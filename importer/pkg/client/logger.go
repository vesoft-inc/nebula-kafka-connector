package client

import (
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/logger"
)

var _ Logger = nebulaLogger{}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type nebulaLogger struct {
	l logger.Logger
}

func newNebulaLogger(l logger.Logger) Logger {
	return nebulaLogger{
		l: l,
	}
}

//revive:disable:empty-lines

func (l nebulaLogger) Info(msg string)  { l.l.Info(msg) }
func (l nebulaLogger) Warn(msg string)  { l.l.Warn(msg) }
func (l nebulaLogger) Error(msg string) { l.l.Error(msg) }
func (l nebulaLogger) Fatal(msg string) { l.l.Fatal(msg) }

//revive:enable:empty-lines
