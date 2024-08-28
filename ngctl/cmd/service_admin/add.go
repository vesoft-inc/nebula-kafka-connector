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
	Short: `Add services into a cluster.`,
	Long: `Add services info a cluster.

Either provides the config file, all services in the config file will be added into the cluster.
Or provides the service info, the service will be added into the cluster.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateServiceFlags(); err != nil {
			return err
		}
		var (
			serviceResource *common.ResourceInfo
			err             error
		)
		if ServiceFlags.configFile == "" {
			serviceResource, err = getServicesDirectly()
		} else {
			serviceResource, err = getServicesWithConfig()
		}
		if err != nil {
			return err
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
			req := meta.NewAddServiceReq(service.IP, uint32(port), serviceType, serviceResource.ClusterName)
			if err := common.MetaClient.AddService(req); err != nil {
				return common.NgctlError("Add service failed", err.Error())
			}
		}
		fmt.Fprintln(common.MetaOutput, "Add services successfully.")
		return nil
	},
}

func getServicesDirectly() (*common.ResourceInfo, error) {
	var flags = ServiceFlags
	// get all services from the command line
	serviceResource := common.ResourceInfo{
		ResourceType:        "services",
		OperationOnResource: "add",
		ResourceList:        make([]common.IPAndPort, 0),
		ClusterName:         flags.clusterName,
	}
	serviceResource.ResourceList = append(serviceResource.ResourceList, common.IPAndPort{
		IP:          flags.host,
		Port:        fmt.Sprintf("%d", flags.port),
		ServiceType: flags.serviceType,
	})
	return &serviceResource, nil
}

func getServicesWithConfig() (*common.ResourceInfo, error) {
	var flags = ServiceFlags
	// get all services from the config file
	serviceList, err := common.DeriveServiceList(flags.clusterName)
	if err != nil {
		return nil, err
	}
	serviceResource := common.ResourceInfo{
		ResourceType:        "services",
		OperationOnResource: "add",
		ResourceList:        make([]common.IPAndPort, 0),
		ClusterName:         flags.clusterName,
	}
	for _, service := range serviceList {
		serviceResource.ResourceList = append(serviceResource.ResourceList, service)
	}
	serviceResource, err = common.ConfirmResourceList(serviceResource)
	if err != nil {
		return nil, common.NgctlError("Failed to confirm resource list", err.Error())
	}
	return &serviceResource, nil

}

func init() {
	addServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "service host")
	addServiceCmd.Flags().Int32VarP(&ServiceFlags.port, "port", "P", -1, "service port")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.clusterName, "cluster", "c", "", "cluster name")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.configFile, "config", "f", "", "config file")
}
