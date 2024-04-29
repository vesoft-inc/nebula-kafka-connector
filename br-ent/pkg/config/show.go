package config

import (
	"fmt"

	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/storage"

	"github.com/spf13/pflag"
)

type ShowConfig struct {
	Backend *agentstorage.Backend
}

func (s *ShowConfig) ParseFlags(flags *pflag.FlagSet) error {
	backend, err := storage.ParseFromFlags(flags)
	if err != nil {
		return fmt.Errorf("parse storage flags failed: %w", err)
	}
	s.Backend = backend
	return nil
}
