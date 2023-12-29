/*
Copyright 2023 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nebula

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/klog/v2"

	"github.com/vesoft-inc/nebula-ng-tools/golang/nrpc"
)

const (
	defaultRPCTimeout        = 200 * time.Millisecond
	defaultReconnectAttempts = 3
	defaultReconnectDelay    = time.Second
)

// ErrorCode reference https://github.com/vesoft-inc/nebula-ng/blob/master/src/common/base/ErrorCode.h
type ErrorCode uint64

const (
	ErrorCodeSucceeded            ErrorCode = 0
	ErrorCodeLeaderChanged        ErrorCode = 21474836486
	ErrorCodeClusterAlreadyExists ErrorCode = 17609365913604
	ErrorCodeServiceAlreadyExists ErrorCode = 18047452577796
	ErrorCodeClusterNotFound      ErrorCode = 17613660880900
	ErrorCodeServiceNotAdded      ErrorCode = 18051747545093
)

var (
	ErrNoAvailableEndpoints = errors.New("metadclient: no available hosts")
	ErrRPCTimeout           = errors.New("metadclient: rpc timeout")
	ErrReconnectFailed      = errors.New("metadclient: reconnect failed")
	ErrClusterNotFound      = errors.New("metadclient: cluster not found")
	ErrLeaderHostNotFound   = errors.New("metadclient: leader host not found")
)

type MetaInterface interface {
	CreateCluster(clusterName string, replica int, zones []string) error

	RemoveCluster(clusterName string) error

	GetCluster(clusterName string) (*ClusterInfo, error)

	ListClusters() ([]ClusterInfo, error)

	InitCluster(clusterName string) error

	AddService(host string, port uint32, serviceType ServiceType, clusterName string) error

	DropService(clusterName string, serviceType ServiceType, host string, port uint32) error

	ListServices(clusterName string) ([]ServiceInfo, error)

	Disconnect() error
}

var _ MetaInterface = (*metaClient)(nil)

type metaClient struct {
	client  *nrpc.Client
	options []Option
}

func NewMetaClient(hosts []string, options ...Option) (MetaInterface, error) {
	if len(hosts) == 0 {
		return nil, ErrNoAvailableEndpoints
	}
	mc, err := newMetaConnection(hosts[0], options...)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func newMetaConnection(endpoint string, options ...Option) (*metaClient, error) {
	rpcClient, err := buildNRpcClient(endpoint, options...)
	if err != nil {
		return nil, err
	}
	mc := &metaClient{
		client:  rpcClient,
		options: options,
	}
	return mc, nil
}

func (m *metaClient) reconnect(endpoint string) error {
	m.client.Close()
	rpcClient, err := buildNRpcClient(endpoint, m.options...)
	if err != nil {
		return err
	}
	m.client = rpcClient
	return nil
}

func (m *metaClient) send(request []byte) ([]byte, error) {
	timeout := defaultRPCTimeout
	opts := loadOptions(m.options...)
	if opts.Timeout > 0 {
		timeout = opts.Timeout
	}
	resp, err := m.client.Send(request, timeout)
	if err != nil {
		nErr, ok := err.(nrpc.Error)
		if ok {
			if nErr.Timeout() {
				return nil, ErrRPCTimeout
			} else if nErr.BadChannel() {
				if err = m.client.Reconnect(defaultReconnectAttempts, defaultReconnectDelay); err != nil {
					return nil, ErrReconnectFailed
				}
				resp, err = m.client.Send(request, timeout)
				if err != nil {
					return nil, err
				}
				return resp, nil
			} else {
				return nil, fmt.Errorf("nrpc error: %v", err)
			}
		}
		return nil, fmt.Errorf("unknown nrpc error: %v", err)
	}
	return resp, nil
}

func (m *metaClient) CreateCluster(clusterName string, replica int, zones []string) error {
	request := NewCreateClusterRequest(clusterName, replica, zones)
	bytes, err := request.Serialize()
	if err != nil {
		return err
	}
	_, err = m.retryOnError(bytes)
	return err
}

func (m *metaClient) RemoveCluster(_ string) error {
	//TODO implement me
	panic("implement me")
}

func (m *metaClient) GetCluster(clusterName string) (*ClusterInfo, error) {
	request := NewListClusterRequest(clusterName)
	bytes, err := request.Serialize()
	if err != nil {
		return nil, err
	}
	deserializer, err := m.retryOnError(bytes)
	if err != nil {
		return nil, err
	}
	clusterResp, err := DeserializeListClusterResponse(deserializer)
	if err != nil {
		return nil, err
	}
	if len(clusterResp.Clusters) == 0 {
		return nil, ErrClusterNotFound
	}
	return &clusterResp.Clusters[0], nil
}

func (m *metaClient) ListClusters() ([]ClusterInfo, error) {
	request := NewListClusterRequest("")
	bytes, err := request.Serialize()
	if err != nil {
		return nil, err
	}
	deserializer, err := m.retryOnError(bytes)
	if err != nil {
		return nil, err
	}
	clusterResp, err := DeserializeListClusterResponse(deserializer)
	if err != nil {
		return nil, err
	}
	return clusterResp.Clusters, nil
}

func (m *metaClient) InitCluster(clusterName string) error {
	request := NewInitClusterRequest(clusterName)
	bytes, err := request.Serialize()
	if err != nil {
		return err
	}
	_, err = m.retryOnError(bytes)
	return err
}

func (m *metaClient) AddService(host string, port uint32, serviceType ServiceType, clusterName string) error {
	request := NewAddServiceRequest(host, port, serviceType, clusterName)
	bytes, err := request.Serialize()
	if err != nil {
		return err
	}
	_, err = m.retryOnError(bytes)
	return err
}

func (m *metaClient) DropService(clusterName string, serviceType ServiceType, host string, port uint32) error {
	request := NewDropServiceRequest(clusterName, serviceType, host, port)
	bytes, err := request.Serialize()
	if err != nil {
		return err
	}
	_, err = m.retryOnError(bytes)
	return err
}

func (m *metaClient) ListServices(clusterName string) ([]ServiceInfo, error) {
	request := NewListServiceRequest(clusterName)
	bytes, err := request.Serialize()
	if err != nil {
		return nil, err
	}
	deserializer, err := m.retryOnError(bytes)
	if err != nil {
		return nil, err
	}
	serviceResp, err := DeserializeListServiceResponse(deserializer)
	if err != nil {
		return nil, err
	}
	return serviceResp.Services, nil
}

func (m *metaClient) Disconnect() error {
	m.client.Close()
	return nil
}

func (m *metaClient) retryOnError(request []byte) (*Deserializer, error) {
	fn := func(req []byte, leaderChanged bool) (*ResponseHeader, *Deserializer, error) {
		resp, err := m.send(request)
		if err != nil {
			return nil, nil, err
		}
		deserializer := NewDeserializer(resp)
		respHeader, err := DeserializeHeader(deserializer, leaderChanged)
		if err != nil {
			return nil, nil, err
		}
		return respHeader, deserializer, nil
	}

	respHeader, deserializer, err := fn(request, true)
	if err != nil {
		return nil, err
	}
	code := ErrorCode(respHeader.Code)
	if code != ErrorCodeSucceeded {
		if code == ErrorCodeLeaderChanged {
			leaderHost := fmt.Sprintf("%s:%d", respHeader.Host, respHeader.Port)
			if leaderHost == "" {
				return nil, ErrLeaderHostNotFound
			}
			klog.Infof("leader changed, reconnect to host: %s", leaderHost)
			// update leader info
			err = m.reconnect(leaderHost)
			if err != nil {
				return nil, err
			}
			respHeader, deserializer, err = fn(request, false)
			if err != nil {
				return nil, err
			}
			code = ErrorCode(respHeader.Code)
			if code != ErrorCodeSucceeded {
				if code == ErrorCodeClusterAlreadyExists ||
					code == ErrorCodeServiceAlreadyExists ||
					code == ErrorCodeClusterNotFound ||
					code == ErrorCodeServiceNotAdded {
					return deserializer, nil
				}
				return nil, fmt.Errorf("metad client retry response code %d, msg: %s", respHeader.Code, respHeader.Msg)
			}
			return deserializer, nil
		} else if code == ErrorCodeClusterAlreadyExists ||
			code == ErrorCodeServiceAlreadyExists ||
			code == ErrorCodeClusterNotFound ||
			code == ErrorCodeServiceNotAdded {
			return deserializer, nil
		}
		return nil, fmt.Errorf("response code: %d, msg: %s", respHeader.Code, respHeader.Msg)
	}
	return deserializer, nil
}
