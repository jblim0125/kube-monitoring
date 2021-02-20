package model

// HashListConfig hash map list config
type HashListConfig struct {
	MaxSize  int `yaml:"maxSize"`
	HashSize int `yaml:"hashSize"`
}

// HashConfig hash list config
type HashConfig struct {
	NodeHashList      HashListConfig `yaml:"nodeHashList"`
	ContainerHashList HashListConfig `yaml:"containerHashList"`
}

// MonitoringThInfo monitoring thread config
type MonitoringThInfo struct {
	Period  int `yaml:"period"`
	NumProc int `yaml:"numProc"`
	NumSave int `yaml:"numSave"`
}

// CollectionConfig Information to be collected
type CollectionConfig struct {
	FileSystems map[string]string `yaml:"fileSystems"`
	NetDevices  []string          `yaml:"networkDevices"`
	DiskDevices []string          `yaml:"diskDevices"`
}

// AdditionalConfig 추가적인 설정들
type AdditionalConfig struct {
}

// MonitoringConfig monitoring config
type MonitoringConfig struct {
	ThreadConfig     MonitoringThInfo `yaml:"thread"`
	HashConfig       HashConfig       `yaml:"hash"`
	CollectionConfig CollectionConfig `yaml:"collection"`
	Additional       AdditionalConfig `yaml:"additional"`
}
