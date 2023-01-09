package config

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
)

type Manager struct {
	Batch               int           `yaml:"batch,omitempty"`
	ReaderConcurrency   int           `yaml:"readerConcurrency,omitempty"`
	ImporterConcurrency int           `yaml:"importerConcurrency,omitempty"`
	StatsInterval       time.Duration `yaml:"statsInterval,omitempty"`
	Hooks               manager.Hooks `yaml:"hooks,omitempty"`
}

func (m *Manager) BuildManager(c client.Client, sources Sources, opts ...manager.Option) (manager.Manager, error) {
	options := make([]manager.Option, 0, 5+len(opts))
	options = append(options,
		manager.WithClient(c),
		manager.WithBatch(m.Batch),
		manager.WithReaderConcurrency(m.ReaderConcurrency),
		manager.WithImporterConcurrency(m.ImporterConcurrency),
		manager.WithStatsInterval(m.StatsInterval),
		manager.WithBeforeHooks(m.Hooks.Before...),
		manager.WithAfterHooks(m.Hooks.After...),
	)
	options = append(options, opts...)

	mgr := manager.NewWithOpts(options...)

	for i := range sources {
		s := sources[i]
		importers, err := s.BuildImporters(c)
		if err != nil {
			return nil, err
		}
		if err := mgr.Import(&s.SourceConfig, importers...); err != nil {
			return nil, err
		}
	}

	return mgr, nil
}
