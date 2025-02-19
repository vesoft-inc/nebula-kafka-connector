package meta

import (
	"context"
	"fmt"

	admin "github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/admin"
	common "github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/common"
)

type (
	CreateServiceGroupReq struct {
		serviceGroupName string
		replica          int
		zones            []string
		owner            string
	}

	AlterServiceGroupReq struct {
		serviceGroupName string
		owner            string
	}

	AddHostReq struct {
		host             string
		serviceGroupName string
		agentPort        uint32
		zoneName         string
	}

	DropHostReq struct {
		host             string
		serviceGroupName string
	}

	ListHostsReq struct {
		serviceGroupName string
	}

	HostInfo struct {
		HostName  string `json:"host"`
		AgentPort uint32 `json:"agent_port"`
	}

	ListHostsResp struct {
		HostInfoList []*HostInfo `json:"host_info_list"`
	}

	AddServiceReq struct {
		host             string
		port             uint32
		serviceType      ServiceType
		serviceGroupName string
	}

	DropServiceReq struct {
		host             string
		port             uint32
		serviceType      ServiceType
		serviceGroupName string
	}

	ServiceType int8

	ServiceGroupInfo struct {
		Id              int64    `json:"id"`
		Name            string   `json:"name"`
		ReplicaRefactor uint32   `json:"replica_refactor"`
		Zones           []string `json:"zones"`
		Owner           string   `json:"owner"`
	}

	ListServiceGroupsReq struct {
		serviceGroupName string
	}

	ListServiceGroupsResp struct {
		ServiceGroups []*ServiceGroupInfo `json:"serviceGroups"`
	}

	ServiceInfo struct {
		Id   int64       `json:"id"`
		Type ServiceType `json:"type"`
		Host string      `json:"host"`
		Port uint32      `json:"port"`
	}

	ListServicesReq struct {
		serviceGroupName string
	}

	ListServicesResp struct {
		Services []*ServiceInfo `json:"services"`
	}

	InitServiceGroupReq struct {
		serviceGroupName string
	}

	DropServiceGroupReq struct {
		serviceGroupName string
		force            bool
	}

	ServiceGroupServiceInfo struct {
		ServiceGroupId int64
		Services       []*ServiceInfo
	}
	AddZoneReq struct {
		zoneName         string
		serviceGroupName string
		replicas         int
		priority         int
	}

	DropZoneReq struct {
		zoneName         string
		serviceGroupName string
	}

	RenameZoneReq struct {
		zoneName         string
		newZoneName      string
		serviceGroupName string
	}

	ZoneInfo struct {
		Name     string   `json:"name"`
		Replicas int      `json:"replicas"`
		Priority int      `json:"priority"`
		Hosts    []string `json:"hosts"`
	}

	ListZonesReq struct {
		serviceGroupName string
	}

	ListZonesResp struct {
		Zones []*ZoneInfo `json:"zones"`
	}
)

const (
	ServiceTypeUnknown ServiceType = iota
	ServiceTypeStoraged
	ServiceTypeGraphd
	ServiceTypeMetad
	ServiceTypeAnalyticd
	ServiceTypeSearch
)

func NewCreateServiceGroupReq(serviceGroupName string, replica int, owner string, zones []string) *CreateServiceGroupReq {
	return &CreateServiceGroupReq{
		serviceGroupName: serviceGroupName,
		replica:          replica,
		zones:            zones,
		owner:            owner,
	}
}

func NewAlterServiceGroupReq(serviceGroupName string, owner string) *AlterServiceGroupReq {
	return &AlterServiceGroupReq{
		serviceGroupName: serviceGroupName,
		owner:            owner,
	}
}

func NewAddHostReq(host string, serviceGroupName string, agentPort uint32, opts ...string) *AddHostReq {
	if len(opts) == 1 {
		return &AddHostReq{
			host:             host,
			serviceGroupName: serviceGroupName,
			agentPort:        9669,
			zoneName:         opts[0],
		}
	}
	return &AddHostReq{
		host:             host,
		serviceGroupName: serviceGroupName,
		agentPort:        agentPort,
		zoneName:         "",
	}
}

func NewDropHostReq(host string, serviceGroupName string) *DropHostReq {
	return &DropHostReq{
		host:             host,
		serviceGroupName: serviceGroupName,
	}
}

func NewShowHostsReq(serviceGroupName string) *ListHostsReq {
	return &ListHostsReq{
		serviceGroupName: serviceGroupName,
	}
}

func NewAddServiceReq(host string, port uint32, serviceType ServiceType, serviceGroupName string) *AddServiceReq {
	return &AddServiceReq{
		host:             host,
		port:             port,
		serviceType:      serviceType,
		serviceGroupName: serviceGroupName,
	}
}

func NewListServiceGroupsReq(serviceGroupName string) *ListServiceGroupsReq {
	return &ListServiceGroupsReq{
		serviceGroupName: serviceGroupName,
	}
}

func NewInitServiceGroupReq(serviceGroupName string) *InitServiceGroupReq {
	return &InitServiceGroupReq{
		serviceGroupName: serviceGroupName,
	}
}

func NewListServicesReq(serviceGroupName string) *ListServicesReq {
	return &ListServicesReq{
		serviceGroupName: serviceGroupName,
	}
}

func NewDropServiceReq(host string, port uint32, serviceType ServiceType, serviceGroupName string) *DropServiceReq {
	return &DropServiceReq{
		host:             host,
		port:             port,
		serviceType:      serviceType,
		serviceGroupName: serviceGroupName,
	}
}

