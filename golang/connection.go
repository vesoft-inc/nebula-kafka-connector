package nebula_ng

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

var defaultConnector = &graphConnector{}
var defaultMsgSize = math.MaxInt64

type graphConnector struct{}

type connection struct {
	mu          sync.Mutex
	graphClient proto.GraphClient
	clientConn  *grpc.ClientConn
	sessionId   int64
	timeout     time.Duration
}

type resultSet struct {
	index    int
	latency  uint64
	result   *proto.ResultTable
	planDesc []byte
}

type rowData struct {
	resultSet *resultSet
	values    []Value
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
	cn.graphClient = proto.NewGraphClient(grpcConn)
	return nil
}

func (cn *connection) authenticate(username, password string) error {
	in := proto.AuthRequest{
		Username: []byte(username),
		Password: []byte(password),
	}
	resp, err := cn.graphClient.Authenticate(context.Background(), &in)
	if err != nil {
		_ = cn.Close()
		return err
	}
	if string(resp.GetGqlStatus().GetCode()) != ErrorSuccessfulCompletion {
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
	in := &proto.ExecuteRequest{
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

	if string(resp.ExecutionOutcome.GetGqlStatus().GetCode()) != (ErrorSuccessfulCompletion) {
		return nil, errServerResponse(
			string(resp.ExecutionOutcome.GetGqlStatus().GetCode()),
			(string(resp.GetExecutionOutcome().GetGqlStatus().GetMessage())))
	}

	return &resultSet{
		index:    0,
		latency:  (resp.LatencyInUs),
		result:   resp.ExecutionOutcome.Result,
		planDesc: resp.ExecutionOutcome.PlanDesc,
	}, nil
}

func (cn *connection) Ping() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	stmt := []byte("RETURN 1")
	in := &proto.ExecuteRequest{
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
	in := &proto.ExecuteRequest{
		SessionId: cn.sessionId,
		Stmt:      []byte("SESSION CLOSE"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), cn.timeout)
	defer cancel()
	_, _ = cn.graphClient.Execute(ctx, in)
	_ = cn.clientConn.Close()
	return nil
}

func (rs *resultSet) HasNext() bool {
	if rs.result == nil {
		return false
	}
	return rs.index < len(rs.result.Records)
}

func (rs *resultSet) Next() (Row, error) {
	if !rs.HasNext() {
		return nil, io.EOF
	}

	values := rs.result.GetRecords()[rs.index].Values
	row := &rowData{
		resultSet: rs,
		values:    make([]Value, 0, len(values)),
	}

	for _, v := range values {
		row.values = append(row.values, &grpcValue{data: v})
	}

	rs.index++
	return row, nil
}

func (rs *resultSet) RowSize() int {
	if rs.result == nil {
		return 0
	}
	return len(rs.result.GetRecords())
}

func (rs *resultSet) Latency() int64 {
	return int64(rs.latency)
}

func (rs *resultSet) Columns() []string {
	if rs.result == nil {
		return nil
	}
	names := rs.result.ColumnNames
	var cols []string
	for _, name := range names {
		cols = append(cols, string(name))
	}
	return cols
}

func (rs *resultSet) ColumnTypes() []ColumnType {
	return nil
}

func (rs *resultSet) PlanDesc() PlanDescer {
	if rs.planDesc == nil {
		return nil
	}
	var planDesc = make(map[string]interface{})
	if err := json.Unmarshal(rs.planDesc, &planDesc); err != nil {
		return nil
	}
	return &plan{
		planDesc: planDesc,
	}
}

func (rd *rowData) Values() []Value {
	return rd.values
}

func (rd *rowData) GetValueByName(name string) (Value, error) {
	names := rd.resultSet.Columns()
	var index int = -1
	for i, n := range names {
		if string(n) == name {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errInternel(fmt.Sprintf("column %s not found", name))
	}
	return rd.values[index], nil
}

func (rd *rowData) GetValueByIndex(index int) (Value, error) {
	if index < 0 || index >= len(rd.values) {
		return nil, errInternel(fmt.Sprintf("index out of range"))
	}
	return rd.values[index], nil
}
