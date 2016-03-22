package main

import (
	"flag"
	"github.com/omakoto/a2h/engine"
	"github.com/omakoto/bashcomp"
)

func main() {
	flag.Parse()
	bashcomp.HandleBashCompletion()

	n := engine.NewConverter()
	n.Convert()
}
