package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/runner"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
	"gopkg.in/yaml.v3"
)

const defaultMetadPort = 9559
const defaultGraphdPort = 9669
const defaultStoragedPort = 9779

type IPAndPort struct {
	// define both members as strings, since we only need them to fill config files
	IP          string
	Port        string
	AgentPort   string
	ServiceType string
}

type Config struct {
	JobSpec     *types.JobSpec
	Version     string      `yaml:"version"`
	InstallPath string      `yaml:"installPath"`
	PackagePath string      `yaml:"packagePath"`
	CaFile      string      `yaml:"caFile,omitempty"`
	CertFile    string      `yaml:"certFile,omitempty"`
	KeyFile     string      `yaml:"keyFile,omitempty"`
	Spec        *ConfigSpec `yaml:"spec"`
}

type ConfigSpec struct {
	Metad         *ServiceSpec        `yaml:"metad"`
	ServiceGroups []*ServiceGroupSpec `yaml:"serviceGroups"`
}

type ServiceGroupSpec struct {
	Name          string       `yaml:"name"`
	ReplicaFactor int          `yaml:"replicaFactor"`
	Graphd        *ServiceSpec `yaml:"graphd"`
	Storaged      *ServiceSpec `yaml:"storaged"`
}

type ServiceSpec struct {
	Config    map[string]any  `yaml:"config,omitempty"`
	Instances []*InstanceSpec `yaml:"instances,omitempty"`
}

type InstanceSpec struct {
	Host      string `yaml:"host"`
	Port      int
	AgentPort int `yaml:"agentPort"`
}

