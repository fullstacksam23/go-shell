package main

import (
	"bufio"
	"fmt"
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

		items := parse(strings.TrimSpace(input))
		command := items[0]
		args := items[1:]

		switch command {
		case "echo":
			builtins.Echo(args)
		case "type":
			builtins.CheckType(args)
		case "pwd":
			builtins.Pwd()
		case "cd":
			builtins.Cd(args)
		case "exit":
			return
		default:
			executor(command, args)
		}

	}
}
