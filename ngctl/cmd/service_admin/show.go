package service_admin

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var showServiceCmd = &cobra.Command{
	Use:   "show",
	Short: "Show services in a srvgrp.",
	Long:  `Show services in a srvgrp.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srvgrp := ServiceFlags.srvgrpName
		if srvgrp == "" {
			return common.NgctlError("srvgrp name is empty", "")
		}
		req := meta.NewListServicesReq(srvgrp)

		resp, err := common.MetaClient.ListServices(req)
		if err != nil {
			return common.NgctlError("Show service failed", err.Error())
		}

		header := []string{"Id", "Type", "Host", "Port"}
		data := make([][]string, 0)
		for _, s := range resp.Services {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.Id))
			if s.Type == meta.ServiceTypeGraphd {
				row = append(row, "graphd")
			} else {
				row = append(row, "storaged")
			}
			row = append(row, s.Host)
			row = append(row, fmt.Sprintf("%d", s.Port))
			data = append(data, row)
		}

		// order by service id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		r, err := common.Format(header, data, common.OutputFormatType(ServiceFlags.output))
		if err != nil {
			return common.NgctlError("Show service failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	showServiceCmd.Flags().StringVarP(&ServiceFlags.srvgrpName, "srvgrp", "s", "", "srvgrp name")
	showServiceCmd.MarkFlagRequired("srvgrp")
}
