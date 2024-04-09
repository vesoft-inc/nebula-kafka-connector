package nebula_ng

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

var defaultConnector = &graphConnector{}
var defaultMsgSize = math.MaxInt64

type graphConnector struct{}

type connection struct {
	mu          sync.Mutex
	graphClient graph.GraphServiceClient
	clientConn  *grpc.ClientConn
	sessionId   int64
	timeout     time.Duration
}

func (c *graphConnector) connect(host *hostAddress, cfg *connConfig) (Conn, error) {
	cn := &connection{}
	start := time.Now()
	if err := cn.open(host.host, host.port, cfg.connectTimeout, nil); err != nil {
		return nil, err
	}
	if err := cn.authenticate(cfg.username, cfg.password); err != nil {
		return nil, err
	}
	if cfg.connectTimeout > 0 && time.Since(start) > cfg.connectTimeout {
		return nil, errConnConnectTimeout(host.host, host.port)
	}

	if cfg.graph != "" {
		// TODO set graph
	}
	cn.timeout = cfg.requestTimeout
	return cn, nil
}

func (cn *connection) open(host string, port int, timeout time.Duration, sslConfig *tls.Config) error {
	var (
		err      error
		grpcConn *grpc.ClientConn
	)
	if sslConfig != nil {
		return fmt.Errorf("ssl is not supported")
	} else {
		timeout := time.Duration(timeout)
		grpcConn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMsgSize), grpc.MaxCallRecvMsgSize(defaultMsgSize)))
		if err != nil {
			return errConnCannotOpen(host, port, err.Error())
		}
	}

	cn.clientConn = grpcConn
	cn.graphClient = graph.NewGraphServiceClient(grpcConn)
	return nil
}

func (cn *connection) authenticate(username, password string) error {
	// TODO just simple auth with password
	authInfo := make(map[string]string)
	authInfo["password"] = password
	bs, err := json.Marshal(authInfo)
	if err != nil {
		return err
	}
	clientInfo := &common.ClientInfo{
		ClientLang:      common.ClientInfo_GO,
		ProtocolVersion: proto.PROTOCOL_VERSION,
		// TODO(yee): set client version to nebula-go version
		ClientVersion: proto.PROTOCOL_VERSION,
	}
	in := graph.AuthRequest{
		Username:   []byte(username),
		AuthInfo:   bs,
		ClientInfo: clientInfo,
	}
	resp, err := cn.graphClient.Authenticate(context.Background(), &in)
	if err != nil {
		_ = cn.Close()
		return err
	}
	if string(resp.GetGqlStatus().GetCode()) != string(ERROR_SUCCESSFUL_COMPLETION) {
		return errServerResponse(string(resp.GetGqlStatus().GetCode()), resp.GetGqlStatus().String())
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
			return nil, errConnRequestTimeout("", 0)
		}
		return nil, err
	}
	if resp.ExecutionOutcome == nil {
		return nil, errInternel("execute failed, response is nil")
	}

	resultResp := resultSet{
		index:     0,
		latency:   (resp.LatencyInUs),
		result:    resp.ExecutionOutcome.Result,
		summary:   resp.ExecutionOutcome.Summary,
		extraInfo: resp.ExecutionOutcome.ExtraInfo,
	}

	if string(resp.ExecutionOutcome.GetGqlStatus().GetCode()) != string(ERROR_SUCCESSFUL_COMPLETION) {
		return &resultResp, errServerResponse(
			string(resp.ExecutionOutcome.GetGqlStatus().GetCode()),
			(string(resp.GetExecutionOutcome().GetGqlStatus().GetMessage())))
	}

	return &resultResp, nil
}

func (cn *connection) Ping() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	stmt := []byte("RETURN 1")
	in := &graph.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      stmt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
	defer cancel()
	_, err := cn.graphClient.Execute(ctx, in)
	return err
}

func (cn *connection) Close() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	// logout via statment, ignore the logout error
	in := &graph.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      []byte("SESSION CLOSE"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
	defer cancel()
	_, _ = cn.graphClient.Execute(ctx, in)
	_ = cn.clientConn.Close()
	return nil
}
