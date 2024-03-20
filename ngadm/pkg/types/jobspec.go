package types

type JobSpec struct {
	Kind           string     `yaml:"kind,omitempty"`
	Version        string     `yaml:"version,omitempty"`
	Rollback       bool       `yaml:"rollback,omitempty"`
	InstallPath    string     `yaml:"installPath,omitempty"`
	CertFile       string     `yaml:"certFile,omitempty"`
	KeyFile        string     `yaml:"keyFile,omitempty"`
	CAFile         string     `yaml:"caFile,omitempty"`
	Spec           ProcessMap `yaml:"spec,omitempty"`
	Info           any        `yaml:"info,omitempty"` //for some thing else just save it
	UtilsProcesses map[string]*Process
}

type ProcessMap struct {
	Metad *MetadSpec `yaml:"metad,omitempty"`
}

type Process struct {
	Name          string         `yaml:"name,omitempty"`
	PackagePath   string         `yaml:"packagePath,omitempty"`
	ConfigPath    string         `yaml:"configPath,omitempty"`
	Config        map[string]any `yaml:"config,omitempty"`
	Hosts         []Agent        `yaml:"hosts,omitempty"`
	StartType     string         `yaml:"startType,omitempty"`     //systemd, shell
	ExecShellPath string         `yaml:"execShellPath,omitempty"` // exec path
	ExecStartPath string         `yaml:"execStartPath,omitempty"` // bin path
	WorkingDir    string         `yaml:"workingDir,omitempty"`    // working dir
}

type MetadSpec struct {
	Process  `yaml:",inline"`
	Clusters []Cluster `yaml:"clusters,omitempty"`
}

type Agent struct {
	Host string `yaml:"host,omitempty"`
}

type Cluster struct {
	ZoneList []string `yaml:"zoneList,omitempty"`
	Name     string   `yaml:"name,omitempty"`
	Replica  int      `yaml:"replica,omitempty"`
	Graphd   Process  `yaml:"graphd,omitempty"`
	Storaged Process  `yaml:"storaged,omitempty"`
}

type StatusItem struct {
	Product string
	Service string
	Host    string
	Port    string
	Status  string
}
