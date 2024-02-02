package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/cmd"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
)

func main() {
	tasks.Init()
	var rootCmd = &cobra.Command{Use: "ngadmin"}
	cmd.AddConfigfileFlag(rootCmd)
	cmd.RegisteInstallCmd(rootCmd)
	cmd.RegisteVerifyCmd(rootCmd)
	cmd.RegisteUninstallCmd(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
