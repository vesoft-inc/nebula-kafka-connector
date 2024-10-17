package restore

import (
	"fmt"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"path/filepath"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type Fix struct {
	r           *Restore
	amg         *clients.AgentManager
	clusters    *utils.NebulaServiceGroups
	metaServiceGroup []*clients.ServiceInfo
	backSuffix  string
}

func NewFixFrom(r *Restore) (*Fix, error) {
	if r.amg == nil {
		return nil, fmt.Errorf("empty agents manager")
	}

	return &Fix{
		r:           r,
		amg:         r.amg,
		backSuffix:  r.backSuffix,
		clusters:    r.clusters,
		metaServiceGroup: r.metaServiceGroup,
	}, nil
}

func (f *Fix) fixServiceData(s *clients.ServiceInfo) error {
	agent, err := f.amg.GetAgent(s.Host)
	if err != nil {
		return fmt.Errorf("get agent for %s failed: %w", s.Host, err)
	}

	for _, d := range s.DataPaths {
		opath := filepath.Join(d, "data")
		bpath := fmt.Sprintf("%s%s", opath, f.backSuffix)

		// check if the old data exist
		exist, err := agent.ExistDir(bpath)
		if err != nil {
			return fmt.Errorf("check %s exist failed: %w", bpath, err)
		}
		if !exist {
			log.WithField("path", bpath).Debug("Origin backup data path does not exist, skip it")
			continue
		}

		// remove the newly downloaded data
		if err = agent.RemoveDir(opath); err != nil {
			return fmt.Errorf("remove new origin dir %s failed: %w", opath, err)
		}

		// move the old data back
		err = agent.MoveDir(bpath, opath)
		if err != nil && !utils.IsNotExist(err) {
			return fmt.Errorf("move data dir back from %s to %s failed: %w", bpath, opath, err)
		}

		log.WithField("origin path", opath).
			WithField("backup path", bpath).
			WithField("origin not exist", utils.IsNotExist(err)).
			Infof("Moveback origin %s data path successfully", s.Host)
	}
	return nil
}

func (f *Fix) getDead(services []*clients.ServiceInfo) ([]*clients.ServiceInfo, error) {
	deadServices := make([]*clients.ServiceInfo, 0)

	for _, service := range services {
		logger := log.WithField("host", service.Host)

		agent, err := f.amg.GetAgent(service.Host)
		if err != nil {
			return nil, fmt.Errorf("get agent %s failed: %w", service.Host, err)
		}

		status, err := agent.ServiceStatus(service.ServiceType, service.InstallPath)
		if err != nil {
			return nil, fmt.Errorf("get service status in host %s failed: %w", service.Host, err)
		}

		if status == clients.ServiceStatusExited {
			logger.WithField("dir", service.InstallPath).WithField("role", service.ServiceType).Debugf("%d:%s is dead.",
				service.ServiceType, service.Host)
			deadServices = append(deadServices, service)
		}
	}

	return deadServices, nil
}

func (f *Fix) startServices(services []*clients.ServiceInfo) error {
	for _, ds := range services {
		name := fmt.Sprintf("%s[%s]", clients.ToName(ds.ServiceType), ds.Host)
		agent, err := f.amg.GetAgent(ds.Host)
		if err != nil {
			return fmt.Errorf("get agent for %s failed: %w", ds.Host, err)
		}
		if err = agent.StartService(ds.ServiceType, ds.InstallPath); err != nil {
			return fmt.Errorf("start %s by agent failed: %w", name, err)
		}
		log.WithField("addr", ds.Host).
			Infof("Start %s by agent successfully.", name)
	}
	return nil
}

func (f *Fix) stopServices(services []*clients.ServiceInfo) error {
	for _, service := range services {
		agent, err := f.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent for %s failed: %w", service.Host, err)
		}

		if err = agent.StopService(service.ServiceType, service.InstallPath); err != nil {
			return fmt.Errorf("stop service %s by agent failed: %w",
				service.Host, err)
		}

		log.WithField("addr", service.Host).
			Info("Stop service successfully.")
	}
	return nil
}

func retry(action func() error, aname string, times int) (err error) {
	for try := 1; try <= times; try++ {
		err = action()
		if err == nil {
			return
		}

		log.WithError(err).Infof("%s failed, try times=%d.", aname, try)
		time.Sleep(time.Second * time.Duration(try))
	}

	return
}

func (f *Fix) fixServices(services []*clients.ServiceInfo) error {
	if len(services) == 0 {
		return nil
	}

	tryTimes := 3

	stopAllServices := func() error {
		return f.stopServices(services)
	}

	// stop all service for data movement
	if err := retry(stopAllServices, "Stop all services", tryTimes); err != nil {
		return err
	}

	waitServicesStopped := func() error {
		time.Sleep(time.Second * 10)
		ds, err := f.getDead(services)
		if err != nil {
			return fmt.Errorf("get services status failed:  %w", err)
		}
		if len(ds) == len(services) {
			return nil
		}

		return fmt.Errorf("not all services are stopped, dead services: %v", ds)
	}

	if err := retry(waitServicesStopped, "Wait for all services stopped", tryTimes); err != nil {
		return err
	}

	fixServicesData := func() error {
		for _, s := range services {
			switch s.ServiceType {
			case meta.ServiceTypeMetad:
				if err := f.fixServiceData(s); err != nil {
					return err
				}
			case meta.ServiceTypeStoraged:
				if err := f.fixServiceData(s); err != nil {
					return err
				}
			default:

			}
		}
		return nil
	}

	// move back data path
	if err := retry(fixServicesData, "Fix data", tryTimes); err != nil {
		return err
	}

	// start all services
	getDeadThenStart := func() error {
		ds, err := f.getDead(services)
		if err != nil {
			return fmt.Errorf("get services failed:  %w", err)
		}
		err = f.startServices(ds)
		if err != nil {
			return fmt.Errorf("start dead services failed: %w", err)
		}
		return nil
	}
	if err := retry(getDeadThenStart, "Get dead services then start", tryTimes); err != nil {
		return err
	}

	return nil
}

func (f *Fix) Fix() error {
	// fix metad cluster
	if err := f.fixServices(f.metaServiceGroup); err != nil {
		return err
	}

	// fix clusters
	for id, cluster := range f.clusters.GetServiceGroups() {
		if id == f.r.cfg.ServiceGroupId {
			if err := f.fixServices(cluster); err != nil {
				return err
			}
		}
	}

	return nil
}
