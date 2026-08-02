package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/builtins"
)

func processCommand(input string) bool {
	items := tokenize(strings.TrimSpace(input))

	cmd, err := parse(items)
	if err != nil {
		fmt.Fprintln(crlfWriter{os.Stderr}, err)
		return false
	}

	background := false
	if len(cmd.Args) > 0 && cmd.Args[len(cmd.Args)-1] == "&" {
		cmd.Args = cmd.Args[:len(cmd.Args)-1]
		background = true
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
	case "jobs":
		builtins.Jobs()
	case "exit":
		return true //return true only when we want to exit shell
	default:

		if err := executor(cmd.Command, cmd.Args, background, stdoutWriter, stderrWriter); err != nil {
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
