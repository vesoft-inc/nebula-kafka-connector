package service_admin

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

func prepareConfigFile(serviceType string, installPath string, host common.IPAndPort, agent types.Agent) error {
	hosts := common.ConfigSpec.Spec.Metad.Hosts
	metaAddrs := []common.IPAndPort{}
	for _, host := range hosts {
		metaAddrs = append(metaAddrs, common.IPAndPort{IP: host.IP, Port: host.Port})
	}
	localCongfigFilePath := installPath + "cluster/etc/nebula-" + serviceType + ".conf"
	backupLocalConfigFilePath := installPath + "cluster/etc/nebula-" + serviceType + fmt.Sprintf("-%s-%s.conf", host.IP, host.Port) + ".conf.bak"
	dstConfigFilePath := installPath + "cluster/etc/"
	err := common.GenerateConfigFile(serviceType, localCongfigFilePath, common.IPAndPort{IP: host.IP, Port: host.Port}, metaAddrs)
	if err != nil {
		return common.NgctlError("Failed to generate metad config file", err.Error())
	}
	if err = common.BackupFile(localCongfigFilePath, backupLocalConfigFilePath); err != nil {
		fmt.Println("Failed to backup the config file: ", err.Error())
	}
	job := runner.NewJob("upload-file")
	exit := make(chan error)
	// dst_ip := host.IP
	dst_agent := agent.Host
	go func() {
		err = job.Run("upload-file", map[string]any{
			"file_path": []string{localCongfigFilePath},
			"dst_path":  []string{dstConfigFilePath},
			"host":      []string{dst_agent},
		}, &common.ConfigSpec)
		if err != nil {
			exit <- err
			return
		}
		exit <- err
	}()
	return nil
}

func ServiceOperation(agent types.Agent, serviceType string, installPath string, operation string) (err error) {
	if operation != "start" && operation != "stop" && operation != "restart" && operation != "status" {
		return common.NgctlError("invalid operation on service", "")
	}
	cmd := exec.Command(installPath+"cluster/scripts/nebula.service", operation, serviceType)
	job := runner.NewJob("start_service")
	workflow := &types.WorkflowSpec{
		Tasks: []*types.TaskSpec{
			{Type: "connect", Params: &tasks.ConnectParams{Host: agent.Host}},
			{Type: "shell", Params: &tasks.ShellParams{
				Host:    agent.Host,
				Command: cmd.String(),
				Sudo:    false,
				CmdID:   "start_service_cmd",
			}},
		},
	}
	err = job.RunWorkflow(workflow)
	return err
}

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: `Start a service on a host`,
	Long:  `Start a service on a host`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := ServiceFlags
		if flags.serviceType == "" || (flags.serviceType != "graphd" && flags.serviceType != "storaged") {
			return common.NgctlError("invalid service type", "")
		}
		if flags.host == "" || common.IsValidIPAddress(flags.host) == false{
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
		// prepare config files for the service
		err = prepareConfigFile(flags.serviceType, common.ConfigSpec.InstallPath, common.IPAndPort{IP: flags.host, Port: port_string}, agent)
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
	startServiceCmd.Flags().Uint32VarP(&ServiceFlags.port, "port", "P", 0, "port")
	startServiceCmd.Flags().StringVarP(&ServiceFlags.configFile, "config", "f", "", "config file path")
	startServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "config")
}
