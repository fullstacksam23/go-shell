package main

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func executor(command string, args []string, writer io.Writer) error {

	isExecutable, _ := pathutil.FindExecutable(command)
	if !isExecutable {
		return errors.New(command + ": not found")
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = writer

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil // suppress exit status
		}
		return err
	}
	return nil
}
