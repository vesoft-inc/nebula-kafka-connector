package meta

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	admin "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/admin"
	common "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/version"
	"google.golang.org/grpc"
)

var defaultMsgSize = math.MaxInt64
var defaultConnectTimeout = 3 * time.Second
var defaultRequestTimeout = 10 * time.Second

type (
	Client interface {
		Close()
		Login() error
		Logout() error
		GetToken() []byte
		UserClient
		ClusterClient
	}

	UserClient interface {
		// change the password of the root with the current password and the new password
		ChangePassword(req *ChangePasswordReq) error
	}

	ClusterClient interface {
		CreateCluster(req *CreateClusterReq) error
		AddService(req *AddServiceReq) error
		DropService(req *DropServiceReq) error
		ShowService(req *ShowServiceReq) (*ShowServiceResp, error)
		InitCluster(req *InitClusterReq) error
		ShowCluster(req *ShowClusterReq) (*ShowClusterResp, error)
	}

	metaClient struct {
		address    string
		client     admin.AdminServiceClient
		clientConn *grpc.ClientConn
		retryTimes int
		timeout    time.Duration
		token      []byte
		user       string
		password   string
	}

	responseHeader interface {
		GetHeader() *common.ResponseHeader
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

func WithUserPassword(user, password string) WithOption {
	return func(client *metaClient) {
		client.user = user
		client.password = password
	}
}

func WithToken(token []byte) WithOption {
	return func(client *metaClient) {
		client.token = token
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
		if err != nil {
			continue
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
	c.client = admin.NewAdminServiceClient(conn)
	return nil
}

func (c *metaClient) Login() error {
	if c.user == "" || c.password == "" {
		return fmt.Errorf("user and password are required")
	}
	token, err := c.authWithPassword(c.user, c.password)
	if err != nil {
		return err
	}
	c.token = token
	return nil

}

func (c *metaClient) authWithPassword(user string, password string) ([]byte, error) {
	info := make(map[string]interface{})
	info["password"] = password
	return c.auth(user, info)
}

func (c *metaClient) auth(user string, authInfo map[string]interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	bs, err := json.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	clientInfo := &common.ClientInfo{
		ClientLang:      common.ClientInfo_GO,
		ProtocolVersion: proto.PROTOCOL_VERSION,
		ClientVersion:   []byte(version.ClientVersion),
	}
	in := &admin.LoginRequest{
		Username:   []byte(user),
		AuthInfo:   bs,
		ClientInfo: clientInfo,
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.Login(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.LoginResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	nebulaErr := nebula.ErrorFromInt(response.Header.Error.Code)
	if nebulaErr != nebula.ERROR_SUCCESSFUL_COMPLETION {
		return nil, nebula.NewNebulaError(
			nebula.ErrorFromInt(response.Header.GetError().GetCode()),
			string(response.Header.GetError().GetMessage()),
		)
	}
	if response.Token == nil {
		return nil, fmt.Errorf("invalid token")
	}
	return response.Token, nil
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
		if nebula.ErrorFromInt(header.GetError().GetCode()) == nebula.ERROR_SUCCESSFUL_COMPLETION {
			return resp, nil
		}
		// if the error is not leader change, then return and do not retry
		if nebula.ErrorFromInt(header.GetError().GetCode()) != nebula.ERROR_LEADER_CHANGED {
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
	errorCode := nebula.ErrorFromInt(header.GetError().GetCode())
	result := &HeaderResponse{
		Code: errorCode,
		Msg:  string(header.GetError().GetMessage()),
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

func (c *metaClient) GetToken() []byte {
	return c.token
}

func (c *metaClient) Logout() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.LogoutRequest{
		Header: &admin.AdminRequestHeader{Token: c.token},
	}
	_, err := c.retry(func() (responseHeader, error) {
		return c.client.Logout(ctx, in)
	})
	if err != nil {
		return err
	}
	return nil
}
