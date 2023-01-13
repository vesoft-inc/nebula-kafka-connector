//go:generate mockgen -source=client.go -destination client_mock.go -package client Client
package client

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

const (
	DefaultUser                     = "root"
	DefaultPassword                 = "nebula"
	DefaultReconnectInitialInterval = time.Second
	DefaultReconnectMaxInterval     = 2 * time.Minute
	DefaultRetry                    = 3
	DefaultRetryInitialInterval     = time.Second
	DefaultRetryMaxInterval         = 2 * time.Minute
	DefaultRetryRandomizationFactor = 0.1
	DefaultRetryMultiplier          = 1.5
	DefaultRetryMaxElapsedTime      = time.Hour
	DefaultConcurrencyPerAddress    = 10
	DefaultQueueSize                = 1000
)

type (
	Client interface {
		Open() error
		Execute(statement string) (ResultSet, error)
		ExecuteChan(statement string) (<-chan ExecuteResult, bool)
		Close() error
	}

	HostAddress struct {
		Host string
		Port int
	}

	defaultClient struct {
		addresses                []string
		user                     string
		password                 string
		reconnectInitialInterval time.Duration
		retry                    int
		retryInitialInterval     time.Duration
		concurrencyPerAddress    int
		queueSize                int
		chExecuteDataQueue       chan executeData
		lock                     sync.RWMutex
		closed                   bool
		done                     chan struct{}
		wgSession                sync.WaitGroup
		wgStatementExecute       sync.WaitGroup
		hostAddresses            []HostAddress
		logger                   logger.Logger
		fnNewSession             NewSessionFunc
	}

	NewSessionFunc func(HostAddress) Session

	Option func(*defaultClient)

	executeData struct {
		statement string
		ch        chan<- ExecuteResult
	}

	ExecuteResult struct {
		ResultSet ResultSet
		Err       error
	}
)

func New(opts ...Option) Client {
	c := &defaultClient{
		user:                     DefaultUser,
		password:                 DefaultPassword,
		reconnectInitialInterval: DefaultReconnectInitialInterval,
		retry:                    DefaultRetry,
		retryInitialInterval:     DefaultRetryInitialInterval,
		concurrencyPerAddress:    DefaultConcurrencyPerAddress,
		queueSize:                DefaultQueueSize,
		done:                     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.chExecuteDataQueue = make(chan executeData, c.queueSize)

	if c.logger == nil {
		c.logger = logger.NopLogger
	}

	if c.fnNewSession == nil {
		WithV5()(c)
	}

	return c
}

func WithAddress(addresses ...string) Option {
	return func(c *defaultClient) {
		for _, addr := range addresses {
			if strings.IndexByte(addr, ',') != -1 {
				c.addresses = append(c.addresses, strings.Split(addr, ",")...)
			} else {
				c.addresses = append(c.addresses, addr)
			}
		}
	}
}

func WithV3() Option {
	return func(c *defaultClient) {
		WithNewSessionFunc(func(hostAddress HostAddress) Session {
			return newSessionV3(hostAddress, c.user, c.password, c.logger)
		})(c)
	}
}

func WithV5() Option {
	return func(c *defaultClient) {
		WithNewSessionFunc(func(hostAddress HostAddress) Session {
			return newSessionV5(hostAddress, c.user, c.password, c.logger)
		})(c)
	}
}

func WithUser(user string) Option {
	return func(c *defaultClient) {
		c.user = user
	}
}

func WithPassword(password string) Option {
	return func(c *defaultClient) {
		c.password = password
	}
}

func WithUserPassword(user, password string) Option {
	return func(c *defaultClient) {
		WithUser(user)(c)
		WithPassword(password)(c)
	}
}

func WithReconnectInitialInterval(interval time.Duration) Option {
	return func(c *defaultClient) {
		if interval > 0 {
			c.reconnectInitialInterval = interval
		}
	}
}

func WithRetry(retry int) Option {
	return func(c *defaultClient) {
		if retry > 0 {
			c.retry = retry
		}
	}
}

func WithRetryInitialInterval(interval time.Duration) Option {
	return func(c *defaultClient) {
		if interval > 0 {
			c.retryInitialInterval = interval
		}
	}
}

func WithConcurrencyPerAddress(concurrencyPerAddress int) Option {
	return func(c *defaultClient) {
		if concurrencyPerAddress > 0 {
			c.concurrencyPerAddress = concurrencyPerAddress
		}
	}
}

func WithQueueSize(queueSize int) Option {
	return func(c *defaultClient) {
		if queueSize > 0 {
			c.queueSize = queueSize
		}
	}
}

func WithLogger(l logger.Logger) Option {
	return func(m *defaultClient) {
		m.logger = l
	}
}

func WithNewSessionFunc(fn NewSessionFunc) Option {
	return func(m *defaultClient) {
		m.fnNewSession = fn
	}
}

func (c *defaultClient) Open() error {
	for _, addr := range c.addresses {
		hostPort := strings.Split(addr, ":")
		if len(hostPort) != 2 {
			return errors.ErrInvalidAddress
		}
		if hostPort[0] == "" {
			return errors.ErrInvalidAddress
		}
		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return errors.ErrInvalidAddress
		}
		hostAddress := HostAddress{Host: hostPort[0], Port: port}
		session, err := c.openSession(hostAddress)
		if err != nil {
			return err
		}
		_ = session.Close()
		c.hostAddresses = append(c.hostAddresses, hostAddress)
	}

	if len(c.hostAddresses) == 0 {
		return errors.ErrNoAddresses
	}

	c.startWorkers()

	return nil
}

