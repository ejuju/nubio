package main

import (
	"os"

	"github.com/ejuju/nubio/pkg/nuage"
)

func main() {
	exitcode := nuage.Run(os.Args[1:]...)
	os.Exit(exitcode)
}
