package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/builtins"
	"golang.org/x/term"
)

func main() {
	SetupAutoComplete()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	line := []rune{}
	lastWasTab := false

	buf := make([]byte, 1)
	stdout := crlfWriter{os.Stdout}
	fmt.Fprint(stdout, "$ ")

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		ch := buf[0]

		switch ch {

		case '\r', '\n':
			fmt.Print("\r\n")

			input := string(line)

			line = line[:0]
			lastWasTab = false

			if processCommand(input) {
				return
			}

			fmt.Print("$ ")

		case '\t':
			prefix, word := splitLastWord(line)
			//there are no spaces so far
			if prefix == "" {

				longestCommonPrefix, matches := autoCompleteCommand(string(word))

				switch len(matches) {

				case 0:
					fmt.Print("\a")
					lastWasTab = false

				case 1:
					suffix := matches[0]
					line = append(line, []rune(suffix)...)
					line = append(line, ' ')
					fmt.Print(suffix + " ")
					lastWasTab = false

				default:
					if len(longestCommonPrefix) > 0 {
						// extend to the common prefix, no trailing space (still ambiguous)
						fmt.Print(string(longestCommonPrefix))
						line = append(line, longestCommonPrefix...)
						lastWasTab = false
					} else if lastWasTab {
						fmt.Print("\r\n")
						for _, m := range matches {
							fmt.Print(string(line) + m + "  ")
						}
						fmt.Print("\b\b\r\n")
						fmt.Print("$ ")
						fmt.Print(string(line))
						lastWasTab = false
					} else {
						fmt.Print("\a")
						lastWasTab = true
					}
				}
				// there are spaces in here
			} else {
				cmd, current, previous := builtins.CompletionContext(line)

				if spec, ok := builtins.GetCompletionSpec(cmd); ok {

					matches, err := builtins.RunCompleter(spec.Script, cmd, current, previous)
					if err == nil {

						if len(matches) == 0 {
							fmt.Print("\a")
							lastWasTab = false
							continue
						}

						if len(matches) == 1 {
							completion := matches[0]

							suffix := completion[len(current):]

							fmt.Print(suffix + " ")

							line = append(line, []rune(suffix)...)
							line = append(line, ' ')

							lastWasTab = false
							continue
						}

						// Multiple matches (same behavior as your existing code)
					}
				}
				//get the path
				path_items := strings.Split(word, "/")
				matches := []string{}
				n := len(path_items)
				if n == 1 {
					matches = autoCompleteFileName(".", word)
				} else {
					matches = autoCompleteFileName(strings.Join(path_items[:n-1], "/"), path_items[n-1])
				}

				if len(matches) == 0 {
					fmt.Print("\a")
					lastWasTab = false
					continue
				} else if len(matches) == 1 {
					suffix := matches[0]
					fmt.Print(suffix)
					line = append(line, []rune(suffix)...)
					//check if it is a file
					if suffix[len(suffix)-1] != '/' {
						fmt.Print(" ")
						line = append(line, ' ')
					}
					lastWasTab = false
				} else {
					if lastWasTab {
						fmt.Print("\r\n")
						for _, m := range matches {
							fmt.Print(string(word) + m + "  ")
						}
						fmt.Print("\b\b\r\n")
						fmt.Print("$ ")
						fmt.Print(string(line))
						lastWasTab = false
					} else {
						lcp := longestCommonPrefixOf(matches)
						if len(lcp) > 0 {
							fmt.Print(lcp)
							line = append(line, []rune(lcp)...)
							lastWasTab = false
						} else {
							fmt.Print("\a")
							lastWasTab = true
						}
					}

				}

			}
		case 127, 8:
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}
			lastWasTab = false

		default:
			line = append(line, rune(ch))
			fmt.Printf("%c", ch)
			lastWasTab = false
		}

	}
}

func processCommand(input string) bool {
	items := tokenize(strings.TrimSpace(input))

	cmd, err := parse(items)
	if err != nil {
		fmt.Fprintln(crlfWriter{os.Stderr}, err)
		return false
	}

	stdoutWriter := io.Writer(crlfWriter{os.Stdout})
	stderrWriter := io.Writer(crlfWriter{os.Stderr})

	var stdoutFile *os.File
	var stderrFile *os.File

	if cmd.StdoutRedirect {
		stdoutFile, err = os.Create(cmd.StdoutFile)
		if err != nil {
			fmt.Fprintln(crlfWriter{os.Stderr}, err)
			return false
		}
		stdoutWriter = stdoutFile
	}
	if cmd.StdoutAppend {
		stdoutFile, err = os.OpenFile(
			cmd.StdoutFile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
		if err != nil {
			fmt.Fprintln(crlfWriter{os.Stderr}, err)
			return false
		}
		stdoutWriter = stdoutFile
	}
	if cmd.StdErrRedirect {
		stderrFile, err = os.Create(cmd.StderrFile)
		if err != nil {
			fmt.Fprintln(crlfWriter{os.Stderr}, err)
			return false
		}
		stderrWriter = stderrFile
	}
	if cmd.StdErrAppend {
		stderrFile, err = os.OpenFile(
			cmd.StderrFile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
		if err != nil {
			fmt.Fprintln(crlfWriter{os.Stderr}, err)
			return false
		}
		stderrWriter = stderrFile
	}

	switch cmd.Command {
	case "echo":
		builtins.Echo(cmd.Args, stdoutWriter)
	case "type":
		builtins.CheckType(cmd.Args, stdoutWriter, stderrWriter)
	case "pwd":
		builtins.Pwd(stdoutWriter)
	case "cd":
		builtins.Cd(cmd.Args, stderrWriter)
	case "complete":
		builtins.Complete(cmd.Args, stdoutWriter, stderrWriter)
	case "exit":
		return true
	default:

		if err := executor(cmd.Command, cmd.Args, stdoutWriter, stderrWriter); err != nil {
			fmt.Fprintln(crlfWriter{os.Stderr}, err)
		}

	}

	if stdoutFile != nil {
		stdoutFile.Close()
	}
	if stderrFile != nil {
		stderrFile.Close()
	}
	return false
}

func splitLastWord(line []rune) (prefix string, word string) {
	s := string(line)
	idx := strings.LastIndex(s, " ")
	if idx == -1 {
		return "", s // no space yet — completing the command itself
	}
	return s[:idx+1], s[idx+1:] // prefix includes the trailing space
}
