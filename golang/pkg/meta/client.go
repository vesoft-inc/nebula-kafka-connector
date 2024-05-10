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
		//if leader change, error is not nil, and return leader info in LoginResponse
		Login() (*LoginResponse, error)
		Logout() error
		UserClient
		ClusterClient
	}

	LoginResponse struct {
		Token  []byte `json:"token"`
		Leader string `json:"leader"`
	}

	UserClient interface {
		// change the password of the root with the current password and the new password
		ChangePassword(req *ChangePasswordReq) error
		CreateUser(req *CreateUserReq) error
		AlterUser(req *AlterUserReq) error
		//if the user list is empty, would show all the users
		ListUsers(req *ListUsersReq) (*ListUsersResp, error)
		DropUser(req *DropUserReq) error
	}

	ClusterClient interface {
		CreateCluster(req *CreateClusterReq) error
		AddService(req *AddServiceReq) error
		DropService(req *DropServiceReq) error
		ListServices(req *ListServicesReq) (*ListServicesResp, error)
		InitCluster(req *InitClusterReq) error
		ListClusters(req *ListClustersReq) (*ListClustersResp, error)
		DropCluster(req *DropClusterReq) error
		CreateBackup(req *CreateBackupReq) (*CreateBackupResp, error)
		DropBackup(req *DropBackupReq) (*DropBackupResp, error)
		Restore(req *RestoreReq) (*RestoreResp, error)
		ShowMeta() (*ShowMetaResp, error)
	}

	metaClient struct {
		address        string
		client         admin.AdminServiceClient
		clientConn     *grpc.ClientConn
		requestTimeout time.Duration
		connectTimeout time.Duration
		token          []byte
		user           string
		password       string
	}

	responseHeader interface {
		GetHeader() *admin.ResponseHeader
	}

	WithOption func(*metaClient)
)

func WithRequestTimeout(timeout time.Duration) WithOption {
	return func(client *metaClient) {
		client.requestTimeout = timeout
	}
}

func WithConnectTimeout(timeout time.Duration) WithOption {
	return func(client *metaClient) {
		client.connectTimeout = timeout
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
			address:        addr,
			requestTimeout: defaultRequestTimeout,
			connectTimeout: defaultConnectTimeout,
		}
		for _, opt := range opts {
			opt(client)
		}
		err = client.open(host, port, client.requestTimeout, nil)
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

func (c *metaClient) Login() (*LoginResponse, error) {
	if c.user == "" || c.password == "" {
		return nil, fmt.Errorf("user and password are required")
	}
	resp, err := c.authWithPassword(c.user, c.password)
	if err != nil {
		return resp, err
	}
	return resp, nil

}

func (c *metaClient) authWithPassword(user string, password string) (*LoginResponse, error) {
	info := make(map[string]interface{})
	info["password"] = password
	return c.auth(user, info)
}

func (c *metaClient) auth(user string, authInfo map[string]interface{}) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	bs, err := json.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	clientInfo := &common.ClientInfo{
		Lang:            common.ClientInfo_GO,
		ProtocolVersion: proto.PROTOCOL_VERSION,
		Version:         []byte(version.ClientVersion),
	}
	in := &admin.LoginRequest{
		Username:   []byte(user),
		AuthInfo:   bs,
		ClientInfo: clientInfo,
	}
	response, err := c.client.Login(ctx, in)
	if err != nil {
		return nil, err
	}

	// retry for leader change
	if nebula.ErrorCode(response.Header.GetStatus().GetCode()) == nebula.ERROR_LEADER_CHANGED {
		leader := response.Header.GetLeader()
		if leader == nil {
			return nil, fmt.Errorf("invalid leader")
		}
		c.Close()
		c.address = fmt.Sprintf("%s:%d", leader.GetHost(), leader.GetPort())
		err := c.open(string(leader.GetHost()), int(leader.GetPort()), c.connectTimeout, nil)
		if err != nil {
			return nil, err
		}
		response, err = c.client.Login(ctx, in)
		if err != nil {
			return nil, err
		}
		if nebula.ErrorCode(response.Header.GetStatus().GetCode()) != nebula.ERROR_SUCCESSFUL_COMPLETION {
			return nil, nebula.NewNebulaError(
				nebula.ErrorCode(string(response.Header.GetStatus().GetCode())),
				string(response.Header.GetStatus().GetMessage()),
			)
		}
	} else if nebula.ErrorCode(response.Header.GetStatus().GetCode()) != nebula.ERROR_SUCCESSFUL_COMPLETION {
		return nil, nebula.NewNebulaError(
			nebula.ErrorCode(string(response.Header.GetStatus().GetCode())),
			string(response.Header.GetStatus().GetMessage()),
		)
	}
	if response.Token == nil {
		return nil, fmt.Errorf("invalid token")
	}
	r := &LoginResponse{
		Token:  response.Token,
		Leader: c.address,
	}
	return r, nil
}

func (c *metaClient) Close() {
	if c.clientConn != nil {
		c.clientConn.Close()
	}
}

// there's no need to retry, because the token is invalid after leader chanage.
func (c *metaClient) execute(fn func() (responseHeader, error)) (responseHeader, error) {
	var (
		resp responseHeader
		err  error
	)
	resp, err = fn()
	if err != nil {
		return nil, err
	}
	header := resp.GetHeader()
	if nebula.ErrorFromBytes(header.GetStatus().GetCode()) == nebula.ERROR_SUCCESSFUL_COMPLETION {
		return resp, nil
	} else {
		return nil, nebula.NewNebulaError(
			nebula.ErrorCode(string(header.GetStatus().GetCode())),
			string(header.GetStatus().GetMessage()),
		)
	}
}

func getResponseHeader(respHeader responseHeader) (*HeaderResponse, error) {
	header := respHeader.GetHeader()
	if header == nil {
		return nil, fmt.Errorf("invalid response")
	}
	leader := header.GetLeader()
	errorCode := nebula.ErrorFromBytes(header.GetStatus().GetCode())
	result := &HeaderResponse{
		Code: errorCode,
		Msg:  string(header.GetStatus().GetMessage()),
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

func responseIsErr(resp responseHeader) error {
	responseHeader, err := getResponseHeader(resp)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetStatus()
	}
	return nil
}

func (c *metaClient) GetToken() []byte {
	return c.token
}

func (c *metaClient) Logout() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.LogoutRequest{
		Header: &admin.RequestHeader{Token: c.token},
	}
	_, err := c.execute(func() (responseHeader, error) {
		return c.client.Logout(ctx, in)
	})
	if err != nil {
		return err
	}
	return nil
}
