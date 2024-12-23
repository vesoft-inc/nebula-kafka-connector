package user_admin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type userFlagsType struct {
	user                string
	password            string
	passwordEncryptType string
	authType            string
	authInfo            string
	output              string
}

var userFlags = userFlagsType{}

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Process user management command",
	Long:  "Execute user management in cli mode",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create user in meta server",
	Long: `ngctl user create --user [username] --password [password] or 
ngctl user create --user [username] --auth-type [authType] --auth-info [authInfo]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if userFlags.password != "" {
			m := make(map[string]string)
			m["password"] = userFlags.password
			m["encry_type"] = userFlags.passwordEncryptType
			bs, err := json.Marshal(m)
			if err != nil {
				return err
			}
			userFlags.authType = "password"
			userFlags.authInfo = string(bs)
		}
		req := meta.NewCreateUserReq(
			userFlags.user,
			userFlags.authType,
			userFlags.authInfo,
		)
		if err := common.MetaClient.CreateUser(req); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Create user successfully.\n")
		return nil
	},
}

var dropUserCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop user in meta server",
	Long:  "ngctl user drop --user [username]",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		req := meta.NewDropUserReq(userFlags.user)
		if err := common.MetaClient.DropUser(req); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Drop user successfully.\n")
		return nil
	},
}

var alterUserCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter user in meta server",
	Long: `ngctl user alter --user [username] --auth-info [authInfo] or
ngctl user alter --user [username] --password [password]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if userFlags.password != "" {
			m := make(map[string]string)
			m["password"] = userFlags.password
			m["encry_type"] = userFlags.passwordEncryptType
			bs, err := json.Marshal(m)
			if err != nil {
				return err
			}
			userFlags.authType = "password"
			userFlags.authInfo = string(bs)
		}
		req := meta.NewAlterUserReq(
			userFlags.user,
			userFlags.authInfo,
			true, //always active for alter user
		)
		if err := common.MetaClient.AlterUser(req); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Alter user successfully.\n")

		return nil
	},
}

var showUserCmd = &cobra.Command{
	Use:   "show",
	Short: "show user in meta server",
	Long:  "ngctl user show --user aa,bb",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var l []string
		if userFlags.user != "" {
			l = append(l, strings.Split(userFlags.user, ",")...)
		}
		req := meta.NewListUsersReq(l)
		resp, err := common.MetaClient.ListUsers(req)
		if err != nil {
			return err
		}
		header := []string{
			"Name",
			"Active",
			"Auth Type",
			"Auth Info",
			"Created Time",
			"Last Updated Time",
			"Last Login Time",
			"Disabled Time",
		}
		data := make([][]string, 0)
		for _, u := range resp.Users {
			row := make([]string, 0)
			row = append(row, u.Name)
			if u.Active {
				row = append(row, "Y")
			} else {
				row = append(row, "N")
			}
			row = append(row, u.AuthType)
			row = append(row, u.AuthInfo)
			row = append(row, common.FormatTime(u.CreatedTime))
			row = append(row, common.FormatTime(u.LastUpdatedTime))
			row = append(row, common.FormatTime(u.LastLoginTime))
			row = append(row, common.FormatTime(u.DisabledTime))
			data = append(data, row)
		}
		// order by user name
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		r, err := common.Format(header, data, common.OutputFormatType(userFlags.output))
		if err != nil {
			return common.NgctlError("Show user failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

var disableUserCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable user in meta server",
	Long:  "ngctl user disable --user aa",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		req := meta.NewAlterUserReq(userFlags.user, "", false)
		if err := common.MetaClient.AlterUser(req); err != nil {
			return err
		}

		fmt.Fprintln(common.MetaOutput, "disable user successfully")
		return nil
	},
}

var enableUserCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable user in meta server",
	Long:  "ngctl user enable --user aa",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		req := meta.NewAlterUserReq(userFlags.user, "", true)
		if err := common.MetaClient.AlterUser(req); err != nil {
			return err
		}

		fmt.Fprintln(common.MetaOutput, "enable user successfully")
		return nil
	},
}

func init() {
	UserCmd.AddCommand(createUserCmd)
	UserCmd.AddCommand(dropUserCmd)
	UserCmd.AddCommand(alterUserCmd)
	UserCmd.AddCommand(showUserCmd)
	UserCmd.AddCommand(disableUserCmd)
	UserCmd.AddCommand(enableUserCmd)

	createUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	createUserCmd.Flags().StringVarP(&userFlags.password, "password", "p", "", "User password")
	createUserCmd.Flags().StringVarP(&userFlags.passwordEncryptType, "encrypt-type", "e", "sha256", "User password encrypt type, options: sha256, sha512, sm3")
	createUserCmd.Flags().StringVar(&userFlags.authType, "auth-type", "", "User auth type")
	createUserCmd.Flags().StringVar(&userFlags.authInfo, "auth-info", "", "User auth info")
	createUserCmd.MarkFlagRequired("user")
	createUserCmd.MarkFlagsRequiredTogether("auth-type", "auth-info")
	createUserCmd.MarkFlagsMutuallyExclusive("auth-type", "password")

	dropUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	dropUserCmd.MarkFlagRequired("user")

	alterUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	alterUserCmd.Flags().StringVarP(&userFlags.password, "password", "p", "", "User password")
	alterUserCmd.Flags().StringVarP(&userFlags.passwordEncryptType, "encrypt-type", "e", "sha256", "User password encrypt type, options: sha256, sha512, sm3")
	alterUserCmd.Flags().StringVar(&userFlags.authInfo, "auth-info", "", "User auth info")
	alterUserCmd.MarkFlagRequired("user")
	alterUserCmd.MarkFlagsMutuallyExclusive("password", "auth-info")
	alterUserCmd.MarkFlagsOneRequired("password", "auth-info")

	showUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "Users, e.g. 'aa,bb'")
	showUserCmd.Flags().StringVarP(&userFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

	disableUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	disableUserCmd.MarkFlagRequired("user")

	enableUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	enableUserCmd.MarkFlagRequired("user")
}
