package storage

import (
	"context"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	storageService struct {
		logx.Logger
		ctx    context.Context
		svcCtx *svc.ServiceContext
	}
)

func NewStorageService(ctx context.Context, svcCtx *svc.ServiceContext) *storageService {
	return &storageService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
