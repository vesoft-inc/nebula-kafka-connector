//go:generate mockgen -source=result_set.go -destination result_set_mock.go -package client ResultSet
package client

import (
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	nebulathrift "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula"
	graphthrift "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula/graph"
)

type (
	ResultSet interface {
		AsStringTable() [][]string
		GetStatus() string
		IsSucceed() bool
		IsSetData() bool
		GetRows() []*nebulathrift.RawRecord
		GetRowSize() int
		GetColNames() []string
		GetColSize() int
		GetLatency() int64
		GetRowValuesByIndex(index int) (*nebula.Record, error)
		IsSetPlanDesc() bool
		GetPlanDesc() *graphthrift.PlanDescription
	}

	defaultResultSet struct {
		*nebula.ResultSet
	}
)

func NewResultSet(rs *nebula.ResultSet) ResultSet {
	return defaultResultSet{
		ResultSet: rs,
	}
}
