package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/host_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/metad_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/srvgrp_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:   "ngctl",
	Short: "Execute ngctl commands in cli mode.",
	Long:  `Execute ngctl commands in cli mode. Use 'ngctl -h' to see usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

func docGen() {
	err := doc.GenMarkdownTree(rootCmd, "doc")
	if err != nil {
		panic(err)
	}
}

func init() {
	// cmds communicating with the metad to manage the metad
	rootCmd.AddCommand(metad_admin.SupermetadCmd)
	// cmds communicating with the metad to manage a specific srvgrp
	rootCmd.AddCommand(srvgrp_admin.SrvgrpCmd)
	// cmds communicating with a agent to mange hosts
	rootCmd.AddCommand(host_admin.HostCmd)
	// cmds communicating with a agent to manage services on a host
	rootCmd.AddCommand(service_admin.ServiceAdminCmd)
	// cmds to show the version of ngctl
	rootCmd.AddCommand(version.VersionCmd)

	rootCmd.PersistentFlags().StringVarP(&common.CachePath, "tokenFile", "", "", "token file path")

	// docGen()
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
