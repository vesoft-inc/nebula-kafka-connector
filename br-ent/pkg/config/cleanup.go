package config

import (
	"crypto/tls"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/storage"
)

func AddCleanupFlags(flags *pflag.FlagSet) {
	flags.String(FlagMetaAddr, "", "Specify meta server")
	flags.String(flagBackupName, "", "Specify backup name")
	flags.String(flagUsername, "", "Username for login metad service")
	flags.String(flagPassword, "", "Password for login metad service")

	// support tls
	flags.Bool(flagEnableSSL, false, "Enable SSL connection")
	flags.String(flagCertPath, "/usr/local/certs/client.crt", "Specify cert file path")
	flags.String(flagKeyPath, "/usr/local/certs/client.key", "Specify key file path")
	flags.String(flagCAPath, "/usr/local/certs/ca.crt", "Specify CA file path")
	flags.Bool(flagInsecureSkipVerify, false, "Skip server side certificate verification")

	cobra.MarkFlagRequired(flags, flagBackupName)
	cobra.MarkFlagRequired(flags, FlagStorage)
}

type CleanupConfig struct {
	MetaAddr   string
	BackupName string
	Backend    *agentstorage.Backend // Backend is associated with the root uri
	TLSConfig  *tls.Config
	Username   string
	Password   string
}

func (c *CleanupConfig) ParseFlags(flags *pflag.FlagSet) error {
	var err error
	c.MetaAddr, err = flags.GetString(FlagMetaAddr)
	if err != nil {
		return err
	}

	c.BackupName, err = flags.GetString(flagBackupName)
	if err != nil {
		return err
	}
	c.Backend, err = storage.ParseFromFlags(flags)
	if err != nil {
		return fmt.Errorf("parse storage flags failed: %w", err)
	}

	c.Username, err = flags.GetString(flagUsername)
	if err != nil {
		return err
	}
	c.Password, err = flags.GetString(flagPassword)
	if err != nil {
		return err
	}

	// support tls
	c.TLSConfig, err = ParseTLSFlag(flags)
	if err != nil {
		return err
	}

	return nil
}
