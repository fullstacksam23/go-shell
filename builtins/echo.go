package builtins

import (
	"fmt"
	"io"
	"strings"
)

func Echo(args []string, stdoutWriter io.Writer) {
	_, err := io.WriteString(stdoutWriter, strings.Join(args, " ")+"\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}
