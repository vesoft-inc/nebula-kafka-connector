package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

var OperationHost string
var KillWait = ""
var ProductInfo = "<all|nebulagraph|metad|graphd|stroaged|other product name>"

func RegisteOperationCmd(rootCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start nebula cluster",
		Long:  `start nebula cluster by config file which is default to .ngadmin.yaml ` + ProductInfo,
		Run: func(cmd *cobra.Command, args []string) {
			err := runOperation(args, "start")
			if err != nil {
				log.Printf("start failed: %v", err)
			} else {
				log.Printf("start success")
			}
		},
	}
	startCmd.Flags().StringVarP(&OperationHost, "host", "H", "", "host to start")

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "stop nebula cluster",
		Long:  `stop nebula cluster by config file which is default to .ngadmin.yaml ` + ProductInfo,
		Run: func(cmd *cobra.Command, args []string) {
			err := runOperation(args, "stop")
			if err != nil {
				log.Printf("stop failed: %v", err)
			} else {
				log.Printf("stop success")
			}
		},
	}
	stopCmd.Flags().StringVarP(&KillWait, "kill-wait", "k", KillWait, "wait for the process to exit, or kill it after the timeout. (support m,s) max 5m")
	stopCmd.Flags().StringVarP(&OperationHost, "host", "H", "", "host to stop")

	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "restart nebula cluster",
		Long:  `restart nebula cluster by config file which is default to .ngadmin.yaml ` + ProductInfo,
		Run: func(cmd *cobra.Command, args []string) {
			err := runOperation(args, "restart")
			if err != nil {
				log.Printf("restart failed: %v", err)
			} else {
				log.Printf("restart success")
			}
		},
	}
	restartCmd.Flags().StringVarP(&OperationHost, "host", "H", "", "host to restart")

	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(restartCmd)
}

func runOperation(args []string, operation string) error {
	component := args[0]
	jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
	if err != nil {
		return err
	}
	job := runner.NewJob("start")
	err = job.Run("operation", map[string]any{
		"host":      OperationHost,
		"operation": operation,
		"component": component,
		"kill-wait": KillWait,
	}, jobSpec)
	return err
}
