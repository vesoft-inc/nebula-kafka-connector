package nebula_ng

import (
	"context"
	"math/rand"
	"time"
)

type (
	Result interface {
		HasNext() bool
		Next() (Row, error)
		Columns() []string
		ColumnTypes() []ColumnType //not support yet
		RowSize() int
		Summary() Summary
		Cursor() []byte
	}

	Row interface {
		Values() []Value
		GetValueByName(name string) (Value, error)
		GetValueByIndex(index int) (Value, error)
	}

	Summary interface {
		BuildTimeUs() int64
		OptimizeTimeUs() int64
		TotalServerTimeUs() int64
		ExplainType() string
		PlanInfo() PlanInfo
		QueryStats() QueryStats
	}
	PlanInfo interface {
		Id() string
		Name() string
		Details() string
		Columns() []string
		TimeMs() float64
		Rows() int64
		MemoryKib() float64
		BlockedMs() float64
		QueuedMs() float64
		ConsumeMs() float64
		ProduceMs() float64
		FinishMs() float64
		Batches() int64
		Concurrency() int64
		OtherStatsJson() []byte
		Children() []PlanInfo
	}
	QueryStats interface {
		NumAffectedNodes() int64
		NumAffectedEdges() int64
	}

	// an internel interface to get the connection
	connector interface {
		connect(host *hostAddress, cfg *connConfig) (Client, error)
	}

	Client interface {
		Execute(stmt string) (Result, error)
		ExecuteContext(ctx context.Context, stmt string) (Result, error)
		Ping() error
		IsClosed() bool
		Close() error
		GetSessionId() (int64, error)
	}

	Pool interface {
		GetClient() (Client, error)
		PutClient(Client) error
		Close() error
	}

	connConfig struct {
		username       string
		password       string
		graph          string
		requestTimeout time.Duration
		connectTimeout time.Duration
		timezone       string
	}

	ColumnType int

	// driverConn is a client wrapper
	driverConn struct {
		cfg           *connConfig
		hostAddresses []*hostAddress
		currentHost   *hostAddress
		connector     connector
		pool          *driverPool
		createAt      time.Time
		conn          Client
		isClosed      bool
		log           Logger
	}
)

const (
	defaultMaxOpenConns   = 100
	defaultMinOpenConns   = 1
	defaultMaxIdleConns   = 5
	defaultMaxLieTime     = 30 * time.Minute
	defaultRequestTimeout = 1 * time.Minute
	defaultConnetTimeout  = 3 * time.Second
	defaultPoolTicker     = 5 * time.Second
)

func NewNebulaClient(addresses, username, password string, opts ...clientOptionsFn) (Client, error) {
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
	for _, o := range opts {
		o(dc)
	}
	if err := dc.open(); err != nil {
		return nil, err
	}

	return dc, nil
}

func NewNebulaPool(addresses, username, password string, opts ...poolOptionsFn) (Pool, error) {
	hostAddresses, err := parseAddresses(addresses)
	if err != nil {
		return nil, err
	}
	connCfg := newConnConfig(username, password)
	ctx, cancel := context.WithCancel(context.Background())
	pool := &driverPool{
		ctx:             ctx,
		stop:            cancel,
		hostAddresses:   hostAddresses,
		connCfg:         connCfg,
		connMap:         make(map[Client]struct{}),
		requestConnChan: make(map[uint64]chan Client),
		requestCount:    0,
		openerCh:        make(chan struct{}, openConnChannelSize),
		connector:       defaultConnector,
		connMaxLifeTime: defaultMaxLieTime,
		maxOpen:         defaultMaxOpenConns,
		minOpen:         defaultMinOpenConns,
		maxIdle:         defaultMaxIdleConns,
		tickerDuration:  defaultPoolTicker,
		log:             DefaultLogger,
	}
	for _, o := range opts {
		o(pool)
	}
	go pool.connectionOpener(ctx)
	return pool, nil
}

func newConnConfig(username, password string) *connConfig {
	return &connConfig{
		username:       username,
		password:       password,
		requestTimeout: defaultRequestTimeout,
		connectTimeout: defaultConnetTimeout,
	}
}

func (dc *driverConn) Execute(stmt string) (Result, error) {
	return dc.ExecuteContext(context.Background(), stmt)
}

func (dc *driverConn) ExecuteContext(ctx context.Context, stmt string) (Result, error) {
	if dc.isClosed {
		return nil, errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	result, err := dc.conn.ExecuteContext(ctx, stmt)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (dc *driverConn) Ping() error {
	if dc.IsClosed() {
		return errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}
	return dc.conn.Ping()
}

func (dc *driverConn) Close() error {
	if dc.IsClosed() {
		return nil
	}
	err := dc.conn.Close()
	if err != nil {
		return err
	}
	dc.isClosed = true
	dc.conn = nil
	return nil
}

func (dc *driverConn) open() error {
	// random select one host
	if dc.currentHost == nil {
		rand.Seed(time.Now().UnixNano())
		hostIndex := rand.Intn(len(dc.hostAddresses))
		dc.currentHost = dc.hostAddresses[hostIndex]
	}

	conn, err := dc.connector.connect(dc.currentHost, dc.cfg)
	if err != nil {
		return err
	}
	dc.createAt = time.Now()
	dc.conn = conn
	return nil
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

func (dc *driverConn) GetSessionId() (int64, error) {
	if dc.IsClosed() {
		return 0, errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}

	return dc.conn.GetSessionId()
}

func (dc *driverConn) IsClosed() bool {
	return dc.conn == nil || dc.isClosed
}
