package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type S3DownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewS3DownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *S3DownloadLogic {
	return &S3DownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *S3DownloadLogic) S3Download(req *types.S3DownloadReq) (resp *types.S3DownloadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).S3Download(req)
}
