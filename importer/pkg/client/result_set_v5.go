//go:generate mockgen -source=result_set.go -destination result_set_mock.go -package client ResultSet
package client

import (
	stderrors "errors"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type defaultResultSetV5 struct {
	*nebula.ResultSet
}

func newResultSetV5(rs *nebula.ResultSet) ResultSet {
	return defaultResultSetV5{
		ResultSet: rs,
	}
}

func (rs defaultResultSetV5) GetError() error {
	if rs.ResultSet.IsSucceed() {
		return nil
	}
	return stderrors.New(rs.ResultSet.GetStatus())
}

func (defaultResultSetV5) IsPermanentError() bool {
	// TODO: SYNTAX_ERROR, SEMANTIC_ERROR
	return false
}

func (defaultResultSetV5) IsRetryMoreError() bool {
	// TODO: RAFT_BUFFER_OVERFLOW
	return false
}
