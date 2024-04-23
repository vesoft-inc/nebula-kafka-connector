package config

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/config"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetComponentConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetComponentConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetComponentConfigLogic {
	return &SetComponentConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetComponentConfigLogic) SetComponentConfig(req *types.SetComponentConfigReq) (resp *types.SetComponentConfigResp, err error) {
	return config.NewConfigService(l.ctx, l.svcCtx).SetComponentConfig(req)
}
