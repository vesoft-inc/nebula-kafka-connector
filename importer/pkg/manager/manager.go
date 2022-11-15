//go:generate mockgen -source=manager.go -destination manager_mock.go -package manager Manager
package manager

import (
	"io"
	"sync"
	"time"

	"github.com/panjf2000/ants"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/importer"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/reader"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/stats"
)

const (
	DefaultReaderConcurrency   = 10
	DefaultImporterConcurrency = 10
	DefaultStatsInterval       = time.Second * 10
)

var sourceOpen = source.Open

type (
	Manager interface {
		ImportNode(name string, sourceConfigs ...*source.Config) error
		ImportEdge(name string, sourceConfigs ...*source.Config) error
		Start() error
		Wait() error
		Stats() *stats.Stats
		Stop() error
	}

	defaultManager struct {
		graph               *spec.Graph
		c                   client.Client
		stats               *stats.ConcurrencyStats
		batch               int
		readerConcurrency   int
		readerWaitGroup     sync.WaitGroup
		readerPool          *ants.Pool
		importerConcurrency int
		importerWaitGroup   sync.WaitGroup
		importerPool        *ants.Pool
		statsInterval       time.Duration
		hooks               *Hooks
		chStart             chan struct{}
		done                chan struct{}
		wgNodes             *WaitGroupMap
		logger              logger.Logger
	}

	Option func(*defaultManager)
)

func New(graph *spec.Graph, c client.Client, opts ...Option) Manager {
	options := make([]Option, 0, 2+len(opts))
	options = append(options, WithGraph(graph), WithClient(c))
	options = append(options, opts...)
	return NewWithOpts(options...)
}

func NewWithOpts(opts ...Option) Manager {
	m := &defaultManager{
		stats:               stats.NewConcurrencyStats(),
		readerConcurrency:   DefaultReaderConcurrency,
		importerConcurrency: DefaultImporterConcurrency,
		statsInterval:       DefaultStatsInterval,
		hooks:               &Hooks{},
		chStart:             make(chan struct{}),
		done:                make(chan struct{}),
		wgNodes:             NewWaitGroups(),
	}

	for _, opt := range opts {
		opt(m)
	}

	m.readerPool, _ = ants.NewPool(m.readerConcurrency)
	m.importerPool, _ = ants.NewPool(m.importerConcurrency)

	if m.logger == nil {
		m.logger = logger.NopLogger
	}

	return m
}

func WithGraph(graph *spec.Graph) Option {
	return func(m *defaultManager) {
		m.graph = graph
	}
}

func WithClient(c client.Client) Option {
	return func(m *defaultManager) {
		m.c = c
	}
}

func WithBatch(batch int) Option {
	return func(m *defaultManager) {
		if batch > 0 {
			m.batch = batch
		}
	}
}
func WithReaderConcurrency(concurrency int) Option {
	return func(m *defaultManager) {
		if concurrency > 0 {
			m.readerConcurrency = concurrency
		}
	}
}

func WithImporterConcurrency(concurrency int) Option {
	return func(m *defaultManager) {
		if concurrency > 0 {
			m.importerConcurrency = concurrency
		}
	}
}

func WithStatsInterval(statsInterval time.Duration) Option {
	return func(m *defaultManager) {
		if statsInterval > 0 {
			m.statsInterval = statsInterval
		}
	}
}

func WithBeforeHooks(hooks ...*Hook) Option {
	return func(m *defaultManager) {
		m.hooks.Before = hooks
	}
}

func WithAfterHooks(hooks ...*Hook) Option {
	return func(m *defaultManager) {
		m.hooks.After = hooks
	}
}

func WithLogger(l logger.Logger) Option {
	return func(m *defaultManager) {
		m.logger = l
	}
}

func (m *defaultManager) ImportNode(name string, sourceConfigs ...*source.Config) error {
	n, ok := m.graph.GetNodeByName(name)
	if !ok {
		err := errors.NewImportError(errors.ErrNodeNotFound, "manager: add import nodes failed").SetNodeName(name)
		m.logError(err, "")
		return err
	}
	i := importer.NewNodeImporter(m.graph, n, m.c)
	return m.importSources(i, sourceConfigs...)
}

func (m *defaultManager) ImportEdge(name string, sourceConfigs ...*source.Config) error {
	e, ok := m.graph.GetEdgeByName(name)
	if !ok {
		err := errors.NewImportError(errors.ErrEdgeNotFound, "manager: add import edges failed").SetEdgeName(name)
		m.logError(err, "")
		return err
	}
	i := importer.NewEdgeImporter(m.graph, e, m.c)
	return m.importSources(i, sourceConfigs...)
}

func (m *defaultManager) Start() error {
	m.logger.Info("manager: starting")

	m.stats.Init()

	if err := m.Before(); err != nil {
		return err
	}

	close(m.chStart)

	go m.loopPrintStats()
	m.logger.Info("manager: start successfully")
	return nil
}

func (m *defaultManager) Wait() error {
	m.logger.Info("manager: wait")

	m.readerWaitGroup.Wait()
	m.importerWaitGroup.Wait()

	if err := m.After(); err != nil {
		return err
	}

	m.logStats()

	m.logger.Info("manager: wait successfully")
	return nil
}

func (m *defaultManager) Stats() *stats.Stats {
	return m.stats.Stats()
}

