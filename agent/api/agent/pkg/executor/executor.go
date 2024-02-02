package executor

import (
	"bytes"
	"context"
	"os/exec"
	"sync"
	"time"
)

const (
	executeDefaultTimeout = time.Second * 60
)

var (
	_ Executor = &defaultExecutor{}

	AsyncExecutor = NewExecutor()
)

type Executor interface {
	Execute(ctx context.Context, cmd string, timeout ...time.Duration) (stdout []byte, stderr []byte, err error)
	ExecuteAsync(ctx context.Context, cmdId, cmd string, timeout ...time.Duration) (pid int, err error)
	GetStatus(cmdId string) (done bool, stdout []byte, stderr []byte, err error)
}

type CommandResult struct {
	Stdout     []byte
	Stderr     []byte
	Err        error
	Done       bool
	FinishTime time.Time
}

type defaultExecutor struct {
	mu      sync.Mutex
	results map[string]CommandResult
}

func NewExecutor() Executor {
	return &defaultExecutor{}
}

func (e *defaultExecutor) Execute(ctx context.Context, command string, timeout ...time.Duration) ([]byte, []byte, error) {
	if len(timeout) == 0 {
		timeout = append(timeout, executeDefaultTimeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	return stdout.Bytes(), stderr.Bytes(), err
}

func (e *defaultExecutor) ExecuteAsync(ctx context.Context, cmdId, command string, timeout ...time.Duration) (pid int, err error) {
	if len(timeout) == 0 {
		timeout = append(timeout, executeDefaultTimeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err = cmd.Start(); err != nil {
		cancel()
		return -1, err
	}

	pid = cmd.Process.Pid

	go func() {
		defer cancel()

		err = cmd.Wait()

		e.mu.Lock()
		e.results[cmdId] = CommandResult{
			Stdout:     stdout.Bytes(),
			Stderr:     stderr.Bytes(),
			Err:        err,
			Done:       true,
			FinishTime: time.Now(),
		}
		e.mu.Unlock()
	}()

	return pid, nil
}

func (e *defaultExecutor) GetStatus(cmdId string) (done bool, stdout []byte, stderr []byte, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	result, ok := e.results[cmdId]
	if !ok {
		return false, nil, nil, nil
	}

	return result.Done, result.Stdout, result.Stderr, result.Err
}
