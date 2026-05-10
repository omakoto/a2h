package main

import (
	"flag"
	"fmt"
	"github.com/omakoto/a2h/engine"
	"github.com/omakoto/bashcomp"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `a2h - Convert ANSI escape sequences to HTML
  https://github.com/omakoto/a2h

Reads from files or stdin and writes a self-contained HTML page to stdout,
preserving colors and text attributes (bold, italic, underline, etc.).

Usage:
  a2h [flags] [file ...]

Examples:
  some-command 2>&1 | a2h > output.html
  a2h output.txt > output.html
  a2h -bg-color '#1e1e1e' -font-size 10pt output.txt > output.html

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()
	bashcomp.HandleBashCompletion()

	n := engine.NewConverter(os.Stdout)
	n.Convert(flag.Args())
}
