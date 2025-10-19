package main

import (
	"os"

	"github.com/gcamargot/kubeselector/cmd/kubeselector"
)

func main() {
	if err := kubeselector.Execute(); err != nil {
		os.Exit(1)
	}
}
