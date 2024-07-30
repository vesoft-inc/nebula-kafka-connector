package service_admin

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: `Add a service into a cluster.`,
	Long:  `Add a service into a cluster.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := ServiceFlags
		readFromConfig := false
		if flags.serviceType == "" {
			readFromConfig = true
		}
		if flags.clusterName == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		if flags.host == "" {
			readFromConfig = true
		}
		if readFromConfig {
			if flags.configFile == "" {
				return common.NgctlError("Neither a valid service nor a config file is provided", "")
			} else {
				configError := common.CheckInConfigFile(flags.configFile)
				if configError != nil {
					return common.NgctlError("Error in config file", configError.Error())
				}
			}
		}
		serviceList, err := common.DeriveServiceList(common.IPAndPort{IP: flags.host, Port: fmt.Sprintf("%d", flags.port), ServiceType: flags.serviceType}, flags.clusterName)
		// Prepare the resource request
		serviceResource := common.ResourceInfo{
			ResourceType:        "services",
			OperationOnResource: "add",
			ResourceList:        make([]common.IPAndPort, 0),
			ClusterName:         flags.clusterName,
		}
		for _, service := range serviceList {
			serviceResource.ResourceList = append(serviceResource.ResourceList, service)
		}
		if flags.host == "" {
			serviceResource, err = common.ConfirmResourceList(serviceResource)
			if err != nil {
				return common.NgctlError("Failed to confirm resource list", err.Error())
			}
		}
		for _, service := range serviceResource.ResourceList {
			port, _ := strconv.Atoi(service.Port)
			var serviceType meta.ServiceType
			if service.ServiceType == "graphd" {
				serviceType = meta.ServiceTypeGraphd
			} else if service.ServiceType == "storaged" {
				serviceType = meta.ServiceTypeStoraged
			} else {
				return common.NgctlError("Invalid service type: "+service.ServiceType, "")
			}
			req := meta.NewAddServiceReq(service.IP, uint32(port), serviceType, flags.clusterName)
			if err := common.MetaClient.AddService(req); err != nil {
				return common.NgctlError("Add service failed", err.Error())
			}
		}
		fmt.Fprintln(common.MetaOutput, "Add services successfully.")
		return nil
	},
}

func init() {
	addServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "service host")
	addServiceCmd.Flags().Uint32VarP(&ServiceFlags.port, "port", "P", 0, "service port")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.clusterName, "cluster", "c", "", "cluster name")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.configFile, "config", "f", "", "config file")
}