func NewAddZoneReq(serviceGroupName string, zoneName string, replicas int, priority int) *AddZoneReq {
	return &AddZoneReq{
		zoneName:         zoneName,
		serviceGroupName: serviceGroupName,
		replicas:         replicas,
		priority:         priority,
	}
}

func NewDropZoneReq(serviceGroupName string, zoneName string) *DropZoneReq {
	return &DropZoneReq{
		zoneName:         zoneName,
		serviceGroupName: serviceGroupName,
	}
}

func NewRenameZoneReq(serviceGroupName string, zoneName string, newZoneName string) *RenameZoneReq {
	return &RenameZoneReq{
		zoneName:         zoneName,
		newZoneName:      newZoneName,
		serviceGroupName: serviceGroupName,
	}
}

func NewListZonesReq(serviceGroupName string) *ListZonesReq {
	return &ListZonesReq{
		serviceGroupName: serviceGroupName,
	}
}

func (c *metaClient) CreateServiceGroup(req *CreateServiceGroupReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	zones := make([][]byte, 0)
	for _, z := range req.zones {
		zones = append(zones, []byte(z))
	}
	in := &admin.CreateServiceGroupRequest{
		Header: &admin.RequestHeader{Token: c.token},
		ServiceGroupDesc: &admin.ServiceGroupDesc{
			ServiceGroup:  []byte(req.serviceGroupName),
			ReplicaFactor: uint32(req.replica),
			Zones:         zones,
			Owner:         []byte(req.owner),
		},
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.CreateServiceGroup(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) AddHost(req *AddHostReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AddHostRequest{
		Header: &admin.RequestHeader{Token: c.token},
		HostInfo: &admin.HostInfo{
			HostName:     []byte(req.host),
			ServiceGroup: []byte(req.serviceGroupName),
			AgentPort:    req.agentPort,
			ZoneName:     []byte(req.zoneName),
		},
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
		Header:       &admin.RequestHeader{Token: c.token},
		HostName:     []byte(req.host),
		ServiceGroup: []byte(req.serviceGroupName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropHost(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) ListHosts(req *ListHostsReq) (*ListHostsResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ListHostsRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
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
			HostName:  string(h.HostName),
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
		Header:       &admin.RequestHeader{Token: c.token},
		Host:         []byte(req.host),
		Port:         req.port,
		Type:         common.ServiceType(req.serviceType),
		ServiceGroup: []byte(req.serviceGroupName),
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
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
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

func (c *metaClient) AlterServiceGroup(req *AlterServiceGroupReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AlterServiceGroupRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
		Owner:        []byte(req.owner),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.AlterServiceGroup(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) InitServiceGroup(req *InitServiceGroupReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.InitStorageRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.InitStorage(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) ListServiceGroups(req *ListServiceGroupsReq) (*ListServiceGroupsResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ShowServiceGroupRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ShowServiceGroup(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err := responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowServiceGroupResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	serviceGroups := make([]*ServiceGroupInfo, 0)
	for _, c := range response.ServiceGroups {
		serviceGroup := &ServiceGroupInfo{
			Id:              c.ServiceGroupId,
			Name:            string(c.Desc.ServiceGroup),
			ReplicaRefactor: c.Desc.ReplicaFactor,
			Zones:           make([]string, 0),
			Owner:           string(c.Desc.Owner),
		}
		for _, z := range c.Desc.Zones {
			serviceGroup.Zones = append(serviceGroup.Zones, string(z))
		}
		serviceGroups = append(serviceGroups, serviceGroup)
	}
	return &ListServiceGroupsResp{
		ServiceGroups: serviceGroups,
	}, nil
}

func (c *metaClient) DropService(req *DropServiceReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropServiceRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		Name:         []byte(req.host),
		Port:         req.port,
		Type:         common.ServiceType(req.serviceType),
		ServiceGroup: []byte(req.serviceGroupName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropService(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewDropServiceGroupReq(serviceGroupName string, force bool) *DropServiceGroupReq {
	return &DropServiceGroupReq{
		serviceGroupName: serviceGroupName,
		force:            force,
	}
}

func (c *metaClient) DropServiceGroup(req *DropServiceGroupReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropServiceGroupRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
		Force:        req.force,
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropServiceGroup(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) AddZone(req *AddZoneReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AddZoneRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
		ZoneDesc: &admin.ZoneDesc{
			ZoneName:    []byte(req.zoneName),
			NumReplicas: uint32(req.replicas),
			Priority:    uint32(req.priority),
		},
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.AddZone(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) DropZone(req *DropZoneReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropZoneRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
		ZoneName:     []byte(req.zoneName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropZone(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) RenameZone(req *RenameZoneReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.RenameZoneRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
		ZoneName:     []byte(req.zoneName),
		NewZoneName:  []byte(req.newZoneName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.RenameZone(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func (c *metaClient) ListZones(req *ListZonesReq) (*ListZonesResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ListZonesRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.serviceGroupName),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ListZones(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err := responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ListZonesResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	zones := make([]*ZoneInfo, 0)
	for _, z := range response.Zones {
		if z.ZoneDesc == nil {
			continue
		}
		zone := &ZoneInfo{
			Name:     string(z.ZoneDesc.ZoneName),
			Replicas: int(z.ZoneDesc.NumReplicas),
			Priority: int(z.ZoneDesc.Priority),
			Hosts:    make([]string, 0),
		}
		for _, h := range z.HostNames {
			zone.Hosts = append(zone.Hosts, string(h))
		}
		zones = append(zones, zone)
	}
	return &ListZonesResp{
		Zones: zones,
	}, nil
}
