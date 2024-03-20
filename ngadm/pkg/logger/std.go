package logger

import (
	"log"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type stdlogger struct {
	progress *types.Progress
}

func NewStdLogger(progress *types.Progress) Logger {
	return &stdlogger{
		progress: progress,
	}
}

func (l *stdlogger) Info(msg string) {
	log.Printf("[INFO] Progress(%d/%d) %s\n", l.progress.Current, l.progress.Total, msg)
}

func (l *stdlogger) Warn(msg string) {
	log.Printf("[WARN] %s\n", msg)
}

func (l *stdlogger) Error(msg string) {
	log.Printf("[ERROR] %s\n", msg)
}

func (l *stdlogger) Fatal(msg string) {
	log.Fatalf("[FATAL] %s\n", msg)
}
