package utils

import (
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type NebulaServiceGroups struct {
	clusters map[int64][]*clients.ServiceInfo // clusterId -> ([storaged, graphd])
}

func NewNebulaServiceGroups(clusters []*clients.ServiceGroupServiceInfo, amg *clients.AgentManager) (*NebulaServiceGroups, error) {
	clusterMap := make(map[int64][]*clients.ServiceInfo)
	for _, cluster := range clusters {
		services := make([]*clients.ServiceInfo, 0)
		for _, service := range cluster.Services {

			services = append(services, &clients.ServiceInfo{
				ServiceId:   service.ServiceId,
				ServiceType: service.ServiceType,
				Host:        service.Host,
				Port:        service.Port,
				InstallPath: service.InstallPath,
				DataPaths:   service.DataPaths,
			})
		}

		clusterMap[cluster.ServiceGroupId] = cluster.Services
	}
	return &NebulaServiceGroups{
		clusters: clusterMap,
	}, nil
}

func (c *NebulaServiceGroups) GetStorages(clusterId int64) []*clients.ServiceInfo {
	var storages []*clients.ServiceInfo
	for _, service := range c.clusters[clusterId] {
		if service.ServiceType == meta.ServiceTypeStoraged {
			storages = append(storages, service)
		}
	}
	return storages
}

func (c *NebulaServiceGroups) GetGraphs(clusterId int64) []*clients.ServiceInfo {
	var graphs []*clients.ServiceInfo
	for _, service := range c.clusters[clusterId] {
		if service.ServiceType == meta.ServiceTypeGraphd {
			graphs = append(graphs, service)
		}
	}
	return graphs
}

func (c *NebulaServiceGroups) GetServices(clusterId int64) []*clients.ServiceInfo {
	return c.clusters[clusterId]
}

func (c *NebulaServiceGroups) GetServiceGroupIds() []int64 {
	var ids []int64
	for id := range c.clusters {
		ids = append(ids, id)
	}
	return ids
}

func (c *NebulaServiceGroups) GetServiceGroups() map[int64][]*clients.ServiceInfo {
	return c.clusters
}
