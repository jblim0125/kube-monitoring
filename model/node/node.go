package node

import (
	"sync"
)

// CoreDetail core info
type CoreDetail struct {
	Idle    float64 `yaml:"idle"`
	Iowait  float64 `yaml:"iowait"`
	Irq     float64 `yaml:"irq"`
	Nice    float64 `yaml:"nice"`
	Softirq float64 `yaml:"softirq"`
	Steal   float64 `yaml:"steal"`
	System  float64 `yaml:"system"`
	User    float64 `yaml:"user"`
}

// Cores cur, prev info
type Cores struct {
	// Persent
	Persent CoreDetail `yaml:"persent"`
	// cur - prev
	Calc CoreDetail `yaml:"calc"`
	// 이전 수집된 값
	Prev CoreDetail `yaml:"prev"`
}

// CPU cpu : sum of core
type CPU struct {
	User  float32          `yaml:"user"`
	Sys   float32          `yaml:"sys"`
	IO    float32          `yaml:"io"`
	Idle  float32          `yaml:"idle"`
	Cores map[string]Cores `yaml:"cores"`
}

// Memory : KiB
type Memory struct {
	Total  uint64 `yaml:"total"`
	Free   uint64 `yaml:"free"`
	Used   uint64 `yaml:"used"`
	Buffer uint64 `yaml:"buffer"`
	Cache  uint64 `yaml:"cache"`
}

// NetDevice info of network device
type NetDevice struct {
	Name          string        `yaml:"name"`
	Status        string        `yaml:"status"`
	Rxbps         uint64        `yaml:"rxkbps"`
	Txbps         uint64        `yaml:"txkbps"`
	RxBytes       uint64        `yaml:"rxBytes"`
	RxPackets     uint64        `yaml:"rxPackets"`
	TxBytes       uint64        `yaml:"txBytes"`
	TxPackets     uint64        `yaml:"txPackets"`
	PrevNetStatus PrevNetStatus `yaml:"prevNetStatus"`
}

// PrevNetStatus prev status
type PrevNetStatus struct {
	RxBytes   uint64 `yaml:"rxBytes"`
	RxPackets uint64 `yaml:"rxPackets"`
	TxBytes   uint64 `yaml:"txBytes"`
	TxPackets uint64 `yaml:"txPackets"`
}

// Network network use
type Network struct {
	Rxbps   uint64      `yaml:"rxkbps"`
	Txbps   uint64      `yaml:"txkbps"`
	Devices []NetDevice `yaml:"devices"`
}

// FileDevice file device usage
type FileDevice struct {
	Device string `yaml:"device"`
	Mount  string `yaml:"mount"`
	Total  uint64 `yaml:"total"`
	Used   uint64 `yaml:"used"`
	Avail  uint64 `yaml:"avail"`
}

// FileSystem file system usage
type FileSystem struct {
	Total   uint64       `yaml:"total"`
	Used    uint64       `yaml:"used"`
	Avail   uint64       `yaml:"avail"`
	Devices []FileDevice `yaml:"devices"`
}

// LoadAverage load average
type LoadAverage struct {
	LoadAverage1M float64 `yaml:"loadAverage1M"`
}

// DiskSecond disk read/write second
type DiskSecond struct {
	ReadTimeTotal       float64 `yaml:"readTimeTotal"`
	ReadCompletedTotal  float64 `yaml:"readComplatedTotal"`
	WriteTimeTotal      float64 `yaml:"writeTimeTotal"`
	WriteCompletedTotal float64 `yaml:"writeComplatedTotal"`
}

// DiskLatency disk latency
type DiskLatency struct {
	Read  float64 `yaml:"read"`
	Write float64 `yaml:"write"`
}

// DiskDevice diskdevice name and info
type DiskDevice struct {
	Device  string      `yaml:"device"`
	Latency DiskLatency `yaml:"latency"`
	Second  DiskSecond  `yaml:"second"`
}

// Disks disk list
type Disks struct {
	Devices []DiskDevice `yaml:"devices"`
}

// Node node info
type Node struct {
	Idx             uint32      `yaml:"idx"`
	key             uint32      `yaml:"hashKey"`
	used            int32       `yaml:"used"`
	cnt             int32       `yaml:"cnt"`
	listIdx         int32       `yaml:"hashListIdx"`
	CID             string      `yaml:"cid"`
	Time            int64       `yaml:"time"`
	Name            string      `yaml:"name"`
	NodeExporterURL string      `yaml:"nodeExporterURL"`
	CPU             CPU         `yaml:"cpu"`
	Mem             Memory      `yaml:"mem"`
	Net             Network     `yaml:"net"`
	File            FileSystem  `yaml:"file"`
	LoadAverage     LoadAverage `yaml:"loadAverage"`
	Disk            Disks       `yaml:"disk"`
}

// Clone deep copy
func (n *Node) Clone() *Node {
	var clone Node
	clone = *n
	clone.CPU.Cores = make(map[string]Cores)
	for k, v := range n.CPU.Cores {
		clone.CPU.Cores[k] = v
	}
	copy(clone.Net.Devices, n.Net.Devices)
	copy(clone.File.Devices, n.File.Devices)
	copy(clone.Disk.Devices, n.Disk.Devices)
	return &clone
}

// Nodes nodes
type Nodes struct {
	List     map[uint32][]*Node
	Lock     []*sync.Mutex
	NodePool []*Node
}
