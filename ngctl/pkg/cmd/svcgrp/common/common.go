package common

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

// create
// init
// alter
// drop
// show
// reset
var (
	svcgrpName  string
	replicas    int
	force       bool
	owner       string
	output      string
	retry       bool
	routePolicy string
)

var CreateCmd = &cobra.Command{
	Use:   "create <svcgrp-name>",
	Short: "Create a service group in the metad",
	Long:  "Create a service group and register in the metad. The service group will be empty after creation. Users need to add hosts and services into a newly created svcgrp",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return common.NgctlError("ServiceGroup name is empty", "")
		}
		svcgrpName = args[0]

		if replicas == 0 {
			return common.NgctlError("Replica number is invalid", "")
		}
		// do not need to specify zones
		req := meta.NewCreateServiceGroupReq(svcgrpName, replicas, owner, nil)
		if err := common.MetaClient.CreateServiceGroup(req); err != nil {
			return common.NgctlError("Create svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Create svcgrp successfully.")
		return nil
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service group in the metad",
	Long:  "Initialize a service group in the metad",
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

		req := meta.NewInitServiceGroupReq(svcgrpName)

		if err := common.MetaClient.InitServiceGroup(req); err != nil {
			return common.NgctlError("Init svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Init svcgrp successfully.")
		return nil
	},
}

var DropCmd = &cobra.Command{
	Use:   "drop <svcgrp-name>",
	Short: "Drop a service group from the metad",
	Long:  "Drop a service group from the metad",
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

		req := meta.NewDropServiceGroupReq(svcgrpName, force)
		if err := common.MetaClient.DropServiceGroup(req); err != nil {
			return common.NgctlError("Drop svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Drop svcgrp successfully.")
		return nil
	},
}

var AlterCmd = &cobra.Command{
	Use:   "alter <svcgrp-name>",
	Short: "Alter the owner of a svcgrp.",
	Long:  `A service group has an owner registered in the metad. Users can alter it.`,
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

		if owner == "" && routePolicy == "" {
			return common.NgctlError("Owner or route policy is empty", "")
		}
		getReq := meta.NewListServiceGroupsReq(svcgrpName)
		getResp, err := common.MetaClient.ListServiceGroups(getReq)
		if err != nil {
			return common.NgctlError("Alter svcgrp failed", err.Error())
		}
		if len(getResp.ServiceGroups) == 0 {
			return common.NgctlError("Service group not found", "")
		}
		r := getResp.ServiceGroups[0]
		var req *meta.AlterServiceGroupReq

		if owner != "" {
			req = meta.NewAlterServiceGroupReq(r.Name, owner, r.RoutePolicy, r.Retry)
		} else {
			req = meta.NewAlterServiceGroupReq(r.Name, "", routePolicy, retry)
		}
		if err := common.MetaClient.AlterServiceGroup(req); err != nil {
			return common.NgctlError("Alter svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Alter svcgrp successfully.")
		return nil
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show [svcgrp-name]",
	Short: "Show details of a service group",
	Long:  "Show details of a service group",
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
			// if there's no any svcgrp name, just all svcgrp
			_ = err
		}
		req := meta.NewListServiceGroupsReq(svcgrpName)
		resp, err := common.MetaClient.ListServiceGroups(req)
		if err != nil {
			return common.NgctlError("Show svcgrp failed", err.Error())
		}

		header := []string{"Id", "Name", "Replica Factor", "Owner", "Route Policy", "Retry"}
		data := make([][]string, 0)
		for _, s := range resp.ServiceGroups {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.Id))
			row = append(row, fmt.Sprintf("%s", s.Name))
			row = append(row, fmt.Sprintf("%d", s.ReplicaRefactor))
			row = append(row, fmt.Sprintf("%s", s.Owner))
			row = append(row, fmt.Sprintf("%s", s.RoutePolicy))
			row = append(row, fmt.Sprintf("%t", s.Retry))
			data = append(data, row)
		}

		// order by svcgrp id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		// printer.FormatTable(headers []string, data [][]string)
		r, err := common.Format(header, data, common.OutputFormatType(output))
		if err != nil {
			return common.NgctlError("Show svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	CreateCmd.Flags().IntVarP(&replicas, "replica-factor", "r", 3, "replica number, default: 3")
	CreateCmd.Flags().StringVarP(&owner, "owner", "o", "", "svcgrp owner")

	DropCmd.Flags().BoolVarP(&force, "force", "f", false, "force drop svcgrp")

	AlterCmd.Flags().StringVarP(&owner, "owner", "o", "", "svcgrp owner")
	AlterCmd.Flags().BoolVarP(&retry, "retry", "r", false, "")
	AlterCmd.Flags().StringVarP(&routePolicy, "route-policy", "p", "", "route policy")
	AlterCmd.MarkFlagsRequiredTogether("route-policy", "retry")
	AlterCmd.MarkFlagsMutuallyExclusive("owner", "route-policy")

	ShowCmd.Flags().StringVarP(&output, "output", "o", "table", "output format. Allowed values: table, json")

}
