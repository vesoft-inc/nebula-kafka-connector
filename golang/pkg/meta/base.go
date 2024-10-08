package meta

import (
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
)

type (
	headerRequest struct {
		requestType string
	}
	HeaderResponse struct {
		Code    errors.ErrorCode
		Msg     string
		NewHost string // only if code is leader changed
		NewPort uint32 // only if code is leader changed
	}
)

func (h *HeaderResponse) GetStatus() error {
	return errors.NewNebulaError(h.Code, "%s", h.Msg)
}

func (h *HeaderResponse) GetErrorCode() errors.ErrorCode {
	return h.Code
}

func (h *HeaderResponse) GetErrorMsg() string {
	return h.GetStatus().Error()
}

func (h *HeaderResponse) IsSucceeded() bool {
	return h.Code == errors.ERROR_SUCCESSFUL_COMPLETION
}
