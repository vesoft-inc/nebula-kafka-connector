package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout meta server",
	Long:  "logout meta server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := common.MetaClient.Logout(); err != nil {
			return err
		}
		if err := common.ClearMetaToken(); err != nil {
			return common.NgctlError("clear meta session failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Logout succeeded.")

		return nil
	},
}

func init() {
}
