package storage

import (
	"fmt"
	"os"
	"strings"
)

// CheckEndpoint check  whether endpoint is combine with ip:port
// example: http://127.0.0.1:9999/xxx
func CheckEndpoint(endpoint string) bool {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	return strings.ContainsAny(endpoint, ":")
}

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("walk to %s failed: %w", path, err)
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}
