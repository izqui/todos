package main

import (
	"bufio"
	"errors"
	"github.com/andrew-d/go-termutil"
	"io/ioutil"
	"os"
	"strings"
)

func ReadFileLines(path string) ([]string, error) {

	f, err := os.OpenFile(path, os.O_CREATE, 0660)
	defer f.Close()

	if err != nil {

		return nil, err
	}

	return ReadLinesFromFile(f)

}

func ReadLinesFromFile(file *os.File) ([]string, error) {

	sc := bufio.NewScanner(file)
	lines := []string{}

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines, sc.Err()
}

func ReadStdin() ([]string, error) {

	if termutil.Isatty(os.Stdin.Fd()) {
		return nil, errors.New("Stdin is empty")
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {

		return nil, err
	}

	str := strings.TrimSuffix(string(bs), "\n")
	return strings.Split(str, "\n"), nil
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
