package main

import (
	"os"

	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
