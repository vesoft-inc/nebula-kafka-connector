package meta

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
    "errors"
    "io/ioutil"
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
    "google.golang.org/grpc/credentials"
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
		BackupRestoreClient
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
		AlterCluster(req *AlterClusterReq) error
		DropCluster(req *DropClusterReq) error
		InitCluster(req *InitClusterReq) error
		ListClusters(req *ListClustersReq) (*ListClustersResp, error)
		AddHost(req *AddHostReq) error
		DropHost(req *DropHostReq) error
		ListHosts(req *ListHostsReq) (*ListHostsResp, error)
		AddService(req *AddServiceReq) error
		DropService(req *DropServiceReq) error
		ListServices(req *ListServicesReq) (*ListServicesResp, error)
	}

	BackupRestoreClient interface {
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
        enableTLS      bool
        ca             string
        cert           string
        key            string
        peerNameVerify bool
        peerName       string
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

func WithTLS(enable bool, ca, cert, key string, peerNameVerify bool, peerName string) WithOption {
    return func(client *metaClient) {
        client.enableTLS = enable
        client.ca = ca
        client.cert = cert
        client.key = key
        client.peerNameVerify = peerNameVerify
        client.peerName = peerName
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

        var tlsCfg *tls.Config
        tlsCfg, err = client.newTLSConfig(host)
        if err != nil {
            continue
        }

        err = client.open(host, port, client.requestTimeout, tlsCfg)
        if err != nil {
            continue
        }

        break
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *metaClient) open(host string, port int, timeout time.Duration, tlsCfg *tls.Config) error {
	var (
		err  error
		conn *grpc.ClientConn
        cred grpc.DialOption
	)

    if tlsCfg != nil {
        cred = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
    } else {
        cred = grpc.WithInsecure()
    }

    duration := time.Duration(timeout)
    conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), cred, grpc.WithBlock(), grpc.WithTimeout(duration),
    grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMsgSize), grpc.MaxCallRecvMsgSize(defaultMsgSize)))

    if err != nil {
        return err
    }

	c.clientConn = conn
	c.client = admin.NewAdminServiceClient(conn)
	return nil
}

func (c *metaClient) newTLSConfig(host string) (*tls.Config, error) {
    if !c.enableTLS {
        return nil, nil
    }

    if c.ca == "" {
        return nil, errors.New("No CA certificate provide")
    }

    peer := c.peerName
    if !c.peerNameVerify {
        peer = ""
    } else if peer == "" {
        peer = host
    }

    tlsCfg := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         peer,
        MinVersion:         tls.VersionTLS13,
    }

    CAs := x509.NewCertPool()
    if ca, err := ioutil.ReadFile(c.ca); err == nil {
        if !CAs.AppendCertsFromPEM(ca) {
            return nil, err
        }
        tlsCfg.RootCAs = CAs
    } else {
        return nil, err
    }

    if c.cert != "" || c.key != "" {
        if cert, err := tls.LoadX509KeyPair(c.cert, c.key); err != nil {
            return nil, err
        } else {
            tlsCfg.Certificates = []tls.Certificate{cert}
        }
    }

    tlsCfg.VerifyPeerCertificate = func(certificates [][]byte, _ [][]*x509.Certificate) error {
        certs := make([]*x509.Certificate, len(certificates))
        for i, data := range certificates {
            cert, err := x509.ParseCertificate(data)
            if err != nil {
                return err
            }
            certs[i] = cert
        }

        opts := x509.VerifyOptions{
            Roots:          tlsCfg.RootCAs,
            DNSName:        tlsCfg.ServerName,
            Intermediates:  x509.NewCertPool(),
        }

        for _, cert := range certs[1:] {
            opts.Intermediates.AddCert(cert)
        }

        _, err := certs[0].Verify(opts)

        return err
    }

    return tlsCfg, nil
}

func (c *metaClient) Login() (*LoginResponse, error) {
	if c.user == "" || c.password == "" {
		return nil, fmt.Errorf("user and password are required")
	}
	resp, err := c.authWithPassword(c.user, c.password)
	if err != nil {
		return nil, err
	}
	c.token = resp.Token
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
