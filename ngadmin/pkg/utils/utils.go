package utils

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

func GetClusterPath(installPath string) string {
	return path.Join(installPath, "cluster/")
}

func GetDownloadPath(installPath string) string {
	return path.Join(installPath, "download/")
}

func GetMetaAddressListString(metaHosts []types.Agent, port string) string {
	metaAddressList := ""
	for _, meta := range metaHosts {

		metaAddressList += fmt.Sprintf("%s:%s,", RemoveAddressPort(meta.Host), port)
	}
	metaAddressList = metaAddressList[:len(metaAddressList)-1]
	return metaAddressList
}

func RemoveAddressPort(address string) string {
	return strings.Split(address, ":")[0]
}

func MergeConfigMap(configMap map[string]string, newConfigMap map[string]string) map[string]string {
	for k, v := range configMap {
		newConfigMap[k] = v
	}
	return newConfigMap
}

func GetUint32Port(portString string) (uint32, error) {
	port, err := strconv.ParseUint(portString, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(port), nil
}

func GetHostIP(host string) string {
	hostArr := strings.Split(host, ":")
	if len(hostArr) == 1 {
		return host
	}
	return hostArr[0]
}

func GetHttpsHost(host string) string {
	if len(host) < 8 || host[:8] != "https://" {
		return "https://" + host
	}
	return host
}
