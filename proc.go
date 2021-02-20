package main

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	queue "github.com/enriquebris/goconcurrentqueue"
	"github.com/sirupsen/logrus"
	"mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	nm "mobigen.com/iris-cloud-monitoring/model/node"
	"mobigen.com/iris-cloud-monitoring/tools"
)

// Proc proc data struct
type Proc struct {
	ID     int
	Log    logrus.FieldLogger
	config model.MonitoringConfig
}

// StartProc node, kube-state-metric data getter
func (m *Monitoring) StartProc(id int) {
	proc := Proc{
		ID:     id,
		Log:    m.logger.WithFields(map[string]interface{}{"proc": id}),
		config: m.config,
	}
	proc.Log.Infof("Proc[ %d ] Started", id)

	for {
		// Get Metrics Data From Queue(sender is gatherer)
		_rcvMsg, err := m.queue[id].Dequeue()
		if err != nil {
			if (err.(*queue.QueueError)).Code() == queue.QueueErrorCodeEmptyQueue {
				time.Sleep(20 * time.Millisecond)
				// proc.Log.Debugf("Proc[ %d ] Q[ %d ] Empty", id, i)
			} else if (err.(*queue.QueueError)).Code() == queue.QueueErrorCodeLockedQueue {
				proc.Log.Debugf("Q[ %d ] Lock", id)
			}
		} else {
			rcvMsg := _rcvMsg.(*model.MetricData)
			if rcvMsg.Type == model.NODE {
				proc.Log.Infof("[ GET >> PROC ] Receive Metric Data. Type[ NODE ] CID[ %s ] Name[ %s ] Time[ %s ]",
					rcvMsg.CID, rcvMsg.Name, tools.PrintTimeFromMilli(rcvMsg.Time))
			} else {
				proc.Log.Infof("[ GET >> PROC ] Receive Metric Data. Type[ KUBE_STATE ] CID[ %s ] Name[ %s ] Time[ %s ]",
					rcvMsg.CID, rcvMsg.Name, tools.PrintTimeFromMilli(rcvMsg.Time))
			}

			// ID, Label, Value 형태로 데이터를 가져 옴.
			descs := proc.MetricDataParseNFilter(rcvMsg)

			if rcvMsg.Type == model.NODE {
				// CPU와 같이 이전 데이터와 비교가 필요한 데이터들의 처리와 저장 규칙(단위)에 맞게 변환하여 저장
				node, err := proc.NodeExporterCalc(rcvMsg, descs)
				if err != nil {
					proc.Log.Errorf("failed to calc that resource usage for node. %s", err.Error())
					m.syncWaitGroup.Done()
					continue
				}
				// Alarm

				// 최종적으로 처리된 시간을 기록
				node.Time = tools.GetMillis()
				proc.Log.WithFields(map[string]interface{}{
					"proc": proc.ID, "cid": node.CID, "name": node.Name,
				}).Infof("[ PROC >> SAVER ] Success. Node Sys Usage Filtering, Parsing, Calc..")
				node.Put()
			} else if rcvMsg.Type == model.KUBESTATE {
				proc.KubeStateMetricCheck(rcvMsg, descs)
			}
			// Save Thread에서 처리할 수 있도록 waitgroup을 처리한다.
			m.syncWaitGroup.Done()
		}
	}
}

