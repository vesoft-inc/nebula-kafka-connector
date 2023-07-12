package meta

import (
	"bytes"
	"fmt"
	"strings"
)

type RequestHeader struct {
	requestType string
	clusterId   int64 // meta admin dont need clusterId
}

type ResponseHeader struct {
	code    uint64
	msg     string
	newHost string // only if code is leader changed
	newPort uint32 // only if code is leader changed
}

func formatTable(headers []string, data [][]string) string {
	var buf bytes.Buffer

	// print header
	for _, header := range headers {
		fmt.Fprintf(&buf, "%-20s", header)
	}
	fmt.Fprintln(&buf)

	// print split line
	for i := 0; i < len(headers)*20; i++ {
		fmt.Fprint(&buf, "-")
	}
	fmt.Fprintln(&buf)

	// print data
	for _, row := range data {
		for _, cell := range row {
			fmt.Fprintf(&buf, "%-20s", cell)
		}
		fmt.Fprintln(&buf)
	}

	return buf.String()
}

func DeserializeHeader(deserializer *Deserializer) *ResponseHeader {
	code, _ := deserializer.DeserializeUINT64()
	msg, _ := deserializer.DeserializeString()
	newHost, _ := deserializer.DeserializeString()
	newPort, _ := deserializer.DeserializeUINT32()
	return &ResponseHeader{
		code:    code,
		msg:     msg,
		newHost: newHost,
		newPort: uint32(newPort),
	}
}

const (
	Unknown = 0
	Storage = 1
	Graph   = 2
	Search  = 3
)

type CreateClusterRequest struct {
	header      RequestHeader
	Clustername string
	Replica     int
	Zones       []string
	IfNotExists bool
}

func NewCreateClusterRequest(clusterName string, replic int, zones []string, ifNotExists bool) *CreateClusterRequest {
	return &CreateClusterRequest{
		header: RequestHeader{
			requestType: "createCluster",
		},
		Clustername: clusterName,
		Replica:     replic,
		Zones:       zones,
		IfNotExists: ifNotExists,
	}
}

func (c *CreateClusterRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(c.header.requestType)
	serializer.SerializeINT64(c.header.clusterId)
	serializer.SerializeString(c.Clustername)
	serializer.SerializeUINT32(uint32(c.Replica))
	serializer.SerializeStringArray(c.Zones)
	serializer.SerializeBool(c.IfNotExists)
	return serializer.GetBytes()
}

type AddServiceRequest struct {
	header      RequestHeader
	Host        string
	Port        uint32
	ServiceType int8
	Clustername string
}

func NewAddServiceRequest(host string, port uint32, serviceType int8, clusterName string) *AddServiceRequest {
	return &AddServiceRequest{
		header: RequestHeader{
			requestType: "addService",
		},
		Host:        host,
		Port:        port,
		ServiceType: serviceType,
		Clustername: clusterName,
	}
}

func (a *AddServiceRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(a.header.requestType)
	serializer.SerializeINT64(a.header.clusterId)
	serializer.SerializeString(a.Host)
	serializer.SerializeUINT32(a.Port)
	serializer.SerializeUINT8(uint8(a.ServiceType))
	serializer.SerializeString(a.Clustername)
	return serializer.GetBytes()
}

type InitClusterRequest struct {
	header      RequestHeader
	clustername string
}

func NewInitClusterRequest(clusterName string) *InitClusterRequest {
	return &InitClusterRequest{
		header: RequestHeader{
			requestType: "initCluster",
		},
		clustername: clusterName,
	}
}

func (i *InitClusterRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(i.header.requestType)
	serializer.SerializeINT64(i.header.clusterId)
	serializer.SerializeString(i.clustername)
	return serializer.GetBytes()
}

type MetaDDLRequest struct {
	header     RequestHeader
	optype     string
	jsonSchema string
}

func NewMetaDDLRequest(optype string, jsonSchema string) *MetaDDLRequest {
	return &MetaDDLRequest{
		header: RequestHeader{
			requestType: "adminDDL",
		},
		optype:     optype,
		jsonSchema: jsonSchema,
	}
}

