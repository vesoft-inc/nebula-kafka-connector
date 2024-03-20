package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

type CreateClusterParams struct {
	ClusterSpec *types.Cluster
	MetaSpec    *types.MetadSpec
}

func NewCreateCluster(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*CreateClusterParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &CreateCluster{
		JobContext:  taskContext,
		taskSpec:    taskSpec,
		clusterSpec: params.ClusterSpec,
		metaSpec:    params.MetaSpec,
	}, nil
}

type CreateCluster struct {
	JobContext  *JobContext
	taskSpec    *types.TaskSpec
	clusterSpec *types.Cluster
	metaSpec    *types.MetadSpec
}

func (d *CreateCluster) Execute() error {
	//1. meta client init & login
	metaClient, err := meta.NewMetaClient(utils.GetMetaAddressListString(d.metaSpec.Hosts, utils.GetConfigPort(d.metaSpec.Config)))
	if err != nil {
		return fmt.Errorf("create meta client failed: %s", err)
	}
	defer metaClient.Close()
	//2. create cluster
	req := meta.NewCreateClusterReq(d.clusterSpec.Name, d.clusterSpec.Replica, d.clusterSpec.ZoneList)
	resp, err := metaClient.CreateCluster(req)
	if err != nil {
		return fmt.Errorf("create cluster failed: %s", err)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	d.JobContext.Logger.Info("create cluster success: " + d.clusterSpec.Name + " " + resp.Msg)
	//3. add graphd & storaged
	for _, host := range d.clusterSpec.Graphd.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.clusterSpec.Graphd.Config))
		if err != nil {
			return fmt.Errorf("get graphd port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Host), port, meta.ServiceTypeGraphd, d.clusterSpec.Name)
		resp, err := metaClient.AddService(req)
		if err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Host + " " + resp.Msg)
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
		resp, err := metaClient.AddService(req)
		if err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		d.JobContext.Logger.Info("add host success: " + host.Host + " " + resp.Msg)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	//4. init cluster
	initReq := meta.NewInitClusterReq(d.clusterSpec.Name)
	initResp, err := metaClient.InitCluster(initReq)
	if err != nil {
		return fmt.Errorf("init cluster failed: %s", err)
	}
	d.JobContext.Logger.Info("init cluster success: " + d.clusterSpec.Name + " " + initResp.Msg)
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
