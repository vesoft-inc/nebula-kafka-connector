/*
Copyright 2024 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ssl

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"time"

	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/cert"
)

const DefaultTimeout = 5 * time.Second

type Option func(ops *Options)

type Options struct {
	EnableMetaTLS bool
	Timeout       time.Duration
	TLSConfig     *tls.Config
}

func ClientOptions(nc *v2alpha1.NebulaCluster, certs *v2alpha1.SSLCertsSpec, enableMetaTLS bool, opts ...Option) ([]Option, error) {
	options := []Option{SetTimeout(DefaultTimeout)}
	if !enableMetaTLS {
		return options, nil
	}
	if certs == nil {
		return nil, errors.New("ssl certs is nil")
	}

	options = append(options, SetMetaTLS(true))
	caCert, clientCert, clientKey, err := getCerts(nc.Namespace, certs)
	if err != nil {
		return nil, fmt.Errorf("get cluster [%s/%s] certs failed: %v", nc.Namespace, nc.Name, err)
	}
	tlsConfig, err := cert.LoadTLSConfig(caCert, clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("load tls config failed: %v", err)
	}
	tlsConfig.ServerName = certs.ServerName
	tlsConfig.InsecureSkipVerify = pointer.BoolDeref(certs.InsecureSkipVerify, false)
	tlsConfig.MaxVersion = tls.VersionTLS12
	klog.V(4).Infof("tls config, ServerName: %s, InsecureSkipVerify: %v, MaxVersion: %d", tlsConfig.ServerName, tlsConfig.InsecureSkipVerify, tlsConfig.MaxVersion)
	options = append(options, SetTLSConfig(tlsConfig))
	options = append(options, opts...)
	return options, nil
}

func LoadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

func SetOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func SetTimeout(duration time.Duration) Option {
	return func(options *Options) {
		options.Timeout = duration
	}
}

func SetTLSConfig(config *tls.Config) Option {
	return func(options *Options) {
		options.TLSConfig = config
	}
}

func SetMetaTLS(e bool) Option {
	return func(options *Options) {
		options.EnableMetaTLS = e
	}
}

func getCerts(namespace string, cert *v2alpha1.SSLCertsSpec) ([]byte, []byte, []byte, error) {
	if os.Getenv("CA_CERT_PATH") != "" &&
		os.Getenv("CLIENT_CERT_PATH") != "" &&
		os.Getenv("CLIENT_KEY_PATH") != "" {
		caCert, err := os.ReadFile(os.Getenv("CA_CERT_PATH"))
		if err != nil {
			return nil, nil, nil, err
		}
		clientCert, err := os.ReadFile(os.Getenv("CLIENT_CERT_PATH"))
		if err != nil {
			return nil, nil, nil, err
		}
		clientKey, err := os.ReadFile(os.Getenv("CLIENT_KEY_PATH"))
		if err != nil {
			return nil, nil, nil, err
		}
		return caCert, clientCert, clientKey, nil
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	client, err := kube.NewClientSet(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	caSecret, err := client.Secret().GetSecret(namespace, cert.CASecret)
	if err != nil {
		return nil, nil, nil, err
	}
	caCert := caSecret.Data[cert.CACert]

	clientSecret, err := client.Secret().GetSecret(namespace, cert.ClientSecret)
	if err != nil {
		return nil, nil, nil, err
	}
	clientCert := clientSecret.Data[cert.ClientCert]
	clientKey := clientSecret.Data[cert.ClientKey]
	return caCert, clientCert, clientKey, nil
}
