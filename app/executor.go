package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/builtins"
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
		pid := cmd.Process.Pid
		full_command := command + " " + strings.Join(args, " ") + " &"

		job := builtins.AddJob(pid, full_command)
		fmt.Fprintf(stdoutWriter, "[%d] %d\n", job.JobID, pid)
		go func() {
			err := cmd.Wait()
			if err != nil {
				job.Status = "Failed"
				return
			}
			job.Status = "Done"
		}()
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
