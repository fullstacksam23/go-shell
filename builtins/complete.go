package builtins

import "io"

func Complete(args []string, stdout, stderr io.Writer) {
	flag := args[0]
	switch flag {
	case "-p":
		msg := []byte("complete: " + args[1] + ": no completion specification\n")
		stderr.Write(msg)
	}
}
