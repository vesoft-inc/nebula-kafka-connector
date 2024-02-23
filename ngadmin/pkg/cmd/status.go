package cmd

import (
	"log"
	"os"
	"sort"
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
	component := "all"
	if len(args) > 0 {
		component = args[0]
	}
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
	err = RenderStatusTableByJob(job)
	return err
}

func RenderStatusTableByJob(job *runner.Job) error {
	tb := table.NewWriter()
	tb.SetOutputMirror(os.Stdout)
	tb.AppendHeader(table.Row{"product", "service", "status", "host", "port"})
	tb.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, Align: text.AlignCenter},
	})
	values := job.Context.ValueMap
	arrs := make([]*types.StatusItem, 0)
	for key, value := range values {
		if strings.HasPrefix(key, "status-") {
			status := value.(types.StatusItem)
			arrs = append(arrs, &status)
		}
	}
	// sort by product and service
	sort.Slice(arrs, func(i, j int) bool {
		return arrs[i].Product < arrs[j].Product
	})
	// sort by product and service
	sortBase := []string{"nebulagraph", "license-manager"}
	index := 0
	for _, base := range sortBase {
		for i := 0; i < len(arrs); i++ {
			if arrs[i].Product == base {
				arrs[i], arrs[index] = arrs[index], arrs[i]
				index++
			}
		}
	}
	for _, status := range arrs {
		color := text.FgBlue
		if status.Status == "running" {
			color = text.FgGreen
		} else if status.Status == "exited" {
			color = text.FgRed
		}
		tb.AppendRow([]interface{}{status.Product, status.Service, text.Colors{color}.Sprint(status.Status), status.Host, status.Port})
		tb.AppendSeparator()
	}
	tb.Render()
	return nil
}
