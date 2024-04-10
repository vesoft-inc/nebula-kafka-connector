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
		ClusterId   int64
		ClusterName string
		Replica     uint32
		Zones       []string
	}

	ShowClusterReq struct {
		clusterName string
	}

	ShowClusterResp struct {
		*HeaderResponse
		Clusters []*ClusterInfo
	}

	ServiceInfo struct {
		ServiceId   int64
		ServiceType ServiceType
		Host        string
		Port        uint32
	}

	ShowServiceReq struct {
		clusterName string
	}

	ShowServiceResp struct {
		*HeaderResponse
		Services []*ServiceInfo
	}

	InitClusterReq struct {
		clustername string
	}
)

const (
	ServiceTypeUnknown ServiceType = iota
	ServiceTypeStoraged
	ServiceTypeGraphd
	ServiceTypeSearch
)

func NewCreateClusterReq(clusterName string, replic int, zones []string) *CreateClusterReq {
	return &CreateClusterReq{
		clusterName: clusterName,
		replica:     replic,
		zones:       zones,
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

func NewShowClusterReq(clusterName string) *ShowClusterReq {
	return &ShowClusterReq{
		clusterName: clusterName,
	}
}

func NewInitClusterReq(clusterName string) *InitClusterReq {
	return &InitClusterReq{
		clustername: clusterName,
	}
}

func NewShowServiceReq(clusterName string) *ShowServiceReq {
	return &ShowServiceReq{
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
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	zones := make([][]byte, 0)
	for _, z := range req.zones {
		zones = append(zones, []byte(z))
	}
	in := &admin.CreateClusterRequest{
		Header: &admin.AdminRequestHeader{Token: c.token},
		ClusterDesc: &admin.ClusterDesc{
			ClusterName:   []byte(req.clusterName),
			ReplicaFactor: uint32(req.replica),
			Zones:         zones,
		},
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.CreateCluster(ctx, in)
	})
	if err != nil {
		return err
	}
	response, ok := resp.(*admin.CreateClusterResponse)
	if !ok {
		return fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetError()
	}
	return nil
}

func (c *metaClient) AddService(req *AddServiceReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	in := &admin.AddServiceRequest{
		Header:      &admin.AdminRequestHeader{Token: c.token},
		Host:        []byte(req.host),
		Port:        req.port,
		Type:        common.ServiceType(req.serviceType),
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.AddService(ctx, in)
	})
	if err != nil {
		return err
	}
	response, ok := resp.(*admin.AddServiceResponse)
	if !ok {
		return fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetError()
	}
	return nil
}
func (c *metaClient) ShowService(req *ShowServiceReq) (*ShowServiceResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.ShowServiceRequest{
		Header:      &admin.AdminRequestHeader{Token: c.token},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ShowService(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowServiceResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
	}
	services := make([]*ServiceInfo, 0)
	for _, s := range response.Services {
		addr := s.Address
		host := addr.GetHost()
		port := addr.GetPort()
		services = append(services, &ServiceInfo{
			ServiceId:   s.ServiceId,
			ServiceType: ServiceType(s.Type),
			Host:        string(host),
			Port:        port,
		})
	}
	if !responseHeader.IsSucceeded() {
		return nil, responseHeader.GetError()
	}
	return &ShowServiceResp{
		HeaderResponse: responseHeader,
		Services:       services,
	}, nil
}
func (c *metaClient) InitCluster(req *InitClusterReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.InitStorageRequest{
		Header:      &admin.AdminRequestHeader{Token: c.token},
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.InitStorage(ctx, in)
	})
	if err != nil {
		return err
	}
	response, ok := resp.(*admin.InitStorageResponse)
	if !ok {
		return fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetError()
	}
	return nil
}
func (c *metaClient) ShowCluster(req *ShowClusterReq) (*ShowClusterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.ShowClusterRequest{
		Header:      &admin.AdminRequestHeader{Token: c.token},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ShowCluster(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowClusterResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
	}
	if !responseHeader.IsSucceeded() {
		return nil, responseHeader.GetError()
	}
	clusters := make([]*ClusterInfo, 0)
	for _, c := range response.Clusters {
		cluster := &ClusterInfo{
			ClusterId:   c.ClusterId,
			ClusterName: string(c.Desc.ClusterName),
			Replica:     c.Desc.ReplicaFactor,
			Zones:       make([]string, 0),
		}
		for _, z := range c.Desc.Zones {
			cluster.Zones = append(cluster.Zones, string(z))
		}
		clusters = append(clusters, cluster)
	}
	return &ShowClusterResp{
		HeaderResponse: responseHeader,
		Clusters:       clusters,
	}, nil
}

func (c *metaClient) DropService(req *DropServiceReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.DropServiceRequest{
		Header:  &admin.AdminRequestHeader{Token: c.token},
		Name:    []byte(req.host),
		Port:    req.port,
		Type:    common.ServiceType(req.serviceType),
		Cluster: []byte(req.clustername),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.DropService(ctx, in)
	})
	if err != nil {
		return err
	}
	response, ok := resp.(*admin.DropServiceResponse)
	if !ok {
		return fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetError()
	}
	return nil
}
