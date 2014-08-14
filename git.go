package main

import (
	"github.com/izqui/functional"
	"os/exec"
	"strings"
)

func GitDirectoryRoot() (string, error) {

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	res, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(res), "\n"), nil
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

func GitGetOwnerRepoFromRepository() (owner, repo string, err error) {

	cmd := exec.Command("git", "ls-remote", "--get-url")
	res, err := cmd.Output()

	if err != nil {
		return "", "", err
	}

	remote := string(res)
	remote = strings.TrimSuffix(remote, ".git\n")
	remote = strings.TrimPrefix(remote, "git@github.com:")
	remote = strings.TrimPrefix(remote, "https://github.com/")

	owner, repo = func(arr []string) (a, b string) { return arr[0], arr[1] }(strings.Split(remote, "/"))
	return
}

func GitAdd(add string) error {

	cmd := exec.Command("git", "add", add)
	_, err := cmd.Output()

	return err
}

func GitBranch() (string, error) {
	c := strings.Split("rev-parse --abbrev-ref HEAD", " ")
	cmd := exec.Command("git", c...)

	res, err := cmd.Output()

	arr := strings.Split(string(res), "\n")
	return arr[0], err
}

func SetupGitHook(path string) {

	bash := "#!/bin/bash"
	script := "git diff --cached --name-only | todos work"

	lns, err := ReadFileLines(path)
	logOnError(err)

	if len(lns) == 0 {
		lns = append(lns, bash)
	}

	//Filter existing script line
	lns = functional.Filter(func(a string) bool { return a != script }, lns).([]string)
	lns = append(lns, script)

	logOnError(WriteFileLines(path, lns, true))
}
