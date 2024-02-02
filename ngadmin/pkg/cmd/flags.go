package cmd

import "github.com/spf13/cobra"

var Configfile string

func AddConfigfileFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&Configfile, "configfile", "c", "ngadmin.yaml", "config file path (default is ngadmin.yaml), if not exists, you must specify a config file with --configfile")
}
