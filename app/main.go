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

		writer := io.Writer(os.Stdout)
		var file *os.File

		if cmd.StdoutRedirect {
			file, err = os.Create(cmd.StdoutFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
			writer = file
		}

		switch cmd.Command {
		case "echo":
			builtins.Echo(cmd.Args, writer)
		case "type":
			builtins.CheckType(cmd.Args, writer)
		case "pwd":
			builtins.Pwd(writer)
		case "cd":
			builtins.Cd(cmd.Args)
		case "exit":
			return
		default:

			if err := executor(cmd.Command, cmd.Args, writer); err != nil {
				fmt.Println(err)
			}

		}

		if cmd.StdoutRedirect {
			file.Close()
		}

	}
}
