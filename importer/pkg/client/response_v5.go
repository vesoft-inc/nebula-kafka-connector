package client

import (
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type defaultResponseV5 struct {
	ResultSet nebula.Result
	respTime  time.Duration
	err       error
}

func newResponseV5(rs nebula.Result, respTime time.Duration, err error) Response {
	return defaultResponseV5{
		ResultSet: rs,
		respTime:  respTime,
		err:       err,
	}
}

func (resp defaultResponseV5) GetLatency() time.Duration {
	return time.Duration(resp.ResultSet.Latency()) * time.Microsecond
}

func (resp defaultResponseV5) GetRespTime() time.Duration {
	return resp.respTime
}

func (resp defaultResponseV5) GetError() error {
	return resp.err
}

func (defaultResponseV5) IsPermanentError() bool {
	// TODO: SYNTAX_ERROR, SEMANTIC_ERROR
	return false
}

func (defaultResponseV5) IsRetryMoreError() bool {
	// TODO: RAFT_BUFFER_OVERFLOW
	return false
}

func (resp defaultResponseV5) IsSucceed() bool {
	return resp.err == nil
}
