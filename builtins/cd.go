package builtins

import (
	"fmt"
	"io"
	"os"
)

func Cd(args []string, stderrWriter io.Writer) {
	if len(args) != 1 {
		io.WriteString(stderrWriter, "cd: requires one argument\n")
		return
	}
	dir := args[0]
	if dir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			io.WriteString(stderrWriter, fmt.Sprintf("cd: %s: No such file or directory\n", dir))
			return
		}
		dir = home
	}

	info, err := os.Stat(dir)
	if err != nil {
		io.WriteString(stderrWriter, fmt.Sprintf("cd: %s: No such file or directory\n", dir))
		return
	}
	if info.IsDir() {
		err = os.Chdir(dir)
		if err != nil {
			io.WriteString(stderrWriter, fmt.Sprintln(err))
			return
		}
	} else {
		io.WriteString(stderrWriter, fmt.Sprintf("cd: %s: No such file or directory\n", dir))
	}
}
