package supercluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Change the password of the user.",
	Long:  `passwd`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheToken, err := common.LoadMetaToken()
		if err != nil || cacheToken == nil {
			return fmt.Errorf("load meta session failed, please login first.")
		}
		client, err := meta.NewMetaClient(cacheToken.Address)
		if err != nil {
			return err
		}
		defer client.Close()
		if _, err := resetPassword(client); err != nil {
			return common.NgctlError("cannot change password, ", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Change password succeeded.")

		return nil
	},
}

func init() {
}
