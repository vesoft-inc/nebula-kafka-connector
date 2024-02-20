package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

func RegisteStatusCmd(rootCmd *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status",
		Short: "view nebulagraph cluster status <metad|graphd|storaged|all>",
		Run: func(cmd *cobra.Command, args []string) {
			err := GetStatus(args)
			if err != nil {
				log.Printf("get status failed: %v", err)
			}
		},
	}
	rootCmd.AddCommand(cmd)
	return cmd
}

func GetStatus(args []string) error {
	component := args[0]
	jobSpec, err := yamlparser.ParseYamlByPath(Configfile)
	if err != nil {
		return err
	}
	job := runner.NewJob("view status")
	err = job.Run("status", map[string]any{
		"component": component,
	}, jobSpec)
	// use go-pretty to print the status
	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"product", "service", "status", "host", "port"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Align: text.AlignCenter},
	})
	values := job.Context.ValueMap
	for key, value := range values {
		if strings.HasPrefix(key, "status-") {
			status := value.(types.StatusItem)
			t.AppendRow([]interface{}{status.Product, status.Service, status.Status, status.Host, status.Port})
		}
	}
	t.Render()
	return err
}
