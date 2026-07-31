package builtins

import (
	"io"
	"os/exec"
	"strings"
)

type CompletionSpec struct {
	Script  string
	Command string
}

var completionStore = map[string]CompletionSpec{}

func Complete(args []string, stdout, stderr io.Writer) {
	if len(args) < 2 {
		stderr.Write([]byte("complete: not enough arguments\n"))
		return
	}

	flag := args[0]
	switch flag {
	case "-p":
		cmd := args[1]
		spec, e := completionStore[cmd]
		if e {
			msg := []byte("complete -C '" + spec.Script + "' " + spec.Command + "\n")
			stdout.Write(msg)
		} else {
			msg := []byte("complete: " + args[1] + ": no completion specification\n")
			stderr.Write(msg)
		}
	case "-C":
		if len(args) < 3 {
			stderr.Write([]byte("complete: not enough arguments\n"))
			return
		}
		script := args[1]
		cmd := args[2]
		spec := CompletionSpec{
			Script:  script,
			Command: cmd,
		}
		completionStore[cmd] = spec
	}
}

func CompletionAvailable(cmd string) bool {
	_, e := completionStore[cmd]
	return e
}
func GetCompletion(cmd string) string {
	spec, e := completionStore[cmd]
	if !e {
		return ""
	}

	command := exec.Command(spec.Script)
	out, err := command.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
