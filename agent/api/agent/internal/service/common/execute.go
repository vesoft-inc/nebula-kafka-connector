package common

import (
	"fmt"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/audit"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/executor"
)

func (s *commonService) CmdExecute(req *types.CmdExecuteReq) (resp *types.CmdExecuteResp, err error) {
	if err = audit.RecordOperation(s.ctx, audit.OpExecuteCmd, fmt.Sprintf("execute cmd `%s`", req.Command)); err != nil {
		return nil, err
	}

	timeout := time.Duration(req.Timeout) * time.Second

	exec := executor.NewExecutor()
	stdout, stderr, err := exec.Execute(s.ctx, req.Command, timeout)
	if err != nil {
		return nil, err
	}

	return &types.CmdExecuteResp{
		Stdout: string(stdout),
		Stderr: string(stderr),
	}, nil
}

func (s *commonService) CmdExecuteAsync(req *types.CmdExecuteAsyncReq) (resp *types.CmdExecuteAsyncResp, err error) {
	if err = audit.RecordOperation(s.ctx, audit.OpExecuteCmd, fmt.Sprintf("execute cmd `%s`", req.Command)); err != nil {
		return nil, err
	}

	timeout := time.Duration(req.Timeout) * time.Second

	pid, err := executor.AsyncExecutor.ExecuteAsync(s.ctx, req.CmdId, req.Command, timeout)
	if err != nil {
		return nil, err
	}

	return &types.CmdExecuteAsyncResp{
		Pid: pid,
	}, nil
}

func (s *commonService) GetCmdExecuteAsyncStatus(req *types.GetCmdExecuteAsyncStatusReq) (resp *types.GetCmdExecuteAsyncStatusResp, err error) {
	done, stdout, stderr, err := executor.AsyncExecutor.GetStatus(req.CmdId)
	if err != nil {
		return nil, err
	}

	return &types.GetCmdExecuteAsyncStatusResp{
		Done:   done,
		Stdout: string(stdout),
		Stderr: string(stderr),
	}, nil
}
