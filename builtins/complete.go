package builtins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
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

func RunCompleter(script, cmd, current, previous, line string) ([]string, error) {
	cmdExec := exec.Command(script, cmd, current, previous)
	cmdExec.Env = append(os.Environ(),
		"COMP_LINE="+line,
		fmt.Sprintf("COMP_POINT=%d", len([]byte(line))),
	)
	out, err := cmdExec.Output()
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			matches = append(matches, line)
		}
	}
	slices.Sort(matches)
	return matches, nil
}
func GetCompletionSpec(cmd string) (CompletionSpec, bool) {
	spec, ok := completionStore[cmd]
	return spec, ok
}

func CompletionContext(line []rune) (cmd, current, previous string) {
	fields := strings.Fields(string(line))

	if len(fields) == 0 {
		return "", "", ""
	}

	cmd = fields[0]

	// Cursor is after a space: "git remote "
	if len(line) > 0 && line[len(line)-1] == ' ' {
		if len(fields) >= 2 {
			previous = fields[len(fields)-1]
		}
		current = ""
		return
	}

	// Cursor is inside a word: "git remote set"
	current = fields[len(fields)-1]

	if len(fields) >= 3 {
		previous = fields[len(fields)-2]
	}

	return
}
