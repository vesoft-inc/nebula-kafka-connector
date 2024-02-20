package types

type JobSpec struct {
	Kind        string     `yaml:"kind,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	Rollback    bool       `yaml:"rollback,omitempty"`
	InstallPath string     `yaml:"installPath,omitempty"`
	Spec        ProcessMap `yaml:"spec,omitempty"`
}

type ProcessMap struct {
	Metad          *MetadSpec `yaml:"metad,omitempty"`
	LicenseManager *Process   `yaml:"license-manager,omitempty"`
}

type Process struct {
	Name          string            `yaml:"name,omitempty"`
	PackagePath   string            `yaml:"packagePath,omitempty"`
	Config        map[string]string `yaml:"config,omitempty"`
	Hosts         []Agent           `yaml:"hosts,omitempty"`
	StartType     string            `yaml:"startType,omitempty"`     //systemd, shell
	ExecShellPath string            `yaml:"execShellPath,omitempty"` // exec path
}

type MetadSpec struct {
	Process  `yaml:",inline"`
	Clusters []Cluster `yaml:"clusters,omitempty"`
}

type Agent struct {
	Host string `yaml:"host,omitempty"`
}

type Cluster struct {
	ZoneList  []string `yaml:"zoneList,omitempty"`
	Name      string   `yaml:"name,omitempty"`
	Partition int      `yaml:"partition,omitempty"`
	Replica   int      `yaml:"replica,omitempty"`
	Graphd    Process  `yaml:"graphd,omitempty"`
	Storaged  Process  `yaml:"storaged,omitempty"`
}

type StatusItem struct {
	Product string
	Service string
	Host    string
	Port    string
	Status  string
}
