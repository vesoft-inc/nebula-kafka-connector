package reader

import "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/source"

type (
	baseReader struct {
		s source.Source
	}
)

func (r *baseReader) Source() source.Source {
	return r.s
}
