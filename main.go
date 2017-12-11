package main

import (
	"flag"
	"github.com/omakoto/a2h/engine"
	"github.com/omakoto/bashcomp"
	"os"
)

func main() {
	flag.Parse()
	bashcomp.HandleBashCompletion()

	n := engine.NewConverter(os.Stdout)
	n.Convert(flag.Args())
}
