package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

type userFlagsType struct {
	user     string
	password string
	authType string
	authInfo string
	disable  bool
}

var userFlags = userFlagsType{}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Process user management command",
	Long:  `Execute user management in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create user in meta server.",
	Long: `ngctl user create --user [username] --password [password] or 
ngctl user create --user [username] --authType [authType] --authInfo [authInfo]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if userFlags.password != "" {
			m := make(map[string]string)
			m["password"] = userFlags.password
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
		if err := metaClient.CreateUser(req); err != nil {
			return err
		}
		fmt.Fprintf(metaOutput, "Create user successfully.\n")
		return nil
	},
}

var dropUserCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop user in meta server.",
	Long:  `ngctl user drop --user [username]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		req := meta.NewDropUserReq(userFlags.user)
		if err := metaClient.DropUser(req); err != nil {
			return err
		}
		fmt.Fprintf(metaOutput, "Drop user successfully.\n")
		return nil
	},
}

var alterUserCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter user in meta server.",
	Long: `ngctl user alter --user [username] --authInfo [authInfo] or
ngctl user alter --user [username] --disable
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if userFlags.password != "" && userFlags.authInfo != "" {
			return fmt.Errorf("password and authInfo are mutually exclusive")
		}
		if userFlags.password != "" {
			m := make(map[string]string)
			m["password"] = userFlags.password
			bs, err := json.Marshal(m)
			if err != nil {
				return err
			}
			userFlags.authType = "password"
			userFlags.authInfo = string(bs)
		}
		if userFlags.authInfo == "" {
			userFlags.authInfo = "{}"
		}
		req := meta.NewAlterUserReq(
			userFlags.user,
			userFlags.authInfo,
			!userFlags.disable,
		)
		if err := metaClient.AlterUser(req); err != nil {
			return err
		}
		fmt.Fprintf(metaOutput, "Alter user successfully.\n")

		return nil
	},
}

var showUserCmd = &cobra.Command{
	Use:   "show",
	Short: "show user in meta server.",
	Long:  `ngctl user show --user aa,bb`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var l []string
		if userFlags.user != "" {
			l = append(l, strings.Split(userFlags.user, ",")...)
		}
		req := meta.NewListUsersReq(l)
		resp, err := metaClient.ListUsers(req)
		if err != nil {
			return err
		}
		header := []string{"Name", "Active", "Auth Type", "Created Time", "Last Updated Time", "Last Login Time"}
		data := make([][]string, 0)
		for _, u := range resp.Users {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%s", u.Name))
			if u.Active {
				row = append(row, "Y")
			} else {
				row = append(row, fmt.Sprintf("%s/%s", "N", formatTime(u.DisabledTime)))
			}
			row = append(row, u.AuthType)
			row = append(row, fmt.Sprintf("%s", formatTime(u.CreatedTime)))
			row = append(row, fmt.Sprintf("%s", formatTime(u.LastUpdatedTime)))
			row = append(row, fmt.Sprintf("%s", formatTime(u.LastLoginTime)))
			data = append(data, row)
		}
		// order by user name
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		fmt.Fprintln(metaOutput, ngctl.FormatTable(header, data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(dropUserCmd)
	userCmd.AddCommand(alterUserCmd)
	userCmd.AddCommand(showUserCmd)

	createUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	createUserCmd.Flags().StringVarP(&userFlags.password, "password", "p", "", "User password")
	createUserCmd.Flags().StringVar(&userFlags.authType, "auth-type", "", "User auth type")
	createUserCmd.Flags().StringVar(&userFlags.authInfo, "auth-info", "", "User auth info")
	createUserCmd.MarkFlagRequired("user")
	createUserCmd.MarkFlagsRequiredTogether("auth-type", "auth-info")
	createUserCmd.MarkFlagsMutuallyExclusive("auth-type", "password")

	dropUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	dropUserCmd.MarkFlagRequired("user")

	alterUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "User name")
	alterUserCmd.Flags().StringVarP(&userFlags.password, "password", "p", "", "User password")
	alterUserCmd.Flags().StringVar(&userFlags.authInfo, "auth-info", "", "User auth info")
	alterUserCmd.Flags().BoolVar(&userFlags.disable, "disable", false, "User disable")
	alterUserCmd.MarkFlagRequired("user")

	showUserCmd.Flags().StringVarP(&userFlags.user, "user", "u", "", "Users, e.g. 'aa,bb'")
}
