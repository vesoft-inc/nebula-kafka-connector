package config

import (
	"crypto/tls"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func AddBackupFlags(flags *pflag.FlagSet) {
	flags.String(FlagMetaAddr, "", "Specify meta server")
	flags.String(FlagAgentsAddr, "", "Specify agents address, eg: 192.168.8.1:6688,192.168.8.2:6688,192.168.8.3:6688")
	flags.String(flagBackupName, "", "Specify backup name")
	flags.Int(flagConcurrency, 5, "Max concurrency for download data and data playback")
	flags.String(flagUsername, "", "Username for login metad service")
	flags.String(flagPassword, "", "Password for login metad service")
	flags.Int64(flagClusterId, 0, "Specify the backup cluster id")
	flags.String(FlagConfig, "", "The config file for ngctl to backup")

	// support tls
	flags.Bool(flagEnableSSL, false, "Enable SSL connection")
	flags.String(flagCertPath, "/usr/local/certs/client.crt", "Specify cert file path")
	flags.String(flagKeyPath, "/usr/local/certs/client.key", "Specify key file path")
	flags.String(flagCAPath, "/usr/local/certs/ca.crt", "Specify CA file path")
	flags.Bool(flagInsecureSkipVerify, false, "Verify the server's certificate chain and host name")
	flags.String(flagServerName, "", "The subject alternative name (SAN) of the peer server to verify")

	cobra.MarkFlagRequired(flags, FlagStorage)
	cobra.MarkFlagRequired(flags, flagClusterId)
	cobra.MarkFlagRequired(flags, FlagMetaAddr)
	cobra.MarkFlagRequired(flags, FlagAgentsAddr)
	cobra.MarkFlagRequired(flags, FlagConfig)
}

type BackupConfig struct {
	BackupName  string
	MetaAddr    string
	AgentsAddr  string
	ClusterId   int64
	Concurrency int
	Backend     *agentstorage.Backend // Backend is associated with the root uri
	TLSConfig   *tls.Config
	Username    string
	Password    string

	Spec *types.JobSpec
}

func (b *BackupConfig) ParseFullFlags(flags *pflag.FlagSet) error {
	var err error
	b.BackupName, err = flags.GetString(flagBackupName)
	if err != nil {
		return err
	}

	b.MetaAddr, err = flags.GetString(FlagMetaAddr)
	if err != nil {
		return err
	}

	b.AgentsAddr, err = flags.GetString(FlagAgentsAddr)
	if err != nil {
		return err
	}

	b.ClusterId, err = flags.GetInt64(flagClusterId)
	if err != nil {
		return err
	}

	b.Concurrency, err = flags.GetInt(flagConcurrency)
	if err != nil {
		return err
	}
	b.Backend, err = storage.ParseFromFlags(flags)
	if err != nil {
		return fmt.Errorf("parse storage flags failed: %w", err)
	}

	b.Username, err = flags.GetString(flagUsername)
	if err != nil {
		return err
	}
	b.Password, err = flags.GetString(flagPassword)
	if err != nil {
		return err
	}

	configPath, err := flags.GetString(FlagConfig)
	if err != nil {
		return err
	}
	b.Spec, err = ParseConfig(configPath)
	if err != nil {
		return err
	}

	// support tls
	b.TLSConfig, err = ParseTLSFlag(flags)
	if err != nil {
		return err
	}

	return nil
}