// MetricDataParseNFilter node-exporter 데이터 파싱, 필터링
func (proc *Proc) MetricDataParseNFilter(metricData *model.MetricData) []model.MetricDesc {
	Log := proc.Log.WithFields(map[string]interface{}{
		"proc": proc.ID,
		"cid":  metricData.CID,
		"name": metricData.Name,
	})

	var descs []model.MetricDesc
	// node_exporter로 부터 수신한 데이터를 개행 문자로 자름.
	lines := bytes.Split(metricData.Data, []byte("\n"))
	for _, line := range lines {
		// #로 시작하는 라인 무시, 문자가 없는 라인 무시
		if bytes.Index(line, []byte("#")) == 0 || len(bytes.TrimSpace(line)) <= 0 {
			continue
		}
		// Name | Labels | Value 로 하나의 행을 자름
		subStrings := proc.MetricsSplitString(string(line))
		if subStrings == nil || len(subStrings) <= 0 || len(subStrings) > 3 {
			Log.Errorf("error. failed to split[ %s ]", string(line))
			continue
		}
		// 필요한 데이터가 아닌 경우 무시
		if metricData.Type == model.NODE {
			if !proc.NodeExporterFilter(subStrings[0]) {
				continue
			}
		} else if metricData.Type == model.KUBESTATE {
			if !proc.KubeStateMetricFilter(subStrings[0]) {
				continue
			}
		}
		// 추가로 스페이스와 좌우 특수 문자 삭제
		for idx := 0; idx < len(subStrings); idx++ {
			subStrings[idx] = strings.TrimSpace(subStrings[idx])
			subStrings[idx] = strings.TrimLeft(subStrings[idx], "{")
			subStrings[idx] = strings.TrimRight(subStrings[idx], "}")
		}

		var desc model.MetricDesc
		if len(subStrings) > 2 {
			desc.FqName = subStrings[0]
			subStrings[1] = strings.ReplaceAll(subStrings[1], "\"", "")
			if metricData.Type == model.NODE {
				desc.Labels = proc.NodeExporterLabelsDecode(desc.FqName, subStrings[1], Log)
			} else if metricData.Type == model.KUBESTATE {
				desc.Labels = proc.MetricsLabelsParser(subStrings[1])
			}
			fValue, err := strconv.ParseFloat(subStrings[2], 64)
			if err != nil {
				Log.Errorf("error. failed to convert(string => float) Name[ %s ] Value[ %s ]",
					subStrings[0], subStrings[2])
				continue
			}
			desc.Value = fValue
		} else {
			desc.FqName = subStrings[0]
			fValue, err := strconv.ParseFloat(subStrings[1], 64)
			if err != nil {
				Log.Errorf("error. failed to convert(string => float) Name[ %s ] Value[ %s ]",
					subStrings[0], subStrings[1])
			}
			desc.Value = fValue
		}
		descs = append(descs, desc)
	}
	Log.Debugf("Success Parsing. CID[ %s ] Name[ %s ]", metricData.CID, metricData.Name)
	return descs
}

// MetricsSplitString string split
func (proc *Proc) MetricsSplitString(line string) []string {
	regex := regexp.MustCompile(`(^\w+)|({.*})|(\s.+)`)
	return regex.FindAllString(line, -1)
}

// NodeExporterFilter 노드 상태 감시를 위해 필요한 키 데이터
func (proc *Proc) NodeExporterFilter(fqName string) bool {
	switch fqName {
	case nm.MetricNameCPU:
	case nm.MetricNameMemTotal:
	case nm.MetricNameMemFree:
	case nm.MetricNameMemBuffer:
	case nm.MetricNameMemCache:
	case nm.MetricNameFileSystemTotal:
	case nm.MetricNameFileSystemFree:
	case nm.MetricNameNetInfo:
	case nm.MetricNameNetRecvBytes:
	case nm.MetricNameNetRecvPackets:
	case nm.MetricNameNetTransmitBytes:
	case nm.MetricNameNetTransmitPackets:
	case nm.MetricNameLoadAverage1M:
	case nm.MetricNameDiskReadTimeTotal:
	case nm.MetricNameDiskReadCompletedTotal:
	case nm.MetricNameDiskWriteTimeTotal:
	case nm.MetricNameDiskWriteCompletedTotal:
	default:
		return false
	}
	return true
}

