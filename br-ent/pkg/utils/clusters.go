package utils

import (
	"fmt"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type NebulaClusters struct {
	clusters map[int64][]*clients.ServiceInfo // clusterId -> ([storaged, graphd])
}

func NewNebulaClusters(clusters []*clients.ClusterServiceInfo, amg *clients.AgentManager) (*NebulaClusters, error) {
	clusterMap := make(map[int64][]*clients.ServiceInfo)
	for _, cluster := range clusters {
		services := make([]*clients.ServiceInfo, 0)
		for _, service := range cluster.Services {
			agent, err := amg.GetAgent(service.Host)
			if err != nil {
				return nil, fmt.Errorf("get agent for %s failed: %w", service.Host, err)
			}
			installPath, err := agent.GetInstallPath(service.ServiceType)
			if err != nil {
				return nil, fmt.Errorf("get install path for %s failed: %w", service.ServiceType, err)
			}

			services = append(services, &clients.ServiceInfo{
				ServiceId:   service.ServiceId,
				ServiceType: service.ServiceType,
				Host:        service.Host,
				Port:        service.Port,
				InstallPath: installPath,
			})
		}

		clusterMap[cluster.ClusterId] = cluster.Services
	}
	return &NebulaClusters{
		clusters: clusterMap,
	}, nil
}

func (c *NebulaClusters) GetStorages(clusterId int64) []*clients.ServiceInfo {
	var storages []*clients.ServiceInfo
	for _, service := range c.clusters[clusterId] {
		if service.ServiceType == meta.ServiceTypeStoraged {
			storages = append(storages, service)
		}
	}
	return storages
}

func (c *NebulaClusters) GetGraphs(clusterId int64) []*clients.ServiceInfo {
	var graphs []*clients.ServiceInfo
	for _, service := range c.clusters[clusterId] {
		if service.ServiceType == meta.ServiceTypeGraphd {
			graphs = append(graphs, service)
		}
	}
	return graphs
}

func (c *NebulaClusters) GetServices(clusterId int64) []*clients.ServiceInfo {
	return c.clusters[clusterId]
}

func (c *NebulaClusters) GetClusterIds() []int64 {
	var ids []int64
	for id := range c.clusters {
		ids = append(ids, id)
	}
	return ids
}

func (c *NebulaClusters) GetClusters() map[int64][]*clients.ServiceInfo {
	return c.clusters
}
