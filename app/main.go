package main

import (
	"fmt"
)

func main() {
	var command string

	fmt.Print("$ ")
	_, err := fmt.Scanln(&command)
	if err != nil {
		fmt.Println("Failed to read input")
	}
	fmt.Printf("%s: command not found", command)
}
