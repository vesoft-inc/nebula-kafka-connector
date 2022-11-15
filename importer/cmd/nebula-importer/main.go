package main

import (
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/cmd"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/cmd/util"
)

func main() {
	command := cmd.NewDefaultImporterCommand()
	if err := util.Run(command); err != nil {
		util.CheckErr(err)
	}
}