// KubeStateMetricFilter 노드 상태 감시를 위해 필요한 키 데이터
func (proc *Proc) KubeStateMetricFilter(fqName string) bool {
	switch fqName {
	// Node
	case ksm.KsmFqNameNodeInfo:
	case ksm.KsmFqNameNodeCreated:
	case ksm.KsmFqNameNodeLabels:
	case ksm.KsmFqNameNodeRole:
	case ksm.KsmFqNameNodeSpecUnschedulable:
	case ksm.KsmFqNameNodeSpecTaint:
	case ksm.KsmFqNameNodeStatusCondition:
	case ksm.KsmFqNameNodeStatusPhase:
	case ksm.KsmFqNameNodeStatusCapacityPods:
	case ksm.KsmFqNameNodeStatusCapacityCPU:
	case ksm.KsmFqNameNodeStatusCapacityMemory:
	case ksm.KsmFqNameNodeStatusAllocatablePods:
	case ksm.KsmFqNameNodeStatusAllocatableCPU:
	case ksm.KsmFqNameNodeStatusAllocatableMemory:
	// Namespace
	case ksm.KsmFqNameNamespaceCreated:
	case ksm.KsmFqNameNamespaceLabels:
	case ksm.KsmFqNameNamespaceStatusPhase:
	// case ksm.KsmFqNameNamespaceStatusCondition:
	// Configmap
	case ksm.KsmFqNameConfigMapInfo:
	case ksm.KsmFqNameConfigMapCreated:
	case ksm.KsmFqNameConfigMapVersion:
	// Cronjob
	case ksm.KsmFqNameCronjobLabels:
	case ksm.KsmFqNameCronjobInfo:
	case ksm.KsmFqNameCronjobCreated:
	case ksm.KsmFqNameCronjobStatus:
	case ksm.KsmFqNameCronjobLastScheduleTime:
	case ksm.KsmFqNameCronjobSpecSuspend:
	case ksm.KsmFqNameCronjobSpecStartingDeadlineSeconds:
	case ksm.KsmFqNameCronjobNextSchedule:
	// Daemonset
	case ksm.KsmFqNameDaemonsetCreated:
	case ksm.KsmFqNameDaemonsetStatusCurrent:
	case ksm.KsmFqNameDaemonsetStatusDesired:
	case ksm.KsmFqNameDaemonsetStatusAvailable:
	case ksm.KsmFqNameDaemonsetStatusMisscheduled:
	case ksm.KsmFqNameDaemonsetStatusReady:
	case ksm.KsmFqNameDaemonsetStatusUnavailable:
	case ksm.KsmFqNameDaemonsetUpdated:
	case ksm.KsmFqNameDaemonsetMetadataGeneration:
	case ksm.KsmFqNameDaemonsetLabels:
	default:
		return false
	}
	return true
}

// NodeExporterLabelsDecode node-exporter의 labels decode
func (proc *Proc) NodeExporterLabelsDecode(fqName, labelStr string, Log logrus.FieldLogger) interface{} {
	switch fqName {
	case nm.MetricNameCPU:
		var labels nm.MetricCPULabels
		kv := proc.MetricsLabelsParser(labelStr)
		labels.CPU = kv["cpu"]
		labels.Mode = kv["mode"]
		// Log.Debugf("CPU: %+v", labels)
		return &labels
	case nm.MetricNameMemTotal:
	case nm.MetricNameMemFree:
	case nm.MetricNameMemBuffer:
	case nm.MetricNameMemCache:
	case nm.MetricNameFileSystemTotal, nm.MetricNameFileSystemFree:
		var labels nm.MetricFileSystemLabels
		kv := proc.MetricsLabelsParser(labelStr)
		labels.Mountpoint = kv["mountpoint"]
		labels.Device = kv["device"]
		labels.FSType = kv["fstype"]
		// Log.Debugf("FileSystemLables : %+v", labels)
		return &labels
	case nm.MetricNameNetInfo, nm.MetricNameNetRecvBytes, nm.MetricNameNetRecvPackets,
		nm.MetricNameNetTransmitBytes, nm.MetricNameNetTransmitPackets:
		var labels nm.MetricNetworkLables
		kv := proc.MetricsLabelsParser(labelStr)
		labels.Device = kv["device"]
		// Log.Debugf("FileSystemLables : %+v", labels)
		return &labels
	case nm.MetricNameLoadAverage1M:
	case nm.MetricNameDiskReadTimeTotal, nm.MetricNameDiskReadCompletedTotal,
		nm.MetricNameDiskWriteTimeTotal, nm.MetricNameDiskWriteCompletedTotal:
		var labels nm.MetricDiskLabel
		kv := proc.MetricsLabelsParser(labelStr)
		labels.Device = kv["device"]
		// Log.Debugf("DiskLables : %+v", labels)
		return &labels
	}
	return nil
}

