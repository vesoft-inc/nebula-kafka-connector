package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "execute meta command in cli mode",
	Long:  `execute meta command in cli mode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello this is nebula-console meta.")
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)
}
