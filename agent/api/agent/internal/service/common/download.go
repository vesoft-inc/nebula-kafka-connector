package common

import (
	"fmt"
	"os"
	"path"

	gopkgmiddleware "github.com/vesoft-inc/go-pkg/middleware"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/audit"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/ecode"
)

func (s *commonService) DownloadFile(req *types.DownloadFileReq) (*types.DownloadFileResp, error) {
	if err := audit.RecordOperation(s.ctx, audit.OpDownloadFile, fmt.Sprintf("download file from %s", req.Path)); err != nil {
		return nil, err
	}

	httpResp, ok := gopkgmiddleware.GetResponseWriter(s.ctx)
	if !ok {
		return nil, ecode.WithInternalServer(fmt.Errorf("unset ResponseWriter"))
	}

	body, err := os.ReadFile(req.Path)
	if err != nil {
		return nil, err
	}

	_, filename := path.Split(req.Path)

	httpResp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	httpResp.Header().Set("Content-Type", "application/octet-stream")

	if _, err = httpResp.Write(body); err != nil {
		return nil, err
	}

	return &types.DownloadFileResp{}, nil
}
