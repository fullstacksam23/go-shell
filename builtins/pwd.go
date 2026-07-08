package builtins

import (
	"fmt"
	"os"
)

func Pwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wd)
}
