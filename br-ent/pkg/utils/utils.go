package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

// DumpMetaToFile serializes meta into JSON and writes it to a file.
func DumpMetaToFile(backupRes *meta.CreateBackupResp, filepath string) error {
	jsonData, err := json.Marshal(backupRes)
	if err != nil {
		return err
	}

	// Create file if it does not exist, truncate it if it does
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// ParseMetaFromFile reads a file and deserializes its content into a CreateBackupResp object.
func ParseMetaFromFile(filename string) (*meta.CreateBackupResp, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var resp *meta.CreateBackupResp
	err = json.Unmarshal(jsonData, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

const (
	LocalTmpDir = "/tmp/br-ent"
)

func EnsureDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("ensure dirs %s failed: %w", dir, err)
	}
	return nil
}

func RemoveDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("remove tmp dirs %s failed: %w", dir, err)
	}
	return nil
}

func IsBackupName(path string) bool {
	return strings.HasPrefix(path, "BACKUP")
}

func UriJoin(elem ...string) (string, error) {
	if len(elem) == 0 {
		return "", fmt.Errorf("empty paths")
	}

	if len(elem) == 1 {
		return elem[0], nil
	}

	u, err := url.Parse(elem[0])
	if err != nil {
		return "", fmt.Errorf("parse base uri %s failed: %w", elem[0], err)
	}

	elem[0] = u.Path
	u.Path = path.Join(elem...)
	return u.String(), nil
}

//func ToRole(r meta.HostRole) pb.ServiceRole {
//	switch r {
//	case meta.HostRole_STORAGE:
//		return pb.ServiceRole_STORAGE
//	case meta.HostRole_GRAPH:
//		return pb.ServiceRole_GRAPH
//	case meta.HostRole_META:
//		return pb.ServiceRole_META
//	default:
//		return pb.ServiceRole_UNKNOWN_ROLE
//	}
//}

func ErrorAggregate(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("%s", strings.Join(messages, ",\n"))
}

func IsNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not exist")
}
