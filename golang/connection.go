package nebula_ng

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
    "io/ioutil"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/version"
	"google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	grpcstatus "google.golang.org/grpc/status"
)

var defaultConnector = &graphConnector{}
var defaultMsgSize = math.MaxInt64

const defaultPingTimeout = 1 * time.Second

type graphConnector struct{}

type connection struct {
	mu          sync.Mutex
	graphClient graph.GraphServiceClient
	clientConn  *grpc.ClientConn
	sessionId   int64
	timeout     time.Duration
	host        string
	port        int
}

func (c *graphConnector) connect(host *hostAddress, cfg *connConfig) (Client, error) {
	cn := &connection{
		host: host.host,
		port: host.port,
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.connectTimeout)
	defer cancel()

    if tlsCfg, err := cn.newTLSConfig(host.host, cfg); err != nil {
        return nil, err
    } else {
        if err := cn.open(host.host, host.port, cfg.connectTimeout, tlsCfg); err != nil {
            return nil, err
        }
    }

	if err := cn.authenticate(ctx, cfg.username, cfg.password); err != nil {
		return nil, err
	}
	cn.timeout = cfg.requestTimeout
	return cn, nil
}

func (cn *connection) open(host string, port int, timeout time.Duration, tlsCfg *tls.Config) error {
	var (
		err      error
		grpcConn *grpc.ClientConn
        cred     grpc.DialOption
	)

    if tlsCfg != nil {
        cred = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
    } else {
        cred = grpc.WithInsecure()
    }
    duration := time.Duration(timeout)
    grpcConn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), cred, grpc.WithBlock(), grpc.WithTimeout(duration),
    grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMsgSize), grpc.MaxCallRecvMsgSize(defaultMsgSize)))
    if err != nil {
        return errConnCannotOpen(host, port, err.Error())
    }

	cn.clientConn = grpcConn
	cn.graphClient = graph.NewGraphServiceClient(grpcConn)
	return nil
}

func (cn *connection) newTLSConfig(host string, cfg *connConfig) (*tls.Config, error) {
    if !cfg.enableTLS {
        return nil, nil
    }

    if cfg.ca == "" {
        return nil, errTLS("No CA certificate provide")
    }

    peer := cfg.peerName
    if !cfg.peerNameVerify {
        peer = ""
    } else if peer == "" {
        peer = host
    }

    tlsCfg := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         peer,
        MinVersion:         tls.VersionTLS12,
    }

    CAs := x509.NewCertPool()
    if ca, err := ioutil.ReadFile(cfg.ca); err == nil {
        if !CAs.AppendCertsFromPEM(ca) {
            return nil, errTLS(err.Error())
        }
        tlsCfg.RootCAs = CAs
    } else {
        return nil, errTLS(err.Error())
    }

    if cfg.cert != "" || cfg.key != "" {
        if cert, err := tls.LoadX509KeyPair(cfg.cert, cfg.key); err != nil {
            return nil, errTLS(err.Error())
        } else {
            tlsCfg.Certificates = []tls.Certificate{cert}
        }
    }

    tlsCfg.VerifyPeerCertificate = func(certificates [][]byte, _ [][]*x509.Certificate) error {
        certs := make([]*x509.Certificate, len(certificates))
        for i, data := range certificates {
            cert, err := x509.ParseCertificate(data)
            if err != nil {
                return errTLS(err.Error())
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

func (cn *connection) authenticate(ctx context.Context, username, password string) error {
	// TODO just simple auth with password
	authInfo := make(map[string]string)
	authInfo["password"] = password
	bs, err := json.Marshal(authInfo)
	if err != nil {
		return err
	}
	clientInfo := &common.ClientInfo{
		Lang:            common.ClientInfo_GO,
		ProtocolVersion: proto.PROTOCOL_VERSION,
		Version:         []byte(version.ClientVersion),
	}
	in := graph.AuthRequest{
		Username:   []byte(username),
		AuthInfo:   bs,
		ClientInfo: clientInfo,
	}
	resp, err := cn.graphClient.Authenticate(ctx, &in)
	if err != nil {
		_ = cn.Close()
		return err
	}
	respErr := resp.GetStatus()
	if string(respErr.GetCode()) != string(ERROR_SUCCESSFUL_COMPLETION) {
		return errServerResponse(string(respErr.GetCode()), string(respErr.GetMessage()))
	}
	cn.sessionId = resp.GetSessionId()
	return nil
}

func (cn *connection) Execute(stmt string) (Result, error) {
	if cn.timeout == 0 {
		return cn.ExecuteContext(context.Background(), stmt)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
		defer cancel()
		return cn.ExecuteContext(ctx, stmt)
	}
}

func (cn *connection) ExecuteContext(ctx context.Context, stmt string) (Result, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	cn.mu.Lock()
	defer cn.mu.Unlock()
	in := &graph.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      []byte(stmt),
	}
	resp, err := cn.graphClient.Execute(ctx, in)
	if err != nil {
		rpcErr, ok := grpcstatus.FromError(err)
		if !ok {
			return nil, err
		}
		switch rpcErr.Code() {
		case grpccodes.DeadlineExceeded, grpccodes.Canceled:
			return nil, errConnRequestTimeout(cn.host, cn.port)
		}
		return nil, err
	}

	resultResp := resultSet{
		index:   0,
		result:  resp.Result,
		summary: resp.Summary,
		cursor:  resp.Cursor,
	}
	if err := cn.isSucceed(resp); err != nil {
		return &resultResp, err
	}

	return &resultResp, nil
}

func (cn *connection) isSucceed(resp *graph.ExecuteResponse) error {
	respErr := resp.GetStatus()
	if string(respErr.GetCode()) != string(ERROR_SUCCESSFUL_COMPLETION) {
		return errServerResponse(string(respErr.GetCode()), string(respErr.GetMessage()))
	}
	return nil
}

func (cn *connection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()

	return cn.PingContext(ctx)

}

func (cn *connection) PingContext(ctx context.Context) error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	stmt := []byte("RETURN 1")
	in := &graph.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      stmt,
	}
	resp, err := cn.graphClient.Execute(ctx, in)
	if err != nil {
		return err
	}
	if err := cn.isSucceed(resp); err != nil {
		return err
	}
	return nil
}

func (cn *connection) Close() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.IsClosed() {
		return nil
	}
	// logout via statment, ignore the logout error
	in := &graph.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      []byte("SESSION CLOSE"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
	defer cancel()
	_, err := cn.graphClient.Execute(ctx, in)
	if err != nil {
		return err
	}

	if err := cn.clientConn.Close(); err != nil {
		return err
	}
	return nil
}

func (cn *connection) GetSessionId() (int64, error) {
	return cn.sessionId, nil
}

func (cn *connection) IsClosed() bool {
	return cn.clientConn.GetState() == connectivity.Shutdown
}
