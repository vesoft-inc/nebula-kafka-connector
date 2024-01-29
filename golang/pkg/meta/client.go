package meta

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/nrpc"
)

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
		address      string
		client       sender
		retryTimes   int
		timeout      time.Duration
		serializer   serializer
		deserializer deserializer
	}

	// interface for nrpc client
	sender interface {
		Send(req []byte, timeout time.Duration) ([]byte, error)
		Close()
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
	rand.Seed(time.Now().UnixNano())
	addr := addrs[rand.Intn(len(addrs))]
	client := &metaClient{
		address:      addr,
		retryTimes:   1,
		timeout:      3 * time.Second,
		serializer:   newDefaultSerializer(),
		deserializer: newDefaultDeserializer(),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func (c *metaClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *metaClient) send(req reqSerializer, resp respDeserializer) error {
	if c.client == nil {
		c.client = nrpc.NewClient(c.address)
	}
	c.serializer.reset()
	c.deserializer.reset()

	bytes := req.serialize(c.serializer)
	response, err := c.client.Send(bytes, c.timeout)
	if err != nil {
		return err
	}
	c.deserializer.setBytes(response)
	if err = resp.deserialize(c.deserializer); err != nil {
		return err
	}

	return nil
}

func (c *metaClient) sendWithRetry(req reqSerializer, resp respDeserializer) error {
	var err error
	for i := 0; i < c.retryTimes+1; i++ {
		if err = c.send(req, resp); err != nil {
			continue
		}
		// if the leader is changed, should update the client, and retry
		// else return without retry
		hresp, ok := resp.(responseHeader)
		if !ok {
			return fmt.Errorf("internal error: resp is not responseHeader")
		}
		header := hresp.getHeader()
		if header.Code == nebula.ErrorSuccessfulCompletion {
			return nil
		}
		if header.Code == nebula.ErrorLeaderChange {
			c.client = nil
			c.address = fmt.Sprintf("%s:%d", header.NewHost, header.NewPort)
			continue
		}
		// other error should not retry
		// and then return the header instead of error
		// e.g. create an existed cluster, then the server will return the error
		return nil
	}
	return err
}

func (c *metaClient) CreateCluster(req *CreateClusterReq) (*CreateClusterResp, error) {
	resp := &CreateClusterResp{}
	if err := c.sendWithRetry(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *metaClient) AddService(req *AddServiceReq) (*AddServiceResp, error) {
	resp := &AddServiceResp{}
	if err := c.sendWithRetry(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *metaClient) ShowCluster(req *ShowClusterReq) (*ShowClusterResp, error) {
	resp := &ShowClusterResp{}
	if err := c.sendWithRetry(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *metaClient) InitCluster(req *InitClusterReq) (*InitClusterResp, error) {
	resp := &InitClusterResp{}
	if err := c.sendWithRetry(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *metaClient) ShowService(req *ShowServiceReq) (*ShowServiceResp, error) {
	resp := &ShowServiceResp{}
	if err := c.sendWithRetry(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
