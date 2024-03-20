package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
)

func RegisteVerifyCmd(rootCmd *cobra.Command) {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "verify nebula cluster",
		Long:  `verify nebula cluster by config file which is default to .ngadm.yaml, check if agent is running`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("verify %s file.", Configfile)

			jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
			if err != nil {
				return err
			}
			job := runner.NewJob("verify")
			err = job.Run("verify", nil, jobSpec)
			if err != nil {
				log.Print("verify failed")
			}
			return err
		},
	}
	rootCmd.AddCommand(verifyCmd)
}
