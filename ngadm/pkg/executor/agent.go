package executor

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vesoft-inc/go-pkg/httpclient"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/utils"
)

var CmdExecutePath = "/api/v1/common/execute"
var CmdExecuteAsyncPath = "/api/v1/common/execute-async"
var CmdExecuteResultPath = "/api/v1/common/execute-async/"
var UploadPath = "/api/v1/common/upload"
var HealthPath = "/health"
var AgentConfigPath = "/api/v1/component-config"

type AgentExecutor struct {
	Host      string
	CertPath  string
	client    httpclient.ObjectClient
	timeout   int
	tlsConfig *tls.Config
}
type AgentResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    AgentData `json:"data"`
}
type AgentData struct {
	Stdout string         `json:"stdout,omitempty"`
	Stderr string         `json:"stderr,omitempty"`
	Err    string         `json:"err,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

type CertConfig struct {
	CAFile  string
	CrtFile string
	KeyFile string
}

var cert CertConfig

func SetCertConfig(certConfig CertConfig) {
	cert = certConfig
}

func NewAgentExecuter(host string, timeout int) (*AgentExecutor, error) {
	tlsConfig, err := GetTLSConfig(cert.CrtFile, cert.KeyFile, cert.CAFile)
	if err != nil {
		return nil, err
	}
	host = utils.GetHttpsHost(host)
	c := httpclient.NewObjectClient(host, httpclient.WithTLSClientConfig(tlsConfig.Clone()))

	agent := &AgentExecutor{
		Host:      host,
		client:    c,
		timeout:   timeout,
		tlsConfig: tlsConfig,
	}
	if err := agent.Health(); err != nil {
		return nil, err
	}
	return agent, nil
}

func (a *AgentExecutor) Health() error {
	var respBody AgentResponse
	err := a.client.Get(HealthPath, &respBody)
	if err != nil {
		return fmt.Errorf("agent %s not heath: %w", a.Host, err)
	}
	if respBody.Code != 0 {
		return fmt.Errorf("get  %s health err ,code: %d, message: %s", a.Host, respBody.Code, respBody.Message)
	}
	return nil
}

func (a *AgentExecutor) SaveAgentConfig(component string, config map[string]any) error {
	var respBody AgentResponse
	err := a.client.Post(AgentConfigPath, map[string]any{
		"component": component,
		"config":    config,
	}, &respBody)
	if err != nil {
		return fmt.Errorf("agent %s get config err: %w", a.Host, err)
	}
	if respBody.Code != 0 {
		return fmt.Errorf("get  %s config err ,code: %d, message: %s", a.Host, respBody.Code, respBody.Message)
	}
	return nil
}

func (a *AgentExecutor) GetAgentConfigg(component string) (map[string]any, error) {
	var respBody AgentResponse
	err := a.client.Get(fmt.Sprintf("%s?component=%s", AgentConfigPath, component), &respBody)
	if err != nil {
		return nil, fmt.Errorf("agent %s get config err: %w", a.Host, err)
	}
	if respBody.Code != 0 {
		return nil, fmt.Errorf("get  %s config err ,code: %d, message: %s", a.Host, respBody.Code, respBody.Message)
	}
	return respBody.Data.Config, nil
}

func (a *AgentExecutor) Shell(cmd string, sudo bool) (stdout string, stderr string, err error) {
	var respBody AgentResponse
	if sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}
	err = a.client.Post(CmdExecutePath, map[string]any{
		"command": cmd,
		"timeout": a.timeout,
	}, &respBody)
	if err != nil {
		return "", "", fmt.Errorf("agent failed to post request: %s", err)
	}
	if respBody.Code != 0 {
		return "", "", fmt.Errorf("code: %d, err: %s, stderr: %s", respBody.Code, respBody.Data.Err, respBody.Data.Stderr)
	}
	if respBody.Data.Err != "" {
		return respBody.Data.Stdout, respBody.Data.Stderr, fmt.Errorf(respBody.Data.Err)
	}

	return respBody.Data.Stdout, respBody.Data.Stderr, nil
}

func (a *AgentExecutor) Upload(src string, dest string) (stdout string, stderr string, err error) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: a.tlsConfig,
		},
	}
	file, err := os.Open(src)
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(src))
	if err != nil {
		return "", "", fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy file to multipart writer: %w", err)
	}

	err = writer.WriteField("path", dest)
	if err != nil {
		return "", "", fmt.Errorf("failed to write uploadPath field: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", "", fmt.Errorf("failed to close multipart writer: %w", err)
	}
	resp, err := client.Post(fmt.Sprintf("%s%s", a.Host, UploadPath), writer.FormDataContentType(), body)
	if err != nil {
		return "", "", fmt.Errorf("failed to post request: %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	var respBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode response body: %w", err)
	}

	code, ok := respBody["code"].(float64)
	if !ok {
		return "", "", fmt.Errorf("unexpected type for code: %T", respBody["code"])
	}
	if code != 0 {
		errorMsg, ok := respBody["message"].(string)
		if !ok {
			return "", "", fmt.Errorf("unexpected type for message: %T", respBody["message"])
		}
		return "", "", fmt.Errorf("error uploading file: %s", errorMsg)
	}

	return "", "", nil
}

func (a *AgentExecutor) Download(src string, dest string) (stdout string, stderr string, err error) {
	return "", "", nil
}

func GetTLSConfig(crtPath string, keyPath string, caPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load client cert and key: %v", err)
	}

	caCert, err := os.ReadFile(caPath)
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
