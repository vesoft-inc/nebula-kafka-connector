package nebula_ng

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/decode"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/grpcutil"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/internal_error"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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
	tlsConfig   *tls.Config
	host        string
	port        int
}

func (c *graphConnector) connect(host *hostAddress, cfg *connConfig) (types.Client, error) {
	cn := &connection{
		host: host.host,
		port: host.port,
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.connectTimeout)
	defer cancel()
	var (
		tlsCfg *tls.Config
		err    error
	)
	if cfg.enableTLS {
		tlsCfg, err = grpcutil.NewTLSConfig(host.host, cfg.ca, cfg.cert, cfg.key, cfg.peerName, cfg.peerNameVerify)
		if err != nil {
			return nil, err
		}
	}

	if err := cn.open(host.host, host.port, cfg.connectTimeout, tlsCfg); err != nil {
		return nil, err
	}

	if err := cn.authenticate(ctx, cfg.username, cfg.password); err != nil {
		return nil, err
	}
	cn.timeout = cfg.requestTimeout
	return cn, nil
}

func (cn *connection) open(host string, port int, timeout time.Duration, tlsCfg *tls.Config) error {
	grpcConn, err := grpcutil.NewGrpcClient(host, port, timeout, tlsCfg)
	if err != nil {
		return err
	}
	cn.clientConn = grpcConn
	cn.graphClient = graph.NewGraphServiceClient(grpcConn)
	return nil
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
		return grpcutil.GetGrpcError(fmt.Sprintf("%s:%d", cn.host, cn.port), err)
	}
	respErr := resp.GetStatus()
	if string(respErr.GetCode()) != string(errors.ERROR_SUCCESSFUL_COMPLETION) {
		return internal_error.ErrServerResponse(string(respErr.GetCode()), string(respErr.GetMessage()))
	}
	cn.sessionId = resp.GetSessionId()
	return nil
}

func (cn *connection) Execute(stmt string) (types.Result, error) {
	if cn.timeout == 0 {
		return cn.ExecuteContext(context.Background(), stmt)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
		defer cancel()
		return cn.ExecuteContext(ctx, stmt)
	}
}

func (cn *connection) ExecuteContext(ctx context.Context, stmt string) (types.Result, error) {
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
		return nil, grpcutil.GetGrpcError(fmt.Sprintf("%s:%d", cn.host, cn.port), err)
	}
	t, err := decode.NewResultTable(resp.Result)
	if err != nil {
		return nil, err
	}
	resultResp := resultSet{
		index:   0,
		table:   t,
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
	if string(respErr.GetCode()) != string(errors.ERROR_SUCCESSFUL_COMPLETION) {
		return internal_error.ErrServerResponse(string(respErr.GetCode()), string(respErr.GetMessage()))
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
