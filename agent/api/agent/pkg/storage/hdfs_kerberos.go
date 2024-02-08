package storage

import (
	"strings"

	krb5 "github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/zeromicro/go-zero/core/logx"
)

type KerberosConfig struct {
	Enable                       bool
	Principle                    string
	KeytabFilePath               string
	ConfigFilePath               string
	KerberosServicePrincipleName string
}

func getKerberosClient(krbConf KerberosConfig) (*krb5.Client, error) {
	ktFromFile, err := keytab.Load(krbConf.KeytabFilePath)
	if err != nil {
		logx.Errorf("unable to get keytab file: %v", err)
		return nil, err
	}
	cfg, err := config.Load(krbConf.ConfigFilePath)
	if err != nil {
		logx.Errorf("unable to get kerberos config: %v", err)
		return nil, err
	}

	principles := strings.Split(krbConf.Principle, "@")
	return krb5.NewWithKeytab(principles[0], principles[1], ktFromFile, cfg), nil
}
