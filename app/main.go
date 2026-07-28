package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/builtins"
)

func main() {
	SetupAutoComplete()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		}
		if err != nil {
			break
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

		if stdoutFile != nil {
			stdoutFile.Close()
		}
		if stderrFile != nil {
			stderrFile.Close()
		}
	}
}
