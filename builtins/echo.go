package builtins

import (
	"fmt"
	"os"
	"strings"
)

func Echo(args []string) {
	redirect := -1
	for i, v := range args {
		if v == ">" || v == "1>" {
			redirect = i
			break
		}
	}
	if redirect == -1 {
		fmt.Println(strings.Join(args, " "))
		return
	}

	if redirect+1 >= len(args) {
		fmt.Println("echo: missing file operand")
		return
	}

	output := strings.Join(args[:redirect], " ")
	fileName := args[redirect+1]

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

}
