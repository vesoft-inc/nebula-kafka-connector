package meta

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/meta"
	"google.golang.org/grpc"
)

var defaultMsgSize = math.MaxInt64
var defaultConnectTimeout = 3 * time.Second
var defaultRequestTimeout = 10 * time.Second

type (
	Client interface {
		Close()
		ClusterClient
	}

	ClusterClient interface {
		CreateCluster(req *CreateClusterReq) (*CreateClusterResp, error)
		AddService(req *AddServiceReq) (*AddServiceResp, error)
		ShowService(req *ShowServiceReq) (*ShowServiceResp, error)
		InitCluster(req *InitClusterReq) (*InitClusterResp, error)
		ShowCluster(req *ShowClusterReq) (*ShowClusterResp, error)
	}

	metaClient struct {
		address    string
		client     meta.MetaServiceClient
		clientConn *grpc.ClientConn
		retryTimes int
		timeout    time.Duration
	}

	responseHeader interface {
		GetHeader() *meta.ResponseHeader
	}

	WithOption func(*metaClient)
)

func WithTimeout(timeout time.Duration) WithOption {
	return func(client *metaClient) {
		client.timeout = timeout
	}
}

func WithRetryTimes(retryTimes int) WithOption {
	return func(client *metaClient) {
		client.retryTimes = retryTimes
	}
}

func NewMetaClient(addresses string, opts ...WithOption) (Client, error) {
	//TODO should verify the address
	// if the address is invalid, then return error
	addrs := strings.Split(addresses, ",")
	if len(addrs) == 0 {
		return nil, fmt.Errorf("invalid address")
	}
	var (
		client *metaClient
		err    error
	)
	for i := 0; i < len(addrs); i++ {
		var port int
		addr := addrs[i]
		if len(strings.Split(addr, ":")) != 2 {
			return nil, fmt.Errorf("invalid address")
		}

		host := strings.Split(addr, ":")[0]
		p := strings.Split(addr, ":")[1]
		port, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid address")
		}
		client = &metaClient{
			address:    addr,
			retryTimes: 1,
			timeout:    defaultRequestTimeout,
		}
		for _, opt := range opts {
			opt(client)
		}
		err = client.open(host, port, defaultConnectTimeout, nil)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *metaClient) open(host string, port int, timeout time.Duration, sslConfig *tls.Config) error {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if sslConfig != nil {
		return fmt.Errorf("ssl is not supported")
	} else {
		timeout := time.Duration(timeout)
		conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMsgSize), grpc.MaxCallRecvMsgSize(defaultMsgSize)))

		if err != nil {
			return err
		}
	}

	c.clientConn = conn
	c.client = meta.NewMetaServiceClient(conn)
	return nil
}

func (c *metaClient) Close() {
	if c.clientConn != nil {
		c.clientConn.Close()
	}
}

func (c *metaClient) retry(fn func() (responseHeader, error)) (responseHeader, error) {
	var (
		resp responseHeader
		err  error
	)
	for i := 0; i < c.retryTimes+1; i++ {
		resp, err = fn()
		if err != nil {
			continue
		}
		header := resp.GetHeader()
		if header.GetOk() {
			return resp, nil
		}
		// if the error is not leader change, then return and do not retry
		if nebula.ErrorFromInt(header.GetCode()) != nebula.ErrorLeaderChange {
			return resp, nil
		}
		newLeader := header.GetLeader()
		if newLeader == nil {
			return nil, fmt.Errorf("invalid leader")
		}
		c.address = fmt.Sprintf("%s:%d", newLeader.GetHost(), newLeader.GetPort())
		c.Close()
		if err := c.open(string(newLeader.GetHost()), int(newLeader.GetPort()),
			defaultConnectTimeout, nil); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getResponseHeader(respHeader responseHeader) (*HeaderResponse, error) {
	header := respHeader.GetHeader()
	if header == nil {
		return nil, fmt.Errorf("invalid response")
	}
	leader := header.GetLeader()
	errorCode := nebula.ErrorFromInt(header.GetCode())
	result := &HeaderResponse{
		OK:   header.GetOk(),
		Code: errorCode,
		Msg:  string(header.GetMessage()),
	}
	if leader == nil {
		result.NewHost = ""
		result.NewPort = 0
	} else {
		result.NewHost = string(leader.GetHost())
		result.NewPort = leader.GetPort()
	}
	return result, nil
}

func (c *metaClient) CreateCluster(req *CreateClusterReq) (*CreateClusterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	zones := make([][]byte, 0)
	for _, z := range req.zones {
		zones = append(zones, []byte(z))
	}
	in := &meta.CreateClusterRequest{
		Header: &meta.RequestHeader{},
		ClusterDesc: &meta.ClusterDesc{
			ClusterName:   []byte(req.clusterName),
			ReplicaFactor: uint32(req.replica),
			Zones:         zones,
		},
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.CreateCluster(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*meta.CreateClusterResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
	}
	return &CreateClusterResp{
		HeaderResponse: responseHeader}, nil
}

func (c *metaClient) AddService(req *AddServiceReq) (*AddServiceResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	in := &meta.AddServiceRequest{
		Header:      &meta.RequestHeader{},
		Host:        []byte(req.host),
		Port:        req.port,
		Type:        meta.ServiceType(req.serviceType),
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.AddService(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*meta.AddServiceResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
	}
	return &AddServiceResp{
		HeaderResponse: responseHeader,
	}, nil
}
func (c *metaClient) ShowService(req *ShowServiceReq) (*ShowServiceResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &meta.ShowServiceRequest{
		Header:      &meta.RequestHeader{},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ShowService(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*meta.ShowServiceResponse)
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
	return &ShowServiceResp{
		HeaderResponse: responseHeader,
		Services:       services,
	}, nil
}
func (c *metaClient) InitCluster(req *InitClusterReq) (*InitClusterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &meta.InitStorageRequest{
		Header:      &meta.RequestHeader{},
		ClusterName: []byte(req.clustername),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.InitStorage(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*meta.InitStorageResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
	}
	return &InitClusterResp{
		HeaderResponse: responseHeader,
	}, nil

}
func (c *metaClient) ShowCluster(req *ShowClusterReq) (*ShowClusterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &meta.ShowClusterRequest{
		Header:      &meta.RequestHeader{},
		ClusterName: []byte(req.clusterName),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ShowCluster(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*meta.ShowClusterResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return nil, err
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
