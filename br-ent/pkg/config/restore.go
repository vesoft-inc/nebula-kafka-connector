package config

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/storage"
)

func AddRestoreFlags(flags *pflag.FlagSet) {
	flags.String(FlagMetaAddr, "", "Specify meta server")
	flags.String(FlagAgentsAddr, "", "Specify agents address, eg: 192.168.8.1:6688,192.168.8.2:6688,192.168.8.3:6688")
	flags.String(flagBackupName, "", "Specify backup name")
	flags.Int(flagConcurrency, 5, "Max concurrency for upload, download and playback data")
	flags.String(flagUsername, "", "Username for login metad service")
	flags.String(flagPassword, "", "Password for login metad service")
	flags.Int64(flagClusterId, 0, "Specify the restore cluster id")
	flags.Int64(flagBackupClusterId, 0, "Specify the backup cluster id")
	flags.String(flagCatalogOwner, "root", "Specify the restore cluster catalog owner")

	// support tls
	flags.Bool(flagEnableSSL, false, "Enable SSL connection")
	flags.String(flagCertPath, "/usr/local/certs/client.crt", "Specify cert file path")
	flags.String(flagKeyPath, "/usr/local/certs/client.key", "Specify key file path")
	flags.String(flagCAPath, "/usr/local/certs/ca.crt", "Specify CA file path")
	flags.Bool(flagInsecureSkipVerify, false, "Skip server side certificate verification")

	cobra.MarkFlagRequired(flags, FlagMetaAddr)
	cobra.MarkFlagRequired(flags, FlagAgentsAddr)
	cobra.MarkFlagRequired(flags, FlagStorage)
	cobra.MarkFlagRequired(flags, flagBackupName)
	cobra.MarkFlagRequired(flags, flagClusterId)
	cobra.MarkFlagRequired(flags, flagBackupClusterId)
}

type RestoreConfig struct {
	BackupName      string
	MetaAddr        string
	AgentsAddr      string
	ClusterId       int64
	BackupClusterId int64
	Concurrency     int
	Backend         *agentstorage.Backend // Backend is associated with the root uri
	TLSConfig       *tls.Config
	Username        string
	Password        string
	CatalogOwner    string
}

func (r *RestoreConfig) ParseFlags(flags *pflag.FlagSet) error {
	var err error

	r.BackupName, err = flags.GetString(flagBackupName)
	if err != nil {
		return err
	}

	r.MetaAddr, err = flags.GetString(FlagMetaAddr)
	if err != nil {
		return err
	}

	r.AgentsAddr, err = flags.GetString(FlagAgentsAddr)
	if err != nil {
		return err
	}

	r.ClusterId, err = flags.GetInt64(flagClusterId)
	if err != nil {
		return err
	}

	r.BackupClusterId, err = flags.GetInt64(flagBackupClusterId)
	if err != nil {
		return err
	}

	r.CatalogOwner, err = flags.GetString(flagCatalogOwner)
	if err != nil {
		return err
	}

	r.Concurrency, err = flags.GetInt(flagConcurrency)
	if err != nil {
		return err
	}
	r.Backend, err = storage.ParseFromFlags(flags)
	if err != nil {
		return fmt.Errorf("parse storage flags failed: %w", err)
	}

	r.Username, err = flags.GetString(flagUsername)
	if err != nil {
		return err
	}
	r.Password, err = flags.GetString(flagPassword)
	if err != nil {
		return err
	}

	// support tls
	r.TLSConfig, err = ParseTLSFlag(flags)
	if err != nil {
		return err
	}

	return nil
}

func parseClusterIdMapping(clusterIdMappingStr string) (map[int64]int64, error) {
	clusterIdMapping := make(map[int64]int64)
	mapping := strings.Split(clusterIdMappingStr, ",")
	for _, m := range mapping {
		ids := strings.Split(m, ":")
		if len(ids) != 2 {
			return nil, fmt.Errorf("invalid cluster id mapping: %s", m)
		}
		oldClusterId, err := strconv.Atoi(ids[0])
		if err != nil {
			return nil, fmt.Errorf("parse old cluster id failed: %w", err)
		}
		newClusterId, err := strconv.Atoi(ids[1])
		if err != nil {
			return nil, fmt.Errorf("parse new cluster id failed: %w", err)
		}
		clusterIdMapping[int64(oldClusterId)] = int64(newClusterId)
	}

	return clusterIdMapping, nil
}
