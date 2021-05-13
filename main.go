package main

import (
	"os"

	"github.com/datreeio/datree/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
