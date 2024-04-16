package meta

import (
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type (
	headerRequest struct {
		requestType string
	}
	HeaderResponse struct {
		Code    nebula.ErrorCode
		Msg     string
		NewHost string // only if code is leader changed
		NewPort uint32 // only if code is leader changed
	}
)

func (h *HeaderResponse) GetStatus() error {
	return nebula.NewNebulaError(h.Code, h.Msg)
}

func (h *HeaderResponse) GetErrorCode() nebula.ErrorCode {
	return h.Code
}

func (h *HeaderResponse) GetErrorMsg() string {
	return h.GetStatus().Error()
}

func (h *HeaderResponse) IsSucceeded() bool {
	return h.Code == nebula.ERROR_SUCCESSFUL_COMPLETION
}
