package configv5

import (
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	configbase "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/base"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/manager"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/reader"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/utils"
)

type (
	Manager struct {
		GraphName          string `yaml:"graphName"`
		configbase.Manager `yaml:",inline"`
	}
)

func (m *Manager) BuildManager(
	wgMap *utils.WaitGroupMap,
	l logger.Logger,
	pool client.Pool,
	sources Sources,
	opts ...manager.Option,
) (manager.Manager, error) {
	options := make([]manager.Option, 0, 8+len(opts))
	options = append(options,
		manager.WithClientPool(pool),
		manager.WithBatch(m.Batch),
		manager.WithReaderConcurrency(m.ReaderConcurrency),
		manager.WithImporterConcurrency(m.ImporterConcurrency),
		manager.WithStatsInterval(m.StatsInterval),
		manager.WithBeforeHooks(m.Hooks.Before...),
		manager.WithAfterHooks(m.Hooks.After...),
		manager.WithLogger(l),
	)
	options = append(options, opts...)

	mgr := manager.NewWithOpts(options...)

	for i := range sources {
		s := sources[i]
		src, brr, err := s.BuildSourceAndReader(reader.WithBatch(m.Batch), reader.WithLogger(l))
		if err != nil {
			return nil, err
		}

		importers, err := s.BuildImporters(wgMap, m.GraphName, pool)
		if err != nil {
			return nil, err
		}
		if err = mgr.Import(src, brr, importers...); err != nil {
			return nil, err
		}
	}

	return mgr, nil
}
