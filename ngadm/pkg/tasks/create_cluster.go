package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

type CreateClusterParams struct {
	ClusterSpec       *types.Cluster
	MetaSpec          *types.MetadSpec
	MetaServerAddress string
	Username          string
	Password          string
}

func NewCreateCluster(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*CreateClusterParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}

	client, err := meta.NewMetaClient(params.MetaServerAddress, meta.WithUserPassword(params.Username, params.Password))
	if err != nil {
		return nil, fmt.Errorf("create meta client failed: %s", err)
	}
	defer client.Close()

	if err = client.Login(); err != nil {
		return nil, fmt.Errorf("login meta server failed: %s", err)
	}

	return &CreateCluster{
		JobContext:  taskContext,
		taskSpec:    taskSpec,
		clusterSpec: params.ClusterSpec,
		metaSpec:    params.MetaSpec,
		metaClient:  client,
	}, nil
}

type CreateCluster struct {
	JobContext  *JobContext
	taskSpec    *types.TaskSpec
	clusterSpec *types.Cluster
	metaSpec    *types.MetadSpec
	metaClient  meta.Client
}

func (d *CreateCluster) Execute() error {
	if d.metaClient == nil {
		return fmt.Errorf("meta client is nil")
	}

	//1. create cluster
	req := meta.NewCreateClusterReq(d.clusterSpec.Name, d.clusterSpec.Replica, d.clusterSpec.ZoneList)
	if err := d.metaClient.CreateCluster(req); err != nil {
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
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Host), port, meta.ServiceTypeGraphd, d.clusterSpec.Name)
		if err := d.metaClient.AddService(req); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Host)
	}
	for _, host := range d.clusterSpec.Storaged.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.clusterSpec.Storaged.Config))
		if err != nil {
			return fmt.Errorf("get storaged port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Host), port, meta.ServiceTypeStoraged, d.clusterSpec.Name)
		if err := d.metaClient.AddService(req); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Host)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	//3. init cluster
	initReq := meta.NewInitClusterReq(d.clusterSpec.Name)
	if err := d.metaClient.InitCluster(initReq); err != nil {
		return fmt.Errorf("init cluster failed: %s", err)
	}
	d.JobContext.Logger.Info("init cluster success: " + d.clusterSpec.Name)
	return nil
}

func (d *CreateCluster) Rollback() error {
	return nil
}

func (d *CreateCluster) String() string {
	return "CreateCluster"
}

func (d *CreateCluster) ifExited() bool {
	select {
	case <-d.JobContext.Sigs:
		return true
	default:
		return false
	}
}
