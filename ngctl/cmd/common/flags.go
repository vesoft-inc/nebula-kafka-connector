package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
)

// the config file used by all cmds of ngctl
var ConfigFile string

// Cache the spec produced by parsing the config file (a yaml file)
var ConfigSpec types.JobSpec

func GetAgentForHost(hostIP string) (agent types.Agent, err error) {
	for _, cluster := range ConfigSpec.Spec.Metad.Clusters {
		for _, h := range cluster.Graphd.Hosts {
			if h.IP == hostIP {
				return h.Agent, nil
			}
		}
		for _, h := range cluster.Storaged.Hosts {
			if h.IP == hostIP {
				return h.Agent, nil
			}
		}
	}
	return agent, fmt.Errorf("host not found")
}

func VerifyConfig(configFilePath string) (configSpec *types.JobSpec, err error) {
	spec, err := yamlparser.ParseYamlByPath(configFilePath)
	if err != nil {
		return nil, err
	}
	job := runner.NewJob("verify")
	err = job.Run("verify", nil, spec)
	if err != nil {
		log.Print("verify config file failed")
	}
	return spec, err
}

func CheckInConfigFile(filepath string) (err error) {
	if filepath == "" && ConfigFile == "" {
		return fmt.Errorf("config file is empty")
	}
	if filepath == "" && ConfigFile != "" {
		// use the exisiting config file
		return nil
	}
	if filepath != "" {
		spec, err := VerifyConfig(filepath)
		if err != nil {
			return err
		}
		// Use the new and verified config file
		ConfigFile = filepath
		ConfigSpec = *spec
	}
	if !strings.HasSuffix(ConfigSpec.InstallPath, "/") {
		ConfigSpec.InstallPath += "/"
	}
	return nil
}

// checking all hosts in the config file, including the metad, storaged, and graphd
func DeriveHostList(hostFromCmdLineOption string, clusterName string, needMetad bool) (hostList []IPAndPort, err error) {
	// hosts for metad
	dict := map[string]IPAndPort{}
	hosts := make([]types.Host, 0)
	if needMetad {
		hosts = append(hosts, ConfigSpec.Spec.Metad.Hosts...)
	}
	for _, cluster := range ConfigSpec.Spec.Metad.Clusters {
		if cluster.Name == clusterName {
			hosts = append(hosts, cluster.Graphd.Hosts...)
			hosts = append(hosts, cluster.Storaged.Hosts...)
		}
	}
	if len(hosts) == 0 {
		return hostList, fmt.Errorf("no valid host found")
	}
	for _, host := range hosts {
		// if host is specified in the command line, return the host directly
		// else, return all the hosts
		if hostFromCmdLineOption != "" {
			if host.IP == hostFromCmdLineOption {
				return []IPAndPort{{IP: host.IP, AgentPort: host.AgentPort}}, nil
			}
		} else {
			ipPort := IPAndPort{IP: host.IP, AgentPort: host.AgentPort}
			if _, ok := dict[host.IP]; !ok {
				dict[host.IP] = ipPort
			}
		}
	}
	for _, ipAndPort := range dict {
		hostList = append(hostList, ipAndPort)
	}
	if len(hostList) == 0 {
		return hostList, fmt.Errorf("no valid host found")
	}
	return hostList, nil
}

// checking all services in the config file, including only storaged and graphd
func DeriveServiceList(serviceFromCmdLineOption IPAndPort, clusterName string) (serviceList []IPAndPort, err error) {
	if serviceFromCmdLineOption.IP != "" && serviceFromCmdLineOption.Port != "" {
		return []IPAndPort{{IP: serviceFromCmdLineOption.IP, Port: serviceFromCmdLineOption.Port, ServiceType: serviceFromCmdLineOption.ServiceType}}, nil
	}
	// graphd and storaged are organized in clusters
	for _, cluster := range ConfigSpec.Spec.Metad.Clusters {
		if cluster.Name == clusterName {
			for _, host := range cluster.Graphd.Hosts {
				serviceList = append(serviceList, IPAndPort{IP: host.IP, Port: host.Port, AgentPort: host.AgentPort, ServiceType: "graphd"})
			}
			for _, host := range cluster.Storaged.Hosts {
				serviceList = append(serviceList, IPAndPort{IP: host.IP, Port: host.Port, AgentPort: host.AgentPort, ServiceType: "storaged"})
			}
		}
	}
	if len(serviceList) == 0 {
		return serviceList, fmt.Errorf("no valid service found")
	}
	return serviceList, nil
}
