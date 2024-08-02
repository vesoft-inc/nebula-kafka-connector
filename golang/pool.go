package nebula_ng

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type (
	driverPool struct {
		ctx      context.Context
		mu       sync.Mutex
		once     sync.Once
		freeConn []Client
		// When driverConn.ExecuteContext is called, it will retry to get a connection.
		// The connMap key is the connection.
		connMap map[Client]struct{}
		// if the max open connection is full, would block
		// and wait for a free connection
		requestConnChan map[uint64]chan Client
		requestCount    uint64
		// open connection channel
		openerCh        chan struct{}
		connector       connector
		maxIdle         int
		maxOpen         int
		minOpen         int
		maxWait         time.Duration
		closed          atomic.Bool
		stop            func()
		hostAddresses   []*hostAddress
		connCfg         *connConfig
		hostIndex       int
		tickerDuration  time.Duration
		connMaxLifeTime time.Duration
		sessionConfig   *sessionConfig
		log             Logger
	}

	hostAddress struct {
		host string
		port int
	}

	sessionConfig struct {
		graph      string
		schema     string
		timezone   string
		parameters map[string]string
	}
)

const (
	openConnChannelSize = 1000000
)

var _ Client = &connection{}
var _ Result = &resultSet{}

func (dp *driverPool) Close() error {
	dp.closed.Store(true)
	dp.stop()
	dp.mu.Lock()
	for dc := range dp.connMap {
		// ignore connection close error
		_ = dc.Close()
	}
	dp.mu.Unlock()
	return nil
}

func (dp *driverPool) openNewConn(address string) (Client, error) {
	dc, err := NewNebulaClient(address, dp.connCfg.username, dp.connCfg.password,
		WithClientConnectTimeout(dp.connCfg.connectTimeout),
		WithClientRequestTimeout(dp.connCfg.requestTimeout),
		WithClientLogger(dp.log),
		withClientConnector(dp.connector),
	)
	if err != nil {
		return nil, err
	}
	// set session config
	var stmt string
	if dp.sessionConfig.graph != "" {
		stmt = fmt.Sprintf("SESSION SET GRAPH %s", dp.sessionConfig.graph)
		if _, err := dc.Execute(stmt); err != nil {
			_ = dc.Close()
			return nil, err
		}
	}
	if dp.sessionConfig.schema != "" {
		stmt = fmt.Sprintf("SESSION SET SCHEMA %s", dp.sessionConfig.schema)
		if _, err := dc.Execute(stmt); err != nil {
			_ = dc.Close()
			return nil, err
		}
	}
	if dp.sessionConfig.timezone != "" {
		stmt = fmt.Sprintf(`SESSION SET TIME ZONE "%s"`, dp.sessionConfig.timezone)
		if _, err := dc.Execute(stmt); err != nil {
			_ = dc.Close()
			return nil, err
		}
	}
	for k, v := range dp.sessionConfig.parameters {
		stmt = fmt.Sprintf("SESSION SET VALUE $%s=%s", k, v)
		if _, err := dc.Execute(stmt); err != nil {
			_ = dc.Close()
			return nil, err
		}
	}

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
	//TODO clean max life conn
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
		delete(dp.connMap, dp.freeConn[i])
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

func (dp *driverPool) GetClient() (Client, error) {
	var (
		dc Client
	)
	timeout, cancel := context.WithTimeout(context.Background(), dp.maxWait)
	defer cancel()

	dp.mu.Lock()
	if len(dp.freeConn) == 0 {
		if len(dp.connMap) < dp.maxOpen {
			dp.openerCh <- struct{}{}
		}

		req := make(chan Client, 1)
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
	// // start ticker
	dp.once.Do(func() {
		go dp.ticker(dp.ctx)
	})
	return dc, nil
}

func (dp *driverPool) PutClient(c Client) error {
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

func (dp *driverPool) putConnLocked(client Client) error {
	dc, ok := client.(*driverConn)
	if !ok {
		return errInternel("invalid client type")
	}
	// if client is closed by user, remove from pool,
	// and then raise an error.
	if dc.IsClosed() {
		delete(dp.connMap, client)
		return errConnIsClosed(dc.currentHost.host, dc.currentHost.port)
	}

	if dp.connMaxLifeTime > 0 && time.Since(dc.createAt) > dp.connMaxLifeTime {
		_ = client.Close()
		delete(dp.connMap, client)
		return nil
	}

	if len(dp.connMap) > dp.maxOpen {
		_ = client.Close()
		delete(dp.connMap, client)
		return nil
	}
	// if there's a conn request, do not put to freeConn
	if len(dp.requestConnChan) > 0 {
		var index uint64
		for i, ch := range dp.requestConnChan {
			ch <- client
			index = i
			break
		}
		delete(dp.requestConnChan, index)
	} else {
		dp.freeConn = append(dp.freeConn, client)
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
			hostAddress := dp.hostAddresses[dp.getHostIndexLocked()]
			address := fmt.Sprintf("%s:%d", hostAddress.host, hostAddress.port)
			dp.mu.Unlock()
			// address :=
			dc, err := dp.openNewConn(address)
			if err != nil {
				dp.log.Error(fmt.Sprintf("open connection failed, err: %s", err.Error()))
				// reset opener
				dp.openerCh <- struct{}{}
				continue
			}
			if dp.closed.Load() {
				_ = dc.Close()
				return
			}
			dp.mu.Lock()
			dp.connMap[dc] = struct{}{}
			_ = dp.putConnLocked(dc)
			dp.mu.Unlock()
		}
	}
}
