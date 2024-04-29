package clients

const (
	CmdExecutePath         = "/api/v1/common/execute"
	S3DownloadPath         = "/api/v1/storage/s3/download"
	S3UploadPath           = "/api/v1/storage/s3/upload"
	HDFSDownloadPath       = "/api/v1/storage/hdfs/download"
	HDFSUploadPath         = "/api/v1/storage/hdfs/upload"
	LocalDownloadPath      = "/api/v1/storage/local/download"
	LocalUploadPath        = "/api/v1/storage/local/upload"
	GetComponentConfigPath = "/api/v1/component-config"
)

type (
	CmdExecuteReq struct {
		Command string `json:"command"`
		Timeout int    `json:"timeout"`
	}

	ExecuteData struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
		Err    string `json:"err"`
	}

	CmdExecuteResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    ExecuteData `json:"data"`
	}

	S3DownloadReq struct {
		Endpoint  string `json:"endpoint"`
		Region    string `json:"region"`
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		Bucket    string `json:"bucket"`
		Path      string `json:"path"`
		LocalPath string `json:"localPath"`
	}

	S3DownloadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	S3UploadReq struct {
		Endpoint  string `json:"endpoint"`
		Region    string `json:"region"`
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		Bucket    string `json:"bucket"`
		Path      string `json:"path"`
		LocalPath string `json:"localPath"`
	}

	S3UploadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	KerberosConfig struct {
		Enable                       bool   `json:"enable"`
		Principle                    string `json:"principle"`
		KeytabFilePath               string `json:"keytabFilePath"`
		ConfigFilePath               string `json:"configFilePath"`
		KerberosServicePrincipleName string `json:"kerberosServicePrincipleName"`
	}

	HDFSDownloadReq struct {
		Address   string         `json:"address"`
		Username  string         `json:"username"`
		Path      string         `json:"path"`
		LocalPath string         `json:"localPath"`
		Kerberos  KerberosConfig `json:"kerberos,optional"`
	}

	HDFSDownloadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	HDFSUploadReq struct {
		Address   string         `json:"address"`
		Username  string         `json:"username"`
		Path      string         `json:"path"`
		LocalPath string         `json:"localPath"`
		Kerberos  KerberosConfig `json:"kerberos,optional"`
	}

	HDFSUploadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	LocalDownloadReq struct {
		Path        string `json:"path"`
		LocalPath   string `json:"localPath"`
		Recursively bool   `json:"recursively"`
	}
	LocalDownloadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	LocalUploadReq struct {
		Path        string `json:"path"`
		LocalPath   string `json:"localPath"`
		Recursively bool   `json:"recursively"`
	}
	LocalUploadResp struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	ComponentConfigData struct {
		Config map[string]interface{} `json:"config"`
	}

	GetComponentConfigResp struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    ComponentConfigData `json:"data"`
	}
)
