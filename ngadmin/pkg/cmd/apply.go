package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

func RegisteApplyCmd(rootCmd *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "apply",
		Short: "apply nebulagraph config with --configfile",
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("install nebulagraph cluster.")
			jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
			if err != nil {
				log.Fatalf("parse yaml file failed: %v", err)
			}
			job := runner.NewJob("install")
			err = job.Run("config", nil, jobSpec)
			if err != nil {
				log.Fatalf("apply config failed: %v", err)
			}
			log.Print("apply config success.")
		},
	}
	rootCmd.AddCommand(cmd)
	return cmd
}
