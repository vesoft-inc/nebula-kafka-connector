package host_admin

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	ngadm_tasks "github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type hostFlagsType struct {
	host          string
	clusterName   string
	agentPort     uint32
	withInstall   bool
	withUninstall bool
	configFile    string
}

var hostFlags hostFlagsType

var HostCmd = &cobra.Command{
	Use:   "host",
	Short: `Process host command`,
	Long:  `Execute host command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func InstallOnHost(hosts []common.IPAndPort, force bool) (err error) {
	var jobName = "install-host"
	log.Print("install nebulagraph cluster.")
	job := runner.NewJob(jobName)
	selectedHosts := make([]string, len(hosts))
	if hosts != nil {
		for i, host := range hosts {
			selectedHosts[i] = host.IP
		}
	}
	err = job.Run(jobName, map[string]any{
		"installAll":   hosts == nil,
		"force":        force,
		"selectedHost": selectedHosts,
	}, &common.ConfigSpec)
	return err
}

// drain, if delete the cluster data
func UninstallOnHost(hosts []common.IPAndPort) (err error) {
	var jobName = "uninstall-host"
	if err != nil {
		log.Println("parse config file failed: ", err)
		return
	}
	selectedHosts := make([]string, len(hosts))
	if hosts != nil {
		for i, host := range hosts {
			selectedHosts[i] = host.IP
		}
	}
	job := runner.NewJob(jobName)
	err = job.Run(jobName, map[string]any{
		"uninstallAll": hosts == nil,
		"selectedHost": selectedHosts,
	}, &common.ConfigSpec)
	if err != nil {
		log.Println("uninstall failed: ", err)
	}
	return
}

func init() {
	HostCmd.AddCommand(addHostCmd)
	HostCmd.AddCommand(dropHostCmd)
	HostCmd.AddCommand(showHostsCmd)

	ngadm_tasks.Init()
}
