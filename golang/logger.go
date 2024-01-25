// Copyright (c) 2022 vesoft inc. All rights reserved.

package nebula_ng

import (
	"log"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type nebulaLogger struct{}
type emptyLogger struct{}

var (
	DefaultLogger = &nebulaLogger{}
	EmptyLogger   = &emptyLogger{}
)

func (l nebulaLogger) Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func (l nebulaLogger) Warn(msg string) {
	log.Printf("[WARNING] %s\n", msg)
}

func (l nebulaLogger) Error(msg string) {
	log.Printf("[ERROR] %s\n", msg)
}

func (l emptyLogger) Info(msg string) {
	return
}

func (l emptyLogger) Warn(msg string) {
	return
}

func (l emptyLogger) Error(msg string) {
	return
}
