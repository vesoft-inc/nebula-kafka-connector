package metad_admin

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/host_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a metad",
	Long:  `Create a metad with the metad install and started on a host`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if metadFlags.configFile == "" {
			return common.NgctlError("config file is empty", "")
		}
		if metadFlags.withInstall {
			return common.NgctlError("withInstall must be set to true", "")
		}
		err := common.CheckInConfigFile(metadFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in config file", err.Error())
		}
		if err := host_admin.InstallOnHost(nil, true); err != nil {
			return common.NgctlError("Failed to install on the host", err.Error())
		}
		dstConfigFilePath := common.ConfigSpec.InstallPath + "srvgrp/etc/"
		// The addrs of all metad hosts as configured
		hosts := common.ConfigSpec.Spec.Metad.Hosts
		metaAddrs := []common.IPAndPort{}
		for _, host := range hosts {
			metaAddrs = append(metaAddrs, common.IPAndPort{IP: host.IP, Port: host.Port})
		}
		// prepare config files
		currentFolder, err := os.Getwd()
		if err != nil {
			return common.NgctlError("Failed to get current folder", err.Error())
		}
		currentFolder = strings.TrimSuffix(currentFolder, "/")
		if metadFlags.serviceConfigFile == "" {
			metadFlags.serviceConfigFile = currentFolder + "/nebula-metad.conf.default"
			err := common.GenerateDefaultConfigFile(metadFlags.serviceConfigFile, "metad")
			if err != nil {
				return common.NgctlError("Failed to generate default config file", err.Error())
			}
		}
		// Generate the config file for metad
		for _, host := range hosts {
			localConfigFilePath := currentFolder + "/nebula-metad.conf"
			backupLocalConfigFilePath := currentFolder + "/nebula-metad-" + host.IP + "-" + fmt.Sprint(host.Port) + ".conf.bak"
			if err := common.GenerateMetadConfigFile(metadFlags.serviceConfigFile, localConfigFilePath, common.IPAndPort{IP: host.IP, Port: host.Port}, metaAddrs); err != nil {
				return common.NgctlError("Failed to generate metad config file", err.Error())
			}
			if err := common.BackupFile(localConfigFilePath, backupLocalConfigFilePath); err != nil {
				fmt.Println("Failed to backup the config file: ", err.Error())
			}
			job := runner.NewJob("upload-file")
			// dst_ip := host.IP
			dst_agent := host.Agent.Host
			err = job.Run("upload-file", map[string]any{
				"file_path": []string{localConfigFilePath},
				"dst_path":  []string{dstConfigFilePath},
				"host":      []string{dst_agent},
			}, &common.ConfigSpec)
			if err != nil {
				return common.NgctlError("Failed to upload", err.Error())
			}
			// Start service
			err = service_admin.ServiceOperation(host.Agent, "metad", common.ConfigSpec.InstallPath, "start")
			if err != nil {
				return common.NgctlError("Failed to start service with the agent"+host.Agent.Host, err.Error())
			}
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&metadFlags.configFile, "config", "f", "", "The config file for ngctl to create the metad")
	createCmd.Flags().StringVarP(&metadFlags.serviceConfigFile, "service_config_file", "F", "", "config file for the service to start")
	createCmd.Flags().BoolVar(&metadFlags.withInstall, "with_install", false, "Install and start the metad on the host")
	createCmd.Flags().Lookup("with_install").NoOptDefVal = "false"
}
