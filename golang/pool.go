package nebula_ng

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type (
	driverPool struct {
		ctx      context.Context
		mu       sync.Mutex
		freeConn []*driverConn
		// When driverConn.ExecuteContext is called, it will retry to get a connection.
		// The connMap key is the connection.
		connMap map[Conn]struct{}
		// if the max open connection is full, would block
		// and wait for a free connection
		requestConnChan map[uint64]chan Conn
		requestCount    uint64
		// open connection channel
		openerCh        chan struct{}
		connector       connector
		maxIdle         int
		maxOpen         int
		minOpen         int
		closed          bool
		stop            func()
		hostAddresses   []*hostAddress
		connCfg         *connConfig
		hostIndex       int
		tickerDuration  time.Duration
		connRetryTimes  int
		connMaxLifeTime time.Duration
		log             Logger
	}

	hostAddress struct {
		host string
		port int
	}
)

const (
	requestConnChannelSize = 10000
	openConnChannelSize    = 10000
)

var _ Conn = &connection{}
var _ Result = &resultSet{}

func newDriverPool(hostAddresses []*hostAddress, cfg *connConfig) *driverPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &driverPool{
		ctx:             ctx,
		stop:            cancel,
		hostAddresses:   hostAddresses,
		connCfg:         cfg,
		connMap:         make(map[Conn]struct{}),
		requestConnChan: make(map[uint64]chan Conn),
		requestCount:    0,
		openerCh:        make(chan struct{}, openConnChannelSize),
		connector:       defaultConnector,
		connMaxLifeTime: defaultMaxLieTime,
		connRetryTimes:  defaultRetryTimes,
		maxOpen:         defaultMaxOpenConns,
		minOpen:         defaultMinOpenConns,
		maxIdle:         defaultMaxIdleConns,
		tickerDuration:  defaultPoolTicker,
		log:             DefaultLogger,
	}
	go pool.connectionOpener(ctx)
	return pool
}

func (dp *driverPool) SetLogger(logger Logger) {
	dp.log = logger
}

func (dp *driverPool) SetMinOpenConnections(min int) {
	if min > dp.maxOpen {
		min = dp.maxOpen
	}
	dp.minOpen = min
}

func (dp *driverPool) SetMaxOpenConnections(max int) {
	dp.maxOpen = max
}

func (dp *driverPool) SetMaxIdleConnections(max int) {
	dp.maxIdle = max
}

func (dp *driverPool) OnOpenClient(fn func(ConnSetter)) {
	fn(dp)
}

func (dp *driverPool) SetRequestTimeout(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dp.connCfg.requestTimeout = timeout
}

func (dp *driverPool) SetConnectTimeout(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dp.connCfg.connectTimeout = timeout
}

func (dp *driverPool) SetMaxLifeTime(timeout time.Duration) {
	if timeout <= 0 {
		timeout = math.MaxInt64
	}
	dp.connMaxLifeTime = timeout
}

func (dp *driverPool) SetGraph(graph string) {
	dp.connCfg.graph = graph
}

func (dp *driverPool) SetRetryTimes(times int) {
	dp.connRetryTimes = times
}

func (dp *driverPool) SetTimeZone(timezone string) {
	dp.connCfg.timezone = timezone
}

func (dp *driverPool) Close() error {
	dp.mu.Lock()
	for dc := range dp.connMap {
		// ignore connection close error
		_ = dc.Close()
	}
	dp.stop()
	dp.closed = true
	dp.mu.Unlock()
	return nil
}

func (dp *driverPool) openNewConnLocked() (*driverConn, error) {
	hostAddress := dp.hostAddresses[dp.getHostIndexLocked()]
	conn, err := dp.connector.connect(hostAddress, dp.connCfg)
	if err != nil {
		return nil, err
	}

	dc := newDriverConn(nil, dp.connCfg)
	dc.pool = dp
	dc.currentHost = hostAddress
	dc.conn = conn
	dc.maxLifeTime = dp.connMaxLifeTime
	dc.retryTimes = dp.connRetryTimes
	dc.createAt = time.Now()
	// init session context
	if dc.cfg.timezone != "" {
		_, err = conn.Execute(fmt.Sprintf(`SESSION SET TIME ZONE "%s"`, dc.cfg.timezone))
		if err != nil {
			return nil, err
		}
	}

	dp.connMap[conn] = struct{}{}
	return dc, nil
}

