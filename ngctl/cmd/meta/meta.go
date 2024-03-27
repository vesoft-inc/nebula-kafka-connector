package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

var (
	metaClient meta.Client
)

func metaClientInit() error {
	if metaClient != nil {
		return nil
	}
	cacheToken, err := ngctl.LoadMetaToken()
	if err != nil || cacheToken == nil {
		return fmt.Errorf("load meta session failed, please login first.")
	}
	metaClient, err = meta.NewMetaClient(cacheToken.Address, meta.WithToken(cacheToken.Token))
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
	Use:   "ngctl",
	Short: "Execute meta command in cli mode.",
	Long: `Execute meta command in cli mode. Use 'ngctl -h' to see usage.
	**Notice:** You should login meta server first`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.Execute()
}
