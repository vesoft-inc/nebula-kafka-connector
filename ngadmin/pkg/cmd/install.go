package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

func RegisteInstallCmd(rootCmd *cobra.Command) *cobra.Command {
	var force bool
	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "install nebulagraph cluster with configfile",
		Long:  `once you installed nebulagraph cluster, you can use ngadmin to manage it by a config file cache, which is default to .ngadmin.yaml`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Print("install nebulagraph cluster.")
			jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
			if err != nil {
				return err
			}
			job := runner.NewJob("install")
			exit := make(chan error)
			go func() {
				err = job.Run("install", map[string]any{
					"force": force,
				}, jobSpec)
				if err != nil {
					exit <- err
					return
				}
				// cache the install config
				err = yamlparser.CacheConfig(jobSpec)
				if err != nil {
					log.Printf("cache config failed: %v", err)
				}
				exit <- err
			}()
			return <-exit
		},
	}
	// add --force
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "force to install")
	rootCmd.AddCommand(installCmd)
	return installCmd
}
