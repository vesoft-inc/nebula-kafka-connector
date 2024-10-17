package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

type CreateServiceGroupParams struct {
	ServiceGroupSpec  *types.ServiceGroup
	MetaSpec          *types.MetadSpec
	MetaServerAddress string
	Username          string
	Password          string
}

func NewCreateServiceGroup(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*CreateServiceGroupParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}

	client, err := meta.NewMetaClient(params.MetaServerAddress, meta.WithUserPassword(params.Username, params.Password))
	if err != nil {
		return nil, fmt.Errorf("create meta client failed: %s", err)
	}
	defer client.Close()

	if _, err = client.Login(); err != nil {
		return nil, fmt.Errorf("login meta server failed: %s", err)
	}

	return &CreateServiceGroup{
		JobContext:  taskContext,
		taskSpec:    taskSpec,
		clusterSpec: params.ServiceGroupSpec,
		metaSpec:    params.MetaSpec,
		metaClient:  client,
	}, nil
}

type CreateServiceGroup struct {
	JobContext  *JobContext
	taskSpec    *types.TaskSpec
	clusterSpec *types.ServiceGroup
	metaSpec    *types.MetadSpec
	metaClient  meta.Client
}

func (d *CreateServiceGroup) Execute() error {
	if d.metaClient == nil {
		return fmt.Errorf("meta client is nil")
	}

	//1. create cluster
	// TODO should modify owner
	req := meta.NewCreateServiceGroupReq(d.clusterSpec.Name, d.clusterSpec.Replica, "", d.clusterSpec.ZoneList)
	if err := d.metaClient.CreateServiceGroup(req); err != nil {
		return fmt.Errorf("create cluster failed: %s", err)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	d.JobContext.Logger.Info("create cluster success: " + d.clusterSpec.Name)
	//2. add graphd & storaged
	for _, host := range d.clusterSpec.Graphd.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.clusterSpec.Graphd.Config))
		if err != nil {
			return fmt.Errorf("get graphd port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Agent.Host), port, meta.ServiceTypeGraphd, d.clusterSpec.Name)
		if err := d.metaClient.AddService(req); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Agent.Host)
	}
	for _, host := range d.clusterSpec.Storaged.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.clusterSpec.Storaged.Config))
		if err != nil {
			return fmt.Errorf("get storaged port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Agent.Host), port, meta.ServiceTypeStoraged, d.clusterSpec.Name)
		if err := d.metaClient.AddService(req); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Agent.Host)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	//3. init cluster
	initReq := meta.NewInitServiceGroupReq(d.clusterSpec.Name)
	if err := d.metaClient.InitServiceGroup(initReq); err != nil {
		return fmt.Errorf("init cluster failed: %s", err)
	}
	d.JobContext.Logger.Info("init cluster success: " + d.clusterSpec.Name)
	return nil
}

func (d *CreateServiceGroup) Rollback() error {
	return nil
}

func (d *CreateServiceGroup) String() string {
	return "CreateServiceGroup"
}

func (d *CreateServiceGroup) ifExited() bool {
	select {
	case <-d.JobContext.Sigs:
		return true
	default:
		return false
	}
}
