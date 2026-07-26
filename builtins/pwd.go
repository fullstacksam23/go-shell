package builtins

import (
	"fmt"
	"io"
	"os"
)

func Pwd(writer io.Writer) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = io.WriteString(writer, wd+"\n")
	if err != nil {
		fmt.Println(err)
	}
}
