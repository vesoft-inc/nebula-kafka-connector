package nebula_ng

import (
	"math"
	"time"
)

type clientOptionsFn func(*driverConn)
type poolOptionsFn func(*driverPool)

func WithClientConnectTimeout(timeout time.Duration) clientOptionsFn {
	return func(ops *driverConn) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.cfg.connectTimeout = timeout
	}
}

func WithClientRequestTimeout(timeout time.Duration) clientOptionsFn {
	return func(ops *driverConn) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.cfg.requestTimeout = timeout
	}
}

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

func WithPoolConnectTime(timeout time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.connCfg.connectTimeout = timeout
	}
}

func WithPoolRequestTimeout(timeout time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if timeout <= 0 {
			timeout = math.MaxInt64
		}
		ops.connCfg.requestTimeout = timeout
	}
}

func WithPoolMinOpenConns(minOpenConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.minOpen = minOpenConns
	}
}

func WithPoolMaxOpenConns(maxOpenConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.maxOpen = maxOpenConns
	}
}

func WithPoolMaxIdleConns(maxIdleConns int) poolOptionsFn {
	return func(ops *driverPool) {
		ops.maxIdle = maxIdleConns
	}
}

func WithPoolTickerDuration(ticker time.Duration) poolOptionsFn {
	return func(ops *driverPool) {
		if ticker <= 0 {
			ticker = math.MaxInt64
		}
		ops.tickerDuration = ticker
	}
}

func WithPoolLogger(logger Logger) poolOptionsFn {
	return func(ops *driverPool) {
		ops.log = logger
	}
}

// used for testing
func withPoolConnector(connector connector) poolOptionsFn {
	return func(ops *driverPool) {
		ops.connector = connector
	}
}
