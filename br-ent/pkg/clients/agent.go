package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vesoft-inc/go-pkg/httpclient"
	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type (
	ServiceName   string
	ServiceStatus string
)

const (
	ServiceNameMetad    ServiceName = "metad"
	ServiceNameStoraged ServiceName = "storaged"
	ServiceNameGraphd   ServiceName = "graphd"
	ServiceNameUnknown  ServiceName = "unknown"

	ServiceStatusRunning ServiceStatus = "Running"
	ServiceStatusExited  ServiceStatus = "Exited"
	ServiceStatusUnknown ServiceStatus = "Unknown"
)

func ToName(t meta.ServiceType) ServiceName {
	switch t {
	case meta.ServiceTypeGraphd:
		return ServiceNameGraphd
	case meta.ServiceTypeStoraged:
		return ServiceNameStoraged
	case meta.ServiceTypeMetad:
		return ServiceNameMetad
	default:
		log.Errorf("Bad format role: %d", t)
	}
	return ServiceNameUnknown
}

var certPath string

func SetCertPath(path string) {
	certPath = path
}

type (
	AgentClient interface {
		UploadFile(b *agentstorage.Backend, sourcePath string, recursively bool) error
		DownloadFile(b *agentstorage.Backend, targetPath string, recursively bool) error
		StopService(serviceType meta.ServiceType, dir string) error
		StartService(serviceType meta.ServiceType, dir string) error
		ServiceStatus(serviceType meta.ServiceType, dir string) (ServiceStatus, error)
		MoveDir(srcPath, dstPath string) error
		RemoveDir(path string) error
		ExistDir(path string) (bool, error)
		GetInstallPath(serviceType meta.ServiceType) (string, error)
		GetDataPaths(serviceType meta.ServiceType, installPath string) ([]string, error)
		DBPlayBack(backupName, installPath, dataPath, serviceMap string) error
	}

	agentClient struct {
		client  httpclient.ObjectClient
		timeout int
	}
)

func NewAgentClient(addr string) (AgentClient, error) {
	if certPath == "" {
		certPath = "certs"
	}
	tlsConfig, err := GetTLSConfig(certPath)
	if err != nil {
		return nil, err
	}
	addr = utils.GetHttpsHost(addr)

	return &agentClient{
		client:  httpclient.NewObjectClient(addr, httpclient.WithTLSClientConfig(tlsConfig)),
		timeout: 10,
	}, nil
}

type AgentManager struct {
	ctx    context.Context
	agents map[string]AgentClient // group by ip or host
}

func NewAgentManager(ctx context.Context, agentsAddr string) (*AgentManager, error) {
	agents := make(map[string]AgentClient)

	agentAddrs := strings.Split(agentsAddr, ",")
	for _, addr := range agentAddrs {
		agent, err := NewAgentClient(addr)
		if err != nil {
			return nil, err
		}

		sAddr := strings.Split(addr, ":")
		if len(sAddr) != 2 {
			return nil, fmt.Errorf("bad format agent address: %s", addr)
		}

		agents[sAddr[0]] = agent
	}

	return &AgentManager{
		ctx:    ctx,
		agents: agents,
	}, nil
}

func (a *AgentManager) GetAgent(host string) (AgentClient, error) {
	agent, ok := a.agents[host]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", host)
	}
	return agent, nil
}

func (a *AgentManager) GetAgents() map[string]AgentClient {
	return a.agents
}

func (ag *agentClient) shell(cmd string, sudo bool) (stdout string, stderr string, err error) {
	var respBody CmdExecuteResp
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	err = ag.client.Post(CmdExecutePath, &CmdExecuteReq{
		Command: cmd,
		Timeout: ag.timeout,
	}, &respBody)

	if err != nil {
		return "", "", fmt.Errorf("agent failed to post request: %s", err)
	}
	if respBody.Code != 0 {
		return "", "", fmt.Errorf("code: %d, err: %s, stderr: %s", respBody.Code, respBody.Data.Err, respBody.Data.Stderr)
	}

	return respBody.Data.Stdout, respBody.Data.Stderr, nil
}

