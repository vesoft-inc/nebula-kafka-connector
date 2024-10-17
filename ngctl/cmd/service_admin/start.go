package service_admin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

func prepareConfigFile(serviceType string, installPath string, host common.IPAndPort, agent types.Agent, configFileTemplate string) error {
	hosts := common.ConfigSpec.Spec.Metad.Hosts
	metaAddrs := []common.IPAndPort{}
	for _, host := range hosts {
		metaAddrs = append(metaAddrs, common.IPAndPort{IP: host.IP, Port: host.Port})
	}

	currentFolder, err := os.Getwd()
	if err != nil {
		return common.NgctlError("Failed to get current folder", err.Error())
	}
	currentFolder = strings.TrimSuffix(currentFolder, "/")
	// the local config file to be rendered
	localCongfigFilePath := currentFolder + "/nebula-" + serviceType + ".conf"
	// the backup config file of the above
	backupLocalConfigFilePath := currentFolder + "/nebula-" + serviceType + fmt.Sprintf("-%s-%s.conf", host.IP, host.Port) + ".bak"
	// the final config file path for the service to be started in the remote host
	dstConfigFilePath := installPath + "srvgrp/etc/"

	if configFileTemplate == "" {
		configFileTemplate = localCongfigFilePath + ".default"
		err := common.GenerateDefaultConfigFile(configFileTemplate, serviceType)
		if err != nil {
			return common.NgctlError("Failed to generate default config file", err.Error())
		}
	}
	// use the external config file template to generate service config files
	err = common.GenerateConfigFile(serviceType, configFileTemplate, localCongfigFilePath, common.IPAndPort{IP: host.IP, Port: host.Port}, metaAddrs)
	if err != nil {
		return common.NgctlError("Failed to generate metad config file", err.Error())
	}
	if err := common.BackupFile(localCongfigFilePath, backupLocalConfigFilePath); err != nil {
		fmt.Println("Failed to backup the config file: ", err.Error())
	}
	job := runner.NewJob("upload-file")
	// dst_ip := host.IP
	dst_agent := agent.Host
	err = job.Run("upload-file", map[string]any{
		"file_path": []string{localCongfigFilePath},
		"dst_path":  []string{dstConfigFilePath},
		"host":      []string{dst_agent},
	}, &common.ConfigSpec)
	if err != nil {
		return common.NgctlError("Failed to upload the config file "+localCongfigFilePath+" to host: "+host.IP, err.Error())
	}
	return nil
}

func ServiceOperation(agent types.Agent, serviceType string, installPath string, operation string) (err error) {
	if operation != "start" && operation != "stop" && operation != "restart" && operation != "status" {
		return common.NgctlError("invalid operation on service", "")
	}
	cmd := exec.Command(installPath+"srvgrp/scripts/nebula.service", operation, serviceType)
	log.Printf("exec cmd: %v on %s", cmd, agent.Host)
	job := runner.NewJob(fmt.Sprintf("%s_service", operation))
	fmt.Println(cmd.String())
	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "connect", Params: &tasks.ConnectParams{Host: agent.Host}},
			{Type: "shell", Params: &tasks.ShellParams{
				Host:    agent.Host,
				Command: cmd.String(),
				Sudo:    false,
				CmdID:   fmt.Sprintf("%s_service_cmd", operation),
			}},
		},
	}
	err = job.RunWorkflow(workflow)
	return err
}

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: `Start a service on a host.`,
	Long:  `Start a service on a host.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := ServiceFlags
		if flags.port < 0 {
			return common.NgctlError("no valid port provided", "")
		}
		if flags.serviceType == "" || (flags.serviceType != "graphd" && flags.serviceType != "storaged") {
			return common.NgctlError("invalid service type", "")
		}
		if flags.host == "" {
			return common.NgctlError("no valid host provided", "")
		}
		if flags.port == 0 {
			return common.NgctlError("port is invalid", "")
		}
		err := common.CheckInConfigFile(flags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in the config file for the install path", err.Error())
		}
		port_string := fmt.Sprintf("%d", flags.port)
		agent, err := common.GetAgentForHost(flags.host)
		if err != nil {
			return common.NgctlError("Failed to get the agent for the service", err.Error())
		}
		err = prepareConfigFile(flags.serviceType, common.ConfigSpec.InstallPath, common.IPAndPort{IP: flags.host, Port: port_string}, agent, flags.serviceConfigFile)
		if err != nil {
			return common.NgctlError("Failed to prepare config file for the service", err.Error())
		}

		// start the service
		err = ServiceOperation(agent, flags.serviceType, common.ConfigSpec.InstallPath, "start")
		if err != nil {
			return common.NgctlError("Failed to start service", err.Error())
		}
		return nil
	},
}

func init() {
	startServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	startServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "host")
	startServiceCmd.Flags().Int32VarP(&ServiceFlags.port, "port", "P", -1, "port")
	startServiceCmd.Flags().StringVarP(&ServiceFlags.configFile, "config", "f", "", "config file for ngctl")
	startServiceCmd.Flags().StringVarP(&ServiceFlags.serviceConfigFile, "service_config_file", "F", "", "config file for the service to start")
	startServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "config")
}
