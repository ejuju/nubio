package main

import (
	"os"

	"github.com/ejuju/nubio/pkg/nubio"
)

func main() {
	exitcode := nubio.Run(os.Args[1:]...)
	os.Exit(exitcode)
}
