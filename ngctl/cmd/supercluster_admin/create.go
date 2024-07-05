package supercluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/host_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a supercluster",
	Long:  `Create a supercluster with the metad install and started on a host`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if superclusterFlags.configFile == "" {
			return common.NgctlError("config file is empty", "")
		}
		if superclusterFlags.withInstall {
			return common.NgctlError("withInstall must be set to true", "")
		}
		err := common.CheckInConfigFile(superclusterFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in config file", err.Error())
		}
		if err := host_admin.InstallOnHost(true); err != nil {
			return common.NgctlError("Failed to install on the host", err.Error())
		}
		dstConfigFilePath := common.ConfigSpec.InstallPath + "cluster/etc/"
		// The addrs of all metad hosts as configured
		hosts := common.ConfigSpec.Spec.Metad.Hosts
		metaAddrs := []common.IPAndPort{}
		for _, host := range hosts {
			metaAddrs = append(metaAddrs, common.IPAndPort{IP: host.IP, Port: host.Port})
		}
		// Generate the config file for metad
		for _, host := range hosts {
			localConfigFilePath := common.ConfigSpec.InstallPath + "cluster/etc/nebula-metad.conf" //
			backupLocalConfigFilePath := common.ConfigSpec.InstallPath + "cluster/etc/nebula-metad" + fmt.Sprintf("-%s-%s.conf", host.IP, host.Port) + ".conf.bak"
			err := common.GenerateMetadConfigFile(localConfigFilePath, common.IPAndPort{IP: host.IP, Port: host.Port}, metaAddrs)
			if err != nil {
				return common.NgctlError("Failed to generate metad config file", err.Error())
			}
			if err := common.BackupFile(localConfigFilePath, backupLocalConfigFilePath); err != nil {
				fmt.Println("Failed to backup the config file: ", err.Error())
			}
			job := runner.NewJob("upload-file")
			exit := make(chan error)
			// dst_ip := host.IP
			dst_agent := host.Agent.Host
			go func() {
				err = job.Run("upload-file", map[string]any{
					"file_path": []string{localConfigFilePath},
					"dst_path":  []string{dstConfigFilePath},
					"host":      []string{dst_agent},
				}, &common.ConfigSpec)
				if err != nil {
					exit <- err
					return
				}
				exit <- err
			}()
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
	createCmd.Flags().StringVar(&superclusterFlags.configFile, "config", "", "The config file to create the supercluster")
	createCmd.Flags().BoolVar(&superclusterFlags.withInstall, "with_install", false, "Install and start the metad on the host")
	createCmd.Flags().Lookup("with_install").NoOptDefVal = "false"
}
