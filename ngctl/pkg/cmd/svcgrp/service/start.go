package service

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var StartServiceCmd = &cobra.Command{
	Use:   "start-service",
	Short: "Start service",
	Long:  "Start service",
	RunE: func(cmd *cobra.Command, args []string) error {
		svcgrpName, err := common.GetResourceName(args)
		if err != nil {
			return err
		}
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(serviceFlags.configFile)
		if err != nil {
			return err
		}
		instances, err := c.GetServiceInstances(svcgrpName,
			serviceFlags.serviceType, serviceFlags.host)
		if err != nil {
			return err
		}
		hosts := make([]string, 0)
		for _, inst := range instances {
			hosts = append(hosts, inst.Host)
		}
		if err := job.NgadmJob.ServiceOperation(c, hosts, serviceFlags.serviceType, "start"); err != nil {
			return common.NgctlError("Failed to start service", err.Error())
		}
		fmt.Fprintf(common.MetaOutput, "Start service successfully.\n")
		return nil
	},
}

func init() {
	StartServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	StartServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "host")
	StartServiceCmd.Flags().StringVarP(&serviceFlags.configFile, "config", "f", "", "config file for ngctl")
	StartServiceCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp start-service <svcgrp-name> [flags]"))
}
