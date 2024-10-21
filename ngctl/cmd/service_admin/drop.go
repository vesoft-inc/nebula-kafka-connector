package service_admin

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropServiceCmd = &cobra.Command{
	Use:   "drop",
	Short: `Drop a service from a svcgrp.`,
	Long:  `Drop a service from a svcgrp.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var flags = ServiceFlags
		if err := validateServiceFlags(); err != nil {
			return err
		}
		var (
			serviceList []common.IPAndPort = make([]common.IPAndPort, 0)
			err         error
		)
		if flags.configFile != "" {
			serviceList, err = common.DeriveServiceList(flags.svcgrpName)
		} else {
			serviceList = append(serviceList, common.IPAndPort{
				IP:          flags.host,
				Port:        strconv.Itoa(int(flags.port)),
				ServiceType: flags.serviceType,
			})
		}
		if err != nil {
			return err
		}
		if len(serviceList) == 0 {
			return common.NgctlError("No service to drop.", "")
		}

		for _, service := range serviceList {
			port, err := strconv.Atoi(service.Port)
			if err != nil {
				return err
			}
			var serviceType meta.ServiceType
			if service.ServiceType == "graphd" {
				serviceType = meta.ServiceTypeGraphd
			} else if service.ServiceType == "storaged" {
				serviceType = meta.ServiceTypeStoraged
			} else {
				return common.NgctlError("Invalid service type.", "")
			}
			req := meta.NewDropServiceReq(service.IP, uint32(port), serviceType, flags.svcgrpName)
			if err := common.MetaClient.DropService(req); err != nil {
				return common.NgctlError("Drop service failed", err.Error())
			}
		}
		fmt.Fprintln(common.MetaOutput, "Drop services successfully.")
		return nil
	},
}

func init() {
	dropServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	dropServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "service host")
	dropServiceCmd.Flags().Int32VarP(&ServiceFlags.port, "port", "P", -1, "service port")
	dropServiceCmd.Flags().StringVarP(&ServiceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
}
