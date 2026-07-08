package pathutil

import "os/exec"

func FindExecutable(command string) (bool, string) {
	path, err := exec.LookPath(command)
	if err != nil {
		return false, ""
	}
	return true, path
}
