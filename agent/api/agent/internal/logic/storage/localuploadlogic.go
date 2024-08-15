package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocalUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLocalUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LocalUploadLogic {
	return &LocalUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LocalUploadLogic) LocalUpload(req *types.LocalUploadReq) (resp *types.LocalUploadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).LocalUpload(req)
}
