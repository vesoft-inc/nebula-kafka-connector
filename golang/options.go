package nebula_ng

import (
	"math"
	"time"
)

type clientOptionsFn func(*driverConn)
type poolOptionsFn func(*driverPool)

// WithClientConnectTimeout sets the timeout for connecting to the server
func WithClientConnectTimeout(timeout time.Duration) clientOptionsFn {
	return func(ops *driverConn) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.cfg.connectTimeout = timeout
	}
}

// WithClientRequestTimeout sets the timeout for a request to the server
func WithClientRequestTimeout(timeout time.Duration) clientOptionsFn {
	return func(ops *driverConn) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.cfg.requestTimeout = timeout
	}
}

// WithClientLogger sets the logger for the client
func WithClientLogger(logger Logger) clientOptionsFn {
	return func(ops *driverConn) {
		ops.log = logger
	}
}

// used for testing
func withClientConnector(connector connector) clientOptionsFn {
	return func(ops *driverConn) {
		ops.connector = connector
	}
}

// Experimental: TLS options
func WithClientTLS(enable bool, ca, cert, key string, peerNameVerify bool, peerName string) clientOptionsFn {
	return func(conn *driverConn) {
		conn.cfg.enableTLS = enable
		conn.cfg.ca = ca
		conn.cfg.cert = cert
		conn.cfg.key = key
		conn.cfg.peerNameVerify = peerNameVerify
		conn.cfg.peerName = peerName
	}
}

// WithPoolConnectTimeout sets the timeout for connecting to the server
func WithPoolConnectTime(timeout time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.connCfg.connectTimeout = timeout
	}
}

// WithPoolRequestTimeout sets the timeout for a request to the server
func WithPoolRequestTimeout(timeout time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.connCfg.requestTimeout = timeout
	}
}

// WithPoolMaxWait sets max wait time
// that the pool will wait (when there are no available connections)
// for a connection to be returned
func WithPoolMaxWait(maxWait time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if maxWait <= 0 {
			maxWait = math.MaxInt64
		}
		ops.maxWait = maxWait
	}
}

// WithPoolMaxLifetime sets the maximum amount of time a connection may be reused
func WithPoolMaxLifetime(maxLifeTime time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if maxLifeTime <= 0 {
			maxLifeTime = math.MaxInt64
		}
		ops.connMaxLifeTime = maxLifeTime
	}
}

// WithPoolMinOpenConns sets the minimum number of open connections
func WithPoolMinOpenConns(minOpenConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.minOpen = minOpenConns
	}
}

// WithPoolMaxOpenConns sets the maximum number of open connections
func WithPoolMaxOpenConns(maxOpenConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.maxOpen = maxOpenConns
	}
}

// WithPoolMaxIdleConns sets the maximum number of idle connections
// would evicte the idle connections until the number of idle connections
func WithPoolMaxIdleConns(maxIdleConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.maxIdle = maxIdleConns
	}
}

// WithPoolTickerDuration sets the duration of the ticker
func WithPoolTickerDuration(ticker time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if ticker <= 0 {
			ticker = math.MaxInt64
		}
		ops.tickerDuration = ticker
	}
}

// WithPoolLogger sets the logger for the pool
func WithPoolLogger(logger Logger) poolOptionsFn {
	return func(ops *driverPool) {
		ops.log = logger
	}
}

// WithPoolGraph sets the graph for the session
func WithPoolGraph(graph string) poolOptionsFn {
	return func(ops *driverPool) {
		ops.sessionConfig.graph = graph
	}
}

// WithPoolSchema sets the schema for the session
func WithPoolSchema(schema string) poolOptionsFn {
	return func(ops *driverPool) {
		ops.sessionConfig.schema = schema
	}
}

// WithPoolTimezone sets the timezone for the session
func WithPoolTimezone(timezone string) poolOptionsFn {
	return func(ops *driverPool) {
		ops.sessionConfig.timezone = timezone
	}
}

// WithPoolParameters sets the parameters for the session
func WithPoolParameters(parameters map[string]string) poolOptionsFn {
	return func(ops *driverPool) {
		ops.sessionConfig.parameters = parameters
	}
}

// WithPoolStrictlyServerHealthy sets the pool to strictly check the server health
// if strictly is false, the pool will be created successfully if any of the servers is healthy
// if strictly is true, the pool will be created successfully only if all of the servers are healthy
// default is false
func WithPoolStrictlyServerHealthy(strictly bool) poolOptionsFn {
	return func(ops *driverPool) {
		ops.strictlyServerHealthy = strictly
	}
}

// used for testing
func withPoolConnector(connector connector) poolOptionsFn {
	return func(ops *driverPool) {
		ops.connector = connector
	}
}

// TLS options
func WithPoolTLS(enable bool, ca, cert, key string, peerNameVerify bool, peerName string) poolOptionsFn {
	return func(pool *driverPool) {
		pool.connCfg.enableTLS = enable
		pool.connCfg.ca = ca
		pool.connCfg.cert = cert
		pool.connCfg.key = key
		pool.connCfg.peerNameVerify = peerNameVerify
		pool.connCfg.peerName = peerName
	}
}

// WithClientAuthInfo sets the authInfo for the client connection
func WithClientAuthInfo(authInfo map[string]string) clientOptionsFn {
	return func(ops *driverConn) {
		ops.cfg.authInfo = authInfo
	}
}

// WithPoolAuthInfo sets the authInfo for the pool connection
func WithPoolAuthInfo(authInfo map[string]string) poolOptionsFn {
	return func(ops *driverPool) {
		ops.connCfg.authInfo = authInfo
	}
}
