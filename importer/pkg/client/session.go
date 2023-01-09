//go:generate mockgen -source=session.go -destination session_mock.go -package client Session
package client

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

type (
	Session interface {
		Open() error
		Execute(statement string) (ResultSet, error)
		Close() error
	}

	defaultSession struct {
		session     *nebula.Session
		hostAddress nebula.HostAddress
		user        string
		password    string
		logger      logger.Logger
	}
)

func newSession(hostAddress nebula.HostAddress, user, password string, l logger.Logger) Session {
	if l == nil {
		l = logger.NopLogger
	}
	return &defaultSession{
		hostAddress: hostAddress,
		user:        user,
		password:    password,
		logger:      l,
	}
}

func (s *defaultSession) Open() error {
	hostAddress := s.hostAddress
	connection := nebula.NewConnection(hostAddress)
	if err := connection.Open(hostAddress, 1000*time.Millisecond, nil); err != nil {
		return err
	}
	authResp, err := connection.Authenticate(s.user, s.password)
	if err != nil {
		return err
	}

	session := nebula.NewSession(
		*authResp.Identifier,
		connection,
		newNebulaLogger(s.logger.With(logger.Field{
			Key:   "address",
			Value: fmt.Sprintf("%s:%d", hostAddress.Host, hostAddress.Port),
		})),
	)

	s.session = session

	return nil
}

func (s *defaultSession) Execute(statement string) (ResultSet, error) {
	rs, err := s.session.Execute(statement)
	if err != nil {
		return nil, err
	}
	return newResultSet(rs), nil
}

func (s *defaultSession) Close() error {
	s.session.Release()
	return nil
}
