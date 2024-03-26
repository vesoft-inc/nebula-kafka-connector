package utils

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func GetUserClusterPath(installPath string, agentPath string) string {
	if agentPath != "" {
		return path.Join(agentPath)
	}
	return path.Join(installPath, "cluster/")
}

func GetClusterPath(installPath string) string {
	return path.Join(installPath, "cluster/")
}
func GetDownloadPath(installPath string) string {
	return path.Join(installPath, "download/")
}
func GetUserDownloadPath(installPath string, downloadPath string) string {
	if downloadPath != "" {
		return path.Join(downloadPath)
	}
	return path.Join(installPath, "download/")
}

func GetUserPackagePath(packagePath string, agentPackagePath string) string {
	if agentPackagePath != "" {
		return agentPackagePath
	}
	return packagePath
}

func GetUtilPath(installPath, utilName string) string {
	return path.Join(installPath, utilName)
}

func GetUserUtilPath(installPath, agentPath, utilName string) string {
	if agentPath != "" {
		return agentPath
	}
	return path.Join(installPath, utilName)
}
func GetMetaAddressListString(metaHosts []types.Agent, port string) string {
	if port == "" {
		port = "9559"
	}
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

func MergeNebulaConfigMap(configMap map[string]any, newConfigMap map[string]string) map[string]string {
	for k, v := range configMap {
		newConfigMap[k] = fmt.Sprintf("%v", v) // convert to string
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

func ParseTimeoutString(timeout string) (time.Duration, error) {
	lastStr := timeout[len(timeout)-1:]
	var timeType time.Duration
	if lastStr == "m" {
		timeType = time.Minute
	} else {
		timeType = time.Second
	}
	numberStr := timeout[:len(timeout)-1]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("parse timeout string error: %s", err)
	}
	return time.Duration(number) * timeType, nil
}

func GetConfigPort(config map[string]any) string {
	if configPort, ok := config["port"]; ok {
		return fmt.Sprintf("%v", configPort)
	}
	if configPort, ok := config["Port"]; ok {
		return fmt.Sprintf("%v", configPort)
	}
	return ""
}

// map component name
func GetAllUtilsProcess(spec *types.JobSpec) []*types.Process {
	allUtils := []*types.Process{}

	for _, process := range spec.UtilsProcesses {
		if process.Config == nil {
			process.Config = map[string]any{}
		}
		allUtils = append(allUtils, process)
	}
	return allUtils
}

func GetAllAgents(spec *types.JobSpec) []*types.Agent {
	allAgents := make(map[string]*types.Agent)
	allHosts := []types.Agent{}
	if spec.Spec.Metad != nil {
		allHosts = append(allHosts, spec.Spec.Metad.Hosts...)
		for _, cluster := range spec.Spec.Metad.Clusters {
			allHosts = append(allHosts, cluster.Graphd.Hosts...)
			allHosts = append(allHosts, cluster.Storaged.Hosts...)
		}
	}
	for _, process := range spec.UtilsProcesses {
		allHosts = append(allHosts, process.Hosts...)
	}
	for _, agent := range allHosts {
		if _, ok := allAgents[agent.Host]; !ok {
			allAgents[agent.Host] = &agent
		} else {
			if agent.SSHConfig != nil {
				allAgents[agent.Host].SSHConfig = agent.SSHConfig
			}
		}
	}
	var agents []*types.Agent
	for _, agent := range allAgents {
		agents = append(agents, agent)
	}
	return agents
}
