package source

type (
	Config struct {
		Path string      `yaml:"path"` // only use for local file
		CSV  *CSVConfig  `yaml:"csv,omitempty"`
		S3   *S3Config   `yaml:"s3,omitempty"`
		OSS  *OSSConfig  `yaml:"oss,omitempty"`
		FTP  *FTPConfig  `yaml:"ftp,omitempty"`
		SFTP *SFTPConfig `yaml:"sftp,omitempty"`
		HDFS *HDFSConfig `yaml:"hdfs,omitempty"`
	}

	CSVConfig struct {
		Delimiter  string `yaml:"delimiter,omitempty"`
		WithHeader bool   `yaml:"withHeader,omitempty"`
	}

	S3Config struct {
		Endpoint  string `yaml:"endpoint,omitempty"`
		Region    string `yaml:"region,omitempty"`
		AccessKey string `yaml:"accessKey,omitempty"`
		SecretKey string `yaml:"secretKey,omitempty"`
		Token     string `yaml:"token,omitempty"`
		Bucket    string `yaml:"bucket,omitempty"`
		Key       string `yaml:"key,omitempty"`
	}

	OSSConfig struct {
		Endpoint  string `yaml:"endpoint,omitempty"`
		AccessKey string `yaml:"accessKey,omitempty"`
		SecretKey string `yaml:"secretKey,omitempty"`
		Bucket    string `yaml:"bucket,omitempty"`
		Key       string `yaml:"key,omitempty"`
	}

	FTPConfig struct {
		Host     string `yaml:"host,omitempty"`
		Port     int    `yaml:"port,omitempty"`
		Username string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
		Path     string `yaml:"path,omitempty"`
	}

	SFTPConfig struct {
		Host       string `yaml:"host,omitempty"`
		Port       int    `yaml:"port,omitempty"`
		SSHUser    string `yaml:"sshUser,omitempty"`
		SSHPwd     string `yaml:"sshPwd,omitempty"`
		SSHKey     string `yaml:"sshKey,omitempty"`
		Passphrase string `yaml:"passphrase,omitempty"`
		Path       string `yaml:"path,omitempty"`
	}

	HDFSConfig struct {
		Host string `yaml:"host,omitempty"`
		Port int    `yaml:"port,omitempty"`
		Path string `yaml:"path,omitempty"`
	}
)

func (c *Config) String() string {
	return c.Path
}
