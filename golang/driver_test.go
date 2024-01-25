package nebula_ng

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type connetorWithErr struct {
	connectTime time.Duration
	executeTime time.Duration
	dummyConnector
	connectTimes int
	err          error
}

type connWithErr struct {
	dummyConn
	executeTime time.Duration
	err         error
}

func (c *connetorWithErr) connect(host *hostAddress, cfg *connConfig) (Conn, error) {
	c.connectTimes++
	return &connWithErr{executeTime: c.executeTime, err: c.err}, nil
}

func (d *connWithErr) ExecuteContext(ctx context.Context, stmt string) (Result, error) {
	time.Sleep(d.executeTime)
	return nil, d.err
}

func TestClientRetry(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	testcases := []struct {
		connectTimes int
		ctx          context.Context
		maxLifeTime  time.Duration
		executionErr error
		err          error
	}{
		{2, timeoutCtx, 150 * time.Millisecond, errConnBroken("", 0), timeoutCtx.Err()},
		{1, context.Background(), 15 * time.Millisecond, nil, nil},
		{3, context.Background(), 15 * time.Millisecond, errConnBroken("", 0), errConnBroken("", 0)},
		{1, context.Background(), 15 * time.Millisecond, fmt.Errorf("not connnection error"), fmt.Errorf("not connnection error")},
	}
	for _, tc := range testcases {
		client := newDriverConn([]*hostAddress{{"127.0.0.1", 9669}}, nil)
		connector := &connetorWithErr{
			connectTime: 3 * time.Millisecond,
			executeTime: 20 * time.Millisecond,
			err:         tc.executionErr,
		}
		client.connector = connector
		client.maxLifeTime = tc.maxLifeTime
		_, err := client.ExecuteContext(tc.ctx, "")
		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		}
		assert.Equal(t, tc.connectTimes, connector.connectTimes)
	}
}

func TestPoolRetry(t *testing.T) {
	// when using pool, user should put the client back to pool
	// means alway create a connection at last.
	// so if retry times is 2, pool would open connection 4 times if all connections are broken
	timeoutCtx, _ := context.WithTimeout(context.Background(), 25*time.Millisecond)
	testcases := []struct {
		connectTimes int
		ctx          context.Context
		maxLifeTime  time.Duration
		executionErr error
		err          error
	}{
		// timeout is 25ms, connectTime is 3ms, executeTime is 20ms
		{3, timeoutCtx, 150 * time.Millisecond, errConnBroken("", 0), timeoutCtx.Err()},
		{1, context.Background(), 15 * time.Millisecond, nil, nil},
		{4, context.Background(), 15 * time.Millisecond, errConnBroken("", 0), errConnBroken("", 0)},
		{1, context.Background(), 15 * time.Millisecond, fmt.Errorf("not connnection error"), fmt.Errorf("not connnection error")},
	}
	for _, tc := range testcases {
		pool := newDriverPool([]*hostAddress{{"127.0.0.1", 9669}}, nil)
		connector := &connetorWithErr{
			connectTime: 3 * time.Millisecond,
			executeTime: 20 * time.Millisecond,
			err:         tc.executionErr,
		}
		pool.connector = connector
		pool.connMaxLifeTime = tc.maxLifeTime
		pool.minOpen = 0
		pool.connCfg = &connConfig{
			requestTimeout: 10 * time.Second,
		}
		client, err := pool.GetClient()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(pool.connMap))
		assert.Equal(t, 0, len(pool.freeConn))
		_, err = client.ExecuteContext(tc.ctx, "")
		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		}
		assert.Equal(t, tc.connectTimes, connector.connectTimes)
	}
}

// after retry, the connection should be put back to pool
// verify the broken connection would be removed from pool
func TestPoolRetry2(t *testing.T) {
	pool := newDriverPool([]*hostAddress{{"127.0.0.1", 9669}}, nil)
	defer pool.Close()
	connector := &connetorWithErr{
		connectTime: 3 * time.Millisecond,
		executeTime: 20 * time.Millisecond,
		err:         errConnBroken("", 0),
	}
	pool.connector = connector
	pool.minOpen = 0
	pool.connCfg = &connConfig{
		requestTimeout: 10 * time.Second,
	}
	client, err := pool.GetClient()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = client.Execute("")
	if err := pool.PutClient(client); err != nil {
		t.Fatal(err)
	}
	pool.mu.Lock()
	assert.Equal(t, 1, len(pool.connMap))
	assert.Equal(t, 1, len(pool.freeConn))
	for _, conn := range pool.freeConn {
		c := conn.conn
		if _, ok := pool.connMap[c]; !ok {
			t.Fatal("conn not in connMap")
		}
	}
	pool.mu.Unlock()

}