// MetricsLabelsParser key=value, key=value
func (proc *Proc) MetricsLabelsParser(labelStr string) map[string]string {
	var labels map[string]string = make(map[string]string)
	subStrs := strings.Split(labelStr, ",")
	for _, subStr := range subStrs {
		// proc.Log.Debugf("labels splitString [ %s ]", subStr)
		kv := strings.Split(subStr, "=")
		// proc.Log.Debugf("key[ %s ] value[ %s ]", kv[0], kv[1])
		labels[kv[0]] = kv[1]
	}
	return labels
}

// NodeExporterCalc cpu, mem, filesystem, network 데이터 기반 계산
func (proc *Proc) NodeExporterCalc(metricData *model.MetricData, descs []model.MetricDesc) (*nm.Node, error) {
	Log := proc.Log.WithFields(map[string]interface{}{
		"proc": proc.ID,
		"cid":  metricData.CID,
		"name": metricData.Name,
	})
	nodeManager := nm.GetInstance()
	node, err := nodeManager.Get(metricData.CID, metricData.Name)
	if err != nil {
		Log.Errorf("error. [ %s ]", err.Error())
		return nil, fmt.Errorf("error. [ %s ]", err.Error())
	}
	for _, desc := range descs {
		switch desc.FqName {
		case nm.MetricNameCPU:
			label := desc.Labels.(*nm.MetricCPULabels)
			proc.NodeCPUUpdate(node, &desc, label, Log)
		case nm.MetricNameMemTotal:
			node.Mem.Total = uint64(desc.Value / 1000)
		case nm.MetricNameMemFree:
			node.Mem.Free = uint64(desc.Value / 1000)
		case nm.MetricNameMemBuffer:
			node.Mem.Buffer = uint64(desc.Value / 1000)
		case nm.MetricNameMemCache:
			node.Mem.Cache = uint64(desc.Value / 1000)
		case nm.MetricNameFileSystemTotal, nm.MetricNameFileSystemFree:
			label := desc.Labels.(*nm.MetricFileSystemLabels)
			proc.NodeFileSystemUpdate(node, &desc, label, Log)
		case nm.MetricNameNetInfo, nm.MetricNameNetRecvBytes, nm.MetricNameNetRecvPackets,
			nm.MetricNameNetTransmitBytes, nm.MetricNameNetTransmitPackets:
			label := desc.Labels.(*nm.MetricNetworkLables)
			proc.NodeNetworkUpdate(node, &desc, label, Log)
		case nm.MetricNameLoadAverage1M:
			proc.NodeLoadAverageUpdate(node, &desc)
		case nm.MetricNameDiskReadTimeTotal, nm.MetricNameDiskReadCompletedTotal,
			nm.MetricNameDiskWriteTimeTotal, nm.MetricNameDiskWriteCompletedTotal:
			label := desc.Labels.(*nm.MetricDiskLabel)
			proc.NodeDiskUpdate(node, &desc, label, Log)
		}
	}
	// CPU
	var tCPUTotal, tUs, tSy, tIo, tIdle float32
	for idx, cpu := range node.CPU.Cores {
		if cpu.Calc.User == 0 && cpu.Calc.System == 0 && cpu.Calc.Idle == 0 {
			continue
		}
		total := cpu.Calc.User + cpu.Calc.Nice + cpu.Calc.System +
			cpu.Calc.Irq + cpu.Calc.Softirq + cpu.Calc.Iowait + cpu.Calc.Idle
		cpu.Persent.User = cpu.Calc.User / total * 100
		cpu.Persent.Nice = cpu.Calc.Nice / total * 100
		cpu.Persent.System = cpu.Calc.System / total * 100
		cpu.Persent.Irq = cpu.Calc.Irq / total * 100
		cpu.Persent.Softirq = cpu.Calc.Softirq / total * 100
		cpu.Persent.Idle = cpu.Calc.Idle / total * 100
		cpu.Persent.Iowait = cpu.Calc.Iowait / total * 100
		cpu.Persent.Steal = cpu.Calc.Steal / total * 100
		node.CPU.Cores[idx] = cpu

		tCPUTotal += 100
		tUs += float32(cpu.Persent.User + cpu.Persent.Nice)
		tSy += float32(cpu.Persent.System + cpu.Persent.Softirq + cpu.Persent.Irq + cpu.Persent.Steal)
		tIo += float32(cpu.Persent.Iowait)
		tIdle += float32(cpu.Persent.Idle)
	}

	node.CPU.User = (tUs / tCPUTotal) * 100
	node.CPU.Sys = (tSy / tCPUTotal) * 100
	node.CPU.IO = (tIo / tCPUTotal) * 100
	node.CPU.Idle = (tIdle / tCPUTotal) * 100

	proc.NodeCPUDebuging(node, Log)

	// Memory
	node.Mem.Used = node.Mem.Total - node.Mem.Free
	Log.Infof("MEM Total[ %d ]GiB Free[ %d ]GiB Use[ %d ]GiB",
		node.Mem.Total/1000/1000, node.Mem.Free/1000/1000, node.Mem.Used/1000/1000)
	// FileSystem
	var tTotal, tUsed, tAvail uint64
	for i, v := range node.File.Devices {
		node.File.Devices[i].Used = v.Total - v.Avail
		tTotal += v.Total
		tAvail += v.Avail
		tUsed += node.File.Devices[i].Used
	}
	node.File.Total = tTotal
	node.File.Used = tUsed
	node.File.Avail = tAvail
	Log.Infof("FS Total[ %d ]GiB Used[ %d ]GiB Avail[ %d ]GiB",
		node.File.Total/1000/1000, node.File.Used/1000/1000, node.File.Avail/1000/1000)
	// Network
	var tRxbps, tTxbps uint64
	for i, v := range node.Net.Devices {
		node.Net.Devices[i].Rxbps = (v.RxBytes * 8) / 5
		node.Net.Devices[i].Txbps = (v.TxBytes * 8) / 5
		Log.Infof("Net Name[ %s ] Rx/Tx[ %d / %d ]kbps Rx/Tx[ %d / %d ]KiB Rx/Tx[ %d / %d ]packets",
			v.Name, node.Net.Devices[i].Rxbps/1000, node.Net.Devices[i].Txbps/1000,
			v.RxBytes, v.TxBytes, v.RxPackets, v.TxPackets)
		tRxbps += node.Net.Devices[i].Rxbps / 1000
		tTxbps += node.Net.Devices[i].Txbps / 1000
	}
	node.Net.Rxbps = tRxbps
	node.Net.Txbps = tTxbps
	Log.Infof("Net Total Rx/Tx[ %d / %d ]kbps",
		node.Net.Rxbps, node.Net.Txbps)

	// Load Average
	Log.Infof("Load Average 1M [ %f ]", node.LoadAverage.LoadAverage1M)
	// Disk Read/Write Latency
	for i, v := range node.Disk.Devices {
		node.Disk.Devices[i].Latency.Read = (v.Second.ReadTimeTotal / v.Second.ReadCompletedTotal) * 100
		node.Disk.Devices[i].Latency.Write = (v.Second.WriteTimeTotal / v.Second.WriteCompletedTotal) * 100
		Log.Infof("Disk Device[ %s ] Read/Write Latency[ %f ]/[ %f ]",
			node.Disk.Devices[i].Device,
			node.Disk.Devices[i].Latency.Read, node.Disk.Devices[i].Latency.Write)
	}

	return node, nil
}

