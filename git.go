package main

import (
	"os/exec"
	"strings"
)

func GitDirectoryRoot() (string, error) {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	res, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.Split(string(res), "\n")[0], nil
}

func GitDiffFiles() ([]string, error) {

	cmd := exec.Command("git", "diff", "--name-only")
	res, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	arr := strings.Split(string(res), "\n")
	return arr[:len(arr)-1], nil
}
