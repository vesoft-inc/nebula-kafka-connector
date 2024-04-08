package executor

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/appleboy/easyssh-proxy"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"golang.org/x/crypto/ssh"
)

type SSHExecutor struct {
	Config        *types.SSHConfig
	easysshConfig *easyssh.MakeConfig
	client        *ssh.Client
}

var clientCacheMap = make(map[string]*SSHExecutor)

// clean used ssh client when exit
func CleanClient() {
	for _, client := range clientCacheMap {
		err := client.client.Close()
		if err != nil {
			fmt.Println("failed to close client: ", err)
		}
	}
	clientCacheMap = make(map[string]*SSHExecutor)
}

func NewSSHExecuter(config *types.SSHConfig) (*SSHExecutor, error) {
	executor := &SSHExecutor{
		Config: config,
	}
	executor.translateConfig()
	session, client, err := executor.easysshConfig.Connect()
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("failed to connect to %s", config.Host)
	}
	defer session.Close()
	executor.client = client
	clientCacheMap[config.Host] = executor
	return executor, nil
}
func (e *SSHExecutor) translateConfig() {
	config := e.Config
	easysshConfig := &easyssh.MakeConfig{
		Server:  config.Host,
		Port:    strconv.Itoa(config.Port),
		User:    config.User,
		Timeout: config.Timeout, // timeout when connecting to remote
	}

	// prefer private key authentication
	if len(config.Key) > 0 {
		easysshConfig.Key = config.Key
		easysshConfig.Passphrase = config.Passphrase
	} else if len(config.Password) > 0 {
		easysshConfig.Password = config.Password
	}
	e.easysshConfig = easysshConfig
}

func (a *SSHExecutor) Health() error {
	_, _, err := a.Shell("echo 1", false)
	return err
}

func (e *SSHExecutor) Shell(cmd string, sudo bool) (stdout string, stderr string, err error) {
	// try to acquire root permission
	if sudo {
		// replace " to \" or bash -c will ignore
		cmd = strings.ReplaceAll(cmd, `"`, `\"`)
		cmd = fmt.Sprintf(`sudo -H bash -c "%s"`, cmd)
	}

	// TODO: sudo uses bash, not sudo uses shell, need to uniform
	// set a basic PATH in case it's empty on login
	cmd = fmt.Sprintf("PATH=$PATH:/usr/bin:/usr/sbin %s", cmd)

	if e.Config.Locale != "" {
		cmd = fmt.Sprintf("export LANG=%s; %s", e.Config.Locale, cmd)
	}
	session, err := e.getSSHSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()
	if err := session.Start(cmd); err != nil {
		return "", "", fmt.Errorf("failed to execute command: %s", err)
	}
	var stderrBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	session.Stderr = &stderrBuf
	session.Stdout = &stdoutBuf
	err = session.Wait()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func (e *SSHExecutor) getSSHSession() (*ssh.Session, error) {
	session, err := e.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}
	return session, nil
}

func (e *SSHExecutor) Upload(src string, dest string) (stdout string, stderr string, err error) {
	return "", "", e.easysshConfig.Scp(src, dest)
}

func (e *SSHExecutor) Download(src string, dest string) (stdout string, stderr string, err error) {
	return "", "", nil
}
