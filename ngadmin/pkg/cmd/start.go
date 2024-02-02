package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

func RegisteStartCmd(rootCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start nebula cluster",
		Long:  `start nebula cluster by config file which is default to .ngadmin.yaml`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("start %s file.", Configfile)

			jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
			if err != nil {
				return err
			}
			job := runner.NewJob("start")
			err = job.Run("start", nil, jobSpec)
			if err != nil {
				log.Print("start failed")
			}
			return err
		},
	}
	rootCmd.AddCommand(startCmd)
}
