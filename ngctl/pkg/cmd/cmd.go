package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/auth"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/down"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/metad"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/up"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/user"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/version"
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
	RootCmd.AddCommand(svcgrp.SvcgrpCmd)
	RootCmd.AddCommand(user.UserCmd)
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(metad.MetadCmd)
	RootCmd.AddCommand(up.UpCmd)
	RootCmd.AddCommand(down.DownCmd)
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.PersistentFlags().StringVarP(&common.CachePath, "tokenFile", "", "", "token file path")
}
