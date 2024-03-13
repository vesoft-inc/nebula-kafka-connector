package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/cache"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type loginFlagsType struct {
	host     string
	port     uint32
	user     string
	password string
}

var loginFlags loginFlagsType

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server.",
	Long:  `login meta server --host [host] --port [port] --user [user] --password [password]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var address string
		address = fmt.Sprintf("%s:%d", loginFlags.host, loginFlags.port)
		_, err := meta.NewMetaClient(address)
		if err != nil {
			return fmt.Errorf("cannot login to %s, err: %v", address, err)
		}
		cache.SaveMetaSession(address)
		// TODO should login?
		// metaclient.Login(loginFlags.User, loginFlags.Pass)

		fmt.Fprintln(metaOutput, "[Warning] Login meta is not implemented yet.")
		fmt.Fprintln(metaOutput, "You can do operations without login now.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginFlags.host, "host", "H", "127.0.0.1", "meta server host")
	loginCmd.Flags().Uint32VarP(&loginFlags.port, "port", "P", 9559, "meta server port")
	loginCmd.Flags().StringVarP(&loginFlags.user, "user", "u", "", "user name")
	loginCmd.Flags().StringVarP(&loginFlags.password, "password", "p", "", "password")
}
