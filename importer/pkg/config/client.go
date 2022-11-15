package config

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
)

type Client struct {
	Addresses                []string      `yaml:"addresses"`
	User                     string        `yaml:"user,omitempty"`
	Password                 string        `yaml:"password,omitempty"`
	ConcurrencyPerAddress    int           `yaml:"concurrencyPerAddress,omitempty"`
	ReconnectInitialInterval time.Duration `yaml:"reconnectInitialInterval,omitempty"`
	Retry                    int           `yaml:"retry,omitempty"`
	RetryInitialInterval     time.Duration `yaml:"retryInitialInterval,omitempty"`
}

func (c *Client) BuildClient(opts ...client.Option) (client.Client, error) {
	options := make([]client.Option, 0, 4+len(opts))
	options = append(
		options,
		client.WithAddress(c.Addresses...),
		client.WithUserPassword(c.User, c.Password),
		client.WithReconnectInitialInterval(c.ReconnectInitialInterval),
		client.WithRetry(c.Retry),
		client.WithRetryInitialInterval(c.RetryInitialInterval),
		client.WithConcurrencyPerAddress(c.ConcurrencyPerAddress),
	)
	options = append(options, opts...)
	cli := client.New(options...)

	if err := cli.Open(); err != nil {
		return nil, err
	}
	return cli, nil
}
