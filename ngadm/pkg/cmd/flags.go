package cmd

import "github.com/spf13/cobra"

var Configfile string

func AddConfigfileFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&Configfile, "configfile", "f", "ngadm.yaml", "config file path (default is ngadm.yaml), if not exists, you must specify a config file with --configfile")
}
