package config

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"

	"github.com/spf13/pflag"
)

const (
	FlagStorage    = "storage"
	FlagMetaAddr   = "meta"
	FlagAgentsAddr = "agents"
	FlagConfig     = "config"

	FlagLogPath  = "log"
	FlagLogDebug = "debug"

	flagBackupName = "name"

	flagCertPath           = "cert-path"
	flagKeyPath            = "key-path"
	flagCAPath             = "ca-path"
	flagEnableSSL          = "enable-ssl"
	flagInsecureSkipVerify = "insecure-skip-verify"
	flagServerName         = "server-name"
	flagConcurrency        = "concurrency"
	flagUsername           = "username"
	flagPassword           = "password"
	flagServiceGroupId          = "clusterId"
	flagCatalogOwner       = "catalog-owner"
	flagForce              = "force"

	CACertPathEnv     = "CA_CERT_PATH"
	ClientCertPathEnv = "CLIENT_CERT_PATH"
	ClientKeyPathEnv  = "CLIENT_KEY_PATH"
)

func AddCommonFlags(flags *pflag.FlagSet) {
	flags.String(FlagLogPath, "br.log", "Specify br detail log path")
	flags.Bool(FlagLogDebug, false, "Output log in debug level or not")
	storage.AddFlags(flags)
}

func ParseTLSFlag(flags *pflag.FlagSet) (*tls.Config, error) {
	enableSSL, err := flags.GetBool(flagEnableSSL)
	if err != nil {
		return nil, err
	}
	if !enableSSL {
		return nil, nil
	}

	certPath, err := flags.GetString(flagCertPath)
	if err != nil {
		return nil, err
	}
	if os.Getenv(ClientCertPathEnv) != "" {
		certPath = os.Getenv(ClientCertPathEnv)
	}
	keyPath, err := flags.GetString(flagKeyPath)
	if err != nil {
		return nil, err
	}
	if os.Getenv(ClientKeyPathEnv) != "" {
		keyPath = os.Getenv(ClientKeyPathEnv)
	}
	caPath, err := flags.GetString(flagCAPath)
	if err != nil {
		return nil, err
	}
	if os.Getenv(CACertPathEnv) != "" {
		caPath = os.Getenv(CACertPathEnv)
	}
	insecureSkipVerify, err := flags.GetBool(flagInsecureSkipVerify)
	if err != nil {
		return nil, err
	}

	serverName, err := flags.GetString(flagServerName)
	if err != nil {
		return nil, err
	}

	caCert, clientCert, clientKey, err := utils.GetCerts(caPath, certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certs: %w", err)
	}
	tlsConfig, err := utils.LoadTLSConfig(caCert, clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load tls config: %w", err)
	}
	tlsConfig.InsecureSkipVerify = insecureSkipVerify
	tlsConfig.ServerName = serverName

	return tlsConfig, nil
}
