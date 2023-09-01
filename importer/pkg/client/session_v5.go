package client

import (
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

type defaultSessionV5 struct {
	session     *nebula.Session
	hostAddress nebula.HostAddress
	user        string
	password    string
	logger      logger.Logger
}

func newSessionV5(hostAddress HostAddress, user, password string, l logger.Logger) Session {
	if l == nil {
		l = logger.NopLogger
	}
	return &defaultSessionV5{
		hostAddress: nebula.HostAddress{
			Host: hostAddress.Host,
			Port: hostAddress.Port,
		},
		user:     user,
		password: password,
		logger:   l,
	}
}

func (s *defaultSessionV5) Open() error {
	hostAddress := s.hostAddress
	connection := nebula.NewConnection(hostAddress)
	if err := connection.Open(hostAddress, 0, nil); err != nil {
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

func (s *defaultSessionV5) Execute(statement string) (Response, error) {
	startTime := time.Now()
	rs, err := s.session.Execute(statement)
	if err != nil {
		return nil, err
	}
	return newResponseV5(rs, time.Since(startTime)), nil
}

func (s *defaultSessionV5) Close() error {
	s.session.Release()
	return nil
}
