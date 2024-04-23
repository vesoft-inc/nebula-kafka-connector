package config

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/config"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetComponentConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetComponentConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetComponentConfigLogic {
	return &GetComponentConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetComponentConfigLogic) GetComponentConfig(req *types.GetComponentConfigReq) (resp *types.GetComponentConfigResp, err error) {
	return config.NewConfigService(l.ctx, l.svcCtx).GetComponentConfig(req)
}
