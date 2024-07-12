package audit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	gopkgmiddleware "github.com/vesoft-inc/go-pkg/middleware"
)

var (
	mu      sync.Mutex
	logFile *os.File
)

type OperationType string

const (
	OpExecuteCmd    OperationType = "execute_cmd"
	OpUploadFile    OperationType = "upload_file"
	OpDownloadFile  OperationType = "download_file"
	OpS3Upload      OperationType = "s3_upload"
	OpS3Download    OperationType = "s3_download"
	OpLocalUpload   OperationType = "local_upload"
	OpLocalDownload OperationType = "local_download"
	OpHDFSUpload    OperationType = "hdfs_upload"
	OpHDFSDownload  OperationType = "hdfs_download"
)

func InitLogFile(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	logFile = file
	return nil
}

func RecordOperation(ctx context.Context, typ OperationType, operation string) error {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return fmt.Errorf("log file is not initialized")
	}

	httpReq, ok := gopkgmiddleware.GetRequest(ctx)
	if !ok {
		return fmt.Errorf("unset KeepRequest")
	}

	logEntry := fmt.Sprintf("start_time: %s,client_addr: %s,type: %s,operation: %s\n",
		time.Now().Format(time.RFC3339), httpReq.RemoteAddr, typ, operation)
	_, err := logFile.WriteString(logEntry)

	return err
}
