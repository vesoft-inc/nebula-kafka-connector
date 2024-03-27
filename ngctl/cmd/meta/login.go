package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

type loginFlagsType struct {
	host           string
	port           uint32
	user           string
	password       string
	promptPassword bool
}

var loginFlags loginFlagsType

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server.",
	Long:  `login meta server --host [host] --port [port] --user [user] --password [password]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginFlags.promptPassword {
			var err error
			pw := promptui.Prompt{
				Label:       "Password",
				AllowEdit:   true,
				Mask:        rune('*'),
				HideEntered: true,
			}
			loginFlags.password, err = pw.Run()
			if err != nil {
				return err
			}
		}
		var address string

		address = fmt.Sprintf("%s:%d", loginFlags.host, loginFlags.port)
		c, err := meta.NewMetaClient(address, meta.WithUserPassword(loginFlags.user, loginFlags.password))
		if err != nil {
			return metaConsoleError(fmt.Sprintf("cannot login to %s", address), err.Error())

		}
		if err := ngctl.SaveMetaToken(address, c.GetToken()); err != nil {
			return metaConsoleError("save meta session failed", err.Error())
		}

		fmt.Fprintln(metaOutput, "Login succeeded.")
		fmt.Fprintf(metaOutput, "Your token will be stored in %s.\n", ngctl.CachePath())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginFlags.host, "host", "H", "127.0.0.1", "meta server host")
	loginCmd.Flags().Uint32VarP(&loginFlags.port, "port", "P", 9559, "meta server port")
	loginCmd.Flags().StringVarP(&loginFlags.user, "user", "u", "root", "user name")
	loginCmd.Flags().StringVarP(&loginFlags.password, "password", "p", "NebulaGraph01", "password")
	loginCmd.Flags().BoolVarP(&loginFlags.promptPassword, "prompt-password", "i", false, "prompt password")
}
