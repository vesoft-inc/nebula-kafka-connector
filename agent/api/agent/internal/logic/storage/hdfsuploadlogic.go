package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/service/storage"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HdfsUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHdfsUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HdfsUploadLogic {
	return &HdfsUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HdfsUploadLogic) HdfsUpload(req *types.HDFSUploadReq) (resp *types.HDFSUploadResp, err error) {
	return storage.NewStorageService(l.ctx, l.svcCtx).HDFSUpload(req)
}
