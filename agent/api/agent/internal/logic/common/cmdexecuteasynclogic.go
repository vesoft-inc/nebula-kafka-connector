package common

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/common"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CmdExecuteAsyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCmdExecuteAsyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CmdExecuteAsyncLogic {
	return &CmdExecuteAsyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CmdExecuteAsyncLogic) CmdExecuteAsync(req *types.CmdExecuteAsyncReq) (resp *types.CmdExecuteAsyncResp, err error) {
	return common.NewCommonService(l.ctx, l.svcCtx).CmdExecuteAsync(req)
}
