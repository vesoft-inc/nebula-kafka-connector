package host

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var RemoveHostCmd = &cobra.Command{
	Use:   "remove-host",
	Short: "Remove a host from a svcgrp",
	Long:  "Remove a host from a svcgrp",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		svcgrpName, err := common.GetResourceName(args)
		if err != nil {
			return err
		}
		req := meta.NewDropHostReq(hostFlags.host, svcgrpName)
		if err := common.MetaClient.DropHost(req); err != nil {
			fmt.Fprintln(common.MetaOutput, fmt.Sprintf("Drop host %s failed: %s", hostFlags.host, err.Error()))
			return err
		}

		fmt.Fprintf(common.MetaOutput, "Drop host %s successfully.\n", hostFlags.host)
		return nil
	},
}

func init() {
	// drop host
	RemoveHostCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp remove-host <svcgrp_name> [flags]"))
	RemoveHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be dropped from a svcgrp")
	RemoveHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 6688, "port of the agent on the host")
}
