package util

import (
	"os/exec"
	"strings"
)

func GetVariableFromGit(variable string) string {
	out, err := exec.Command("git", "config", "--get", variable).Output()
	if err != nil {
		panic(err)
	}

	return string(out)[:len(out)-1]
}

func GetDiffFiles() {
	out, err := exec.Command("git", "diff", "--name-only").Output()
	if err != nil {
		panic(err)
	}

	println(strings.Split(string(out)[:len(out)-1], "\n")[0])
}