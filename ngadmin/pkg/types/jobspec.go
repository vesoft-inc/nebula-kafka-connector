package types

type JobSpec struct {
	Kind        string     `yaml:"kind,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	Rollback    bool       `yaml:"rollback,omitempty"`
	InstallPath string     `yaml:"installPath,omitempty"`
	Spec        ProcessMap `yaml:"spec,omitempty"`
}

type ProcessMap struct {
	Metad *MetadSpec `yaml:"metad,omitempty"`
}

type Process struct {
	PackagePath string            `yaml:"packagePath,omitempty"`
	Config      map[string]string `yaml:"config,omitempty"`
	Hosts       []Agent           `yaml:"hosts,omitempty"`
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
