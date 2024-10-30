package passwd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type passwdFlagsType struct {
	host           string
	port           uint32
	user           string
	currentPasswd  string
	newPasswd      string
	enableTLS      bool
	ca             string
	cert           string
	key            string
	peerNameVerify bool
	peerName       string
}

var passwdFlags passwdFlagsType

var PasswdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Change the password of the user",
	Long:  "Change the password of the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		//passwd no need to login metad
		//user can reset password first without login
		cacheToken, err := common.LoadMetaToken()
		var addr string
		if err != nil || cacheToken == nil {
			if passwdFlags.host == "" || passwdFlags.port == 0 {
				return common.NgctlError("Please specify the host and port.", "")
			}
			addr = fmt.Sprintf("%s:%d", passwdFlags.host, passwdFlags.port)
		} else {
			addr = cacheToken.Address
		}
		var client meta.Client
		if cacheToken != nil {
			client, err = meta.NewMetaClient(addr, meta.WithTLS(
				cacheToken.EnableTLS,
				cacheToken.CA,
				cacheToken.Cert,
				cacheToken.Key,
				cacheToken.PeerNameVerify,
				cacheToken.PeerName),
			)
		} else {
			client, err = meta.NewMetaClient(addr, meta.WithTLS(
				passwdFlags.enableTLS,
				passwdFlags.ca,
				passwdFlags.cert,
				passwdFlags.key,
				passwdFlags.peerNameVerify,
				passwdFlags.peerName),
			)
		}
		if err != nil {
			return err
		}
		defer client.Close()
		var old, new string
		if cacheToken != nil && passwdFlags.host == "" {
			old, new, err = common.GetPromptPassword()
			if err != nil {
				return common.NgctlError("cannot change password, ", err.Error())
			}
		} else {
			old = passwdFlags.currentPasswd
			new = passwdFlags.newPasswd
		}
		if old == "" || new == "" {
			return common.NgctlError("password cannot be empty", "")
		}
		if passwdFlags.user == "" {
			return common.NgctlError("user cannot be empty", "")
		}
		if err := common.ResetPassword(client, passwdFlags.user, old, new); err != nil {
			return common.NgctlError("reset password failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Change password succeeded.")

		return nil
	},
}

func init() {
	PasswdCmd.Flags().StringVarP(&passwdFlags.host, "host", "H", "127.0.0.1", "meta server host")
	PasswdCmd.Flags().Uint32VarP(&passwdFlags.port, "port", "P", 9559, "meta server port")
	PasswdCmd.Flags().StringVarP(&passwdFlags.user, "user", "u", "root", "user name")
	PasswdCmd.Flags().StringVarP(&passwdFlags.currentPasswd, "current_password", "c", "", "current password")
	PasswdCmd.Flags().StringVarP(&passwdFlags.newPasswd, "new_password", "p", "", "new password")
	PasswdCmd.Flags().BoolVarP(&passwdFlags.enableTLS, "enable-tls", "", false, "Enable TLS")
	PasswdCmd.Flags().StringVarP(&passwdFlags.ca, "ca", "", "", "Certificate of trusted CA, in PEM format")
	PasswdCmd.Flags().StringVarP(&passwdFlags.cert, "cert", "", "", "Certificate of meta client, in PEM format")
	PasswdCmd.Flags().StringVarP(&passwdFlags.key, "key", "", "", "Private key of meta client, in PEM format")
	PasswdCmd.Flags().BoolVarP(&passwdFlags.peerNameVerify, "peer-name-verify", "", false, "Enable peer name verification")
	PasswdCmd.Flags().StringVarP(&passwdFlags.peerName, "peer-name", "", "", "Peer name to override the default, i.e. domain name")
}
