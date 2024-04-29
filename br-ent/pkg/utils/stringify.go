package utils

import (
	"fmt"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func StringifyBackup(b *meta.CreateBackupResp) string {
	m := map[string]string{
		"backup name":  b.BackupName,
		"created time": time.Unix(b.CreateTime/1000, 0).Local().String(),
	}

	//s := make([]string, 0, len(b.MetaBackups))
	//for _, f := range b.MetaBackups {
	//	s = append(s, f)
	//}
	//m["meta files"] = strings.Join(s, ",")

	return fmt.Sprintf("%v", m)
}
