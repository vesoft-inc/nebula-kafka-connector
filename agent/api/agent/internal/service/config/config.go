package config

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	configService struct {
		logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
)

func NewConfigService(ctx context.Context, svcCtx *svc.ServiceContext) *configService {
	return &configService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
