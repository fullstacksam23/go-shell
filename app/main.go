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
	fmt.Print("$ ")

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
			matches := autocompleteSuffix(string(line))

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
				if lastWasTab {

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
		fmt.Println(err)
		return false
	}

	stdoutWriter := io.Writer(os.Stdout)
	stderrWriter := io.Writer(os.Stderr)

	var stdoutFile *os.File
	var stderrFile *os.File

	if cmd.StdoutRedirect {
		stdoutFile, err = os.Create(cmd.StdoutFile)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			return false
		}
		stdoutWriter = stdoutFile
	}
	if cmd.StdErrRedirect {
		stderrFile, err = os.Create(cmd.StderrFile)
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
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
	case "exit":
		return true
	default:

		if err := executor(cmd.Command, cmd.Args, stdoutWriter, stderrWriter); err != nil {
			fmt.Println(err)
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
