package srvgrp_admin

import (
	"github.com/spf13/cobra"
)

type srvgrpFlagsType struct {
	srvgrpName string
	replicas   int
	force      bool
	owner      string
	output     string
}

var srvgrpFlags srvgrpFlagsType

var SrvgrpCmd = &cobra.Command{
	Use:   "srvgrp",
	Short: "Run commands managing a service group",
	Long:  `A service group may contain multiple graphd and storaged services. Users could run commands to create, initialize, alternate, show, backup or drop a srvgrp.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SrvgrpCmd.PersistentFlags().StringVarP(&srvgrpFlags.srvgrpName, "srvgrp", "s", "", "Cluster name")
	SrvgrpCmd.AddCommand(createSrvgrpCmd)
	SrvgrpCmd.AddCommand(initSrvgrpCmd)
	SrvgrpCmd.AddCommand(showSrvgrpCmd)
	showSrvgrpCmd.Flags().StringVarP(&srvgrpFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

	//disable br
	// SrvgrpCmd.AddCommand(brCmd)
	createSrvgrpCmd.Flags().IntVarP(&srvgrpFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createSrvgrpCmd.Flags().StringVarP(&srvgrpFlags.owner, "owner", "o", "", "srvgrp owner")

	SrvgrpCmd.AddCommand(dropSrvgrpCmd)
	dropSrvgrpCmd.Flags().BoolVarP(&srvgrpFlags.force, "force", "f", false, "force drop srvgrp")

	SrvgrpCmd.AddCommand(alterSrvgrpCmd)
	alterSrvgrpCmd.Flags().StringVarP(&srvgrpFlags.owner, "owner", "o", "", "srvgrp owner")
}
