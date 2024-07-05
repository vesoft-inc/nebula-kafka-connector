package meta

import (
	"context"
	"fmt"

	admin "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/admin"
	common "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
)

type (
	CreateClusterReq struct {
		clusterName string
		replica     int
		zones       []string
		owner       string
	}

	AlterClusterReq struct {
		clusterName string
		owner       string
	}

	AddHostReq struct {
		host		string
		clustername string
		agentPort	uint32
	}

	DropHostReq struct {
		host        string
		clustername string
	}

	ListHostsReq struct {
		clustername string
	}

	HostInfo struct {
		HostName string `json:"host"`
		AgentPort uint32 `json:"agent_port"`
	}

	ListHostsResp struct {
		HostInfoList	[]*HostInfo `json:"host_info_list"`
	}

	AddServiceReq struct {
		host        string
		port        uint32
		serviceType ServiceType
		clustername string
	}

	DropServiceReq struct {
		host        string
		port        uint32
		serviceType ServiceType
		clustername string
	}

	ServiceType int8

	ClusterInfo struct {
		Id              int64    `json:"id"`
		Name            string   `json:"name"`
		ReplicaRefactor uint32   `json:"replica_refactor"`
		Zones           []string `json:"zones"`
		Owner           string   `json:"owner"`
	}

	ListClustersReq struct {
		clusterName string
	}

	ListClustersResp struct {
		Clusters []*ClusterInfo `json:"clusters"`
	}

	ServiceInfo struct {
		Id   int64       `json:"id"`
		Type ServiceType `json:"type"`
		Host string      `json:"host"`
		Port uint32      `json:"port"`
	}

	ListServicesReq struct {
		clusterName string
	}

	ListServicesResp struct {
		Services []*ServiceInfo `json:"services"`
	}

	InitClusterReq struct {
		clustername string
	}

	DropClusterReq struct {
		clustername string
		force       bool
	}

	ClusterServiceInfo struct {
		ClusterId int64
		Services  []*ServiceInfo
	}
)

const (
	ServiceTypeUnknown ServiceType = iota
	ServiceTypeStoraged
	ServiceTypeGraphd
	ServiceTypeMetad
	ServiceTypeSearch
)

func NewCreateClusterReq(clusterName string, replic int, owner string, zones []string) *CreateClusterReq {
	return &CreateClusterReq{
		clusterName: clusterName,
		replica:     replic,
		zones:       zones,
		owner:       owner,
	}
}

func NewAlterClusterReq(clusterName string, owner string) *AlterClusterReq {
	return &AlterClusterReq{
		clusterName: clusterName,
		owner:       owner,
	}
}

func NewAddHostReq(host string, clusterName string, agentPort uint32) *AddHostReq {
	return &AddHostReq{
		host:        host,
		clustername: clusterName,
		agentPort:   agentPort,
	}
}

func NewDropHostReq(host string, clusterName string) *DropHostReq {
	return &DropHostReq{
		host:        host,
		clustername: clusterName,
	}
}

func NewShowHostsReq(clusterName string) *ListHostsReq {
	return &ListHostsReq{
		clustername: clusterName,
	}
}

func NewAddServiceReq(host string, port uint32, serviceType ServiceType, clusterName string) *AddServiceReq {
	return &AddServiceReq{
		host:        host,
		port:        port,
		serviceType: serviceType,
		clustername: clusterName,
	}
}

func NewListClustersReq(clusterName string) *ListClustersReq {
	return &ListClustersReq{
		clusterName: clusterName,
	}
}

func NewInitClusterReq(clusterName string) *InitClusterReq {
	return &InitClusterReq{
		clustername: clusterName,
	}
}

func NewListServicesReq(clusterName string) *ListServicesReq {
	return &ListServicesReq{
		clusterName: clusterName,
	}
}

func NewDropServiceReq(host string, port uint32, serviceType ServiceType, clusterName string) *DropServiceReq {
	return &DropServiceReq{
		host:        host,
		port:        port,
		serviceType: serviceType,
		clustername: clusterName,
	}
}

