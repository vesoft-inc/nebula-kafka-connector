package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
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
		if loginFlags.password == "" {
			var err error
			pw := promptui.Prompt{
				Label:       "Password",
				AllowEdit:   true,
				Mask:        rune(' '),
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
			return metaConsoleError(fmt.Sprintf("cannot create client to %s", address), err.Error())
		}
		if err := c.Login(); err != nil {
			//should reset password for first login
			e, ok := err.(*nebula.NebulaError)
			if !ok {
				return metaConsoleError("Login failed", err.Error())
			}
			if e.Code() != nebula.ERROR_AUTH_NEED_CHANGE_PASSWORD {
				return metaConsoleError("Login failed", err.Error())
			}
			// reset password and re-login
			fmt.Fprintln(metaOutput, "Please reset the password for the first login.")
			var newPassword string
			newPassword, err = resetPassword(c)
			if err != nil {
				return metaConsoleError("Cannot reset password", err.Error())
			}
			fmt.Fprintf(metaOutput, "Reset password succeeded for %s, re-login with the new password.\n", loginFlags.user)
			c.Close()
			c, err = meta.NewMetaClient(address, meta.WithUserPassword(loginFlags.user, newPassword))
			if err != nil {
				return err
			}
			if err := c.Login(); err != nil {
				return metaConsoleError("Login failed", err.Error())
			}
		}
		if err := ngctl.SaveMetaToken(address, c.GetToken()); err != nil {
			return metaConsoleError("Save meta session failed", err.Error())
		}

		fmt.Fprintln(metaOutput, "Login succeeded.")
		fmt.Fprintf(metaOutput, "Your token will be stored in %s.\n", ngctl.CachePath())

		return nil
	},
}

func resetPassword(c meta.Client) (string, error) {
	currentPassword := promptui.Prompt{
		Label:     "Current password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	currentPasswordStr, err := currentPassword.Run()
	if err != nil {
		return "", err
	}
	newPassword := promptui.Prompt{
		Label:     "New password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	newPasswordStr, err := newPassword.Run()
	if err != nil {
		return "", err
	}
	confirmPassword := promptui.Prompt{
		Label:     "Retype new password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	confirmPasswordStr, err := confirmPassword.Run()
	if err != nil {
		return "", err
	}
	if newPasswordStr != confirmPasswordStr {
		return "", fmt.Errorf("Sorry, the passwords you entered do not match.")
	}
	req := meta.NewChangePasswordReq(
		loginFlags.user,
		currentPasswordStr,
		newPasswordStr,
	)
	if err := c.ChangePassword(req); err != nil {
		return "", err
	}
	return newPasswordStr, nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginFlags.host, "host", "H", "127.0.0.1", "meta server host")
	loginCmd.Flags().Uint32VarP(&loginFlags.port, "port", "P", 9559, "meta server port")
	loginCmd.Flags().StringVarP(&loginFlags.user, "user", "u", "root", "user name")
	loginCmd.Flags().StringVarP(&loginFlags.password, "password", "p", "", "password")
}
