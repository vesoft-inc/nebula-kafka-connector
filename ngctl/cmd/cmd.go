package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/down"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/host_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/login"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/logout"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/metad_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/passwd"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/svcgrp_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/up"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/user_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/version"
)

var RootCmd = &cobra.Command{
	Use:   "ngctl",
	Short: "Execute ngctl commands in cli mode",
	Long:  "Execute ngctl commands in cli mode. Use 'ngctl -h' to see usage",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// cmds communicating with the metad to manage the metad
	RootCmd.AddCommand(metad_admin.MetadCmd)
	// cmds communicating with the metad to manage a specific svcgrp
	RootCmd.AddCommand(svcgrp_admin.SvcgrpCmd)
	// cmds communicating with a agent to mange hosts
	RootCmd.AddCommand(host_admin.HostCmd)
	// cmds communicating with a agent to manage services on a host
	RootCmd.AddCommand(service_admin.ServiceAdminCmd)
	// cmds to show the version of ngctl
	RootCmd.AddCommand(version.VersionCmd)

	RootCmd.AddCommand(up.UpCmd)
	RootCmd.AddCommand(down.DownCmd)
	RootCmd.AddCommand(login.LoginCmd)
	RootCmd.AddCommand(logout.LogoutCmd)
	RootCmd.AddCommand(passwd.PasswdCmd)
	RootCmd.AddCommand(user_admin.UserCmd)

	RootCmd.PersistentFlags().StringVarP(&common.CachePath, "tokenFile", "", "", "token file path")
}
