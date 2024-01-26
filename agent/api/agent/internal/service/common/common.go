package common

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	commonService struct {
		logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
)

func NewCommonService(ctx context.Context, svcCtx *svc.ServiceContext) *commonService {
	return &commonService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
