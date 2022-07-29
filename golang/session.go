// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng_go

import (
	"fmt"
	"sync"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/generated_code/v5.0.0/nebula/graph"
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
	// timezoneInfo
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
func (session *Session) Execute(stmt string) (*graph.ExecutionResponse, error) {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.connection == nil {
		return nil, fmt.Errorf("failed to execute: Session has been released")
	}

	resp, err := session.connection.execute(session.sessionID, stmt)
	if err != nil {
		return nil, err
	}
	return resp, err
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