// NodeCPUUpdate node CPU usage update
func (proc *Proc) NodeCPUUpdate(node *nm.Node, desc *model.MetricDesc, label *nm.MetricCPULabels, Log logrus.FieldLogger) error {
	if cpu, ok := node.CPU.Cores[label.CPU]; ok {
		// Update
		// Debuging
		Log.Debugf("Update CPU[ %-2s / %-7s ]", label.CPU, label.Mode)
		switch label.Mode {
		case nm.MetricCPUIdle:
			// calc
			if cpu.Prev.Idle > 0 {
				cpu.Calc.Idle = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Idle))
			}
			// save
			cpu.Prev.Idle = desc.Value
		case nm.MetricCPUIOWait:
			if cpu.Prev.Iowait > 0 {
				cpu.Calc.Iowait = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Iowait))
			}
			cpu.Prev.Iowait = desc.Value
		case nm.MetricCPUIRQ:
			if cpu.Prev.Irq > 0 {
				cpu.Calc.Irq = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Irq))
			}
			cpu.Prev.Irq = desc.Value
		case nm.MetricCPUNice:
			if cpu.Prev.Nice > 0 {
				cpu.Calc.Nice = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Nice))
			}
			cpu.Prev.Nice = desc.Value
		case nm.MetricCPUSoftIRQ:
			if cpu.Prev.Softirq > 0 {
				cpu.Calc.Softirq = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Softirq))
			}
			cpu.Prev.Softirq = desc.Value
		case nm.MetricCPUSteal:
			if cpu.Prev.Steal > 0 {
				cpu.Calc.Steal = math.Min(100, math.Max(0, desc.Value-cpu.Prev.Steal))
			}
			cpu.Prev.Steal = desc.Value
		case nm.MetricCPUSystem:
			if cpu.Prev.System > 0 {
				cpu.Calc.System = math.Min(100, math.Max(0, desc.Value-cpu.Prev.System))
			}
			cpu.Prev.System = desc.Value
		case nm.MetricCPUUser:
			if cpu.Prev.User > 0 {
				cpu.Calc.User = math.Min(100, math.Max(0, desc.Value-cpu.Prev.User))
			}
			cpu.Prev.User = desc.Value
		}
		node.CPU.Cores[label.CPU] = cpu
	} else {
		// Add
		// Log.Debugf("CPU Add[ %s ][ %s ][ %f ]", label.CPU, label.Mode, desc.Value)
		cpu := nm.Cores{}
		switch label.Mode {
		case nm.MetricCPUIdle:
			cpu.Prev.Idle = desc.Value
		case nm.MetricCPUIOWait:
			cpu.Prev.Iowait = desc.Value
		case nm.MetricCPUIRQ:
			cpu.Prev.Irq = desc.Value
		case nm.MetricCPUNice:
			cpu.Prev.Nice = desc.Value
		case nm.MetricCPUSoftIRQ:
			cpu.Prev.Softirq = desc.Value
		case nm.MetricCPUSteal:
			cpu.Prev.Steal = desc.Value
		case nm.MetricCPUSystem:
			cpu.Prev.System = desc.Value
		case nm.MetricCPUUser:
			cpu.Prev.User = desc.Value
		}
		if len(node.CPU.Cores) <= 0 {
			node.CPU.Cores = make(map[string]nm.Cores)
		}
		node.CPU.Cores[label.CPU] = cpu
		// 	Log.Debugf("CPU[ %s ] %+v", label.CPU, cpu)
	}
	return nil
}

