package utils

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func GetUserCluster(installPath string, agentPath string) string {
	if agentPath != "" {
		return path.Join(agentPath, "/")
	}
	return path.Join(installPath, "/")
}

func GetCluster(installPath string) string {
	return path.Join(installPath, "/")
}
func GetDownloadPath(installPath string) string {
	return path.Join(installPath, "/")
}
func GetUserDownloadPath(installPath string, downloadPath string) string {
	if downloadPath != "" {
		return path.Join(downloadPath, "/")
	}
	return path.Join(installPath, "/")
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
		return path.Join(agentPath, utilName)
	}
	return path.Join(installPath, utilName)
}
func GetMetaAddressListString(metaHosts []types.Host, port string) string {
	if port == "" {
		port = "9559"
	}
	metaAddressList := ""
	for _, meta := range metaHosts {

		metaAddressList += fmt.Sprintf("%s:%s,", RemoveAddressPort(meta.Agent.Host), port)
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
func GetHostPort(host string) int {
	hostArr := strings.Split(host, ":")
	if len(hostArr) == 1 {
		return 0
	}
	port, err := GetUint32Port(hostArr[1])
	if err != nil {
		return 0
	}
	return int(port)
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
		for _, host := range spec.Spec.Metad.Hosts {
			allHosts = append(allHosts, host.Agent)
		}
		for _, cluster := range spec.Spec.Metad.ServiceGroups {
			for _, host := range cluster.Graphd.Hosts {
				allHosts = append(allHosts, host.Agent)
			}
			for _, host := range cluster.Storaged.Hosts {
				allHosts = append(allHosts, host.Agent)
			}
		}
	}
	for _, process := range spec.UtilsProcesses {
		for _, host := range process.Hosts {
			allHosts = append(allHosts, host.Agent)
		}
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

func MergeAgentConfig(agentConfig map[string]any, newConfig map[string]any) map[string]any {
	if agentConfig == nil {
		agentConfig = map[string]any{}
	}
	for k, v := range newConfig {
		(agentConfig)[k] = v
	}
	return agentConfig
}
