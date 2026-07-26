package builtins

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/pathutil"
)

func CheckType(args []string, stdoutWriter, stderrWriter io.Writer) {
	if len(args) > 1 {
		io.WriteString(stderrWriter, fmt.Sprintln("Type: command accepts a single parameter"))
		return
	}

	if len(args) != 1 {
		io.WriteString(stderrWriter, fmt.Sprintln("Type: Requires one parameter"))
		return
	}
	command := args[0]

	var data string
	_, ok := BuiltIn[command]
	if ok {
		data = command + " is a shell builtin"
	} else {
		found, path := pathutil.FindExecutable(command)
		if found {
			data = fmt.Sprintf("%s is %s", command, path)
		} else {
			data = command + ": not found"
		}
	}

	_, err := io.WriteString(stdoutWriter, data+"\n")
	if err != nil {
		io.WriteString(stderrWriter, fmt.Sprintln(err))
	}
}
