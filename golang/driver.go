package nebula_ng

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	internel_error "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/internal_error"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

type (
	// an internel interface to get the connection
	connector interface {
		connect(host *hostAddress, cfg *connConfig) (types.Client, error)
	}
	connConfig struct {
		username       string
		password       string
		graph          string
		requestTimeout time.Duration
		connectTimeout time.Duration
		timezone       string
		enableTLS      bool
		cert           string
		key            string
		ca             string
		peerName       string
		peerNameVerify bool
	}
	// driverConn is a client wrapper
	driverConn struct {
		cfg           *connConfig
		hostAddresses []*hostAddress
		currentHost   *hostAddress
		connector     connector
		pool          *driverPool
		createAt      time.Time
		conn          types.Client
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
	defaultTicker         = 5 * time.Second
	defaultMaxWait        = 1 * time.Minute
)

func NewNebulaClient(addresses, username, password string, opts ...clientOptionsFn) (types.Client, error) {
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

func NewNebulaPool(addresses, username, password string, opts ...poolOptionsFn) (types.Pool, error) {
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
		connMap:         make(map[types.Client]struct{}),
		requestConnChan: make(map[uint64]chan types.Client),
		requestCount:    0,
		openerCh:        make(chan struct{}, openConnChannelSize),
		connector:       defaultConnector,
		connMaxLifeTime: defaultMaxLieTime,
		maxOpen:         defaultMaxOpenConns,
		minOpen:         defaultMinOpenConns,
		maxIdle:         defaultMaxIdleConns,
		maxWait:         defaultMaxWait,
		sessionConfig:   &sessionConfig{},
		tickerDuration:  defaultTicker,
		log:             DefaultLogger,
	}
	for _, o := range opts {
		o(pool)
	}
	var (
		successed = 0
		dc        types.Client
	)
	for _, h := range pool.hostAddresses {
		if !pool.strictlyServerHealthy && successed > 0 {
			break
		}
		address := fmt.Sprintf("%s:%d", h.host, h.port)
		dc, err = pool.openNewConn(address)
		if err != nil {
			if pool.strictlyServerHealthy {
				break
			}
		} else {
			successed++
			pool.putNewConn(dc)
		}
	}
	if successed == 0 || (pool.strictlyServerHealthy && successed != len(pool.hostAddresses)) {
		return nil, err
	}

	go pool.connectionOpener(ctx)
	// start ticker
	go pool.ticker(pool.ctx)
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

func (dc *driverConn) Execute(stmt string) (types.Result, error) {
	return dc.ExecuteContext(context.Background(), stmt)
}

func (dc *driverConn) ExecuteContext(ctx context.Context, stmt string) (types.Result, error) {
	if dc.isClosed {
		return nil, internel_error.ErrConnIsClosed(dc.currentHost.host, dc.currentHost.port)
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
		return internel_error.ErrConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}
	return dc.conn.Ping()
}

func (dc *driverConn) PingContext(ctx context.Context) error {
	if dc.IsClosed() {
		return internel_error.ErrConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}
	return dc.conn.PingContext(ctx)
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
		return internel_error.ErrInternel("invalid connection type")
	}
	dc.conn = newConn.conn
	dc.currentHost = newConn.currentHost
	dc.createAt = newConn.createAt
	dc.isClosed = newConn.isClosed
	return nil
}

func (dc *driverConn) GetSessionId() (int64, error) {
	if dc.IsClosed() {
		return 0, internel_error.ErrConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}

	return dc.conn.GetSessionId()
}

func (dc *driverConn) IsClosed() bool {
	return dc.conn == nil || dc.isClosed
}
