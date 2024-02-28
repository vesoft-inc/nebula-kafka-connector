package meta

import (
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type (
	headerRequest struct {
		requestType string
		clusterId   int64 // meta admin dont need clusterId
	}
	HeaderResponse struct {
		OK      bool
		Code    nebula.ErrorCode
		Msg     string
		NewHost string // only if code is leader changed
		NewPort uint32 // only if code is leader changed
	}
)

func (h *HeaderResponse) GetErrorCode() nebula.ErrorCode {
	return h.Code
}

func (h *HeaderResponse) GetErrorMsg() string {
	return h.Msg
}

func (h *HeaderResponse) IsSucceeded() bool {
	return h.OK
}
