package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func executor(command string, args []string) {

	isExecutable, _ := pathutil.FindExecutable(command)
	if !isExecutable {
		fmt.Println(command + ": not found")
		return
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
