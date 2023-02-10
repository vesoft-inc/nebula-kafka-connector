package reader

import (
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/spec"
)

const (
	DefaultBatchSize = 128
)

type (
	RecordReader interface {
		source.Sizer
		Read() (int, spec.Record, error)
	}
)

func NewRecordReader(s source.Source) RecordReader {
	// TODO: support other source formats
	return NewCSVReader(s)
}
