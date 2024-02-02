package logger

import "github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type NewLogger func(*types.Progress) Logger
