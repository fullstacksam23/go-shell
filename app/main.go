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

	for {

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		fmt.Print("$ ")
		input, err := ReadLine()

		term.Restore(int(os.Stdin.Fd()), oldState)

		if err != nil {
			fmt.Println(err)
			return
		}

		items := tokenize(strings.TrimSpace(input))

		cmd, err := parse(items)
		if err != nil {
			fmt.Println(err)
			continue
		}

		stdoutWriter := io.Writer(os.Stdout)
		stderrWriter := io.Writer(os.Stderr)

		var stdoutFile *os.File
		var stderrFile *os.File

		if cmd.StdoutRedirect {
			stdoutFile, err = os.Create(cmd.StdoutFile)
			if err != nil {
				fmt.Println(err)
				continue
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
				continue
			}
			stdoutWriter = stdoutFile
		}
		if cmd.StdErrRedirect {
			stderrFile, err = os.Create(cmd.StderrFile)
			if err != nil {
				fmt.Println(err)
				continue
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
				continue
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
			return
		default:

			if err := executor(cmd.Command, cmd.Args, stdoutWriter, stderrWriter); err != nil {
				fmt.Println(err)
			}

		}

		if cmd.StdoutRedirect || cmd.StdoutAppend {
			stdoutFile.Close()
		}
		if cmd.StdErrRedirect {
			stderrFile.Close()
		}
	}
}

func ReadLine() (string, error) {
	var line []byte
	buf := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return "", err
		}

		switch buf[0] {

		case '\r', '\n':
			fmt.Println()
			return string(line), nil

		case '\t':
			// autocomplete
			completed := autocomplete(string(line))

			// erase current line
			fmt.Print("\r$ ")
			fmt.Print(strings.Repeat(" ", len(line)))
			fmt.Print("\r$ ")

			line = []byte(completed)
			fmt.Print(completed)

		case 127: // backspace
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}

		default:
			line = append(line, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}
}
