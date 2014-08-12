package main

import (
	"bufio"
	"os"
)

//TODO: Read file based on regex directly

func ReadFileLines(path string) ([]string, error) {

	f, err := os.Open(path)
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

func WriteFileLines(path string, lines []string) error {

	f, err := os.OpenFile(path, os.O_RDWR, 0660)
	defer f.Close()

	if err != nil {

		return err
	}

	sc := bufio.NewWriter(f)

	for _, l := range lines {
		sc.WriteString(l + "\n")
	}

	return sc.Flush()
}
