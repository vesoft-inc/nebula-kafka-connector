package client

import (
	stderrors "errors"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type defaultResponseV5 struct {
	*nebula.ResultSet
	respTime time.Duration
}

func newResponseV5(rs *nebula.ResultSet, respTime time.Duration) Response {
	return defaultResponseV5{
		ResultSet: rs,
		respTime:  respTime,
	}
}

func (resp defaultResponseV5) GetLatency() time.Duration {
	return time.Duration(resp.ResultSet.GetLatency()) * time.Microsecond
}

func (resp defaultResponseV5) GetRespTime() time.Duration {
	return resp.respTime
}

func (resp defaultResponseV5) GetError() error {
	if resp.ResultSet.IsSucceed() {
		return nil
	}
	return stderrors.New(resp.ResultSet.GetStatus())
}

func (defaultResponseV5) IsPermanentError() bool {
	// TODO: SYNTAX_ERROR, SEMANTIC_ERROR
	return false
}

func (defaultResponseV5) IsRetryMoreError() bool {
	// TODO: RAFT_BUFFER_OVERFLOW
	return false
}
