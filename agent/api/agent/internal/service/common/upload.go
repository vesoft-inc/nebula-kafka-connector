package common

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"

	gopkgmiddleware "github.com/vesoft-inc/go-pkg/middleware"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/ecode"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/utils"
)

func (s *commonService) UploadFile(req *types.UploadFileReq) (*types.UploadFileResp, error) {
	httpReq, ok := gopkgmiddleware.GetRequest(s.ctx)
	if !ok {
		return nil, ecode.WithInternalServer(fmt.Errorf("unset KeepRequest"))
	}

	if err := httpReq.ParseMultipartForm(500 << 20); err != nil {
		if err != nil {
			return nil, err
		}
	}

	file, header, err := httpReq.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = validateFile(header); err != nil {
		return nil, err
	}
	if _, err = os.Stat(req.Path); os.IsNotExist(err) {
		if err = os.MkdirAll(req.Path, 0755); err != nil {
			return nil, err
		}
	}

	filePath := path.Join(req.Path, header.Filename)
	if _, err = utils.SaveFormFile(header, filePath); err != nil {
		return nil, err
	}
	return &types.UploadFileResp{
		Data: filePath,
	}, nil
}

func validateFile(fh *multipart.FileHeader) error {
	if fh.Size > 500*1024*1024 {
		return fmt.Errorf("file size is too large")
	}
	return nil
}
