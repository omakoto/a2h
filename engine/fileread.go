package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadFilesFromArgs(args []string, f func(line []byte) bool) {
	if args == nil {
		args = os.Args[1:]
	}
	if len(args) == 0 {
		args = []string{"-"}
	}
	readFiles(args, f)
}

func readFiles(names []string, f func(line []byte) bool) {
	for _, name := range names {
		openAndReadFile(name, f)
	}
}

func openAndReadFile(name string, f func(line []byte) bool) {
	var file *os.File
	var err error

	if name == "-" {
		file = os.Stdin
	} else {
		file, err = os.OpenFile(name, os.O_RDONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s: %s\n", name, err)
			return
		}
		defer file.Close()

		fstat, err := file.Stat()
		check(err, "Stat failed")
		if fstat.Mode().IsDir() {
			fmt.Fprintf(os.Stderr, "Skipping directory %s...\n", name)
			return
		}
	}

	readFile(file, f)
}

func readFile(file *os.File, f func(line []byte) bool) {
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if line != nil {
			f(line)
		}
		if err == io.EOF {
			return
		}
		check(err, "ReadBytes failed")
	}
}