// NodeFileSystemUpdate node file system usage update
func (proc *Proc) NodeFileSystemUpdate(node *nm.Node, desc *model.MetricDesc, label *nm.MetricFileSystemLabels, Log logrus.FieldLogger) error {
	isExist := false
	for idx, fs := range node.File.Devices {
		if fs.Device == label.Device && fs.Mount == label.Mountpoint {
			Log.Debugf("Update FileSystem Device[ %s ] MountPoint[ %s ]", label.Device, label.Mountpoint)
			isExist = true
			if desc.FqName == nm.MetricNameFileSystemTotal {
				node.File.Devices[idx].Total = uint64(desc.Value / 1000)
			} else if desc.FqName == nm.MetricNameFileSystemFree {
				node.File.Devices[idx].Avail = uint64(desc.Value / 1000)
			}
			break
		}
	}
	if !isExist {
		if val, ok := proc.config.CollectionConfig.FileSystems[label.Device]; ok {
			if val == label.Mountpoint {
				Log.Debugf("Add New FileSystem Device[ %s ] MountPoint[ %s ]", label.Device, label.Mountpoint)
				fs := nm.FileDevice{
					Device: label.Device,
					Mount:  label.Mountpoint,
				}
				if desc.FqName == nm.MetricNameFileSystemTotal {
					fs.Total = uint64(desc.Value / 1000)
				} else if desc.FqName == nm.MetricNameFileSystemFree {
					fs.Avail = uint64(desc.Value / 1000)
				}
				node.File.Devices = append(node.File.Devices, fs)
			}
		}
	}
	return nil
}

