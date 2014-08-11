package main

import (
	"os/exec"
	"strings"
)

func GitDirectoryRoot() (string, error) {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	res, err := cmd.Output()

	return string(res), err
}

func GitDiffFiles() ([]string, error) {

	cmd := exec.Command("git", "diff", "--show-names")
	res, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return strings.Split(string(res), "\n"), err
}
