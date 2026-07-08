package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var BuiltIn = map[string]struct{}{
	"echo": {},
	"type": {},
	"exit": {},
}

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
			echo(args)
		case "type":
			checkType(args)
		case "exit":
			return
		default:
			fmt.Printf("%s: command not found\n", command)
		}

	}
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func checkType(args []string) {
	if len(args) > 1 {
		fmt.Println("Type command accepts a single parameter")
		return
	}
	command := args[0]
	_, ok := BuiltIn[command]
	if ok {
		fmt.Println(command + " is a shell builtin")
		return
	}
	found, path := findExecutable(command)
	if found {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Println(command + ": not found")
	}
}

func findExecutable(command string) (bool, string) {
	pathEnv := os.Getenv("PATH")
	extensionsEnv := os.Getenv("PATHEXT")

	directories := filepath.SplitList(pathEnv)
	extensions := filepath.SplitList(extensionsEnv)

	for _, dir := range directories {
		if dir == "" {
			continue
		}
		for _, ext := range extensions {
			filename := command + ext

			fullPath := filepath.Join(dir, filename)

			if isExecutable(fullPath) {
				return true, fullPath
			}
		}
	}
	return false, "not found"
}

func isExecutable(fullPath string) bool {
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	mode := info.Mode()
	return mode&0111 != 0
}
