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

	return &CreateServiceGroup{
		JobContext:       taskContext,
		params:           params,
		taskSpec:         taskSpec,
		serviceGroupSpec: params.ServiceGroupSpec,
		metaSpec:         params.MetaSpec,
	}, nil
}

type CreateServiceGroup struct {
	JobContext       *JobContext
	params           *CreateServiceGroupParams
	taskSpec         *types.TaskSpec
	serviceGroupSpec *types.ServiceGroup
	metaSpec         *types.MetadSpec
}

func (d *CreateServiceGroup) Execute() error {
	client, err := meta.NewMetaClient(d.params.MetaServerAddress,
		meta.WithUserPassword(d.params.Username, d.params.Password))
	if err != nil {
		return fmt.Errorf("create meta client failed: %s", err)
	}
	defer client.Close()
	if _, err = client.Login(); err != nil {
		return fmt.Errorf("login meta server failed: %s", err)
	}

	//1. create cluster
	// TODO should modify owner
	req := meta.NewCreateServiceGroupReq(d.serviceGroupSpec.Name, d.serviceGroupSpec.Replica, "", d.serviceGroupSpec.ZoneList)
	if err := client.CreateServiceGroup(req); err != nil {
		return fmt.Errorf("create cluster failed: %s", err)
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	d.JobContext.Logger.Info("create cluster success: " + d.serviceGroupSpec.Name)
	//2. add graphd & storaged
	var addedHost = make(map[string]struct{})
	for _, host := range d.serviceGroupSpec.Graphd.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		if err := d.addHost(addedHost, client, host.IP, d.serviceGroupSpec.Name, utils.GetHostPort(host.Agent.Host)); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.serviceGroupSpec.Graphd.Config))
		if err != nil {
			return fmt.Errorf("get graphd port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Agent.Host), port, meta.ServiceTypeGraphd, d.serviceGroupSpec.Name)
		if err := client.AddService(req); err != nil {
			return fmt.Errorf("add service failed: %s", err)
		}
		d.JobContext.Logger.Info(fmt.Sprintf("add %s service success: %s:%d", "graphd", host.IP, port))
	}
	for _, host := range d.serviceGroupSpec.Storaged.Hosts {
		if d.ifExited() {
			return fmt.Errorf("exited signal received")
		}
		if err := d.addHost(addedHost, client, host.IP, d.serviceGroupSpec.Name, utils.GetHostPort(host.Agent.Host)); err != nil {
			return fmt.Errorf("add host failed: %s", err)
		}
		port, err := utils.GetUint32Port(utils.GetConfigPort(d.serviceGroupSpec.Storaged.Config))
		if err != nil {
			return fmt.Errorf("get storaged port failed: %s", err)
		}
		req := meta.NewAddServiceReq(utils.GetHostIP(host.Agent.Host), port, meta.ServiceTypeStoraged, d.serviceGroupSpec.Name)
		if err := client.AddService(req); err != nil {
			return fmt.Errorf("add service failed: %s", err)
		}
		d.JobContext.Logger.Info(fmt.Sprintf("add %s service success: %s:%d", "storaged", host.IP, port))
	}
	if d.ifExited() {
		return fmt.Errorf("exited signal received")
	}
	//3. init cluster
	initReq := meta.NewInitServiceGroupReq(d.serviceGroupSpec.Name)
	if err := client.InitServiceGroup(initReq); err != nil {
		return fmt.Errorf("init cluster failed: %s", err)
	}
	d.JobContext.Logger.Info("init cluster success: " + d.serviceGroupSpec.Name)
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

func (d *CreateServiceGroup) addHost(added map[string]struct{}, client meta.Client, host string, serviceGroup string, agentPort int) error {
	if _, ok := added[host]; ok {
		return nil
	}
	req := meta.NewAddHostReq(host, serviceGroup, uint32(agentPort), "")
	if err := client.AddHost(req); err != nil {
		return err
	}
	added[host] = struct{}{}
	return nil
}
