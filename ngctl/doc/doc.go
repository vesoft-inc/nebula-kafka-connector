package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, ".")
	if err != nil {
		panic(err)
	}
}
