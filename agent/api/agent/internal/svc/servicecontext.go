package svc

import (
	"database/sql"
	"errors"

	"github.com/vesoft-inc/go-pkg/httpclient"
	"github.com/vesoft-inc/go-pkg/response"
	"github.com/vesoft-inc/go-pkg/validator"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/config"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/ecode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config          config.Config
	ResponseHandler response.Handler
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		ResponseHandler: createResponseHandler(c),
	}
}

func createResponseHandler(c config.Config) response.Handler { //nolint:gocritic
	detailsType := response.StandardHandlerDetailsNone
	if c.Debug.Enable {
		detailsType = response.StandardHandlerDetailsFull
	}
	return response.NewStandardHandler(response.StandardHandlerParams{
		GetErrCode: func(err error) *ecode.ErrCode {
			return ecode.TakeCodePriority(func() *ecode.ErrCode {
				if _, ok := err.(validator.ValidationErrors); ok {
					return ecode.ErrParam
				}
				return nil
			}, func() *ecode.ErrCode {
				if errors.Is(err, sql.ErrNoRows) {
					return ecode.ErrNotFound
				}
				return nil
			}, func() *ecode.ErrCode {
				if e, ok := httpclient.AsResponseError(err); ok {
					return ecode.GetErrCodeByHTTPStatus(e.GetResponse().StatusCode())
				}
				return nil
			}, func() *ecode.ErrCode {
				return ecode.ErrInternalServer
			})
		},
		Errorf:      logx.Errorf,
		DetailsType: detailsType,
	})
}