func verifyConfig(configFilePath string) (configSpec *types.JobSpec, err error) {
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

func NewConfigFromFile(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return NewConfigFromBytes(data)
}

func NewConfigFromBytes(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	for _, rule := range validateRules {
		if err := rule(&c); err != nil {
			return nil, err
		}
	}
	c.formatPort()
	if err := c.convertToJobSpec(); err != nil {
		return nil, err
	}
	if err := yamlparser.CheckMetaSpec(c.JobSpec.Spec.Metad); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) formatPort() {
	var metadPort, storagedPort, graphdPort int
	if c.Spec.Metad.Config == nil {
		metadPort = defaultGraphdPort
	} else {
		if port, ok := c.Spec.Metad.Config["port"]; ok {
			metadPort = port.(int)
		} else {
			metadPort = defaultMetadPort
		}
	}
	for _, inst := range c.Spec.Metad.Instances {
		inst.Port = metadPort
	}
	for _, svcgrp := range c.Spec.ServiceGroups {
		if svcgrp.Graphd.Config == nil {
			graphdPort = defaultGraphdPort
		} else {
			if port, ok := svcgrp.Graphd.Config["port"]; ok {
				graphdPort = port.(int)
			} else {
				graphdPort = defaultGraphdPort
			}
		}
		if svcgrp.Storaged.Config == nil {
			storagedPort = defaultStoragedPort
		} else {
			if port, ok := svcgrp.Storaged.Config["port"]; ok {
				storagedPort = port.(int)
			} else {
				storagedPort = defaultStoragedPort
			}
		}
		for _, inst := range svcgrp.Graphd.Instances {
			inst.Port = graphdPort
		}
		for _, inst := range svcgrp.Storaged.Instances {
			inst.Port = storagedPort
		}
	}
}

func (c *Config) convertToJobSpec() error {
	jobSpec := &types.JobSpec{
		Kind:        "NebulaCluster",
		Version:     "1.0",
		Rollback:    false,
		InstallPath: c.InstallPath,
		CertFile:    c.CertFile,
		KeyFile:     c.KeyFile,
		CAFile:      c.CaFile,
	}

	if c.Spec == nil {
		return fmt.Errorf("spec is empty")
	}
	if c.Spec.Metad == nil {
		return fmt.Errorf("metad is empty")
	}

	jobMetadProcess := c.convertServiceToJobSpec(c.Spec.Metad)
	jobMetaSpec := &types.MetadSpec{
		ServiceGroups: make([]types.ServiceGroup, 0),
	}
	jobMetaSpec.Process.Hosts = jobMetadProcess.Hosts
	jobMetaSpec.Process.Config = jobMetadProcess.Config
	jobMetaSpec.Process.InstallPath = jobMetadProcess.InstallPath
	jobMetaSpec.Process.PackagePath = jobMetadProcess.PackagePath

	for _, svcgrp := range c.Spec.ServiceGroups {
		jobGraphdProcess := c.convertServiceToJobSpec(svcgrp.Graphd)
		jobStoragedProcess := c.convertServiceToJobSpec(svcgrp.Storaged)
		jobServiceGroup := types.ServiceGroup{
			Name:     svcgrp.Name,
			Replica:  svcgrp.ReplicaFactor,
			Graphd:   *jobGraphdProcess,
			Storaged: *jobStoragedProcess,
		}
		jobMetaSpec.ServiceGroups = append(jobMetaSpec.ServiceGroups, jobServiceGroup)
	}
	jobSpec.Spec.Metad = jobMetaSpec
	c.JobSpec = jobSpec
	return nil
}

func (c *Config) convertServiceToJobSpec(svc *ServiceSpec) *types.Process {
	p := &types.Process{
		Config:      svc.Config,
		InstallPath: c.InstallPath,
		PackagePath: c.PackagePath,
	}
	for _, inst := range svc.Instances {
		host := types.Host{
			IP:   inst.Host,
			Port: strconv.Itoa(inst.Port),
			Agent: types.Agent{
				Host: fmt.Sprintf("%s:%d", inst.Host, inst.AgentPort),
			},
		}
		p.Hosts = append(p.Hosts, host)
	}
	return p
}

func (c *Config) GetServiceInstances(serviceGroup string, serviceType string, host string) ([]*InstanceSpec, error) {
	var instances []*InstanceSpec
	if serviceType == "metad" {
		instances = c.Spec.Metad.Instances
	} else {
		var svcgrp *ServiceGroupSpec
		for _, sg := range c.Spec.ServiceGroups {
			if sg.Name == serviceGroup {
				svcgrp = sg
				break
			}
		}
		if svcgrp == nil {
			return nil, fmt.Errorf("service group %s not found", serviceGroup)
		}
		switch serviceType {
		case "graphd":
			instances = svcgrp.Graphd.Instances
		case "storaged":
			instances = svcgrp.Storaged.Instances
		default:
			return nil, fmt.Errorf("invalid service type %s", serviceType)
		}
	}
	if len(instances) == 0 {
		return nil, fmt.Errorf("no service instance found")
	}
	if host == "" {
		return instances, nil
	}
	for _, inst := range instances {
		if inst.Host == host {
			return []*InstanceSpec{inst}, nil
		}
	}
	return nil, fmt.Errorf("host not found")
}

func (c *Config) GetInstanceFromHost(host string) (*InstanceSpec, error) {
	for _, inst := range c.Spec.Metad.Instances {
		if inst.Host == host {
			return inst, nil
		}
	}
	for _, svcgrp := range c.Spec.ServiceGroups {
		for _, inst := range svcgrp.Graphd.Instances {
			if inst.Host == host {
				return inst, nil
			}
		}
		for _, inst := range svcgrp.Storaged.Instances {
			if inst.Host == host {
				return inst, nil
			}
		}
	}
	return nil, fmt.Errorf("host not found")
}

func (c *Config) GetMetadAddress() string {
	metadAddresses := make([]string, 0)
	for _, inst := range c.Spec.Metad.Instances {
		metadAddresses = append(metadAddresses, fmt.Sprintf("%s:%d", inst.Host, inst.Port))
	}
	return strings.Join(metadAddresses, ",")
}

func getPortFromAddress(address string) string {
	if strings.Contains(address, ":") {
		return strings.Split(address, ":")[1]
	}
	return ""
}
