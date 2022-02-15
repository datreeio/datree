package main

import (
	"errors"
	"os"

	"github.com/datreeio/datree/cmd"
	"github.com/datreeio/datree/cmd/test"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if errors.Is(err, test.InvocationError) {
			os.Exit(9)
		}
		os.Exit(1)
	}
}
