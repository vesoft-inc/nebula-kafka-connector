// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng

import (
	"fmt"
	"sync"
)

type timezoneInfo struct {
	offset int32
	name   []byte
}

type Session struct {
	sessionID  int64
	connection *connection
	log        Logger
	mu         sync.Mutex
	timezoneInfo
}

// TODO(Aiee) used for demo only
func NewSession(sessionID int64, connection *connection, log Logger) *Session {
	return &Session{
		sessionID:  sessionID,
		connection: connection,
		log:        log,
	}
}

// Execute returns the result of the given query as a ResultSet
func (session *Session) Execute(stmt string) (*ResultSet, error) {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.connection == nil {
		return nil, fmt.Errorf("failed to execute: Session has been released")
	}

	resp, err := session.connection.execute(session.sessionID, stmt)
	if err != nil {
		return nil, err
	}

	resSet, err := genResultSet(resp, session.timezoneInfo)
	if err != nil {
		return nil, err
	}
	return resSet, nil
}

// Release logs out and releases connection hold by session.
func (session *Session) Release() {
	if session == nil {
		return
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.connection == nil {
		session.log.Warn("Session has been released")
		return
	}
	if err := session.connection.signOut(session.sessionID); err != nil {
		session.log.Warn(fmt.Sprintf("Sign out failed, %s", err.Error()))
	}
	session.connection = nil
}
