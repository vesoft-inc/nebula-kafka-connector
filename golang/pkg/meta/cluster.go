package meta

type (
	CreateClusterReq struct {
		*headerRequest
		ifNotExists bool
		Clustername string
		Replica     int
		Zones       []string
	}

	CreateClusterResp struct {
		*HeaderResponse
	}

	AddServiceReq struct {
		*headerRequest
		Host        string
		Port        uint32
		ServiceType ServiceType
		Clustername string
	}

	AddServiceResp struct {
		*HeaderResponse
	}

	ServiceType int8

	ClusterInfo struct {
		ClusterId   int64
		ClusterName string
		Replica     uint32
		Zones       []string
	}

	ShowClusterReq struct {
		*headerRequest
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
		*headerRequest
		clusterName string
	}

	ShowServiceResp struct {
		*HeaderResponse
		Services []*ServiceInfo
	}

	InitClusterReq struct {
		*headerRequest
		clustername string
	}

	InitClusterResp struct {
		*HeaderResponse
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
		headerRequest: &headerRequest{
			requestType: "createCluster",
		},
		Clustername: clusterName,
		Replica:     replic,
		Zones:       zones,
		//forbid to create cluster if it exists
		ifNotExists: false,
	}
}

func (req *CreateClusterReq) serialize(s serializer) []byte {
	s.serializeHeader(req.headerRequest)
	s.serializeString(req.Clustername)
	s.serializeUINT32(uint32(req.Replica))
	s.serializeStringArray(req.Zones)
	s.serializeBool(req.ifNotExists)
	return s.getBytes()
}

func (resp *CreateClusterResp) deserialize(d deserializer) error {
	var err error
	resp.HeaderResponse, err = d.deserializeHeader()
	if err != nil {
		return err
	}

	return nil
}

func NewAddServiceReq(host string, port uint32, serviceType ServiceType, clusterName string) *AddServiceReq {
	return &AddServiceReq{
		headerRequest: &headerRequest{
			requestType: "addService",
		},
		Host:        host,
		Port:        port,
		ServiceType: serviceType,
		Clustername: clusterName,
	}
}

func (req *AddServiceReq) serialize(s serializer) []byte {
	s.serializeHeader(req.headerRequest)
	s.serializeString(req.Host)
	s.serializeUINT32(req.Port)
	s.serializeUINT8(uint8(req.ServiceType))
	s.serializeString(req.Clustername)
	return s.getBytes()
}

func (resp *AddServiceResp) deserialize(d deserializer) error {
	var err error
	resp.HeaderResponse, err = d.deserializeHeader()
	if err != nil {
		return err
	}

	return nil
}

func NewShowClusterReq(clusterName string) *ShowClusterReq {
	return &ShowClusterReq{
		headerRequest: &headerRequest{
			requestType: "showCluster",
		},
		clusterName: clusterName,
	}
}

func (req *ShowClusterReq) serialize(s serializer) []byte {
	s.serializeHeader(req.headerRequest)
	s.serializeString(req.clusterName)
	return s.getBytes()
}

func (resp *ShowClusterResp) deserialize(d deserializer) error {
	var err error
	resp.HeaderResponse, err = d.deserializeHeader()
	if err != nil {
		return err
	}
	//should handle leader change outside
	if resp.GetErrorCode() != 0 {
		return nil
	}
	clusterLen, err := d.deserializeUINT32()
	if err != nil {
		return err
	}
	for i := 0; i < int(clusterLen); i++ {
		clusterId, err := d.deserializeUINT64()
		if err != nil {
			return err
		}
		clusterName, err := d.deserializeString()
		if err != nil {
			return err
		}
		replica, err := d.deserializeUINT32()
		if err != nil {
			return err
		}
		zones, err := d.deserializeStringArray()
		if err != nil {
			return err
		}
		resp.Clusters = append(resp.Clusters, &ClusterInfo{
			ClusterId:   int64(clusterId),
			ClusterName: clusterName,
			Replica:     replica,
			Zones:       zones,
		})
	}
	return nil
}

func NewInitClusterReq(clusterName string) *InitClusterReq {
	return &InitClusterReq{
		headerRequest: &headerRequest{
			requestType: "initCluster",
		},
		clustername: clusterName,
	}
}

func (req *InitClusterReq) serialize(s serializer) []byte {
	s.serializeHeader(req.headerRequest)
	s.serializeString(req.clustername)
	return s.getBytes()
}

func (resp *InitClusterResp) deserialize(d deserializer) error {
	var err error
	resp.HeaderResponse, err = d.deserializeHeader()
	if err != nil {
		return err
	}

	return nil
}

func NewShowServiceReq(clusterName string) *ShowServiceReq {
	return &ShowServiceReq{
		headerRequest: &headerRequest{
			requestType: "showService",
		},
		clusterName: clusterName,
	}
}

func (req *ShowServiceReq) serialize(s serializer) []byte {
	s.serializeHeader(req.headerRequest)
	s.serializeString(req.clusterName)
	return s.getBytes()
}

func (resp *ShowServiceResp) deserialize(d deserializer) error {
	var err error
	resp.HeaderResponse, err = d.deserializeHeader()
	if err != nil {
		return err
	}
	//should handle leader change outside
	if resp.GetErrorCode() != 0 {
		return nil
	}

	serviceLen, err := d.deserializeUINT32()
	if err != nil {
		return err
	}
	for i := 0; i < int(serviceLen); i++ {
		serviceId, err := d.deserializeUINT64()
		if err != nil {
			return err
		}
		t, err := d.deserializeINT8()
		if err != nil {
			return err
		}
		serviceType := ServiceType(t)
		host, err := d.deserializeString()
		if err != nil {
			return err
		}
		port, err := d.deserializeUINT32()
		if err != nil {
			return err
		}
		resp.Services = append(resp.Services, &ServiceInfo{
			ServiceId:   int64(serviceId),
			ServiceType: ServiceType(serviceType),
			Host:        host,
			Port:        port,
		})
	}
	return nil
}
