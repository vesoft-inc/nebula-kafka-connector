/*
Copyright 2023 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain s copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nebula

type ServiceType int8

const (
	Unknown        ServiceType = 0
	StorageService ServiceType = 1
	GraphService   ServiceType = 2
	MetaService    ServiceType = 3
)

type RequestHeader struct {
	RequestType string
	ClusterId   int64 // meta admin don't need clusterId
}

type LeaderHost struct {
	Host string
	Port uint32
}

type ResponseHeader struct {
	Code uint64
	Msg  string
	LeaderHost
}

func DeserializeHeader(deserializer *Deserializer, leaderChanged bool) (*ResponseHeader, error) {
	v, err := deserializer.DeserializeBool()
	if err != nil {
		return nil, err
	}
	var code uint64
	var msg string
	if !v {
		code, err = deserializer.DeserializeUint64()
		if err != nil {
			return nil, err
		}
		msg, err = deserializer.DeserializeString()
		if err != nil {
			return nil, err
		}
	}
	respHeader := &ResponseHeader{
		Code: code,
		Msg:  msg,
	}
	if leaderChanged {
		host, err := deserializer.DeserializeString()
		if err != nil {
			return nil, err
		}
		port, err := deserializer.DeserializeUint32()
		if err != nil {
			return nil, err
		}
		respHeader.Host = host
		respHeader.Port = port
	}

	return respHeader, nil
}

type CreateClusterRequest struct {
	Header      RequestHeader
	ClusterName string
	Replica     int
	Zones       []string
}

func NewCreateClusterRequest(clusterName string, replica int, zones []string) *CreateClusterRequest {
	return &CreateClusterRequest{
		Header: RequestHeader{
			RequestType: "createCluster",
		},
		ClusterName: clusterName,
		Replica:     replica,
		Zones:       zones,
	}
}

func (s *CreateClusterRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	if err := serializer.SerializeUint32(uint32(s.Replica)); err != nil {
		return nil, err
	}
	if err := serializer.SerializeStringArray(s.Zones); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type InitClusterRequest struct {
	Header      RequestHeader
	ClusterName string
}

func NewInitClusterRequest(clusterName string) *InitClusterRequest {
	return &InitClusterRequest{
		Header: RequestHeader{
			RequestType: "initCluster",
		},
		ClusterName: clusterName,
	}
}

func (s *InitClusterRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ListClusterRequest struct {
	Header      RequestHeader
	ClusterName string
}

func NewListClusterRequest(clusterName string) *ListClusterRequest {
	return &ListClusterRequest{
		Header: RequestHeader{
			RequestType: "showCluster",
		},
		ClusterName: clusterName,
	}
}

func (s *ListClusterRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type AddServiceRequest struct {
	Header      RequestHeader
	Host        string
	Port        uint32
	ServiceType ServiceType
	ClusterName string
}

func NewAddServiceRequest(host string, port uint32, serviceType ServiceType, clusterName string) *AddServiceRequest {
	return &AddServiceRequest{
		Header: RequestHeader{
			RequestType: "addService",
		},
		Host:        host,
		Port:        port,
		ServiceType: serviceType,
		ClusterName: clusterName,
	}
}

func (s *AddServiceRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.Host); err != nil {
		return nil, err
	}
	if err := serializer.SerializeUint32(s.Port); err != nil {
		return nil, err
	}
	if err := serializer.SerializeUint8(uint8(s.ServiceType)); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type DropServiceRequest struct {
	Header      RequestHeader
	ClusterName string
	ServiceType ServiceType
	Host        string
	Port        uint32
}

func NewDropServiceRequest(clusterName string, serviceType ServiceType, host string, port uint32) *DropServiceRequest {
	return &DropServiceRequest{
		Header: RequestHeader{
			RequestType: "dropService",
		},
		ClusterName: clusterName,
		ServiceType: serviceType,
		Host:        host,
		Port:        port,
	}
}

func (s *DropServiceRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	if err := serializer.SerializeUint8(uint8(s.ServiceType)); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.Host); err != nil {
		return nil, err
	}
	if err := serializer.SerializeUint32(s.Port); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ListServiceRequest struct {
	Header      RequestHeader
	ClusterName string
}

func NewListServiceRequest(clusterName string) *ListServiceRequest {
	return &ListServiceRequest{
		Header: RequestHeader{
			RequestType: "showService",
		},
		ClusterName: clusterName,
	}
}

func (s *ListServiceRequest) Serialize() ([]byte, error) {
	serializer := NewSerializer()
	if err := serializer.SerializeString(s.Header.RequestType); err != nil {
		return nil, err
	}
	if err := serializer.SerializeInt64(s.Header.ClusterId); err != nil {
		return nil, err
	}
	if err := serializer.SerializeString(s.ClusterName); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type ServiceInfo struct {
	ServiceId   int64
	ServiceType ServiceType
	Host        string
	Port        uint32
}

type ListServiceResponse struct {
	Services []ServiceInfo
}

func DeserializeListServiceResponse(deserializer *Deserializer) (*ListServiceResponse, error) {
	resp := &ListServiceResponse{}
	length, err := deserializer.DeserializeUint32()
	if err != nil {
		return nil, err
	}
	for s := 0; s < int(length); s++ {
		serviceId, dErr := deserializer.DeserializeInt64()
		if dErr != nil {
			return nil, dErr
		}
		serviceType, dErr := deserializer.DeserializeInt8()
		if dErr != nil {
			return nil, dErr
		}
		host, dErr := deserializer.DeserializeString()
		if dErr != nil {
			return nil, dErr
		}
		port, dErr := deserializer.DeserializeUint32()
		if dErr != nil {
			return nil, dErr
		}
		resp.Services = append(resp.Services, ServiceInfo{
			ServiceId:   serviceId,
			ServiceType: ServiceType(serviceType),
			Host:        host,
			Port:        port,
		})
	}
	return resp, nil
}

type ClusterInfo struct {
	ClusterId   int64
	ClusterName string
	Replica     uint32
	Zones       []string
}

type ListClusterResponse struct {
	Clusters []ClusterInfo
}

func DeserializeListClusterResponse(deserializer *Deserializer) (*ListClusterResponse, error) {
	resp := &ListClusterResponse{}
	length, err := deserializer.DeserializeUint32()
	if err != nil {
		return nil, err
	}
	for s := 0; s < int(length); s++ {
		clusterId, dErr := deserializer.DeserializeUint64()
		if dErr != nil {
			return nil, dErr
		}
		clusterName, dErr := deserializer.DeserializeString()
		if dErr != nil {
			return nil, dErr
		}
		replica, dErr := deserializer.DeserializeUint32()
		if dErr != nil {
			return nil, dErr
		}
		zones, dErr := deserializer.DeserializeStringArray()
		if dErr != nil {
			return nil, dErr
		}
		resp.Clusters = append(resp.Clusters, ClusterInfo{
			ClusterId:   int64(clusterId),
			ClusterName: clusterName,
			Replica:     replica,
			Zones:       zones,
		})
	}
	return resp, nil
}
