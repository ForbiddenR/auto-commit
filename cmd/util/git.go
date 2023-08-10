package util

import (
	"os/exec"
)

func GetVariableFromGit(variable string) string {
	out, err := exec.Command("git", "config", "--get", variable).Output()
	if err != nil {
		panic(err)
	}

	return string(out)[:len(out)-1]
}
