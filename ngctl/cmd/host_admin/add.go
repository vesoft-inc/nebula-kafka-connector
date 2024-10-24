package host_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var addHostCmd = &cobra.Command{
	Use:   "add",
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
		if err := validateAddorDropFlags(); err != nil {
			return err
		}
		req := meta.NewAddHostReq(hostFlags.host, hostFlags.svcgrpName, hostFlags.agentPort)
		if err := common.MetaClient.AddHost(req); err != nil {
			return common.NgctlError("Add host failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Add hosts successfully.")
		return nil
	},
}

func init() {
	// add host
	addHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a svcgrp")
	addHostCmd.Flags().StringVarP(&hostFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
	addHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 6688, "agent port")
}
