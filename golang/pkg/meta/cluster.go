package meta

type (
	CreateClusterReq struct {
		clusterName string
		replica     int
		zones       []string
	}

	CreateClusterResp struct {
		*HeaderResponse
		ClusterId int64
	}

	AddServiceReq struct {
		host        string
		port        uint32
		serviceType ServiceType
		clustername string
	}

	AddServiceResp struct {
		*HeaderResponse
	}

	DropServiceReq struct {
		host        string
		port        uint32
		serviceType ServiceType
		clustername string
	}
	DropServiceResp struct {
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
