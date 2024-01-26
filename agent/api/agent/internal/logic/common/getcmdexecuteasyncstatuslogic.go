package common

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/common"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCmdExecuteAsyncStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCmdExecuteAsyncStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCmdExecuteAsyncStatusLogic {
	return &GetCmdExecuteAsyncStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCmdExecuteAsyncStatusLogic) GetCmdExecuteAsyncStatus(req *types.GetCmdExecuteAsyncStatusReq) (resp *types.GetCmdExecuteAsyncStatusResp, err error) {
	return common.NewCommonService(l.ctx, l.svcCtx).GetCmdExecuteAsyncStatus(req)
}
