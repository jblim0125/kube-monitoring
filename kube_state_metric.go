package main

import (
	"fmt"
	"strings"

	logrus "github.com/sirupsen/logrus"
	"mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	"mobigen.com/iris-cloud-monitoring/tools"
)

func (proc *Proc) KubeStateMetricCheck(metricData *model.MetricData, descs []model.MetricDesc) error {
	Log := proc.Log.WithFields(map[string]interface{}{
		"proc": proc.ID,
		"cid":  metricData.CID,
		"name": metricData.Name,
	})
	ksmManager := ksm.GetInstance()
	for _, desc := range descs {
		// Log.Debugf("FqName[ %s ] Labels[ %+v ] Value[ %f ]", desc.FqName, desc.Labels, desc.Value)
		resType := proc.GetTypeOfFqName(desc.FqName)
		// Log.Debugf("Type[ %s ]", resType)

		switch resType {
		case "configmap":
			proc.KsmUpdateConfigmap(ksmManager, metricData, desc, Log)
		case "cronjob":
			proc.KsmUpdateCronjob(ksmManager, metricData, desc, Log)
		case "daemonset":
			proc.KsmUpdateDaemonset(ksmManager, metricData, desc, Log)
		case "deployment":
			proc.KsmUpdateDeployment(ksmManager, metricData, desc, Log)
		case "ingress":
			proc.KsmUpdateIngress(ksmManager, metricData, desc, Log)
		case "job":
			proc.KsmUpdateJob(ksmManager, metricData, desc, Log)
		case "namespace":
			proc.KsmUpdateNamespace(ksmManager, metricData, desc, Log)
		case "node":
			proc.KsmUpdateNode(ksmManager, metricData, desc, Log)
		case "persistentvolume":
			proc.KsmUpdatePersistentVolume(ksmManager, metricData, desc, Log)
		case "persistentvolumeclaim":
			proc.KsmUpdatePersistentVolumeClaim(ksmManager, metricData, desc, Log)
		case "pod":
			proc.KsmUpdatePod(ksmManager, metricData, desc, Log)
		case "replicaset":
			proc.KsmUpdateReplicaset(ksmManager, metricData, desc, Log)
		case "replicationcontroller":
			proc.KsmUpdateReplicationController(ksmManager, metricData, desc, Log)
		case "secret":
			proc.KsmUpdateSecret(ksmManager, metricData, desc, Log)
		case "service":
			proc.KsmUpdateService(ksmManager, metricData, desc, Log)
		case "statefulset":
			proc.KsmUpdateStatefulset(ksmManager, metricData, desc, Log)
		default:
			Log.Errorf("Error. Not Supported FqName[ %s ]", desc.FqName)
		}
	}
	return nil
}

// GetTypeOfFqName fqname split and return resource type(node, namespace, pods)
func (proc *Proc) GetTypeOfFqName(input string) string {
	ret := strings.Split(input, "_")
	return ret[1]
}

// KsmUpdateConfigmap configmap
func (proc *Proc) KsmUpdateConfigmap(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	Log.Debugf("Update Configmap State Metric")
	return nil
}

// KsmUpdateCronjob cronjob
func (proc *Proc) KsmUpdateCronjob(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	Log.Debugf("Update Cronjob State Metric")
	return nil
}

// KsmUpdateDeamonset daemonset
func (proc *Proc) KsmUpdateDaemonset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Daemonset State Metric")
	return nil
}

// KsmUpdateDeployment deployment
func (proc *Proc) KsmUpdateDeployment(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Deployment State Metric")
	return nil
}

// KsmUpdateIngress ingress
func (proc *Proc) KsmUpdateIngress(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Ingress State Metric")
	return nil
}

// KsmUpdateJob job
func (proc *Proc) KsmUpdateJob(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Job State Metric")
	return nil
}

