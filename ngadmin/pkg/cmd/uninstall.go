package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

var Drain bool = false

func RegisteUninstallCmd(rootCmd *cobra.Command) {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall nebula cluster",
		Long:  `uninstall nebula cluster by config file which is default to .ngadmin.yaml, check if agent is running`,
		Run: func(cmd *cobra.Command, args []string) {
			jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
			if err != nil {
				fmt.Println("parse config file failed: ", err)
				return
			}
			job := runner.NewJob("uninstall")
			err = job.Run("uninstall", map[string]any{
				"kill-wait": KillWait,
				"drain":     Drain,
			}, jobSpec)
			if err != nil {
				fmt.Println("uninstall failed: ", err)
			}
		},
	}
	uninstallCmd.Flags().BoolP("drain", "d", Drain, "delete nebulagraph cluster data")
	uninstallCmd.Flags().StringVarP(&KillWait, "kill-wait", "k", KillWait, "wait for the process to exit, or kill it after the timeout. (support m,s) max 5m")
	rootCmd.AddCommand(uninstallCmd)
}