func (c *metaClient) CreateCluster(req *CreateClusterReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	zones := make([][]byte, 0)
	for _, z := range req.zones {
		zones = append(zones, []byte(z))
	}
	in := &admin.CreateClusterRequest{
		Header: &admin.RequestHeader{Token: c.token},
		ClusterDesc: &admin.ClusterDesc{
			ClusterName:   []byte(req.clusterName),
			ReplicaFactor: uint32(req.replica),
			Zones:         zones,
			Owner:         []byte(req.owner),
		},
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.CreateCluster(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c* metaClient) AddHost(req *AddHostReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AddHostRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		HostInfo:    &admin.HostInfo{HostName: []byte(req.host), ClusterName: []byte(req.clustername), AgentPort: req.agentPort},
		Force: 	 false,
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.AddHost(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) DropHost(req *DropHostReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.RemoveHostRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		HostName:   []byte(req.host),
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropHost(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c* metaClient) ListHosts(req *ListHostsReq) (*ListHostsResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ListHostsRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ListHosts(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err := responseIsErr(resp); err != nil {
        return nil, err
    }
    response, ok := resp.(*admin.ListHostsResponse)
    if !ok {
        return nil, fmt.Errorf("invalid response")
    }
	hostInfoList := make([]*HostInfo, 0)
	for _, h := range response.HostInfo {
		hostInfoList = append(hostInfoList, &HostInfo{
			HostName: string(h.HostName),
			AgentPort: h.AgentPort,
		})
	}
	return &ListHostsResp{
		HostInfoList: hostInfoList,
	}, nil
}

func (c *metaClient) AddService(req *AddServiceReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	in := &admin.AddServiceRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		Host:        []byte(req.host),
		Port:        req.port,
		Type:        common.ServiceType(req.serviceType),
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.AddService(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}
func (c *metaClient) ListServices(req *ListServicesReq) (*ListServicesResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ShowServiceRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ShowService(ctx, in)
	})
	if err != nil {
		return nil, err
	}

	if err := responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowServiceResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	services := make([]*ServiceInfo, 0)
	for _, s := range response.Services {
		addr := s.Address
		host := addr.GetHost()
		port := addr.GetPort()
		services = append(services, &ServiceInfo{
			Id:   s.ServiceId,
			Type: ServiceType(s.Type),
			Host: string(host),
			Port: port,
		})
	}
	return &ListServicesResp{
		Services: services,
	}, nil
}

func (c *metaClient) AlterCluster(req *AlterClusterReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AlterClusterRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		ClusterName: []byte(req.clusterName),
		Owner:       []byte(req.owner),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.AlterCluster(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) InitCluster(req *InitClusterReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.InitStorageRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.InitStorage(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) ListClusters(req *ListClustersReq) (*ListClustersResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ShowClusterRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ShowCluster(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err := responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowClusterResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	clusters := make([]*ClusterInfo, 0)
	for _, c := range response.Clusters {
		cluster := &ClusterInfo{
			Id:              c.ClusterId,
			Name:            string(c.Desc.ClusterName),
			ReplicaRefactor: c.Desc.ReplicaFactor,
			Zones:           make([]string, 0),
			Owner:           string(c.Desc.Owner),
		}
		for _, z := range c.Desc.Zones {
			cluster.Zones = append(cluster.Zones, string(z))
		}
		clusters = append(clusters, cluster)
	}
	return &ListClustersResp{
		Clusters: clusters,
	}, nil
}

func (c *metaClient) DropService(req *DropServiceReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropServiceRequest{
		Header:  &admin.RequestHeader{Token: c.token},
		Name:    []byte(req.host),
		Port:    req.port,
		Type:    common.ServiceType(req.serviceType),
		Cluster: []byte(req.clustername),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropService(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewDropClusterReq(clusterName string, force bool) *DropClusterReq {
	return &DropClusterReq{
		clustername: clusterName,
		force:       force,
	}
}

func (c *metaClient) DropCluster(req *DropClusterReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropClusterRequest{
		Header:  &admin.RequestHeader{Token: c.token},
		Cluster: []byte(req.clustername),
		Force:   req.force,
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropCluster(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}
