package host

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var AddHostCmd = &cobra.Command{
	Use:   "add-host",
	Short: "Add a host into a svcgrp",
	Long:  "Add a host into a svcgrp. A host is identified by its IP address. The port of the deployed agent is also needed",
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

		req := meta.NewAddHostReq(hostFlags.host, svcgrpName, hostFlags.agentPort)
		if err := common.MetaClient.AddHost(req); err != nil {
			return common.NgctlError("Add host failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Add hosts successfully.")
		return nil
	},
}

func init() {
	// add host
	AddHostCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp add-host <svcgrp_name> [flags]"))
	AddHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a svcgrp")
	AddHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 6688, "agent port")
}
