package restore

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func (r *Restore) stopServiceGroupGraph(clusterId int64) error {
	for id, cluster := range r.clusters.GetServiceGroups() {
		if id == clusterId {
			for _, service := range cluster {
				if service.ServiceType == meta.ServiceTypeGraphd {
					agent, err := r.amg.GetAgent(service.Host)
					if err != nil {
						return fmt.Errorf("get agent for graphd %s failed: %w", service.Host, err)
					}

					if err = agent.StopService(service.ServiceType, service.InstallPath); err != nil {
						return fmt.Errorf("stop graphd service %s by agent failed: %w",
							service.Host, err)
					}

					log.WithField("addr", service.Host).
						Info("Stop graphd service successfully.")
				}
			}
		}
	}

	return nil
}

func (r *Restore) startServiceGroupGraph(clusterId int64) error {
	for id, cluster := range r.clusters.GetServiceGroups() {
		if id == clusterId {
			for _, service := range cluster {
				if service.ServiceType == meta.ServiceTypeGraphd {
					agent, err := r.amg.GetAgent(service.Host)
					if err != nil {
						return fmt.Errorf("get agent for graphd %s failed: %w", service.Host, err)
					}

					if err = agent.StartService(service.ServiceType, service.InstallPath); err != nil {
						return fmt.Errorf("start graphd service %s by agent failed: %w",
							service.Host, err)
					}

					log.WithField("addr", service.Host).
						Info("Start graphd service successfully.")
				}
			}
		}
	}

	return nil
}

func (r *Restore) stopServiceGroupStorage(clusterId int64) error {
	for id, cluster := range r.clusters.GetServiceGroups() {
		if id == clusterId {
			for _, service := range cluster {
				if service.ServiceType == meta.ServiceTypeStoraged {
					agent, err := r.amg.GetAgent(service.Host)
					if err != nil {
						return fmt.Errorf("get agent for storaged %s failed: %w", service.Host, err)
					}

					if err = agent.StopService(service.ServiceType, service.InstallPath); err != nil {
						return fmt.Errorf("stop storaged service %s by agent failed: %w",
							service.Host, err)
					}

					log.WithField("addr", service.Host).
						Info("Stop storaged service successfully.")
				}
			}
		}
	}

	return nil
}

func (r *Restore) startServiceGroupStorage(clusterId int64) error {
	for id, cluster := range r.clusters.GetServiceGroups() {
		if id == clusterId {
			for _, service := range cluster {
				if service.ServiceType == meta.ServiceTypeStoraged {
					agent, err := r.amg.GetAgent(service.Host)
					if err != nil {
						return fmt.Errorf("get agent for storaged %s failed: %w", service.Host, err)
					}

					if err = agent.StartService(service.ServiceType, service.InstallPath); err != nil {
						return fmt.Errorf("start storaged service %s by agent failed: %w",
							service.Host, err)
					}

					log.WithField("addr", service.Host).
						Info("Start storaged service successfully.")
				}
			}
		}
	}

	return nil
}

func (r *Restore) stopAllServiceGroupGraph() error {
	for id := range r.clusters.GetServiceGroups() {
		if err := r.stopServiceGroupGraph(id); err != nil {
			return fmt.Errorf("stop graphd cluster %d failed: %w", id, err)
		}
	}
	return nil
}

func (r *Restore) startAllServiceGroupGraph() error {
	for id := range r.clusters.GetServiceGroups() {
		if err := r.startServiceGroupGraph(id); err != nil {
			return fmt.Errorf("start graphd cluster %d failed: %w", id, err)
		}
	}
	return nil
}

func (r *Restore) stopAllServiceGroupStorage() error {
	for id := range r.clusters.GetServiceGroups() {
		if err := r.stopServiceGroupStorage(id); err != nil {
			return fmt.Errorf("stop storaged cluster %d failed: %w", id, err)
		}
	}
	return nil
}

func (r *Restore) startAllServiceGroupStorage() error {
	for id := range r.clusters.GetServiceGroups() {
		if err := r.startServiceGroupStorage(id); err != nil {
			return fmt.Errorf("start storaged cluster %d failed: %w", id, err)
		}
	}
	return nil
}

func (r *Restore) stopAllServiceGroups() error {
	if err := r.stopAllServiceGroupGraph(); err != nil {
		return fmt.Errorf("stop graphd cluster failed: %w", err)
	}

	if err := r.stopAllServiceGroupStorage(); err != nil {
		return fmt.Errorf("stop storaged cluster failed: %w", err)
	}

	return nil
}

func (r *Restore) stopStorageds() error {
	//if err := r.stopAllServiceGroupStorage(); err != nil {
	//	return fmt.Errorf("stop storaged cluster failed: %w", err)
	//}

	storages := r.clusters.GetStorages(r.clusterId)
	for _, service := range storages {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent for storaged %s failed: %w", service.Host, err)
		}

		if err = agent.StopService(service.ServiceType, service.InstallPath); err != nil {
			return fmt.Errorf("stop storaged service %s by agent failed: %w",
				service.Host, err)
		}

		log.WithField("addr", service.Host).
			Info("Stop storaged service successfully.")
	}
	//log.Info("Waiting for lm to clear quota...")
	//time.Sleep(time.Second * 30)

	return nil
}

func (r *Restore) stopGraphds() error {
	graphs := r.clusters.GetGraphs(r.clusterId)
	for _, service := range graphs {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent for graphd %s failed: %w", service.Host, err)
		}

		if err = agent.StopService(service.ServiceType, service.InstallPath); err != nil {
			return fmt.Errorf("stop graphd service %s by agent failed: %w",
				service.Host, err)
		}

		log.WithField("addr", service.Host).
			Info("Stop graphd service successfully.")
	}

	return nil

}
