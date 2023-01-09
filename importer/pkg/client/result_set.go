//go:generate mockgen -source=result_set.go -destination result_set_mock.go -package client ResultSet
package client

import (
	stderrors "errors"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type (
	ResultSet interface {
		GetStatus() string
		IsSucceed() bool
		GetLatency() int64
		GetError() error
		IsPermanentError() bool
		IsRetryMoreError() bool
	}

	defaultResultSet struct {
		*nebula.ResultSet
	}
)

func newResultSet(rs *nebula.ResultSet) ResultSet {
	return defaultResultSet{
		ResultSet: rs,
	}
}

func (rs defaultResultSet) GetError() error {
	if rs.ResultSet.IsSucceed() {
		return nil
	}
	return stderrors.New(rs.ResultSet.GetStatus())
}

func (defaultResultSet) IsPermanentError() bool {
	// TODO: SYNTAX_ERROR, SEMANTIC_ERROR
	return false
}

func (defaultResultSet) IsRetryMoreError() bool {
	// TODO: RAFT_BUFFER_OVERFLOW
	return false
}
