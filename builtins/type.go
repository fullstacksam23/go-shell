package builtins

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func CheckType(args []string) {
	if len(args) > 1 {
		fmt.Println("Type command accepts a single parameter")
		return
	}

	if len(args) != 1 {
		return
	}
	command := args[0]

	_, ok := BuiltIn[command]
	if ok {
		fmt.Println(command + " is a shell builtin")
		return
	}
	found, path := pathutil.FindExecutable(command)
	if found {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Println(command + ": not found")
	}
}
