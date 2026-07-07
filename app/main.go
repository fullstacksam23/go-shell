package main

import (
	"fmt"
)

func main() {
	var command string

	for {
		fmt.Print("$ ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println("Failed to read input")
			return
		}
		if command == "exit" {
			return
		}
		fmt.Printf("%s: command not found\n", command)
	}
}