func (m *defaultManager) Stop() (err error) {
	m.logger.Info("manager: stop")
	defer func() {
		if err != nil {
			err = errors.NewImportError(err, "manager: stop failed").SetGraphName(m.graph.Name)
			m.logError(err, "")
		} else {
			m.logger.Info("manager: stop successfully")
		}
	}()
	close(m.done)
	return m.Wait()
}

func (m *defaultManager) Before() error {
	m.logger.Info("manager: exec before hook")
	return m.execHooks(BeforeHook)
}

func (m *defaultManager) After() error {
	m.logger.Info("manager: exec after hook")
	return m.execHooks(AfterHook)
}

func (m *defaultManager) execHooks(name HookName) error {
	var hooks []*Hook
	switch name {
	case BeforeHook:
		hooks = m.hooks.Before
	case AfterHook:
		hooks = m.hooks.After
	}
	if len(hooks) == 0 {
		return nil
	}
	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		for _, statement := range hook.Statements {
			if statement == "" {
				continue
			}
			rs, err := m.c.Execute(statement)
			if err != nil {
				err = errors.NewImportError(err,
					"manager: exec failed in %s hook", name,
				).SetGraphName(m.graph.Name).SetStatement(statement)
				m.logError(err, "")
				return err
			}
			if !rs.IsSucceed() {
				err = errors.NewImportError(err,
					"manager: exec failed in %s hook, %s", name, rs.GetStatus(),
				).SetGraphName(m.graph.Name).SetStatement(statement)
				m.logError(err, "")
				return err
			}
		}
		if hook.Wait != 0 {
			time.Sleep(hook.Wait)
		}
	}
	return nil
}

func (m *defaultManager) importSources(i importer.Importer, sourceConfigs ...*source.Config) error {
	for _, sourceConfig := range sourceConfigs {
		if err := m.importSource(i, sourceConfig); err != nil {
			return err
		}
	}
	return nil
}

func (m *defaultManager) importSource(i importer.Importer, sourceConfig *source.Config) error {
	log := m.logger.With(logger.Field{Key: "source", Value: sourceConfig.String()})
	s, err := sourceOpen(sourceConfig)
	if err != nil {
		err = errors.NewImportError(err, "manager: open import source failed").SetGraphName(m.graph.Name)
		m.logError(err, "")
		return err
	}

	nBytes, err := s.Size()
	if err != nil {
		_ = s.Close()
		err = errors.NewImportError(err, "manager: get size of import source failed").SetGraphName(m.graph.Name)
		m.logError(err, "")
		return err
	}
	m.stats.AddTotalBytes(nBytes)

	rr := reader.NewRecordReader(s)
	bcr := reader.NewBatchRecordReader(rr, m.batch)

	node := i.Node()
	if node != nil {
		m.wgNodes.Add(node.Name, 1)
	}
	edge := i.Edge()

	m.readerWaitGroup.Add(1)
	cleanup := func() {
		if node != nil {
			m.wgNodes.Done(node.Name)
		}
		m.readerWaitGroup.Done()
		s.Close()
	}

	go func() {
		err = m.readerPool.Submit(func() {
			<-m.chStart
			defer cleanup()

			if edge != nil {
				m.wgNodes.Wait(edge.Src.Name)
				m.wgNodes.Wait(edge.Dst.Name)
			}
			_ = m.loopImport(i, bcr)
		})
		if err != nil {
			cleanup()
			m.logError(err, "manager: submit reader failed")
		}
	}()

	log.Info("manager: add import source successfully")
	return nil
}

func (m *defaultManager) loopImport(i importer.Importer, r reader.BatchRecordReader) error {
	for {
		select {
		case <-m.done:
			return nil
		default:
			nBytes, records, err := r.ReadBatch()
			if err != nil {
				if err != io.EOF {
					err = errors.NewImportError(err, "manager: read batch failed").SetGraphName(m.graph.Name)
					m.logError(err, "")
					return err
				}
				return nil
			}
			m.submitImporterTask(i, nBytes, records)
		}
	}
}

func (m *defaultManager) submitImporterTask(i importer.Importer, nBytes int, records spec.Records) {
	m.importerWaitGroup.Add(1)
	if err := m.importerPool.Submit(func() {
		defer m.importerWaitGroup.Done()
		result, err := i.Import(records...)
		if err != nil {
			m.logError(err, "manager: import failed")
			m.onFailed(nBytes, records)
		} else {
			m.onSucceeded(nBytes, records, result)
		}
	}); err != nil {
		m.importerWaitGroup.Done()
		m.logError(err, "manager: submit importer failed")
	}
}

func (m *defaultManager) loopPrintStats() {
	if m.statsInterval <= 0 {
		return
	}
	ticker := time.NewTicker(m.statsInterval)
	m.logStats()
	for {
		select {
		case <-ticker.C:
			m.logStats()
		case <-m.done:
			return
		}
	}
}

func (m *defaultManager) logStats() {
	m.logger.Info(m.Stats().String())
}

func (m *defaultManager) onFailed(nBytes int, records spec.Records) {
	m.stats.Failed(int64(nBytes), int64(len(records)))
}

func (m *defaultManager) onSucceeded(nBytes int, records spec.Records, result *importer.ImportResp) {
	m.stats.Succeeded(int64(nBytes), int64(len(records)), result.Latency, result.ReqTime)
}

func (m *defaultManager) logError(err error, msg string, fields ...logger.Field) { //nolint:unparam
	e := errors.AsOrNewImportError(err)
	fields = append(fields, logger.MapToFields(e.Fields())...)
	m.logger.SkipCaller(1).WithError(e.Cause()).Error(msg, fields...)
}