// NodeNetworkUpdate node network usage update
func (proc *Proc) NodeNetworkUpdate(node *nm.Node, desc *model.MetricDesc, label *nm.MetricNetworkLables, Log logrus.FieldLogger) error {

	isExist := false
	for idx, net := range node.Net.Devices {
		if net.Name == label.Device {
			Log.Debugf("Update Network Device[ %s ] Type[ %s ]", label.Device, desc.FqName)
			isExist = true
			switch desc.FqName {
			case nm.MetricNameNetInfo:
			case nm.MetricNameNetRecvBytes:
				// calc
				if net.PrevNetStatus.RxBytes != 0 {
					if net.PrevNetStatus.RxBytes > uint64(desc.Value) {
						// ( UINT64_MAX - prev ) + cur;
						node.Net.Devices[idx].RxBytes = (math.MaxUint64 - net.PrevNetStatus.RxBytes) + uint64(desc.Value)
					} else {
						// cur - prev;
						node.Net.Devices[idx].RxBytes = uint64(desc.Value) - net.PrevNetStatus.RxBytes
					}
				}
				// Save
				node.Net.Devices[idx].PrevNetStatus.RxBytes = uint64(desc.Value)
			case nm.MetricNameNetRecvPackets:
				if net.PrevNetStatus.RxPackets != 0 {
					if net.PrevNetStatus.RxPackets > uint64(desc.Value) {
						node.Net.Devices[idx].RxPackets = (math.MaxUint64 - net.PrevNetStatus.RxPackets) + uint64(desc.Value)
					} else {
						node.Net.Devices[idx].RxPackets = uint64(desc.Value) - net.PrevNetStatus.RxPackets
					}
				}
				// Save
				node.Net.Devices[idx].PrevNetStatus.RxPackets = uint64(desc.Value)
			case nm.MetricNameNetTransmitBytes:
				if net.PrevNetStatus.TxBytes != 0 {
					if net.PrevNetStatus.TxBytes > uint64(desc.Value) {
						node.Net.Devices[idx].TxBytes = (math.MaxUint64 - net.PrevNetStatus.TxBytes) + uint64(desc.Value)
					} else {
						node.Net.Devices[idx].TxBytes = uint64(desc.Value) - net.PrevNetStatus.TxBytes
					}
				}
				// Save
				node.Net.Devices[idx].PrevNetStatus.TxBytes = uint64(desc.Value)
			case nm.MetricNameNetTransmitPackets:
				if net.PrevNetStatus.TxPackets != 0 {
					if net.PrevNetStatus.TxPackets > uint64(desc.Value) {
						node.Net.Devices[idx].TxPackets = (math.MaxUint64 - net.PrevNetStatus.TxPackets) + uint64(desc.Value)
					} else {
						node.Net.Devices[idx].TxPackets = uint64(desc.Value) - net.PrevNetStatus.TxPackets
					}
				}
				// Save
				node.Net.Devices[idx].PrevNetStatus.TxPackets = uint64(desc.Value)
			}
			break
		}
	}
	if !isExist {
		isExist := false
		// 수집 대상인지 확인
		for _, name := range proc.config.CollectionConfig.NetDevices {
			if name == label.Device {
				isExist = true
				break
			}
		}
		if isExist {
			Log.Debugf("Add New Network Device[ %s ]", label.Device)
			nd := nm.NetDevice{
				Name:   label.Device,
				Status: "up",
			}
			switch desc.FqName {
			case nm.MetricNameNetInfo:
			case nm.MetricNameNetRecvBytes:
				nd.PrevNetStatus.RxBytes = uint64(desc.Value / 1000)
			case nm.MetricNameNetRecvPackets:
				nd.PrevNetStatus.RxPackets = uint64(desc.Value / 1000)
			case nm.MetricNameNetTransmitBytes:
				nd.PrevNetStatus.TxBytes = uint64(desc.Value / 1000)
			case nm.MetricNameNetTransmitPackets:
				nd.PrevNetStatus.TxPackets = uint64(desc.Value / 1000)
			}
			node.Net.Devices = append(node.Net.Devices, nd)
		}
	}
	return nil
}