func (ag *agentClient) UploadFile(b *agentstorage.Backend, sourcePath string, recursively bool) error {
	switch b.Type() {
	case agentstorage.S3Type:
		req := &S3UploadReq{
			Endpoint:  b.S3.Endpoint,
			Region:    b.S3.Region,
			AccessKey: b.S3.AccessKey,
			SecretKey: b.S3.SecretKey,
			Bucket:    b.S3.Bucket,
			Path:      b.S3.Path,
			LocalPath: sourcePath,
		}

		var respBody S3UploadResp
		if err := ag.client.Post(S3UploadPath, req, &respBody); err != nil {
			return fmt.Errorf("upload s3 file failed: %w", err)
		}

		return nil
	case agentstorage.HDFSType:
		req := &HDFSUploadReq{
			Address:   b.HDFS.Address,
			Username:  b.HDFS.Username,
			Path:      b.HDFS.Path,
			LocalPath: sourcePath,
		}

		if b.HDFS.Kerberos.Enable {
			req.Kerberos = KerberosConfig{
				Enable:                       b.HDFS.Kerberos.Enable,
				Principle:                    b.HDFS.Kerberos.Principle,
				KeytabFilePath:               b.HDFS.Kerberos.KeytabFilePath,
				ConfigFilePath:               b.HDFS.Kerberos.ConfigFilePath,
				KerberosServicePrincipleName: b.HDFS.Kerberos.KerberosServicePrincipleName,
			}
		}

		var respBody HDFSUploadResp
		if err := ag.client.Post(HDFSUploadPath, req, &respBody); err != nil {
			return fmt.Errorf("upload hdfs file failed: %w", err)
		}

		return nil
	case agentstorage.LocalType:
		req := &LocalUploadReq{
			Path:        b.Local.Path,
			LocalPath:   sourcePath,
			Recursively: recursively,
		}

		var respBody LocalUploadResp
		if err := ag.client.Post(LocalUploadPath, req, &respBody); err != nil {
			return fmt.Errorf("upload local file failed: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown storage backend type")
	}
}

func (ag *agentClient) DownloadFile(b *agentstorage.Backend, targetPath string, recursively bool) error {
	switch b.Type() {
	case agentstorage.S3Type:
		req := &S3DownloadReq{
			Endpoint:  b.S3.Endpoint,
			Region:    b.S3.Region,
			AccessKey: b.S3.AccessKey,
			SecretKey: b.S3.SecretKey,
			Bucket:    b.S3.Bucket,
			Path:      b.S3.Path,
			LocalPath: targetPath,
		}

		var respBody S3DownloadResp
		if err := ag.client.Post(S3DownloadPath, req, &respBody); err != nil {
			return fmt.Errorf("download s3 file failed: %w", err)
		}
		return nil
	case agentstorage.HDFSType:
		req := &HDFSDownloadReq{
			Address:   b.HDFS.Address,
			Username:  b.HDFS.Username,
			Path:      b.HDFS.Path,
			LocalPath: targetPath,
		}

		if b.HDFS.Kerberos.Enable {
			req.Kerberos = KerberosConfig{
				Enable:                       b.HDFS.Kerberos.Enable,
				Principle:                    b.HDFS.Kerberos.Principle,
				KeytabFilePath:               b.HDFS.Kerberos.KeytabFilePath,
				ConfigFilePath:               b.HDFS.Kerberos.ConfigFilePath,
				KerberosServicePrincipleName: b.HDFS.Kerberos.KerberosServicePrincipleName,
			}
		}

		var respBody HDFSDownloadResp
		if err := ag.client.Post(HDFSDownloadPath, req, &respBody); err != nil {
			return fmt.Errorf("download hdfs file failed: %w", err)
		}

		return nil
	case agentstorage.LocalType:
		req := &LocalDownloadReq{
			Path:        b.Local.Path,
			LocalPath:   targetPath,
			Recursively: recursively,
		}

		var respBody LocalDownloadResp
		if err := ag.client.Post(LocalDownloadPath, req, &respBody); err != nil {
			return fmt.Errorf("download local file failed: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown storage backend type")
	}
}

func (ag *agentClient) StopService(serviceType meta.ServiceType, dir string) error {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service stop %s", dir, ToName(serviceType))

	stdout, stderr, err := ag.shell(cmdStr, false)
	if err != nil {
		return fmt.Errorf("stop service failed %s, %s, %w", stdout, stderr, err)
	}

	return nil
}

func (ag *agentClient) StartService(serviceType meta.ServiceType, dir string) error {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service start %s", dir, ToName(serviceType))
	stdout, stderr, err := ag.shell(cmdStr, false)
	if err != nil {
		return fmt.Errorf("start service failed %s, %s, %w", stdout, stderr, err)

	}
	return nil
}

func (ag *agentClient) ServiceStatus(serviceType meta.ServiceType, dir string) (ServiceStatus, error) {
	cmdStr := fmt.Sprintf("cd %s && scripts/nebula.service status %s", dir, ToName(serviceType))
	stdout, stderr, err := ag.shell(cmdStr, false)
	if err != nil {
		return ServiceStatusUnknown, fmt.Errorf("get service status failed %s, %s, %w", stdout, stderr, err)
	}

	// an example: [INFO] nebula-graphd(46b2aac66): Exited
	if strings.Contains(stdout, "Exit") {
		return ServiceStatusExited, nil
	}

	// an example: [INFO] nebula-metad(46b2aac66): Running as 25859, Listening on 29559
	if strings.Contains(stdout, "Run") {
		return ServiceStatusRunning, nil
	}
	return ServiceStatusUnknown, nil
}

func (ag *agentClient) MoveDir(srcPath, dstPath string) error {
	stdout, stderr, err := ag.shell(fmt.Sprintf("mv %s %s", srcPath, dstPath), false)
	if err != nil {
		return fmt.Errorf("move dir failed %s, %s, %w", stdout, stderr, err)
	}
	return nil
}

func (ag *agentClient) RemoveDir(path string) error {
	stdout, stderr, err := ag.shell(fmt.Sprintf("rm -rf %s", path), false)
	if err != nil {
		return fmt.Errorf("remove dir failed %s, %s, %w", stdout, stderr, err)
	}

	return nil
}

func (ag *agentClient) ExistDir(path string) (bool, error) {
	stdout, stderr, err := ag.shell(fmt.Sprintf("[ -d %s ] && echo 'exist' || echo 'not exist'", path), false)
	if err != nil {
		return false, fmt.Errorf("check dir exist failed %s, %s, %w", stdout, stderr, err)
	}

	if strings.Contains(stdout, "exist") {
		return true, nil
	}
	return false, nil
}

func (ag *agentClient) getComponentConfig(serviceType meta.ServiceType) (map[string]interface{}, error) {
	var respBody GetComponentConfigResp
	err := ag.client.Get(fmt.Sprintf("%s?component=%s", GetComponentConfigPath, ToName(serviceType)), &respBody)
	if err != nil {
		return nil, fmt.Errorf("get component config failed: %w", err)
	}
	if respBody.Code != 0 {
		return nil, fmt.Errorf("get component config failed: %s", respBody.Message)
	}

	return respBody.Data.Config, nil
}

func (ag *agentClient) GetInstallPath(serviceType meta.ServiceType) (string, error) {
	config, err := ag.getComponentConfig(serviceType)
	if err != nil {
		return "", err
	}

	installPath, ok := config["install_path"].(string)
	if !ok {
		return "", fmt.Errorf("install path not found")
	}

	return installPath, nil
}

func (ag *agentClient) GetDataPaths(serviceType meta.ServiceType, installPath string) ([]string, error) {
	cmdStr := fmt.Sprintf("cd %s && cat etc/nebula-%s.conf | grep data_path | grep -v '^\\s*#'", installPath, ToName(serviceType))
	stdout, stderr, err := ag.shell(cmdStr, false)
	if err != nil {
		return nil, fmt.Errorf("get data paths failed %s, %s, %w", stdout, stderr, err)
	}

	stdout = strings.TrimSpace(stdout)

	datapath := strings.Split(stdout, "=")
	if len(datapath) != 2 {
		return nil, fmt.Errorf("data dataPath not valid: %s", datapath)
	}

	var paths []string
	for _, dataPath := range strings.Split(datapath[1], ",") {
		if dataPath == "" {
			continue
		}
		if dataPath[0] == '/' {
			paths = append(paths, dataPath)
		} else {
			dataPath = filepath.Join(installPath, dataPath)
			paths = append(paths, dataPath)
		}
	}

	return paths, nil
}

func (ag *agentClient) DBPlayBack(backupName, installPath, dataPath, serviceMap string) error {
	cmdStr := fmt.Sprintf("cd %s && bin/db_playback --v=2 --db_path=%s --backup_name=%s --service_map=%s", installPath, dataPath, backupName, serviceMap)

	log.Infof("db playback cmd: %s", cmdStr)

	stdout, stderr, err := ag.shell(cmdStr, false)
	if err != nil {
		return fmt.Errorf("db playback failed %s, %s, %w", stdout, stderr, err)
	}

	fmt.Println("db playback success, log:", stderr)

	return nil
}

func GetTLSConfig(certPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath+"/client.crt", certPath+"/client.key")
	if err != nil {
		return nil, fmt.Errorf("unable to load client cert and key: %v", err)
	}

	caCert, err := os.ReadFile(certPath + "/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("unable to read ca cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("unable to append ca cert to pool")
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}, nil
}
