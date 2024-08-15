package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocalDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLocalDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LocalDownloadLogic {
	return &LocalDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LocalDownloadLogic) LocalDownload(req *types.LocalDownloadReq) (resp *types.LocalDownloadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).LocalDownload(req)
}
