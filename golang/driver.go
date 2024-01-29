package nebula_ng

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type (
	Result interface {
		HasNext() bool
		Next() (Row, error)
		Latency() int64 // in us
		Columns() []string
		ColumnTypes() []ColumnType //not support yet
		RowSize() int
		PlanDesc() PlanDescer
	}

	Row interface {
		Values() []Value
		GetValueByName(name string) (Value, error)
		GetValueByIndex(index int) (Value, error)
	}

	// an internel interface to get the connection
	connector interface {
		connect(host *hostAddress, cfg *connConfig) (Conn, error)
	}

	Client interface {
		Conn
		ConnSetter
		SetLogger(logger Logger)
	}

	ConnSetter interface {
		SetRequestTimeout(timeout time.Duration)
		SetConnectTimeout(timeout time.Duration)
		SetMaxLifeTime(timeout time.Duration)
		SetRetryTimes(times int)
		SetGraph(graph string)
	}

	Conn interface {
		Execute(stmt string) (Result, error)
		ExecuteContext(ctx context.Context, stmt string) (Result, error)
		Ping() error
		Close() error
	}

	Pool interface {
		GetClient() (Conn, error)
		SetMinOpenConnections(min int)
		SetMaxOpenConnections(max int)
		SetMaxIdleConnections(max int)
		SetLogger(logger Logger)
		PutClient(Conn) error
		OnOpenClient(func(ConnSetter))
		Close() error
	}

	PlanDescer interface {
		GetHeader() []string
		GetPlanPrintFormat() string
		MakePlanByRow() (rightSepToTailWidth []int, rows [][]interface{})
		GetBuildTimeInUs() int64
		GetOptimizeTimeInUs() int64
	}

	connConfig struct {
		username       string
		password       string
		graph          string
		requestTimeout time.Duration
		connectTimeout time.Duration
	}

	ColumnType int

	// driverConn is a client wrapper
	// support reconnect, timeout, etc.
	driverConn struct {
		cfg           *connConfig
		hostAddresses []*hostAddress
		currentHost   *hostAddress
		connector     connector
		pool          *driverPool
		createAt      time.Time
		maxLifeTime   time.Duration
		retryTimes    int
		conn          Conn
		isClosed      bool
		log           Logger
	}
)

const (
	defaultMaxOpenConns   = 100
	defaultMinOpenConns   = 1
	defaultMaxIdleConns   = 5
	defaultMaxLieTime     = 30 * time.Minute
	defaultRetryTimes     = 2
	defaultRequestTimeout = 1 * time.Minute
	defaultConnetTimeout  = 3 * time.Second
	defaultPoolTicker     = 5 * time.Second
)

func NewNebulaClient(addresses, username, password string) (Client, error) {
	hostAddresses, err := parseAddresses(addresses)
	if err != nil {
		return nil, err
	}
	cfg := newConnConfig(username, password)
	dc := &driverConn{
		hostAddresses: hostAddresses,
		connector:     defaultConnector,
		cfg:           cfg,
	}

	return dc, nil
}

func NewNebulaPool(addresses, username, password string) (Pool, error) {
	hostAddresses, err := parseAddresses(addresses)
	if err != nil {
		return nil, err
	}
	connCfg := newConnConfig(username, password)

	return newDriverPool(hostAddresses, connCfg), nil
}

func newConnConfig(username, password string) *connConfig {
	return &connConfig{
		username:       username,
		password:       password,
		requestTimeout: defaultRequestTimeout,
		connectTimeout: defaultConnetTimeout,
	}
}

func newDriverConn(hostAddresses []*hostAddress, cfg *connConfig) *driverConn {
	return &driverConn{
		hostAddresses: hostAddresses,
		cfg:           cfg,
		connector:     defaultConnector,
		maxLifeTime:   defaultMaxLieTime,
		retryTimes:    defaultRetryTimes,
		createAt:      time.Now(),
		log:           DefaultLogger,
	}
}

func (dc *driverConn) Execute(stmt string) (Result, error) {
	return dc.ExecuteContext(context.Background(), stmt)
}

func (dc *driverConn) ExecuteContext(ctx context.Context, stmt string) (Result, error) {
	if dc.isClosed {
		return nil, errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}
	return dc.retryExecuteLocked(ctx, stmt)
}

func (dc *driverConn) retryExecuteLocked(ctx context.Context, stmt string) (Result, error) {
	var (
		result Result
		err    error
	)

	for i := 0; i < dc.retryTimes+1; i++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if dc.pool == nil && dc.conn == nil {
			conn, err := dc.newDirect()
			if err != nil {
				// TODO warning
				continue
			}
			dc.conn = conn
		}

		if dc.maxLifeTime > 0 && time.Since(dc.createAt) > dc.maxLifeTime {
			if dc.pool == nil {
				_ = dc.Close()
				dc.conn = nil
			} else {
				_ = dc.Close()
				_ = dc.replaceFromPool()
			}
			continue
		}

		result, err = dc.conn.ExecuteContext(ctx, stmt)
		if err == nil {
			return result, nil
		}
		if !isConnectionError(err) {
			break
		}
		if dc.pool == nil {
			_ = dc.conn.Close()
			dc.conn = nil
		} else {
			_ = dc.Close()
			_ = dc.replaceFromPool()
		}
	}
	return nil, err
}

func (dc *driverConn) Ping() error {
	if dc.pool == nil && dc.conn == nil {
		conn, err := dc.newDirect()
		if err != nil {
			return err
		}
		dc.conn = conn
	}
	return dc.conn.Ping()
}

func (dc *driverConn) Close() error {
	if dc.isClosed {
		return nil
	}
	if dc.conn == nil {
		return nil
	}
	err := dc.conn.Close()
	if err != nil {
		return err
	}
	dc.isClosed = true
	return nil
}

func (dc *driverConn) SetLogger(logger Logger) {
	dc.log = logger
}

func (dc *driverConn) SetRequestTimeout(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dc.cfg.requestTimeout = timeout
}

func (dc *driverConn) SetConnectTimeout(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dc.cfg.connectTimeout = timeout
}

func (dc *driverConn) SetMaxLifeTime(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dc.maxLifeTime = timeout
}

func (dc *driverConn) SetGraph(graph string) {
	dc.cfg.graph = graph
}

func (dc *driverConn) SetRetryTimes(times int) {
	dc.retryTimes = times
}

func (dc *driverConn) newDirect() (Conn, error) {
	// random select one host
	rand.Seed(time.Now().UnixNano())
	hostIndex := rand.Intn(len(dc.hostAddresses))
	dc.currentHost = dc.hostAddresses[hostIndex]
	conn, err := dc.connector.connect(dc.currentHost, dc.cfg)
	if err != nil {
		return nil, err
	}
	dc.createAt = time.Now()
	return conn, nil
}

func (dc *driverConn) replaceFromPool() error {
	if dc.pool == nil {
		return nil
	}
	// return the old connection
	// pool would delete it if invalid.
	_ = dc.pool.PutClient(dc)
	conn, err := dc.pool.GetClient()
	if err != nil {
		return err
	}
	newConn, ok := conn.(*driverConn)
	if !ok {
		// not reachable
		return errInternel("invalid connection type")
	}
	dc.conn = newConn.conn
	dc.currentHost = newConn.currentHost
	dc.createAt = newConn.createAt
	dc.isClosed = newConn.isClosed
	return nil
}
