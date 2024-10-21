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
	Long:  `A service group may contain multiple graphd and storaged services. Users could run commands to create, initialize, alternate, show, backup or drop a svcgrp.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SvcgrpCmd.PersistentFlags().StringVarP(&svcgrpFlags.svcgrpName, "svcgrp", "s", "", "Cluster name")
	SvcgrpCmd.AddCommand(createsvcgrpCmd)
	SvcgrpCmd.AddCommand(initsvcgrpCmd)
	SvcgrpCmd.AddCommand(showsvcgrpCmd)
	showsvcgrpCmd.Flags().StringVarP(&svcgrpFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

	//disable br
	// svcgrpCmd.AddCommand(brCmd)
	createsvcgrpCmd.Flags().IntVarP(&svcgrpFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createsvcgrpCmd.Flags().StringVarP(&svcgrpFlags.owner, "owner", "o", "", "svcgrp owner")

	SvcgrpCmd.AddCommand(dropsvcgrpCmd)
	dropsvcgrpCmd.Flags().BoolVarP(&svcgrpFlags.force, "force", "f", false, "force drop svcgrp")

	SvcgrpCmd.AddCommand(altersvcgrpCmd)
	altersvcgrpCmd.Flags().StringVarP(&svcgrpFlags.owner, "owner", "o", "", "svcgrp owner")
}
