package svcgrp

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp/host"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp/pkg"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp/service"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/svcgrp/zone"
)

var SvcgrpCmd = &cobra.Command{
	Use:   "svcgrp",
	Short: "Run commands managing a service group",
	Long:  "A service group may contain multiple graphd and storaged services. Users could run commands to create, initialize, alternate, show, backup or drop a svcgrp",
}

func init() {
	commonCmds := []*cobra.Command{
		common.CreateCmd,
		common.DropCmd,
		common.AlterCmd,
		common.InitCmd,
		common.ShowCmd,
	}

	hostCmds := []*cobra.Command{
		host.AddHostCmd,
		host.RemoveHostCmd,
		host.ShowHostsCmd,
	}

	pkgCmds := []*cobra.Command{
		pkg.InstallCmd,
		pkg.UninstallCmd,
	}

	serviceCmds := []*cobra.Command{
		service.AddServiceCmd,
		service.DropServiceCmd,
		service.ShowServiceCmd,
		service.StartServiceCmd,
		service.StopServiceCmd,
		service.ConfigServiceCmd,
	}

	zoneCmds := []*cobra.Command{
		zone.AddCmd,
		zone.RemoveCmd,
		zone.ShowCmd,
		zone.RenameCmd,
	}

	cmds := make([]*cobra.Command, 0)
	cmds = append(cmds, commonCmds...)
	cmds = append(cmds, hostCmds...)
	cmds = append(cmds, pkgCmds...)
	cmds = append(cmds, serviceCmds...)
	cmds = append(cmds, zoneCmds...)

	SvcgrpCmd.AddCommand(cmds...)

}
