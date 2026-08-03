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
		full_command := command + " " + strings.Join(args, " ")

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

func executePipeline(commands []*ParsedCommand, stdoutWriter, stderrWriter io.Writer) error {
	cmds := make([]*exec.Cmd, len(commands))

	for i, p := range commands {
		cmds[i] = exec.Command(p.Command, p.Args...)
	}

	var prev io.ReadCloser

	for i, cmd := range cmds {

		// stdin
		if prev != nil {
			cmd.Stdin = prev
		}

		// stdout
		if i != len(cmds)-1 {
			next, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			prev = next
		} else {
			cmd.Stdout = stdoutWriter
		}

		cmd.Stderr = stderrWriter
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		cmd.Wait()
	}
	return nil
}
