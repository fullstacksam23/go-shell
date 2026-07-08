package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/builtins"
	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func main() {

	for {
		fmt.Print("$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		items := strings.Fields(input)
		command := items[0]
		args := items[1:]

		switch command {
		case "echo":
			builtins.Echo(args)
		case "type":
			builtins.CheckType(args)
		case "pwd":
			builtins.Pwd()
		case "exit":
			return
		default:
			executor(command, args)
		}

	}
}

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
