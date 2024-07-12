package clients

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type NebulaMeta struct {
	leaderAddr string
	client     meta.Client
	tlsConfig  *tls.Config
}

func NewMeta(addrStrs, username, password string, tlsConfig *tls.Config) (*NebulaMeta, error) {
	// TODO: support tls
	client, err := meta.NewMetaClient(addrStrs, meta.WithUserPassword(username, password))
	if err != nil {
		return nil, fmt.Errorf("create meta client failed: %v", err)
	}

	resp, err := client.Login()
	if err != nil {
		return nil, fmt.Errorf("login meta failed: %v", err)
	}

	leaderAddr := strings.Split(resp.Leader, ",")
	m := &NebulaMeta{
		leaderAddr: leaderAddr[0],
		tlsConfig:  tlsConfig,
		client:     client,
	}

	return m, nil
}

func (m *NebulaMeta) LeaderAddr() string {
	return m.leaderAddr
}

func (m *NebulaMeta) CreateFullBackup(backupName string, clusterId int64) (*meta.CreateBackupResp, error) {
	resp, err := m.client.CreateBackup(meta.NewCreateBackupReq(backupName, []int64{clusterId}))
	if err != nil {
		return nil, err
	}

	if resp.NewHost != "" {
		m.leaderAddr = fmt.Sprintf("%s:%d", resp.NewHost, resp.NewPort)
	}

	return resp, nil
}

func (m *NebulaMeta) Restore(req *meta.RestoreReq) (*meta.RestoreResp, error) {
	resp, err := m.client.Restore(req)
	if err != nil {
		return nil, err
	}

	if resp.NewHost != "" {
		m.leaderAddr = fmt.Sprintf("%s:%d", resp.NewHost, resp.NewPort)
	}

	return resp, nil
}

func (m *NebulaMeta) DropBackup(req *meta.DropBackupReq) (*meta.DropBackupResp, error) {
	resp, err := m.client.DropBackup(req)
	if err != nil {
		return nil, err
	}

	if resp.NewHost != "" {
		m.leaderAddr = fmt.Sprintf("%s:%d", resp.NewHost, resp.NewPort)
	}

	return resp, nil
}

func (m *NebulaMeta) ShowMeta() (*meta.ShowMetaResp, error) {
	resp, err := m.client.ShowMeta()
	if err != nil {
		return nil, err
	}

	if resp.NewHost != "" {
		m.leaderAddr = fmt.Sprintf("%s:%d", resp.NewHost, resp.NewPort)
	}
	return resp, nil
}

type ServiceInfo struct {
	ServiceId   int64
	ServiceType meta.ServiceType
	Host        string
	Port        uint32
	InstallPath string
	DataPaths   []string
}

type ClusterServiceInfo struct {
	ClusterId int64
	Services  []*ServiceInfo
}

func (m *NebulaMeta) ListClusters(amg *AgentManager, clusterId int64) ([]*ClusterServiceInfo, error) {
	clusterResp, err := m.client.ListClusters(meta.NewListClustersReq(""))
	if err != nil {
		return nil, err
	}

	clusters := make([]*ClusterServiceInfo, 0)
	for _, c := range clusterResp.Clusters {
		if c.Id != clusterId {
			continue
		}
		cluster := &ClusterServiceInfo{
			ClusterId: c.Id,
			Services:  make([]*ServiceInfo, 0),
		}

		serviceResp, err := m.client.ListServices(meta.NewListServicesReq(c.Name))
		if err != nil {
			return nil, err
		}
		for _, s := range serviceResp.Services {
			sInfo := &ServiceInfo{
				ServiceId:   s.Id,
				ServiceType: s.Type,
				Host:        s.Host,
				Port:        s.Port,
			}
			if s.Type == meta.ServiceTypeStoraged {
				agent, err := amg.GetAgent(s.Host)
				if err != nil {
					return nil, fmt.Errorf("get agent %s failed: %w", s.Host, err)
				}
				installPath, err := agent.GetInstallPath(s.Type)
				if err != nil {
					return nil, fmt.Errorf("get metad %s install path failed: %w", s.Host, err)
				}
				dataPaths, err := agent.GetDataPaths(s.Type, installPath)
				if err != nil {
					return nil, fmt.Errorf("get storaged %s data path failed: %w", s.Host, err)
				}

				sInfo.InstallPath = installPath
				sInfo.DataPaths = dataPaths
			}

			cluster.Services = append(cluster.Services, sInfo)
		}
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}
