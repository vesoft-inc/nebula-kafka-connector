//go:generate mockgen -source=session.go -destination session_mock.go -package client NebulaSession
package client

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

type (
	NebulaSession interface {
		Open() error
		Execute(statement string) (*nebula.ResultSet, error)
		Close() error
	}

	defaultNebulaSession struct {
		session     *nebula.Session
		hostAddress nebula.HostAddress
		user        string
		password    string
		logger      logger.Logger
	}
)

func newNebulaSession(hostAddress nebula.HostAddress, user, password string, l logger.Logger) NebulaSession {
	if l == nil {
		l = logger.NopLogger
	}
	return &defaultNebulaSession{
		hostAddress: hostAddress,
		user:        user,
		password:    password,
		logger:      l,
	}
}

func (s *defaultNebulaSession) Open() error {
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

func (s *defaultNebulaSession) Execute(statement string) (*nebula.ResultSet, error) {
	return s.session.Execute(statement)
}

func (s *defaultNebulaSession) Close() error {
	s.session.Release()
	return nil
}
