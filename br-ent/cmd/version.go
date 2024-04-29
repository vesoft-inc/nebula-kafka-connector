package cmd

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/version"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of br-ent tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(`%s,V-%d.%d.%d
   GitSha: %s
   GitRef: %s
please run "help" subcommand for more infomation.`,
				version.VerName, version.VerMajor, version.VerMinor, version.VerPatch,
				version.GitSha,
				version.GitRef)

			return nil
		},
	}
	return versionCmd
}
