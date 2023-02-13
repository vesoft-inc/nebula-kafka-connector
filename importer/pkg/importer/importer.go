//go:generate mockgen -source=importer.go -destination importer_mock.go -package importer Importer
package importer

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/client"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

type (
	Importer interface {
		Wait()
		Import(records ...spec.Record) (*ImportResp, error)
		Done()
	}

	ImportResp struct {
		Latency  time.Duration
		RespTime time.Duration
	}

	ImportResult struct {
		Resp *ImportResp
		Err  error
	}

	Option func(*defaultImporter)

	defaultImporter struct {
		builder spec.StatementBuilder
		pool    client.Pool
		fnWait  func()
		fnDone  func()
	}
)

func New(builder spec.StatementBuilder, pool client.Pool, opts ...Option) Importer {
	options := make([]Option, 0, 2+len(opts))
	options = append(options, WithStatementBuilder(builder), WithClientPool(pool))
	options = append(options, opts...)
	return NewWithOpts(options...)
}

func NewWithOpts(opts ...Option) Importer {
	i := &defaultImporter{}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func WithStatementBuilder(builder spec.StatementBuilder) Option {
	return func(i *defaultImporter) {
		i.builder = builder
	}
}

func WithClientPool(p client.Pool) Option {
	return func(i *defaultImporter) {
		i.pool = p
	}
}

func WithWaitFunc(fn func()) Option {
	return func(i *defaultImporter) {
		i.fnWait = fn
	}
}

func WithDoneFunc(fn func()) Option {
	return func(i *defaultImporter) {
		i.fnDone = fn
	}
}

func (i *defaultImporter) Wait() {
	if i.fnWait != nil {
		i.fnWait()
	}
}

func (i *defaultImporter) Import(records ...spec.Record) (*ImportResp, error) {
	statement, err := i.builder.Build(records...)
	if err != nil {
		return nil, err
	}

	resp, err := i.pool.Execute(statement)
	if err != nil {
		return nil, errors.NewImportError(err).
			SetStatement(statement)
	}
	if !resp.IsSucceed() {
		return nil, errors.NewImportError(err, "the execute error is %s ", resp.GetError()).
			SetStatement(statement)
	}

	return &ImportResp{
		RespTime: resp.GetRespTime(),
		Latency:  resp.GetLatency(),
	}, nil
}

func (i *defaultImporter) Done() {
	if i.fnDone != nil {
		i.fnDone()
	}
}
