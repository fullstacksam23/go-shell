package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/builtins"
)

func main() {

	for {
		fmt.Print("$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
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
				return
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
