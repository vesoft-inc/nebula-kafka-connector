package common

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

var (
	MetaClient meta.Client
)

func MetaClientInit() error {
	if MetaClient != nil {
		return nil
	}
	cacheToken, err := LoadMetaToken()
	if err != nil || cacheToken == nil {
		return fmt.Errorf("load meta session failed, please login first.")
	}
	MetaClient, err = meta.NewMetaClient(cacheToken.Leader, meta.WithToken(cacheToken.Token))
	if err != nil {
		return err
	}
	return nil
}

func MetaClientClose() {
	if MetaClient != nil {
		MetaClient.Close()
	}
}

func NgctlError(message string, err string) error {
	if err != "" {
		return fmt.Errorf("%s, err: %s", message, err)
	} else {
		return fmt.Errorf("%s", message)
	}
}

// common.MetaOutput is the output of meta command
// using os.Stdout by default, and could use other output for testing
var MetaOutput io.Writer = os.Stdout

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
