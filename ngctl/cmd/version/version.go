package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitTag    string
	GitCommit string
	BuildTime string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ngctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", GitTag)
		fmt.Println("Git commit:", GitCommit)
		fmt.Println("Build time:", BuildTime)
	},
}
