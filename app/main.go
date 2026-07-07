package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
			fmt.Println(echo(args))
		case "exit":
			return
		default:
			fmt.Printf("%s: command not found\n", command)
		}

	}
}

func echo(args []string) string {
	return strings.Join(args, " ")
}
