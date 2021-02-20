package node

//fqName : for filtering
const (
	// MetricNameCPU : CPU
	MetricNameCPU = "node_cpu_seconds_total"
	// MetricNameMemTotal
	MetricNameMemTotal = "node_memory_MemTotal_bytes"
	// MetricNameMemFree
	MetricNameMemFree = "node_memory_MemFree_bytes"
	// MetricNameMemBuffer
	MetricNameMemBuffer = "node_memory_Buffers_bytes"
	// MetricNameMemCache
	MetricNameMemCache = "node_memory_Cached_bytes"
	// MetricNameFileSystemTotal filesystem size
	MetricNameFileSystemTotal = "node_filesystem_size_bytes"
	// MetricNameFileSystemFree filesystem free
	MetricNameFileSystemFree = "node_filesystem_free_bytes"
	// MetricNameNetInfo
	MetricNameNetInfo = "node_network_info"
	// MetricNameNetRecvBytes
	MetricNameNetRecvBytes = "node_network_receive_bytes_total"
	// MetricNameNetRecvPackets
	MetricNameNetRecvPackets = "node_network_receive_packets_total"
	// MetricNameNetTransmitBytes
	MetricNameNetTransmitBytes = "node_network_transmit_bytes_total"
	// MetricNameNetTransmitPackets
	MetricNameNetTransmitPackets = "node_network_transmit_packets_total"

	// MetricNameLoadAverage1M
	MetricNameLoadAverage1M = "node_load1"
	// MetricNameDiskReadTimeTotal
	MetricNameDiskReadTimeTotal = "node_disk_read_time_seconds_total"
	// MetricNameDiskReadCompletedTotal
	MetricNameDiskReadCompletedTotal = "node_disk_reads_completed_total"
	// MetricNameDiskWriteTimeTotal
	MetricNameDiskWriteTimeTotal = "node_disk_write_time_seconds_total"
	// MetricNameDiskWriteCompletedTotal
	MetricNameDiskWriteCompletedTotal = "node_disk_writes_completed_total"
)

const (
	// MetricCPUIdle "idle"
	MetricCPUIdle = "idle"
	// MetricCPUIOWait "iowait"
	MetricCPUIOWait = "iowait"
	// MetricCPUIRQ "irq"
	MetricCPUIRQ = "irq"
	// MetricCPUNice "nice"
	MetricCPUNice = "nice"
	// MetricCPUSoftIRQ "softirq"
	MetricCPUSoftIRQ = "softirq"
	// MetricCPUSteal "steal"
	MetricCPUSteal = "steal"
	// MetricCPUSystem "system"
	MetricCPUSystem = "system"
	// MetricCPUUser "user"
	MetricCPUUser = "user"
)

// MetricCPULabels cpu labels
// {device="/dev/mapper/centos-root",fstype="xfs",mountpoint="/"}
type MetricCPULabels struct {
	CPU  string `yaml:"cpu" json:"cpu"`
	Mode string `yaml:"mode" json:"mode"`
}

// MetricFileSystemLabels filesystem labels
// {device="/dev/mapper/centos-root",fstype="xfs",mountpoint="/"}
type MetricFileSystemLabels struct {
	Device     string `yaml:"device" json:"device"`
	FSType     string `yaml:"fstype" json:"fstype"`
	Mountpoint string `yaml:"mountpoint" json:"mountpoint"`
}

// MetricNetworkLables network labels
// {address="02:42:02:f4:de:63",broadcast="ff:ff:ff:ff:ff:ff",device="docker0",duplex="",ifalias="",operstate="down"} 1
type MetricNetworkLables struct {
	Device string `yaml:"device,omitempty" json:"device,omitempty"`
}

// MetricDiskLabel disk labels
// {device="dm-0"}
type MetricDiskLabel struct {
	Device string `yaml:"device,omitempty" json:"device,omitempty"`
}
