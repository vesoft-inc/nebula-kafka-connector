package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/cmd"
)

func main() {
	var rootCmd = &cobra.Command{Use: "ngadm"}
	cmd.AddConfigfileFlag(rootCmd)
	cmd.RegisteInstallCmd(rootCmd)
	cmd.RegisteVerifyCmd(rootCmd)
	cmd.RegisteUninstallCmd(rootCmd)
	cmd.RegisteOperationCmd(rootCmd)
	cmd.RegisteStatusCmd(rootCmd)
	cmd.RegisteApplyCmd(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
