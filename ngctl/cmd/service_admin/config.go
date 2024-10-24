package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var configServiceCmd = &cobra.Command{
	Use:   "config",
	Short: "Config service",
	Long:  "Config service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(serviceFlags.configFile)
		if err != nil {
			return err
		}
		instances, err := c.GetServiceInstances(serviceFlags.svcgrpName, serviceFlags.serviceType, serviceFlags.host)
		if err != nil {
			return err
		}
		if len(instances) == 0 {
			return common.NgctlError("no service instance found", "")
		}
		hosts := make([]string, 0)
		for _, inst := range instances {
			hosts = append(hosts, inst.Host)
		}

		if err := job.NgadmJob.UpdateConfig(c, serviceFlags.svcgrpName, hosts, serviceFlags.serviceType); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Config service successfully.\n")
		return nil
	},
}

func init() {
	configServiceCmd.Flags().StringVarP(&serviceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
	configServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	configServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "host")
	configServiceCmd.Flags().StringVarP(&serviceFlags.configFile, "config", "f", "", "config file for ngctl")
}
