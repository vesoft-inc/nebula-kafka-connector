package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HdfsDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHdfsDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HdfsDownloadLogic {
	return &HdfsDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HdfsDownloadLogic) HdfsDownload(req *types.HDFSDownloadReq) (resp *types.HDFSDownloadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).HDFSDownload(req)
}