// KsmUpdateNamespace namespace
func (proc *Proc) KsmUpdateNamespace(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Namespace State Metric")
	descLabels := desc.Labels.(map[string]string)
	if nmName, ok := descLabels["namespace"]; ok {
		var nm *ksm.Namespace
		nm, err := ksmManager.GetNamespace(metricData.CID, nmName)
		if err != nil {
			nm, err = ksmManager.AllocNamespace()
			if err != nil {
				return fmt.Errorf("ERROR. Failed to Update Namespace[ %s ]. [ %s ]", nmName, err.Error())
			}
			nm.CID = metricData.CID
			nm.Name = nmName
			nm.Add()
		}
		// Namespace 상태 갱신 시간 업데이트
		nm.RefreshTime = tools.GetMillis()
		switch desc.FqName {
		case ksm.KsmFqNameNamespaceCreated:
			if nm.Created == 0 {
				Log.Debugf("Namespace[ %s ] Created Time[ %s ]", nm.Name, tools.PrintTimeFromSec(int64(desc.Value)))
				nm.Created = int64(desc.Value) * 1000
			} else if nm.Created != int64(desc.Value)*1000 {
				Log.Debugf("Namespace[ %s ] Update Created Time. Before [ %s ] Currnet [ %s ]", nm.Name,
					tools.PrintTimeFromMilli(nm.Created), tools.PrintTimeFromSec(int64(desc.Value)))
				nm.Created = int64(desc.Value) * 1000
			}
		case ksm.KsmFqNameNamespaceLabels:
			if !proc.CompareMapData(nm.Labels, descLabels) {
				Log.Debugf("Namespace[ %s ] Changed Labels.", nm.Name)
				nm.UpdateTime = tools.GetMillis()
				nm.Labels = nil
				nm.Labels = make(map[string]string)
				proc.CopyMapData(nm.Labels, descLabels)
				Log.Debugf("Namespace[ %s ] Update Labels.[ %+v ]", nm.Name, nm.Labels)
			}
		case ksm.KsmFqNameNamespaceStatusPhase:
			phase := descLabels["phase"]
			if int(desc.Value) == 1 {
				if nm.Phase != phase {
					nm.UpdateTime = tools.GetMillis()
					Log.Debugf("Namespace[ %s ] Changed Phase[ %s ].", nm.Name, phase)
					nm.Phase = phase
				}
			}
		// case ksm.KsmFqNameNamespaceStatusCondition:
		default:
			Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
		}
		nm.Put()
	} else {
		return fmt.Errorf("ERROR. Not Found Namespace[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}

	return nil
}

// KsmUpdateNode node
func (proc *Proc) KsmUpdateNode(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	descLabels := desc.Labels.(map[string]string)
	if nodeName, ok := descLabels["node"]; ok {
		var node *ksm.Node
		node, err := ksmManager.GetNode(metricData.CID, nodeName)
		if err != nil {
			node, err = ksmManager.AllocNode()
			if err != nil {
				return fmt.Errorf("ERROR. Failed to Update Node[ %s ]. [ %s ]", nodeName, err.Error())
			}
			node.CID = metricData.CID
			node.Name = nodeName
			node.Add()
		}
		// Node 상태 갱신 시간 업데이트
		node.RefreshTime = tools.GetMillis()
		switch desc.FqName {
		case ksm.KsmFqNameNodeInfo:
			if !proc.CompareMapData(node.NodeInfo, descLabels) {
				Log.Debugf("Node[ %s ] Changed NodeInfo.", node.Name)
				node.UpdateTime = tools.GetMillis()
				node.NodeInfo = nil
				node.NodeInfo = make(map[string]string)
				proc.CopyMapData(node.NodeInfo, descLabels)
				Log.Debugf("Node[ %s ] Update NodeInfo.[ %+v ]", node.Name, node.NodeInfo)
			}
		case ksm.KsmFqNameNodeCreated:
			if node.Created == 0 {
				Log.Debugf("Node[ %s ] Created Time[ %s ]", node.Name, tools.PrintTimeFromSec(int64(desc.Value)))
				node.Created = int64(desc.Value) * 1000
			} else if node.Created != int64(desc.Value)*1000 {
				Log.Debugf("Node[ %s ] Update Created Time. Before [ %s ] Currnet [ %s ]", node.Name,
					tools.PrintTimeFromMilli(node.Created), tools.PrintTimeFromSec(int64(desc.Value)))
				node.Created = int64(desc.Value) * 1000
			}
		case ksm.KsmFqNameNodeLabels:
			if !proc.CompareMapData(node.Labels, descLabels) {
				Log.Debugf("Node[ %s ] Changed Labels.", node.Name)
				node.UpdateTime = tools.GetMillis()
				node.Labels = nil
				node.Labels = make(map[string]string)
				proc.CopyMapData(node.Labels, descLabels)
				Log.Debugf("Node[ %s ] Update Labels.[ %+v ]", node.Name, node.Labels)
			}
		case ksm.KsmFqNameNodeRole:
			if !proc.CompareMapData(node.Roles, descLabels) {
				Log.Debugf("Node[ %s ] Changed Role.", node.Name)
				node.UpdateTime = tools.GetMillis()
				node.Roles = nil
				node.Roles = make(map[string]string)
				proc.CopyMapData(node.Roles, descLabels)
				Log.Debugf("Node[ %s ] Update Role.[ %+v ]", node.Name, node.Roles)
			}
		case ksm.KsmFqNameNodeSpecUnschedulable:
			var unschedulable bool
			if int(desc.Value) == 1 {
				unschedulable = true
			} else {
				unschedulable = false
			}
			if node.Spec.Unschedulable != unschedulable {
				Log.Debugf("Node[ %s ] Changed Unschedulable. To[ %t ]", node.Name, unschedulable)
				node.UpdateTime = tools.GetMillis()
				node.Spec.Unschedulable = unschedulable
			}
		case ksm.KsmFqNameNodeSpecTaint:
			// 데이터를 비교하고 동일한 경우 갱신시간을 업데이트한다.
			taintEquals, err := proc.CompareTaint(node, descLabels)
			if err != nil {
				return fmt.Errorf("ERROR. Can't Compare Node Taints[ %s ]", err.Error())
			}
			// 변경된 부분을 변경하거나 추가한다.
			if !taintEquals {
				Log.Debugf("Node[ %s ] Changed Taint.", node.Name)
				node.UpdateTime = tools.GetMillis()
				proc.UpdateOrInsertTaint(node, descLabels)
				Log.Debugf("Node[ %s ] Update/Insert Taint.[ %+v ]", node.Name, descLabels)
			}
		case ksm.KsmFqNameNodeStatusCondition:
			if proc.UpdateNodeCondition(node, desc, Log) {
				Log.Debugf("Node[ %s ] Changed Condition.[ %+v ]", node.Name, node.Status.Condition)
				// TODO :: Alarm??
				if node.Status.Condition.NetworkUnavailable != ksm.KsmStatusFalse {
					Log.Debugf("Node[ %s ] Network Condition Unnormal[ %s ].",
						node.Name, node.Status.Condition.NetworkUnavailable)
				}
				if node.Status.Condition.MemoryPressure != ksm.KsmStatusFalse {
					Log.Debugf("Node[ %s ] Memory Condition Unnormal[ %s ].",
						node.Name, node.Status.Condition.MemoryPressure)
				}
				if node.Status.Condition.DiskPressure != ksm.KsmStatusFalse {
					Log.Debugf("Node[ %s ] Disk Condition Unnormal[ %s ].",
						node.Name, node.Status.Condition.DiskPressure)
				}
				if node.Status.Condition.PIDPressure != ksm.KsmStatusFalse {
					Log.Debugf("Node[ %s ] PID Condition Unnormal[ %s ].",
						node.Name, node.Status.Condition.PIDPressure)
				}
				if node.Status.Condition.Ready != ksm.KsmStatusTrue {
					Log.Debugf("Node[ %s ] Ready Condition Unnormal[ %s ].",
						node.Name, node.Status.Condition.Ready)
				}
			}
		case ksm.KsmFqNameNodeStatusPhase:
			Log.Debugf("Node[ %s ] Status Phase.[ %+v ][ %f ]", node.Name, descLabels, desc.Value)
		case ksm.KsmFqNameNodeStatusCapacityPods:
			if node.Status.Capacity.Pods != int32(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Capacity Pods.[ %d ]", node.Name, int32(desc.Value))
				node.Status.Capacity.Pods = int32(desc.Value)
			}
		case ksm.KsmFqNameNodeStatusCapacityCPU:
			if node.Status.Capacity.Cores != int32(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Capacity Cores.[ %d ]", node.Name, int32(desc.Value))
				node.Status.Capacity.Cores = int32(desc.Value)
			}
		case ksm.KsmFqNameNodeStatusCapacityMemory:
			if node.Status.Capacity.MemoryBytes != int64(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Capacity Mem Bytes.[ %d ]GiB", node.Name, int64(desc.Value)/1000/1000/1000)
				node.Status.Capacity.MemoryBytes = int64(desc.Value)
			}
		case ksm.KsmFqNameNodeStatusAllocatablePods:
			if node.Status.Allocatable.Pods != int32(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Allocatable Pods.[ %d ]", node.Name, int32(desc.Value))
				node.Status.Allocatable.Pods = int32(desc.Value)
			}
		case ksm.KsmFqNameNodeStatusAllocatableCPU:
			if node.Status.Allocatable.Cores != int32(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Allocatable Cores.[ %d ]", node.Name, int32(desc.Value))
				node.Status.Allocatable.Cores = int32(desc.Value)
			}
		case ksm.KsmFqNameNodeStatusAllocatableMemory:
			if node.Status.Allocatable.MemoryBytes != int64(desc.Value) {
				Log.Debugf("Node[ %s ] Changed Allocatable Mem Bytes.[ %d ]GiB", node.Name, int64(desc.Value)/1000/1000/1000)
				node.Status.Allocatable.MemoryBytes = int64(desc.Value)
			}
		default:
			Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
		}
		node.Put()
	} else {
		return fmt.Errorf("ERROR. Not Found Node Name[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}
	return nil
}

// KsmUpdatePersistentvolume persistentvolume
func (proc *Proc) KsmUpdatePersistentVolume(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	Log.Debugf("Update PV State Metric")
	return nil
}

// KsmUpdatePersistentVolumeClaim persistentvolumeclaim
func (proc *Proc) KsmUpdatePersistentVolumeClaim(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update PVC State Metric")
	return nil
}

// KsmUpdatePod pod
func (proc *Proc) KsmUpdatePod(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Pod State Metric")
	return nil
}

// KsmUpdateReplicaset replicaset
func (proc *Proc) KsmUpdateReplicaset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update RS State Metric")
	return nil
}

// KsmUpdateReplicationController replicationcontroller
func (proc *Proc) KsmUpdateReplicationController(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update RC State Metric")
	return nil
}

// KsmUpdateSecret secret
func (proc *Proc) KsmUpdateSecret(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Secret State Metric")
	return nil
}

// KsmUpdateService service
func (proc *Proc) KsmUpdateService(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Service State Metric")
	return nil
}

// KsmUpdateStatefulset statefulset
func (proc *Proc) KsmUpdateStatefulset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	Log.Debugf("Update Statefulset State Metric")
	return nil
}

// CompareMapData compare map[string]string
func (proc *Proc) CompareMapData(left map[string]string, right map[string]string) bool {
	for lk, lv := range left {
		if rv, ok := right[lk]; ok {
			if lv != rv {
				return false
			}
			continue
		}
		return false
	}
	for rk, rv := range right {
		if rk == "node" {
			continue
		}
		if lv, ok := left[rk]; ok {
			if lv != rv {
				return false
			}
			continue
		}
		return false
	}
	return true
}

// CopyMapData copy map[string]string
func (proc *Proc) CopyMapData(left map[string]string, right map[string]string) {
	for rightKey, rightValue := range right {
		if rightKey == "node" {
			continue
		}
		left[rightKey] = rightValue
	}
}

// CompareTaint compare taint
func (proc *Proc) CompareTaint(node *ksm.Node, descLabels map[string]string) (bool, error) {
	if len(node.Spec.Taints) <= 0 {
		return false, nil
	}
	for idx, taint := range node.Spec.Taints {
		key, ok := descLabels["key"]
		if !ok {
			return false, fmt.Errorf("not found key from kube state metric taint[ %+v ]", descLabels)
		}
		value, ok := descLabels["value"]
		if !ok {
			return false, fmt.Errorf("not found value from kube state metric taint[ %+v ]", descLabels)
		}
		effect, ok := descLabels["effect"]
		if !ok {
			return false, fmt.Errorf("not found effect from kube state metric taint[ %+v ]", descLabels)
		}
		if taint.Key == key && taint.Value == value && taint.Effect == effect {
			node.Spec.Taints[idx].RefreshTime = tools.GetMillis()
			return true, nil
		}
	}
	return false, nil
}

// UpdateOrInsertTaint update or insert node taint
func (proc *Proc) UpdateOrInsertTaint(node *ksm.Node, descLabels map[string]string) {
	// Key가 동일한 경우를 처리
	for idx, taint := range node.Spec.Taints {
		// map에서 키/밸류에 대한 오류 체크는 compare단계에서 이미 진행하여 생략한다.
		key := descLabels["key"]
		value := descLabels["value"]
		effect := descLabels["effect"]
		if taint.Key == key {
			node.Spec.Taints[idx].RefreshTime = tools.GetMillis()
			node.Spec.Taints[idx].Value = value
			node.Spec.Taints[idx].Effect = effect
			return
		}
	}
	// Key가 없는 경우를 처리
	taint := ksm.NodeTaint{
		RefreshTime: tools.GetMillis(),
		Key:         descLabels["key"],
		Value:       descLabels["value"],
		Effect:      descLabels["effect"],
	}
	node.Spec.Taints = append(node.Spec.Taints, taint)
}

// UpdateNodeCondition update node condition
func (proc *Proc) UpdateNodeCondition(node *ksm.Node, desc model.MetricDesc, Log logrus.FieldLogger) bool {
	isChanged := false
	descLabels := desc.Labels.(map[string]string)
	condition := descLabels["condition"]
	status := descLabels["status"]
	if int(desc.Value) == 1 {
		switch condition {
		case ksm.KsmNodeStatusConditionNetworkUnavailable:
			if node.Status.Condition.NetworkUnavailable != status {
				isChanged = true
			}
			node.Status.Condition.NetworkUnavailable = status
		case ksm.KsmNodeStatusConditionMemoryPressure:
			if node.Status.Condition.MemoryPressure != status {
				isChanged = true
			}
			node.Status.Condition.MemoryPressure = status
		case ksm.KsmNodeStatusConditionDiskPressure:
			if node.Status.Condition.DiskPressure != status {
				isChanged = true
			}
			node.Status.Condition.DiskPressure = status
		case ksm.KsmNodeStatusConditionPIDPressure:
			if node.Status.Condition.PIDPressure != status {
				isChanged = true
			}
			node.Status.Condition.PIDPressure = status
		case ksm.KsmNodeStatusConditionReady:
			if node.Status.Condition.Ready != status {
				isChanged = true
			}
			node.Status.Condition.Ready = status
		default:
			Log.Errorf("ERROR. Not Supported Condition Type[ %s ]", condition)
		}
	}
	return isChanged
}
