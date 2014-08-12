package main

import (
	"bufio"
	"os"
)

//TODO: Read file based on regex directly

func ReadFileLines(path string) ([]string, error) {

	f, err := os.OpenFile(path, os.O_CREATE, 0660)
	defer f.Close()

	if err != nil {

		return nil, err
	}

	sc := bufio.NewScanner(f)
	lines := []string{}

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines, sc.Err()
}

func WriteFileLines(path string, lines []string, exec bool) error {

	var mode os.FileMode
	if exec {
		mode = 0755
	} else {
		mode = 0660
	}

	f, err := os.OpenFile(path, os.O_RDWR, mode)
	defer f.Close()

	if err != nil {

		return err
	}

	sc := bufio.NewWriter(f)

	for _, l := range lines {
		sc.WriteString(l + "\n")
	}

	os.Chmod(path, mode)
	return sc.Flush()
}
