package metad

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var showMetadCmd = &cobra.Command{
	Use:   "show",
	Short: "Show metad",
	Long:  "Show metad",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := common.MetaClient.ShowMeta()
		if err != nil {
			return common.NgctlError("Show metad failed", err.Error())
		}
		header := []string{"ClusterId", "Git Sha", "Peers", "Agents"}
		data := make([][]string, 0)
		row := make([]string, 0)
		row = append(row, fmt.Sprintf("%d", resp.ClusterId))
		row = append(row, resp.GitInfoSha)
		row = append(row, strings.Join(resp.Peers, ","))
		row = append(row, strings.Join(resp.Agents, ","))
		data = append(data, row)
		r, err := common.Format(header, data, common.OutputFormatType(metadFlags.output))
		if err != nil {
			return common.NgctlError("Show metad failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	showMetadCmd.Flags().StringVarP(&metadFlags.output, "output", "o", "table", "output format. Allowed values: table, json")
}
