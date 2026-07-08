package builtins

import (
	"fmt"
	"os"
)

func Cd(args []string) {
	if len(args) != 1 {
		fmt.Println("cd: requires one argument")
		fmt.Println(args)
		fmt.Println(len(args))
		return
	}
	dir := args[0]
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
		return
	}
	if info.IsDir() {
		err = os.Chdir(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
	}
}