func (c *defaultClient) Execute(statement string) (ResultSet, error) {
	if c.IsClosed() {
		return nil, ErrClosed
	}
	c.wgStatementExecute.Add(1)
	defer c.wgStatementExecute.Done()

	ch := make(chan ExecuteResult, 1)
	data := executeData{
		statement: statement,
		ch:        ch,
	}
	c.chExecuteDataQueue <- data
	result := <-ch
	return result.ResultSet, result.Err
}

func (c *defaultClient) ExecuteChan(statement string) (<-chan ExecuteResult, bool) {
	if c.IsClosed() {
		return nil, false
	}
	c.wgStatementExecute.Add(1)
	defer c.wgStatementExecute.Done()

	ch := make(chan ExecuteResult, 1)
	data := executeData{
		statement: statement,
		ch:        ch,
	}
	select {
	case c.chExecuteDataQueue <- data:
		return ch, true
	default:
		return nil, false
	}
}

func (c *defaultClient) Close() error {
	c.lock.Lock()
	c.closed = true
	c.lock.Unlock()

	c.wgStatementExecute.Wait()
	close(c.done)
	c.wgSession.Wait()
	close(c.chExecuteDataQueue)
	return nil
}

func (c *defaultClient) IsClosed() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.closed
}

func (c *defaultClient) startWorkers() {
	for _, hostAddress := range c.hostAddresses {
		hostAddress := hostAddress
		for i := 0; i < c.concurrencyPerAddress; i++ {
			c.wgSession.Add(1)
			go func() {
				defer c.wgSession.Done()
				c.worker(hostAddress)
			}()
		}
	}
}

func (c *defaultClient) worker(hostAddress HostAddress) {
	for {
		select {
		case <-c.done:
			return
		default:
			exp := backoff.NewExponentialBackOff()
			exp.InitialInterval = c.reconnectInitialInterval
			exp.MaxInterval = DefaultReconnectMaxInterval
			exp.RandomizationFactor = DefaultRetryRandomizationFactor
			exp.Multiplier = DefaultRetryMultiplier

			var session Session
			_ = backoff.Retry(func() error {
				var err error
				session, err = c.openSession(hostAddress)
				if err != nil {
					c.logger.WithError(err).Error("open session failed")
				}
				return err
			}, exp)
			c.loopSession(session)
		}
	}
}

func (c *defaultClient) openSession(hostAddress HostAddress) (Session, error) {
	session := c.fnNewSession(hostAddress)
	if err := session.Open(); err != nil {
		return nil, err
	}
	return session, nil
}

func (c *defaultClient) loopSession(session Session) {
	defer func() {
		_ = session.Close()
	}()
	for {
		select {
		case data, ok := <-c.chExecuteDataQueue:
			if !ok {
				continue
			}

			exp := backoff.NewExponentialBackOff()
			exp.InitialInterval = c.retryInitialInterval
			exp.MaxInterval = DefaultRetryMaxInterval
			exp.MaxElapsedTime = DefaultRetryMaxElapsedTime
			exp.Multiplier = DefaultRetryMultiplier
			exp.RandomizationFactor = DefaultRetryRandomizationFactor

			var (
				err   error
				rs    ResultSet
				retry = c.retry
			)

			// There are three cases of retry
			// * Case 1: retry no more
			// * Case 2. retry as much as possible
			// * Case 3: retry with limit times
			_ = backoff.Retry(func() error {
				rs, err = session.Execute(data.statement)
				if err == nil && rs.IsSucceed() {
					return nil
				}
				retryErr := err
				if rs != nil {
					retryErr = rs.GetError()

					// Case 1: retry no more
					if rs.IsPermanentError() {
						// stop the retry
						return backoff.Permanent(retryErr)
					}

					// Case 2. retry as much as possible
					if rs.IsRetryMoreError() {
						retry = c.retry
						return retryErr
					}
				}

				// Case 3: retry with limit times
				if retry <= 0 {
					// stop the retry
					return backoff.Permanent(retryErr)
				}
				retry--
				return retryErr
			}, exp)
			if err != nil {
				c.logger.WithError(err).Error("execute statement failed")
			}

			data.ch <- ExecuteResult{
				ResultSet: rs,
				Err:       err,
			}
		case <-c.done:
			return
		}
	}
}
