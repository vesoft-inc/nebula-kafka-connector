package common

import (
	"fmt"
	"log"
	"net"
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

func IsValidIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}
