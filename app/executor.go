package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func executor(command string, args []string, background bool, stdoutWriter, stderrWriter io.Writer) error {

	isExecutable, _ := pathutil.FindExecutable(command)
	if !isExecutable {
		return errors.New(command + ": not found")
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderrWriter
	cmd.Stdout = stdoutWriter

	if background {
		err := cmd.Start()
		if err != nil {
			return err
		}
		fmt.Fprintf(stdoutWriter, "[1] %d\n", cmd.Process.Pid)
		go cmd.Wait()
	} else {
		err := cmd.Run()
		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return nil // suppress exit status
			}
			return err
		}
	}
	return nil
}
