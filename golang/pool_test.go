package nebula_ng

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

type dummyConnector struct{}

type dummyConn struct {
	id int
}

func (c *dummyConnector) connect(host *hostAddress, cfg *connConfig) (Conn, error) {
	return &dummyConn{}, nil
}

func (d *dummyConn) Execute(stmt string) (Result, error) {
	return nil, nil
}

func (d *dummyConn) ExecuteContext(ctx context.Context, stmt string) (Result, error) {
	return nil, nil
}
func (d *dummyConn) Ping() error {
	return nil
}
func (d *dummyConn) Close() error {
	return nil
}

func TestCleanIdle(t *testing.T) {
	testcases := []struct {
		maxIdle  int
		conns    []*driverConn
		expected int
	}{
		{5,
			[]*driverConn{
				{conn: &dummyConn{0}},
				{conn: &dummyConn{1}},
				{conn: &dummyConn{2}},
				{conn: &dummyConn{3}},
				{conn: &dummyConn{4}},
			},
			5,
		},
		{4,
			[]*driverConn{
				{conn: &dummyConn{0}},
				{conn: &dummyConn{1}},
				{conn: &dummyConn{2}},
				{conn: &dummyConn{3}},
				{conn: &dummyConn{4}},
			},
			4,
		},
		{2,
			[]*driverConn{
				{conn: &dummyConn{0}},
				{conn: &dummyConn{1}},
			},
			2,
		},
	}

	for _, tc := range testcases {
		connMap := make(map[Conn]struct{})
		for _, conn := range tc.conns {
			c := conn.conn
			connMap[c] = struct{}{}
		}
		pool := &driverPool{
			maxIdle:  tc.maxIdle,
			freeConn: tc.conns,
			connMap:  connMap,
		}
		pool.clearIdleConn()
		assert.Equal(t, tc.expected, len(pool.freeConn))
		assert.Equal(t, tc.expected, len(pool.connMap))
	}
}

func TestPool(t *testing.T) {
	var (
		maxOpen = 10
		minOpen = 5
		maxIdle = 6
	)
	pool := newDriverPool([]*hostAddress{{"127.0.0.1", 9669}}, nil)
	defer pool.Close()
	pool.connector = &dummyConnector{}
	pool.SetMaxOpenConnections(maxOpen)
	pool.SetMinOpenConnections(minOpen)
	pool.SetMaxIdleConnections(maxIdle)
	pool.tickerDuration = 100 * time.Millisecond
	pool.connCfg = &connConfig{
		requestTimeout: 10 * time.Second,
	}

	c, err := pool.GetClient()
	if err != nil {
		t.Fatal(err)
	}
	<-time.After(10 * time.Millisecond)
	if err := pool.PutClient(c); err != nil {
		t.Fatal(err)
	}
	//test min open
	assert.Equal(t, pool.minOpen, len(pool.connMap))
	assert.Equal(t, pool.minOpen, len(pool.freeConn))
	//test max open
	var wg sync.WaitGroup
	for i := 0; i < maxOpen; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := pool.GetClient()
			if err != nil {
				t.Log(err)
			}
			<-time.After(10 * time.Millisecond)
			err = pool.PutClient(c)
			if err != nil {
				t.Log(err)
			}
		}()
	}
	wg.Wait()

	pool.mu.Lock()
	assert.Equal(t, maxOpen, len(pool.connMap))
	assert.Equal(t, maxOpen, len(pool.freeConn))
	pool.mu.Unlock()
	//test idle conn
	<-time.After(200 * time.Millisecond)

	pool.mu.Lock()
	assert.Equal(t, maxIdle, len(pool.freeConn))
	assert.Equal(t, maxIdle, len(pool.connMap))
	pool.mu.Unlock()
}

func TestPoolPut(t *testing.T) {
	testcases := []struct {
		isClosed bool
		exceed   bool
		errMsg   string
	}{
		{false, false, ""},
		{true, false, "[99005]: connection to 127.0.0.1:9669 is closed"},
		{false, true, ""},
	}
	for _, tc := range testcases {
		pool := newDriverPool([]*hostAddress{{"127.0.0.1", 9669}}, nil)
		pool.connector = &dummyConnector{}
		pool.minOpen = 0
		pool.connCfg = &connConfig{
			requestTimeout: 10 * time.Second,
		}
		c, err := pool.GetClient()
		if err != nil {
			t.Fatal(err)
		}
		if tc.isClosed {
			c.Close()
		}
		if tc.exceed {
			pool.maxOpen = 0
		}
		err = pool.PutClient(c)
		if tc.errMsg != "" {
			assert.EqualError(t, err, tc.errMsg)
		} else {
			assert.NoError(t, err)
		}
		pool.mu.Lock()
		if tc.exceed || tc.isClosed {
			assert.Equal(t, 0, len(pool.freeConn))
		} else {
			assert.Equal(t, 1, len(pool.freeConn))
		}
		pool.mu.Unlock()
		pool.Close()
	}
}

func TestPoolGet(t *testing.T) {
	testcases := []struct {
		concurrency int
		maxOpen     int
		runTimes    int
	}{
		{10, 10, 10},
		{20, 10, 30},
		{10, 30, 20},
	}
	for _, tc := range testcases {
		pool := newDriverPool([]*hostAddress{{"127.0.0.1", 9669}}, nil)
		pool.connector = &dummyConnector{}
		pool.minOpen = 0
		pool.maxOpen = tc.maxOpen
		pool.connCfg = &connConfig{
			requestTimeout: 10 * time.Second,
		}
		var eg errgroup.Group
		for i := 0; i < tc.concurrency; i++ {
			eg.Go(func() error {
				for j := 0; j < tc.runTimes; j++ {
					c, err := pool.GetClient()
					if err != nil {
						return err
					}
					<-time.After(10 * time.Millisecond)
					if err := pool.PutClient(c); err != nil {
						return err
					}
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			t.Fatal(err)
		}
	}
}
