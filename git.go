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

	return strings.TrimSuffix(string(res), "\n"), nil
}

func GitDiffFiles() ([]string, error) {

	cmd := exec.Command("git", "diff", "--cached", "--name-only")
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
