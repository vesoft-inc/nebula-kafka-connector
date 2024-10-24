package svcgrp_admin

import (
	"github.com/spf13/cobra"
)

type svcgrpFlagsType struct {
	svcgrpName string
	replicas   int
	force      bool
	owner      string
	output     string
}

var svcgrpFlags svcgrpFlagsType

var SvcgrpCmd = &cobra.Command{
	Use:   "svcgrp",
	Short: "Run commands managing a service group",
	Long:  "A service group may contain multiple graphd and storaged services. Users could run commands to create, initialize, alternate, show, backup or drop a svcgrp",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SvcgrpCmd.PersistentFlags().StringVarP(&svcgrpFlags.svcgrpName, "svcgrp", "s", "", "Service group name")
	SvcgrpCmd.AddCommand(createsvcgrpCmd)
	SvcgrpCmd.AddCommand(initsvcgrpCmd)
	SvcgrpCmd.AddCommand(showsvcgrpCmd)
	//disable br
	// svcgrpCmd.AddCommand(brCmd)
	SvcgrpCmd.AddCommand(dropsvcgrpCmd)
	SvcgrpCmd.AddCommand(altersvcgrpCmd)

}