func (c *MetaDDLRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(c.header.requestType)
	serializer.SerializeINT64(c.header.clusterId)
	serializer.SerializeString(c.optype)
	serializer.SerializeString(c.jsonSchema)
	return serializer.GetBytes()
}

type ListClusterRequest struct {
	header      RequestHeader
	clusterName string
}

func NewListClusterRequest(clusterName string) *ListClusterRequest {
	return &ListClusterRequest{
		header: RequestHeader{
			requestType: "showCluster",
		},
		clusterName: clusterName,
	}
}

func (l *ListClusterRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(l.header.requestType)
	serializer.SerializeINT64(l.header.clusterId)
	serializer.SerializeString(l.clusterName)
	return serializer.GetBytes()
}

type ListServiceRequest struct {
	header      RequestHeader
	clusterName string
}

func NewListServiceRequest(clusterName string) *ListServiceRequest {
	return &ListServiceRequest{
		header: RequestHeader{
			requestType: "showService",
		},
		clusterName: clusterName,
	}
}

func (l *ListServiceRequest) Serialize() []byte {
	serializer := NewSerializer()
	serializer.SerializeString(l.header.requestType)
	serializer.SerializeINT64(l.header.clusterId)
	serializer.SerializeString(l.clusterName)
	return serializer.GetBytes()
}

type ServiceInfo struct {
	serviceId   int64
	serviceType int8
	host        string
	port        uint32
}

type ListServiceResponse struct {
	info []ServiceInfo
}

func formatServiceType(serviceType int8) string {
	switch serviceType {
	case Storage:
		return "storage"
	case Graph:
		return "graph"
	case Search:
		return "search"
	default:
		return "unknown"
	}
}

func (l *ListServiceResponse) Format() string {
	var data [][]string
	for _, info := range l.info {
		var line []string
		line = append(line, fmt.Sprintf("%d", info.serviceId))
		line = append(line, formatServiceType(info.serviceType))
		line = append(line, info.host)
		line = append(line, fmt.Sprintf("%d", info.port))
		data = append(data, line)
	}

	return formatTable([]string{"service id", "service type", "host", "port"}, data)
}

func DeserializeListServiceResponse(deserializer *Deserializer) *ListServiceResponse {
	resp := &ListServiceResponse{}
	infoLen, _ := deserializer.DeserializeUINT32()
	for i := 0; i < int(infoLen); i++ {
		serviceId, _ := deserializer.DeserializeINT64()
		serviceType, _ := deserializer.DeserializeINT8()
		host, _ := deserializer.DeserializeString()
		port, _ := deserializer.DeserializeUINT32()
		resp.info = append(resp.info, ServiceInfo{
			serviceId:   int64(serviceId),
			serviceType: int8(serviceType),
			host:        host,
			port:        port,
		})
	}
	return resp
}

type ClusterInfo struct {
	clusterId   int64
	clusterName string
	replica     uint32
	zones       []string
}

type ListClusterResponse struct {
	info []ClusterInfo
}

func DeserializeListClusterResponse(deserializer *Deserializer) *ListClusterResponse {
	resp := &ListClusterResponse{}
	infoLen, _ := deserializer.DeserializeUINT32()
	for i := 0; i < int(infoLen); i++ {
		clusterId, _ := deserializer.DeserializeUINT64()
		clusterName, _ := deserializer.DeserializeString()
		replica, _ := deserializer.DeserializeUINT32()
		zones, _ := deserializer.DeserializeStringArray()
		resp.info = append(resp.info, ClusterInfo{
			clusterId:   int64(clusterId),
			clusterName: clusterName,
			replica:     replica,
			zones:       zones,
		})
	}
	return resp
}

func (r *ListClusterResponse) Format() string {
	var data [][]string
	for _, info := range r.info {
		var line []string
		line = append(line, fmt.Sprintf("%d", info.clusterId))
		line = append(line, info.clusterName)
		line = append(line, fmt.Sprintf("%d", info.replica))
		line = append(line, strings.Join(info.zones, ","))
		data = append(data, line)
	}

	return formatTable([]string{"cluster id", "cluster name", "replica", "zones"}, data)
}
