package client

import (
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

type defaultSessionV5 struct {
	session     nebula.Conn
	hostAddress string
	user        string
	password    string
	logger      logger.Logger
}

func newSessionV5(hostAddress, user, password string, l logger.Logger) Session {
	if l == nil {
		l = logger.NopLogger
	}
	return &defaultSessionV5{
		hostAddress: hostAddress,
		user:        user,
		password:    password,
		logger:      l,
	}
}

func (s *defaultSessionV5) Open() error {
	hostAddress := s.hostAddress
	connection, err := nebula.NewNebulaClient(hostAddress, s.user, s.password)
	if err != nil {
		return err
	}
	if err := connection.Ping(); err != nil {
		return err
	}

	s.session = connection
	return nil
}

func (s *defaultSessionV5) Execute(statement string) (Response, error) {
	startTime := time.Now()
	rs, err := s.session.Execute(statement)
	if err != nil {
		return nil, err
	}
	return newResponseV5(rs, time.Since(startTime), nil), nil
}

func (s *defaultSessionV5) Close() error {
	s.session.Close()
	return nil
}
