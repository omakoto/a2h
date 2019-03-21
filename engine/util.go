package engine

import (
	"fmt"
	"os"
)

func check(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprint(os.Stderr, message)
	os.Exit(1)
}
