package client

import (
	"fmt"
	"strings"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

type defaultResultSetV3 struct {
	*nebula.ResultSet
}

func newResultSetV3(rs *nebula.ResultSet) ResultSet {
	return defaultResultSetV3{
		ResultSet: rs,
	}
}

func (rs defaultResultSetV3) GetError() error {
	if rs.ResultSet.IsSucceed() {
		return nil
	}
	errorCode := rs.ResultSet.GetErrorCode()
	errorMsg := rs.ResultSet.GetErrorMsg()
	return fmt.Errorf("%d:%s", errorCode, errorMsg)
}

func (rs defaultResultSetV3) IsPermanentError() bool {
	switch rs.ResultSet.GetErrorCode() { //nolint:exhaustive
	default:
		return false
	case nebula.ErrorCode_E_SYNTAX_ERROR:
	case nebula.ErrorCode_E_SEMANTIC_ERROR:
	}
	return true
}

func (rs defaultResultSetV3) IsRetryMoreError() bool {
	errorMsg := rs.ResultSet.GetErrorMsg()
	// TODO: compare with E_RAFT_BUFFER_OVERFLOW
	// Can not get the E_RAFT_BUFFER_OVERFLOW inside storage now.
	return strings.Contains(errorMsg, "raft buffer is full")
}
