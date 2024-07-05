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

func InstallOnHost(force bool) (err error) {
	var jobName = "install-host"
	log.Print("install nebulagraph cluster.")
	job := runner.NewJob(jobName)
	exit := make(chan error)
	go func() {
		err = job.Run(jobName, map[string]any{
			"force": force,
		}, &common.ConfigSpec)
		if err != nil {
			exit <- err
			return
		}
		exit <- err
	}()
	return <-exit
}

// drain, if delete the cluster data
func UninstallOnHost() (err error) {
	var jobName = "uninstall-host"
	if err != nil {
		log.Println("parse config file failed: ", err)
		return
	}
	job := runner.NewJob(jobName)
	err = job.Run(jobName, map[string]any{}, &common.ConfigSpec)
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
