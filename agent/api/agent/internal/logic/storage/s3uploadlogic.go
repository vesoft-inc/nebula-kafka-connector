package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type S3UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewS3UploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *S3UploadLogic {
	return &S3UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *S3UploadLogic) S3Upload(req *types.S3UploadReq) (resp *types.S3UploadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).S3Upload(req)
}
