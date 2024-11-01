package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/errors"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type ResetMetaPasswordParams struct {
	MetaServerAddress string
	Username          string
	Password          string
	NewPassword       string
	TimeoutSec        int
}

type ResetMetaPassword struct {
	JobContext *JobContext
	params     *ResetMetaPasswordParams
}

func NewResetMetaPassword(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ResetMetaPasswordParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}

	return &ResetMetaPassword{
		JobContext: taskContext,
		params:     params,
	}, nil
}

func (d *ResetMetaPassword) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.params.TimeoutSec)*time.Second)
	defer cancel()
	var err error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("reset meta password timeout, err: %s", err)
		default:
			changeFailed, err := d.reset()
			if err == nil {
				return nil
			}
			if changeFailed {
				return err
			}
		}
	}
}

func (d *ResetMetaPassword) reset() (changeFailed bool, err error) {
	client, err := meta.NewMetaClient(d.params.MetaServerAddress,
		meta.WithUserPassword(d.params.Username, d.params.Password))
	if err != nil {
		return false, fmt.Errorf("create meta client failed: %s", err)
	}
	defer client.Close()
	changePasswdReq := meta.NewChangePasswordReq(d.params.Username, d.params.Password, d.params.NewPassword)
	if err := client.ChangePassword(changePasswdReq); err != nil {
		ngerr, ok := err.(*errors.NebulaError)
		if !ok {
			return false, fmt.Errorf("change password failed: %s", err)
		}
		if ngerr.Code() == errors.ERROR_LEADER_CHANGED {
			return false, err
		} else {
			return true, err
		}
	}
	return false, nil
}

func (d *ResetMetaPassword) Rollback() error {
	return nil
}

func (d *ResetMetaPassword) String() string {
	return "ResetMetaPassword"
}
