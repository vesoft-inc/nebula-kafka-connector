package common

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/common"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CmdExecuteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCmdExecuteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CmdExecuteLogic {
	return &CmdExecuteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CmdExecuteLogic) CmdExecute(req *types.CmdExecuteReq) (resp *types.CmdExecuteResp, err error) {
	return common.NewCommonService(l.ctx, l.svcCtx).CmdExecute(req)
}
