package builtins

import (
	"fmt"
	"io"
	"strings"
)

func Echo(args []string, writer io.Writer) {
	_, err := io.WriteString(writer, strings.Join(args, " ")+"\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}
