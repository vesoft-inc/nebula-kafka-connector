package login

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type loginFlagsType struct {
	host           string
	port           uint32
	user           string
	password       string
	enableTLS      bool
	ca             string
	cert           string
	key            string
	peerNameVerify bool
	peerName       string
}

var loginFlags loginFlagsType

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server",
	Long:  "ngctl login --host [host] --port [port] --user [user] --password [password]",
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
		c, err := meta.NewMetaClient(address, meta.WithUserPassword(loginFlags.user, loginFlags.password), meta.WithTLS(loginFlags.enableTLS, loginFlags.ca, loginFlags.cert, loginFlags.key, loginFlags.peerNameVerify, loginFlags.peerName))
		if err != nil {
			return common.NgctlError(fmt.Sprintf("cannot create client to %s", address), err.Error())
		}
		resp, err := c.Login()
		if err != nil {
			//should reset password for first login
			e, ok := err.(*errors.NebulaError)
			if !ok {
				return common.NgctlError("Login failed", err.Error())
			}
			if e.Code() != errors.ERROR_AUTH_NEED_CHANGE_PASSWORD {
				return common.NgctlError("Login failed", err.Error())
			}
			// reset password and re-login
			fmt.Fprintln(common.MetaOutput, "Please reset the password for the first login.")
			currentPassword, newPassword, err := common.GetPromptPassword()
			if err != nil {
				return common.NgctlError("Cannot reset password", err.Error())
			}
			if err := common.ResetPassword(
				c, loginFlags.user, currentPassword, newPassword); err != nil {
				return common.NgctlError("Reset password failed", err.Error())
			}
			fmt.Fprintf(
				common.MetaOutput,
				"Reset password succeeded for %s, re-login with the new password.\n",
				loginFlags.user,
			)
			c.Close()
			c, err = meta.NewMetaClient(address, meta.WithUserPassword(loginFlags.user, newPassword), meta.WithTLS(loginFlags.enableTLS, loginFlags.ca, loginFlags.cert, loginFlags.key, loginFlags.peerNameVerify, loginFlags.peerName))
			if err != nil {
				return err
			}
			resp, err = c.Login()
			if err != nil {
				return common.NgctlError("Login failed", err.Error())
			}
		}
		if err := common.SaveMetaToken(address, resp.Leader, resp.Token, loginFlags.enableTLS, loginFlags.ca, loginFlags.cert, loginFlags.key, loginFlags.peerNameVerify, loginFlags.peerName); err != nil {
			return common.NgctlError("Save meta session failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Login succeeded.")
		fmt.Fprintf(common.MetaOutput, "Your token will be stored in <%s>.\n", common.GetCachePath())

		return nil
	},
}

func init() {
	LoginCmd.Flags().StringVarP(&loginFlags.host, "host", "H", "127.0.0.1", "meta server host")
	LoginCmd.Flags().Uint32VarP(&loginFlags.port, "port", "P", 9559, "meta server port")
	LoginCmd.Flags().StringVarP(&loginFlags.user, "user", "u", "root", "user name")
	LoginCmd.Flags().StringVarP(&loginFlags.password, "password", "p", "", "password")
	// LoginCmd.Flags().BoolVarP(&loginFlags.enableTLS, "enable-tls", "", false, "Enable TLS")
	// LoginCmd.Flags().StringVarP(&loginFlags.ca, "ca", "", "", "Certificate of trusted CA, in PEM format")
	// LoginCmd.Flags().StringVarP(&loginFlags.cert, "cert", "", "", "Certificate of meta client, in PEM format")
	// LoginCmd.Flags().StringVarP(&loginFlags.key, "key", "", "", "Private key of meta client, in PEM format")
	// LoginCmd.Flags().BoolVarP(&loginFlags.peerNameVerify, "peer-name-verify", "", false, "Enable peer name verification")
	// LoginCmd.Flags().StringVarP(&loginFlags.peerName, "peer-name", "", "", "Peer name to override the default, i.e. domain name")
}
