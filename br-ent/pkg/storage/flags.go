package storage

import (
	"fmt"
	"strings"

	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagStorage = "storage"

	flagS3Endpoint  = "s3.endpoint"
	flagS3Region    = "s3.region"
	flagS3AccessKey = "s3.access_key"
	flagS3SecretKey = "s3.secret_key"
)

func AddFlags(flags *pflag.FlagSet) {
	flags.String(flagStorage, "",
		`backup target url, format: <SCHEME>://<PATH>.
    <SCHEME>: a string indicating which backend type. optional: local, s3.
    now only s3-compatible is supported.
    example:
    for local - "local:///the/local/path/to/backup"
    for s3  - "s3://example/url/to/the/backup"
    `)
	if err := cobra.MarkFlagRequired(flags, flagStorage); err != nil {
		log.Errorf("Failed to mark flag %s required: %v.", flagStorage, err)
	}
	AddS3Flags(flags)
	AddLocalFlags(flags)
}

func AddS3Flags(flags *pflag.FlagSet) {
	flags.String(flagS3Region, "", "S3 Option: set region or location to upload or download backup")
	flags.String(flagS3Endpoint, "",
		"S3 Option: set the S3 endpoint URL, please specify the http or https scheme explicitly")
	flags.String(flagS3AccessKey, "", "S3 Option: set access key id")
	flags.String(flagS3SecretKey, "", "S3 Option: set secret key for access id")
}

func AddLocalFlags(flags *pflag.FlagSet) {
	// There is no need extra flags for local storage other than local uri
}

func ParseFromFlags(flags *pflag.FlagSet) (*agentstorage.Backend, error) {
	s, err := flags.GetString(flagStorage)
	if err != nil {
		return nil, err
	}
	s = strings.TrimRight(s, "/ ") // trim tailing space and / in passed in storage uri

	t := agentstorage.ParseType(s)
	b := &agentstorage.Backend{}
	switch t {
	case agentstorage.S3Type:
		region, err := flags.GetString(flagS3Region)
		if err != nil {
			return nil, err
		}
		endpoint, err := flags.GetString(flagS3Endpoint)
		if err != nil {
			return nil, err
		}
		accessKey, err := flags.GetString(flagS3AccessKey)
		if err != nil {
			return nil, err
		}
		secretKey, err := flags.GetString(flagS3SecretKey)
		if err != nil {
			return nil, err
		}
		if err := b.SetUri(s); err != nil {
			return nil, err
		}
		b.S3.Region = region
		b.S3.Endpoint = endpoint
		b.S3.AccessKey = accessKey
		b.S3.SecretKey = secretKey
	default:
		return nil, fmt.Errorf("bad format backend: %d", t)
	}

	log.WithField("type", t).WithField("uri", s).Debugln("Parse storage flag.")
	return b, nil
}