func (dp *driverPool) getHostIndexLocked() int {
	if len(dp.hostAddresses) == 1 {
		return 0
	}
	dp.hostIndex++
	if dp.hostIndex >= len(dp.hostAddresses) {
		dp.hostIndex = 0
	}
	return dp.hostIndex
}

// ticker is used to check the connection status
// maxIdleTime
// maxIdle connections
// minOpen connections
func (dp *driverPool) ticker(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	// immediately run the first time
	dp.clearIdleConn()
	dp.openMinConn()

	ticker := time.NewTicker(dp.tickerDuration)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dp.clearIdleConn()
			dp.openMinConn()
		}
	}
}

// alway put a connection to freeConn at the end,
// so the first connection will be max idle time
func (dp *driverPool) clearIdleConn() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	total := len(dp.freeConn)
	if total == 0 {
		return
	}
	var index int

	if total > dp.maxIdle {
		index = total - dp.maxIdle
	}
	for i := 0; i < index; i++ {
		_ = dp.freeConn[i].Close()
		delete(dp.connMap, dp.freeConn[i].conn)
	}
	dp.freeConn = dp.freeConn[index:]
}

func (dp *driverPool) openMinConn() {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	needOpen := dp.minOpen - len(dp.connMap)
	if needOpen <= 0 {
		return
	}

	for i := 0; i < needOpen; i++ {
		// it may cost a lot of time to open a new connection
		// maybe we can use a channel to open connection in parallel
		dp.openerCh <- struct{}{}
	}
}

func (dp *driverPool) GetClient() (Conn, error) {
	var (
		dc *driverConn
	)
	timeout, cancel := context.WithTimeout(context.Background(), dp.connCfg.requestTimeout)
	defer cancel()
	dp.mu.Lock()
	if len(dp.freeConn) == 0 {
		if len(dp.connMap) < dp.maxOpen {
			dp.openerCh <- struct{}{}
		}

		req := make(chan Conn, 1)
		if dp.requestCount == math.MaxUint64 {
			dp.requestCount = 0
		} else {
			dp.requestCount++
		}
		dp.requestConnChan[dp.requestCount] = req
		// should unlock and wait for the new connection
		dp.mu.Unlock()
		select {
		case <-timeout.Done():
			return nil, errInternel("cannot get the valid connection")
		case conn := <-req:
			dc = conn.(*driverConn)
		}
	} else {
		dc = dp.freeConn[len(dp.freeConn)-1]
		dp.freeConn = dp.freeConn[:len(dp.freeConn)-1]
		dp.mu.Unlock()
	}

	// start ticker
	var once sync.Once
	once.Do(func() {
		go dp.ticker(dp.ctx)
	})

	return dc, nil
}

func (dp *driverPool) PutClient(c Conn) error {
	if c == nil {
		return errInternel("connection is nil")
	}

	dc, ok := c.(*driverConn)
	if !ok {
		// never happen from nebula client
		return errInternel("invalid client type")
	}

	dp.mu.Lock()
	defer dp.mu.Unlock()
	return dp.putConnLocked(dc)

}

func (dp *driverPool) putConnLocked(dc *driverConn) error {
	if dc.isClosed {
		delete(dp.connMap, dc.conn)
		return errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}

	if dc.maxLifeTime > 0 && time.Since(dc.createAt) > dc.maxLifeTime {
		_ = dc.Close()
		delete(dp.connMap, dc.conn)
		return nil
	}

	if len(dp.connMap) > dp.maxOpen {
		_ = dc.Close()
		delete(dp.connMap, dc.conn)
		return nil
	}
	// if there's a conn request, do not put to freeConn
	if len(dp.requestConnChan) > 0 {
		var index uint64
		for i, ch := range dp.requestConnChan {
			ch <- dc
			index = i
			break
		}
		delete(dp.requestConnChan, index)
	} else {
		dp.freeConn = append(dp.freeConn, dc)
	}
	return nil
}

func (dp *driverPool) connectionOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-dp.openerCh:
			dp.mu.Lock()
			dc, err := dp.openNewConnLocked()
			if err != nil {
				// reset opener
				dp.openerCh <- struct{}{}
				dp.mu.Unlock()
				continue
			}
			_ = dp.putConnLocked(dc)
			dp.mu.Unlock()
		}
	}
}
