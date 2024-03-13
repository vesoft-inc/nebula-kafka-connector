package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/cache"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

var (
	metaClient meta.Client
)

func metaClientInit() error {
	if metaClient != nil {
		return nil
	}
	cacheSession, err := cache.LoadMetaSession()
	if err != nil {
		return fmt.Errorf("load meta session failed: %s", err)
	}
	metaClient, err = meta.NewMetaClient(cacheSession.Address)
	if err != nil {
		return err
	}
	return nil
}

func metaClientClose() {
	if metaClient != nil {
		metaClient.Close()
	}
}

func metaConsoleError(message string, err string) error {
	if err != "" {
		return fmt.Errorf("%s, err: %s", message, err)
	} else {
		return fmt.Errorf("%s", message)
	}
}

// metaOutput is the output of meta command
// using os.Stdout by default, and could use other output for testing
var metaOutput io.Writer = os.Stdout

var rootCmd = &cobra.Command{
	Use:   "meta-ctl",
	Short: "Execute meta command in cli mode.",
	Long: `Execute meta command in cli mode. Use 'meta-ctl -h' to see usage.
	**Notice:** You should login meta server first`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.Execute()
}