// NodeLoadAverageUpdate node load average 1m update
func (proc *Proc) NodeLoadAverageUpdate(node *nm.Node, desc *model.MetricDesc) error {
	node.LoadAverage.LoadAverage1M = desc.Value
	return nil
}

// NodeDiskUpdate node disk latency update
func (proc *Proc) NodeDiskUpdate(node *nm.Node, desc *model.MetricDesc, label *nm.MetricDiskLabel, Log logrus.FieldLogger) error {

	isExist := false
	for idx, disk := range node.Disk.Devices {
		if disk.Device == label.Device {
			Log.Debugf("Update Disk Device[ %s ] Type[ %s ]", label.Device, desc.FqName)
			isExist = true
			switch desc.FqName {
			case nm.MetricNameDiskReadTimeTotal:
				node.Disk.Devices[idx].Second.ReadTimeTotal = desc.Value
			case nm.MetricNameDiskReadCompletedTotal:
				node.Disk.Devices[idx].Second.ReadCompletedTotal = desc.Value
			case nm.MetricNameDiskWriteTimeTotal:
				node.Disk.Devices[idx].Second.WriteTimeTotal = desc.Value
			case nm.MetricNameDiskWriteCompletedTotal:
				node.Disk.Devices[idx].Second.WriteCompletedTotal = desc.Value
			}
			break
		}
	}
	if !isExist {
		// 수집 대상인지 확인
		isExist := false
		for _, name := range proc.config.CollectionConfig.DiskDevices {
			if name == label.Device {
				isExist = true
				break
			}
		}
		if isExist {
			Log.Debugf("Add New Disk Device[ %s ]", label.Device)
			dd := nm.DiskDevice{
				Device: label.Device,
			}
			switch desc.FqName {
			case nm.MetricNameDiskReadTimeTotal:
				dd.Second.ReadTimeTotal = desc.Value
			case nm.MetricNameDiskReadCompletedTotal:
				dd.Second.ReadCompletedTotal = desc.Value
			case nm.MetricNameDiskWriteTimeTotal:
				dd.Second.WriteTimeTotal = desc.Value
			case nm.MetricNameDiskWriteCompletedTotal:
				dd.Second.WriteCompletedTotal = desc.Value
			}
			node.Disk.Devices = append(node.Disk.Devices, dd)
		}
	}
	return nil
}

// NodeCPUDebuging cpu debuging
func (proc *Proc) NodeCPUDebuging(node *nm.Node, Log logrus.FieldLogger) {
	// Sort
	CPUIDs := []string{}
	for id := range node.CPU.Cores {
		CPUIDs = append(CPUIDs, id)
	}
	sort.Slice(CPUIDs, func(i, j int) bool {
		numA, _ := strconv.Atoi(CPUIDs[i])
		numB, _ := strconv.Atoi(CPUIDs[j])
		return numA < numB
	})

	// Print
	for _, id := range CPUIDs {
		if cpu, ok := node.CPU.Cores[id]; ok {
			if cpu.Persent.User == 0 && cpu.Persent.System == 0 && cpu.Persent.Nice == 0 && cpu.Persent.Idle == 0 {
				continue
			}
			Log.Debugf("CPU[ %2s ] us[ %5.2f ] sy[ %5.2f ] ni[ %5.2f ] id[ %5.2f ] wa[ %5.2f ] hi[ %5.2f ] si[ %5.2f ] st[ %5.2f ]",
				id, cpu.Persent.User, cpu.Persent.System, cpu.Persent.Nice, cpu.Persent.Idle,
				cpu.Persent.Iowait, cpu.Persent.Irq, cpu.Persent.Softirq, cpu.Persent.Steal)
		}
	}
	// Total
	if !math.IsNaN(float64(node.CPU.User)) && !math.IsNaN(float64(node.CPU.Sys)) &&
		!math.IsNaN(float64(node.CPU.Idle)) && !math.IsNaN(float64(node.CPU.IO)) {
		Log.Infof("CPU Total us[ %5.2f ] sy[ %5.2f ] id[ %5.2f ] io[ %5.2f ]",
			node.CPU.User, node.CPU.Sys, node.CPU.Idle, node.CPU.IO)
	}
}
