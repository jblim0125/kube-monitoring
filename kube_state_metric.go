package main

import (
	"fmt"
	"strings"

	logrus "github.com/sirupsen/logrus"
	"mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	"mobigen.com/iris-cloud-monitoring/tools"
)

// KubeStateMetricCheck kube state metric data parsing and local data update
func (proc *Proc) KubeStateMetricCheck(metricData *model.MetricData, descs []model.MetricDesc) error {
	Log := proc.Log.WithFields(map[string]interface{}{
		"proc": proc.ID,
		"cid":  metricData.CID,
		"name": metricData.Name,
	})
	var errCnt int = 0
	ksmManager := ksm.GetInstance()
	for _, desc := range descs {
		// Log.Debugf("FqName[ %s ] Labels[ %+v ] Value[ %f ]", desc.FqName, desc.Labels, desc.Value)
		resType := proc.GetTypeOfFqName(desc.FqName)
		// Log.Debugf("Type[ %s ]", resType)

		var err error

		switch resType {
		case "configmap":
			err = proc.KsmUpdateConfigmap(ksmManager, metricData, desc, Log)
		case "cronjob":
			err = proc.KsmUpdateCronjob(ksmManager, metricData, desc, Log)
		case "daemonset":
			err = proc.KsmUpdateDaemonset(ksmManager, metricData, desc, Log)
		case "deployment":
			err = proc.KsmUpdateDeployment(ksmManager, metricData, desc, Log)
		case "ingress":
			err = proc.KsmUpdateIngress(ksmManager, metricData, desc, Log)
		case "job":
			err = proc.KsmUpdateJob(ksmManager, metricData, desc, Log)
		case "namespace":
			err = proc.KsmUpdateNamespace(ksmManager, metricData, desc, Log)
		case "node":
			err = proc.KsmUpdateNode(ksmManager, metricData, desc, Log)
		case "persistentvolume":
			err = proc.KsmUpdatePersistentVolume(ksmManager, metricData, desc, Log)
		case "persistentvolumeclaim":
			err = proc.KsmUpdatePersistentVolumeClaim(ksmManager, metricData, desc, Log)
		case "pod":
			err = proc.KsmUpdatePod(ksmManager, metricData, desc, Log)
		case "replicaset":
			err = proc.KsmUpdateReplicaset(ksmManager, metricData, desc, Log)
		case "replicationcontroller":
			err = proc.KsmUpdateReplicationController(ksmManager, metricData, desc, Log)
		case "secret":
			err = proc.KsmUpdateSecret(ksmManager, metricData, desc, Log)
		case "service":
			err = proc.KsmUpdateService(ksmManager, metricData, desc, Log)
		case "statefulset":
			err = proc.KsmUpdateStatefulset(ksmManager, metricData, desc, Log)
		default:
			Log.Errorf("Error. Not Supported FqName[ %s ]", desc.FqName)
		}
		if err != nil {
			Log.Errorf("%s", err.Error())
			errCnt++
		}
	}
	if errCnt <= 0 {
		return nil
	}
	return fmt.Errorf("ERROR. Failed Processing Count[ %d ]", errCnt)
}

// GetTypeOfFqName fqname split and return resource type(node, namespace, pods)
func (proc *Proc) GetTypeOfFqName(input string) string {
	ret := strings.Split(input, "_")
	return ret[1]
}

// KsmUpdateConfigmap configmap
func (proc *Proc) KsmUpdateConfigmap(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("error. failed to parse configmap. can not found namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["configmap"]
	if !isExist {
		return fmt.Errorf("error. failed to parse configmap. can not found name.[ %+v ]", desc)
	}

	var configmap *ksm.Configmap
	configmap, err := ksmManager.GetConfigmap(metricData.CID, namespace, name)
	if err != nil {
		configmap, err = ksmManager.AllocConfigmap()
		if err != nil {
			return fmt.Errorf("error. failed to update configmap[ %s-%s ] [ %s ]", namespace, name, err.Error())
		}
		configmap.CID = metricData.CID
		configmap.Namespace = namespace
		configmap.Name = name
		configmap.Add()
	}
	// 갱신 시간 업데이트
	configmap.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameConfigMapInfo:
	case ksm.KsmFqNameConfigMapCreated:
		if configmap.Created == 0 {
			Log.Debugf("Configmap[ %s ] NS[ %s ] Created Time[ %s ]", configmap.Name,
				configmap.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			configmap.UpdateTime = tools.GetMillis()
			configmap.Created = int64(desc.Value) * 1000
		} else if configmap.Created != int64(desc.Value)*1000 {
			Log.Debugf("Configmap[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", configmap.Name,
				configmap.Namespace, tools.PrintTimeFromMilli(configmap.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			configmap.UpdateTime = tools.GetMillis()
			configmap.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameConfigMapVersion:
		if configmap.Metadata.ResourceVersion == 0 {
			Log.Debugf("Configmap[ %s ] NS[ %s ] Resource Version[ %d ]", configmap.Name,
				configmap.Namespace, int64(desc.Value))
			configmap.UpdateTime = tools.GetMillis()
			configmap.Metadata.ResourceVersion = int64(desc.Value)
		} else if configmap.Metadata.ResourceVersion != int64(desc.Value) {
			Log.Debugf("Configmap[ %s ] NS[ %s ] Changed Resource Version Before [ %d ] Currnet [ %d ]", configmap.Name,
				configmap.Namespace, configmap.Metadata.ResourceVersion, int64(desc.Value))
			configmap.UpdateTime = tools.GetMillis()
			configmap.Metadata.ResourceVersion = int64(desc.Value)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	configmap.Put()
	return nil
}

// KsmUpdateCronjob cronjob
func (proc *Proc) KsmUpdateCronjob(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Cronjob. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["cronjob"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Cronjob. Can Not Found Name.[ %+v ]", desc)
	}

	var cronjob *ksm.Cronjob
	cronjob, err := ksmManager.GetCronjob(metricData.CID, namespace, name)
	if err != nil {
		cronjob, err = ksmManager.AllocCronjob()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Cronjob[ %s-%s ]. [ %s ]", namespace, name, err.Error())
		}
		cronjob.CID = metricData.CID
		cronjob.Namespace = namespace
		cronjob.Name = name
		cronjob.Add()
	}
	// 갱신 시간 업데이트
	cronjob.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameCronjobLabels:
		if !proc.CompareMapData(cronjob.Labels, descLabels, []string{"namespace", "cronjob"}) {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Labels.", cronjob.Name, cronjob.Namespace)
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Labels = nil
			cronjob.Labels = make(map[string]string)
			proc.CopyMapData(cronjob.Labels, descLabels, []string{"namespace", "cronjob"})
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Update Labels.[ %+v ]", cronjob.Name, cronjob.Namespace, cronjob.Labels)
		}
	case ksm.KsmFqNameCronjobInfo:
		if !proc.CompareMapData(cronjob.Info, descLabels, []string{"namespace", "cronjob"}) {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Info.", cronjob.Name, cronjob.Namespace)
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Labels = nil
			cronjob.Labels = make(map[string]string)
			proc.CopyMapData(cronjob.Info, descLabels, []string{"namespace", "cronjob"})
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Update Info.[ %+v ]", cronjob.Name, cronjob.Namespace, cronjob.Info)
		}
	case ksm.KsmFqNameCronjobCreated:
		if cronjob.Created == 0 {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Created Time[ %s ]", cronjob.Name, cronjob.Namespace,
				tools.PrintTimeFromSec(int64(desc.Value)))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Created = int64(desc.Value) * 1000
		} else if cronjob.Created != int64(desc.Value)*1000 {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", cronjob.Name,
				cronjob.Namespace, tools.PrintTimeFromMilli(cronjob.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameCronjobStatus:
		var status bool
		if int(desc.Value) == 0 {
			status = false
		} else {
			status = true
		}
		if cronjob.Status.Active != status {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Status Active[ %t ]", cronjob.Name,
				cronjob.Namespace, cronjob.Status.Active)
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Status.Active = status
		}
	case ksm.KsmFqNameCronjobLastScheduleTime:
		if cronjob.Status.LastScheduleTime != int64(desc.Value)*1000 {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Last ScheduleTime. Before [ %s ] Currnet [ %s ]", cronjob.Name,
				cronjob.Namespace, tools.PrintTimeFromMilli(cronjob.Status.LastScheduleTime), tools.PrintTimeFromSec(int64(desc.Value)))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Status.LastScheduleTime = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameCronjobSpecSuspend:
		if cronjob.Spec.Suspend != int32(desc.Value) {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Spec Suspend. [ %d ]", cronjob.Name,
				cronjob.Namespace, int(desc.Value))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Spec.Suspend = int32(desc.Value)
		}
	case ksm.KsmFqNameCronjobSpecStartingDeadlineSeconds:
		if cronjob.Spec.StartingDeadlineSeconds != int64(desc.Value)*1000 {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Spec Starting Deadline Seconds. [ %d ]", cronjob.Name,
				cronjob.Namespace, int(desc.Value))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Spec.StartingDeadlineSeconds = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameCronjobNextSchedule:
		if cronjob.Spec.NextScheduletime != int64(desc.Value)*1000 {
			Log.Debugf("Cronjob[ %s ] NS[ %s ] Changed Spec NextScheduleTime. [ %s ]", cronjob.Name,
				cronjob.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			cronjob.UpdateTime = tools.GetMillis()
			cronjob.Spec.NextScheduletime = int64(desc.Value) * 1000
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	cronjob.Put()
	return nil
}

// KsmUpdateDaemonset daemonset
func (proc *Proc) KsmUpdateDaemonset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Daemonset State Metric")

	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Daemonset. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["daemonset"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Daemonset. Can Not Found Name.[ %+v ]", desc)
	}

	var daemonset *ksm.Daemonset
	daemonset, err := ksmManager.GetDaemonset(metricData.CID, namespace, name)
	if err != nil {
		daemonset, err = ksmManager.AllocDaemonset()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Daemonset[ %s-%s ]. [ %s ]", namespace, name, err.Error())
		}
		daemonset.CID = metricData.CID
		daemonset.Namespace = namespace
		daemonset.Name = name
		daemonset.Add()
	}
	// 갱신 시간 업데이트
	daemonset.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameDaemonsetCreated:
		if daemonset.Created == 0 {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Created Time[ %s ]", daemonset.Name, daemonset.Namespace,
				tools.PrintTimeFromSec(int64(desc.Value)))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Created = int64(desc.Value) * 1000
		} else if daemonset.Created != int64(desc.Value)*1000 {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", daemonset.Name,
				daemonset.Namespace, tools.PrintTimeFromMilli(daemonset.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameDaemonsetStatusCurrent:
		if daemonset.Status.Current != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Current Scheduled.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Current = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetStatusDesired:
		if daemonset.Status.Desired != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Desired.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Desired = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetStatusAvailable:
		if daemonset.Status.Available != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Available.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Available = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetStatusMisscheduled:
		if daemonset.Status.Misscheduled != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Misscheduled.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Misscheduled = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetStatusReady:
		if daemonset.Status.Ready != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Ready.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Ready = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetStatusUnavailable:
		if daemonset.Status.Unavailable != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Unavailable.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Unavailable = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetUpdated:
		if daemonset.Status.Scheduled != int16(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Number Scheduled.[ %d ]", daemonset.Name, daemonset.Namespace, int(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Status.Scheduled = int16(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetMetadataGeneration:
		if daemonset.Metadata.Generation != int64(desc.Value) {
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Changed Meta Generation.[ %d ]", daemonset.Name, daemonset.Namespace, int32(desc.Value))
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Metadata.Generation = int64(desc.Value)
		}
	case ksm.KsmFqNameDaemonsetLabels:
		if !proc.CompareMapData(daemonset.Labels, descLabels, []string{"namespace", "daemonset"}) {
			Log.Infof("Daemonset[ %s ] NS[ %s ] Changed Labels.", daemonset.Name, daemonset.Namespace)
			daemonset.UpdateTime = tools.GetMillis()
			daemonset.Labels = nil
			daemonset.Labels = make(map[string]string)
			proc.CopyMapData(daemonset.Labels, descLabels, []string{"namespace", "daemonset"})
			Log.Debugf("Daemonset[ %s ] NS[ %s ] Update Labels.[ %+v ]",
				daemonset.Name, daemonset.Namespace, daemonset.Labels)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	daemonset.Put()
	return nil
}

// KsmUpdateDeployment deployment
func (proc *Proc) KsmUpdateDeployment(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["deployment"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found Name.[ %+v ]", desc)
	}

	var deploy *ksm.Deployment
	deploy, err := ksmManager.GetDeployment(metricData.CID, namespace, name)
	if err != nil {
		deploy, err = ksmManager.AllocDeployment()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Deploy[ %s-%s ]. [ %s ]", namespace, name, err.Error())
		}
		deploy.CID = metricData.CID
		deploy.Namespace = namespace
		deploy.Name = name
		deploy.Add()
	}
	// 갱신 시간 업데이트
	deploy.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameDeploymentCreated:
		if deploy.Created == 0 {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Created Time[ %s ]", deploy.Name,
				deploy.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Created = int64(desc.Value) * 1000
		} else if deploy.Created != int64(desc.Value)*1000 {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", deploy.Name,
				deploy.Namespace, tools.PrintTimeFromMilli(deploy.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameDeploymentStatusReplicas:
		if deploy.Status.Replicas != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Replicas.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Status.Replicas = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentStatusReplicasAvailable:
		if deploy.Status.ReplicasAvailable != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Replicas Available.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Status.ReplicasAvailable = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentStatusReplicasUnavailable:
		if deploy.Status.ReplicasUnavailable != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Replicas Unavailable.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Status.ReplicasUnavailable = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentStatusReplicasUpdated:
		if deploy.Status.ReplicasUpdated != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Replicas Updated.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Status.ReplicasUpdated = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentStatusObserved:
		if deploy.Status.ObservedGeneration != int64(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Observed Generation.[ %d ]", deploy.Name, deploy.Namespace, int64(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Status.ObservedGeneration = int64(desc.Value)
		}
	case ksm.KsmFqNameDeploymentStatusCondition:
		condition, isExist := descLabels["condition"]
		if !isExist {
			deploy.Put()
			return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found contition.[ %+v ]", desc)
		}
		status, isExist := descLabels["status"]
		if !isExist {
			deploy.Put()
			return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found status.[ %+v ]", desc)
		}
		if int(desc.Value) == 1 {
			switch condition {
			case ksm.KsmDeployStatusConditionAvailable:
				if deploy.Status.Condition.Available != status {
					Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Condition Available.[ %s ]",
						deploy.Name, deploy.Namespace, status)
					deploy.UpdateTime = tools.GetMillis()
					deploy.Status.Condition.Available = status
				}
			case ksm.KsmDeployStatusConditionProgressing:
				if deploy.Status.Condition.Progressing != status {
					Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Condition Progressing.[ %s ]",
						deploy.Name, deploy.Namespace, status)
					deploy.UpdateTime = tools.GetMillis()
					deploy.Status.Condition.Progressing = status
				}
			case ksm.KsmDeployStatusConditionReplicaFailure:
				if deploy.Status.Condition.ReplicaFailure != status {
					Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Condition ReplicaFailure.[ %s ]",
						deploy.Name, deploy.Namespace, status)
					deploy.UpdateTime = tools.GetMillis()
					deploy.Status.Condition.ReplicaFailure = status
				}
			default:
				Log.Errorf("ERROR. Not supported Conditon[ %s ][ %+v ]", condition, desc)
			}
		}
	case ksm.KsmFqNameDeploymentSpecReplicas:
		if deploy.Spec.Replicas != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Spec Replicas.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Spec.Replicas = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentSpecPaused:
		if deploy.Spec.Paused != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Spec Pause.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Spec.Paused = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentSpecStrategyRollingupdateMaxUnavailable:
		if deploy.Spec.Strategy.MaxUnavailable != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Spec RollingUpdate Max Unavailable.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Spec.Strategy.MaxUnavailable = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentSpecStrategyRollingupdateMaxSurge:
		if deploy.Spec.Strategy.MaxSurge != int16(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Spec RollingUpdate Maxsurge.[ %d ]", deploy.Name, deploy.Namespace, int(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Spec.Strategy.MaxSurge = int16(desc.Value)
		}
	case ksm.KsmFqNameDeploymentMetadataGeneration:
		if deploy.Metadata.Generation != int64(desc.Value) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Meta Generation.[ %d ]", deploy.Name, deploy.Namespace, int64(desc.Value))
			deploy.UpdateTime = tools.GetMillis()
			deploy.Metadata.Generation = int64(desc.Value)
		}
	case ksm.KsmFqNameDeploymentLabels:
		if !proc.CompareMapData(deploy.Labels, descLabels, []string{"namespace", "deployment"}) {
			Log.Debugf("Deploy[ %s ] NS[ %s ] Changed Labels.", deploy.Name, deploy.Namespace)
			deploy.UpdateTime = tools.GetMillis()
			deploy.Labels = nil
			deploy.Labels = make(map[string]string)
			proc.CopyMapData(deploy.Labels, descLabels, []string{"namespace", "deployment"})
			Log.Debugf("Deploy[ %s ] NS[ %s ] Update Labels.[ %+v ]", deploy.Name, deploy.Namespace, deploy.Labels)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	deploy.Put()
	return nil
}

// KsmUpdateIngress ingress
func (proc *Proc) KsmUpdateIngress(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Ingress State Metric")
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["ingress"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingressh. Can Not Found Name.[ %+v ]", desc)
	}

	var ingress *ksm.Ingress
	ingress, err := ksmManager.GetIngress(metricData.CID, namespace, name)
	if err != nil {
		ingress, err = ksmManager.AllocIngress()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Ingress[ %s-%s ]. [ %s ]", namespace, name, err.Error())
		}
		ingress.CID = metricData.CID
		ingress.Namespace = namespace
		ingress.Name = name
		ingress.Add()
	}
	// 갱신 시간 업데이트
	ingress.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameIngressInfo:
	case ksm.KsmFqNameIngressLabels:
		if !proc.CompareMapData(ingress.Labels, descLabels, []string{"namespace", "ingress"}) {
			Log.Debugf("Ingress[ %s ] NS[ %s ] Changed Labels.", ingress.Name, ingress.Namespace)
			ingress.UpdateTime = tools.GetMillis()
			ingress.Labels = nil
			ingress.Labels = make(map[string]string)
			proc.CopyMapData(ingress.Labels, descLabels, []string{"namespace", "ingress"})
			Log.Debugf("Ingress[ %s ] NS[ %s ] Update Labels.[ %+v ]", ingress.Name, ingress.Namespace, ingress.Labels)
		}
	case ksm.KsmFqNameIngressCreated:
		if ingress.Created == 0 {
			Log.Debugf("Ingress[ %s ] NS[ %s ] Created Time[ %s ]", ingress.Name,
				ingress.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			ingress.UpdateTime = tools.GetMillis()
			ingress.Created = int64(desc.Value) * 1000
		} else if ingress.Created != int64(desc.Value)*1000 {
			Log.Debugf("Ingress[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", ingress.Name,
				ingress.Namespace, tools.PrintTimeFromMilli(ingress.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			ingress.UpdateTime = tools.GetMillis()
			ingress.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameIngressMetadataResourceVersion:
		if ingress.Metadata.ResourceVersion != int64(desc.Value) {
			Log.Debugf("Ingress[ %s ] NS[ %s ] Changed Resource Version.[ %d ]", ingress.Name, ingress.Namespace, int64(desc.Value))
			ingress.UpdateTime = tools.GetMillis()
			ingress.Metadata.ResourceVersion = int64(desc.Value)
		}
	case ksm.KsmFqNameIngressPath:
		err := proc.UpdateOrInsertIngressPath(ingress, descLabels, Log)
		if err != nil {
			ingress.Put()
			return err
		}
	case ksm.KsmFqNameIngressTLS:
		err := proc.UpdateOrInsertIngressTLS(ingress, descLabels, Log)
		if err != nil {
			ingress.Put()
			return err
		}
	default:
		Log.Errorf("Error. Not Supported FqName[ %s ]", desc.FqName)
	}
	ingress.Put()
	return nil
}

// KsmUpdateJob job
func (proc *Proc) KsmUpdateJob(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Job State Metric")
	return nil
}

// KsmUpdateNamespace namespace
func (proc *Proc) KsmUpdateNamespace(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Namespace State Metric")
	descLabels := desc.Labels.(map[string]string)
	nmName, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse. Can Not Found Namespace[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}

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
			nm.UpdateTime = tools.GetMillis()
			nm.Created = int64(desc.Value) * 1000
		} else if nm.Created != int64(desc.Value)*1000 {
			Log.Debugf("Namespace[ %s ] Update Created Time. Before [ %s ] Currnet [ %s ]", nm.Name,
				tools.PrintTimeFromMilli(nm.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			nm.UpdateTime = tools.GetMillis()
			nm.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameNamespaceLabels:
		if !proc.CompareMapData(nm.Labels, descLabels, []string{"namespace"}) {
			Log.Debugf("Namespace[ %s ] Changed Labels.", nm.Name)
			nm.UpdateTime = tools.GetMillis()
			nm.Labels = nil
			nm.Labels = make(map[string]string)
			proc.CopyMapData(nm.Labels, descLabels, []string{"namespace"})
			Log.Debugf("Namespace[ %s ] Update Labels.[ %+v ]", nm.Name, nm.Labels)
		}
	case ksm.KsmFqNameNamespaceStatusPhase:
		phase := descLabels["phase"]
		if int(desc.Value) == 1 {
			if nm.Phase != phase {
				Log.Debugf("Namespace[ %s ] Changed Phase[ %s ].", nm.Name, phase)
				nm.UpdateTime = tools.GetMillis()
				nm.Phase = phase
			}
		}
	// case ksm.KsmFqNameNamespaceStatusCondition:
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	nm.Put()

	return nil
}

// KsmUpdateNode node
func (proc *Proc) KsmUpdateNode(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	descLabels := desc.Labels.(map[string]string)
	nodeName, isExist := descLabels["node"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed To Parse. Can Not Found Node Name[ %+v ]", desc)
	}
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
		if !proc.CompareMapData(node.NodeInfo, descLabels, []string{"node"}) {
			Log.Debugf("Node[ %s ] Changed NodeInfo.", node.Name)
			node.UpdateTime = tools.GetMillis()
			node.NodeInfo = nil
			node.NodeInfo = make(map[string]string)
			proc.CopyMapData(node.NodeInfo, descLabels, []string{"node"})
			Log.Debugf("Node[ %s ] Update NodeInfo.[ %+v ]", node.Name, node.NodeInfo)
		}
	case ksm.KsmFqNameNodeCreated:
		if node.Created == 0 {
			node.UpdateTime = tools.GetMillis()
			Log.Debugf("Node[ %s ] Created Time[ %s ]", node.Name, tools.PrintTimeFromSec(int64(desc.Value)))
			node.Created = int64(desc.Value) * 1000
		} else if node.Created != int64(desc.Value)*1000 {
			node.UpdateTime = tools.GetMillis()
			Log.Debugf("Node[ %s ] Update Created Time. Before [ %s ] Currnet [ %s ]", node.Name,
				tools.PrintTimeFromMilli(node.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			node.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameNodeLabels:
		if !proc.CompareMapData(node.Labels, descLabels, []string{"node"}) {
			Log.Debugf("Node[ %s ] Changed Labels.", node.Name)
			node.UpdateTime = tools.GetMillis()
			node.Labels = nil
			node.Labels = make(map[string]string)
			proc.CopyMapData(node.Labels, descLabels, []string{"node"})
			Log.Debugf("Node[ %s ] Update Labels.[ %+v ]", node.Name, node.Labels)
		}
	case ksm.KsmFqNameNodeRole:
		if !proc.CompareMapData(node.Roles, descLabels, []string{"node"}) {
			Log.Debugf("Node[ %s ] Changed Role.", node.Name)
			node.UpdateTime = tools.GetMillis()
			node.Roles = nil
			node.Roles = make(map[string]string)
			proc.CopyMapData(node.Roles, descLabels, []string{"node"})
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
			node.Put()
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
			node.UpdateTime = tools.GetMillis()
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
		node.UpdateTime = tools.GetMillis()
		Log.Debugf("Node[ %s ] Status Phase.[ %+v ][ %f ]", node.Name, descLabels, desc.Value)
	case ksm.KsmFqNameNodeStatusCapacityPods:
		if node.Status.Capacity.Pods != int32(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.Pods = int32(desc.Value)
			Log.Debugf("Node[ %s ] Changed Capacity Pods.[ %d ]", node.Name, int32(desc.Value))
		}
	case ksm.KsmFqNameNodeStatusCapacityCPU:
		if node.Status.Capacity.Cores != int32(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.Cores = int32(desc.Value)
			Log.Debugf("Node[ %s ] Changed Capacity Cores.[ %d ]",
				node.Name, node.Status.Capacity.Cores)
		}
	case ksm.KsmFqNameNodeStatusCapacityMemory:
		if node.Status.Capacity.MemoryBytes != int64(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.MemoryBytes = int64(desc.Value)
			Log.Debugf("Node[ %s ] Changed Capacity Mem Bytes.[ %d ]GiB",
				node.Name, node.Status.Capacity.MemoryBytes/1000/1000/1000)
		}
	case ksm.KsmFqNameNodeStatusAllocatablePods:
		if node.Status.Allocatable.Pods != int32(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.Pods = int32(desc.Value)
			Log.Debugf("Node[ %s ] Changed Allocatable Pods.[ %d ]",
				node.Name, node.Status.Allocatable.Pods)
		}
	case ksm.KsmFqNameNodeStatusAllocatableCPU:
		if node.Status.Allocatable.Cores != int32(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.Cores = int32(desc.Value)
			Log.Debugf("Node[ %s ] Changed Allocatable Cores.[ %d ]",
				node.Name, node.Status.Allocatable.Cores)
		}
	case ksm.KsmFqNameNodeStatusAllocatableMemory:
		if node.Status.Allocatable.MemoryBytes != int64(desc.Value) {
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.MemoryBytes = int64(desc.Value)
			Log.Debugf("Node[ %s ] Changed Allocatable Mem Bytes.[ %d ]GiB",
				node.Name, node.Status.Allocatable.MemoryBytes/1000/1000/1000)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	node.Put()
	return nil
}

// KsmUpdatePersistentVolume persistentvolume
func (proc *Proc) KsmUpdatePersistentVolume(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	pvName, isExist := descLabels["persistentvolume"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse. Can Not Found persistentvolume[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}

	var pv *ksm.PersistentVolume
	pv, err := ksmManager.GetPersistentVolume(metricData.CID, pvName)
	if err != nil {
		pv, err = ksmManager.AllocPersistentVolume()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Persistentvolume[ %s ]. [ %s ]", pvName, err.Error())
		}
		pv.CID = metricData.CID
		pv.Name = pvName
		pv.Add()
	}
	// 갱신 시간 업데이트
	pv.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNamePersistentvolumeLabels:
		if !proc.CompareMapData(pv.Labels, descLabels, []string{"persistentvolume"}) {
			Log.Debugf("PV[ %s ] Changed Labels.", pv.Name)
			pv.UpdateTime = tools.GetMillis()
			pv.Labels = nil
			pv.Labels = make(map[string]string)
			proc.CopyMapData(pv.Labels, descLabels, []string{"persistentvolume"})
			Log.Debugf("PV[ %s ] Update Labels.[ %+v ]", pv.Name, pv.Labels)
		}
	case ksm.KsmFqNamePersistentvolumeStatusPhase:
		if int(desc.Value) == 1 {
			phase, isExist := descLabels["phase"]
			if !isExist {
				pv.Put()
				return fmt.Errorf("ERROR. Failed to Parse Persistentvolume. Can Not Found phase.[ %+v ]", descLabels)
			}
			if pv.Phase != phase {
				pv.UpdateTime = tools.GetMillis()
				pv.Phase = phase
				Log.Debugf("PV[ %s ] Changed Status.[ %s ]", pv.Name, pv.Phase)
			}
		}
	case ksm.KsmFqNamePersistentvolumeInfo:
		if !proc.CompareMapData(pv.Infos, descLabels, []string{"persistentvolume"}) {
			Log.Debugf("PV[ %s ] Changed Info.", pv.Name)
			pv.UpdateTime = tools.GetMillis()
			pv.Infos = nil
			pv.Infos = make(map[string]string)
			proc.CopyMapData(pv.Infos, descLabels, []string{"persistentvolume"})
			Log.Debugf("PV[ %s ] Update Infos.[ %+v ]", pv.Name, pv.Infos)
		}
	case ksm.KsmFqNamePersistentvolumeCapacityBytes:
		if pv.CapacityBytes != int64(desc.Value) {
			pv.UpdateTime = tools.GetMillis()
			pv.CapacityBytes = int64(desc.Value)
			Log.Debugf("PV[ %s ] Changed CapacityBytes.[ %d ]MiB", pv.Name, pv.CapacityBytes/1000/1000)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}

	pv.Put()
	return nil
}

// KsmUpdatePersistentVolumeClaim persistentvolumeclaim
func (proc *Proc) KsmUpdatePersistentVolumeClaim(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse. Can Not Found namespace[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}
	pvcName, isExist := descLabels["persistentvolumeclaim"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse. Can Not Found persistentvolumeclaim[ %s ][ %+v ]", desc.FqName, desc.Labels)
	}

	var pvc *ksm.PersistentVolumeClaim
	pvc, err := ksmManager.GetPersistentVolumeClaim(metricData.CID, namespace, pvcName)
	if err != nil {
		pvc, err = ksmManager.AllocPersistentVolumeClaim()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Persistentvolumeclaim[ %s ]. [ %s ]", pvcName, err.Error())
		}
		pvc.CID = metricData.CID
		pvc.Name = pvcName
		pvc.Namespace = namespace
		pvc.Add()
	}
	// 갱신 시간 업데이트
	pvc.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNamePersistentvolumeclaimLabels:
		if !proc.CompareMapData(pvc.Labels, descLabels, []string{"namespace", "persistentvolumeclaim"}) {
			Log.Debugf("PVC[ %s ] NS[ %s ] Changed Labels.",
				pvc.Name, pvc.Namespace)
			pvc.UpdateTime = tools.GetMillis()
			pvc.Labels = nil
			pvc.Labels = make(map[string]string)
			proc.CopyMapData(pvc.Labels, descLabels, []string{"namespace", "persistentvolumeclaim"})
			Log.Debugf("PVC[ %s ] NS[ %s ] Update Labels.[ %+v ]",
				pvc.Name, pvc.Namespace, pvc.Labels)
		}
	case ksm.KsmFqNamePersistentvolumeclaimInfo:
		if !proc.CompareMapData(pvc.Infos, descLabels, []string{"namespace", "persistentvolumeclaim"}) {
			Log.Debugf("PVC[ %s ] NS[ %s ] Changed Infos.", pvc.Name, pvc.Namespace)
			pvc.UpdateTime = tools.GetMillis()
			pvc.Infos = nil
			pvc.Infos = make(map[string]string)
			proc.CopyMapData(pvc.Infos, descLabels, []string{"namespace", "persistentvolumeclaim"})
			Log.Debugf("PVC[ %s ] NS[ %s ] Update Infos.[ %+v ]",
				pvc.Name, pvc.Namespace, pvc.Infos)
		}
	case ksm.KsmFqNamePersistentvolumeclaimStatusPhase:
		if int(desc.Value) == 1 {
			phase, isExist := descLabels["phase"]
			if !isExist {
				pvc.Put()
				return fmt.Errorf("ERROR. Failed to Parse. Can Not Found phase[ %s ][ %+v ]",
					desc.FqName, desc.Labels)
			}
			if pvc.Phase != phase {
				pvc.UpdateTime = tools.GetMillis()
				pvc.Phase = phase
				Log.Debugf("PVC[ %s ] NS[ %s ] Changed Status.[ %+v ]",
					pvc.Name, pvc.Namespace, pvc.Phase)
			}
		}
	case ksm.KsmFqNamePersistentvolumeclaimResourceRequestsBytes:
		if pvc.RequestBytes != int64(desc.Value) {
			pvc.UpdateTime = tools.GetMillis()
			pvc.RequestBytes = int64(desc.Value)
			Log.Debugf("PVC[ %s ] NS[ %s ] Changed RequestBytes.[ %s ]",
				pvc.Name, pvc.Namespace, pvc.RequestBytes)
		}
	case ksm.KsmFqNamePersistentvolumeclaimAccessMode:
		err := proc.UpdateOrInsertPVCAccessModes(pvc, descLabels, Log)
		if err != nil {
			pvc.Put()
			return err
		}
	case ksm.KsmFqNamePersistentvolumeclaimStatusCondition:
		if int(desc.Value) == 1 {
			err := proc.UpdatePVCCondition(pvc, descLabels, Log)
			if err != nil {
				pvc.Put()
				return err
			}
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	pvc.Put()
	return nil
}

// KsmUpdatePod pod
func (proc *Proc) KsmUpdatePod(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod. Can Not Found namespace[ %s ][ %+v ]",
			desc.FqName, desc.Labels)
	}
	podName, isExist := descLabels["pod"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod. Can Not Found pod[ %s ][ %+v ]",
			desc.FqName, desc.Labels)
	}

	var pod *ksm.Pod
	pod, err := ksmManager.GetPod(metricData.CID, namespace, podName)
	if err != nil {
		pod, err = ksmManager.AllocPod()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update POD[ %s ]. [ %s ]", podName, err.Error())
		}
		pod.CID = metricData.CID
		pod.Namespace = namespace
		pod.Name = podName
		pod.Add()
	}
	// 갱신 시간 업데이트
	pod.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNamePodInfo:
		// Node 정보는 resource, pod_info에서만 확인 가능
		node := descLabels["node"]
		if pod.Node != node {
			pod.Node = node
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Node[ %s ].", pod.Name, pod.Namespace, pod.Node)
		}
		if !proc.CompareMapData(pod.Infos, descLabels, []string{"namespace", "pod"}) {
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Infos.", pod.Name, pod.Namespace)
			pod.UpdateTime = tools.GetMillis()
			pod.Infos = nil
			pod.Infos = make(map[string]string)
			proc.CopyMapData(pod.Infos, descLabels, []string{"namespace", "pod"})
			Log.Debugf("POD[ %s ] NS[ %s ] Update Infos.[ %+v ]", pod.Name, pod.Namespace, pod.Infos)
		}
	case ksm.KsmFqNamePodStartTime:
		if pod.StartTime != int64(desc.Value)*1000 {
			pod.UpdateTime = tools.GetMillis()
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Start Time. Before [ %s ] Currnet [ %s ]", pod.Name,
				pod.Namespace, tools.PrintTimeFromMilli(pod.StartTime), tools.PrintTimeFromSec(int64(desc.Value)))
			pod.StartTime = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNamePodCompletionTime:
	case ksm.KsmFqNamePodOwner:
		if !proc.CompareMapData(pod.Owner, descLabels, []string{"namespace", "pod"}) {
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Owner.", pod.Name, pod.Namespace)
			pod.UpdateTime = tools.GetMillis()
			pod.Owner = nil
			pod.Owner = make(map[string]string)
			proc.CopyMapData(pod.Owner, descLabels, []string{"namespace", "pod"})
			Log.Debugf("POD[ %s ] NS[ %s ] Update Owner.[ %+v ]", pod.Name, pod.Namespace, pod.Owner)
		}
	case ksm.KsmFqNamePodLabels:
		if !proc.CompareMapData(pod.Labels, descLabels, []string{"namespace", "pod"}) {
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Labels.", pod.Name, pod.Namespace)
			pod.UpdateTime = tools.GetMillis()
			pod.Labels = nil
			pod.Labels = make(map[string]string)
			proc.CopyMapData(pod.Labels, descLabels, []string{"namespace", "pod"})
			Log.Debugf("POD[ %s ] NS[ %s ] Update Labels.[ %+v ]", pod.Name, pod.Namespace, pod.Labels)
		}
	case ksm.KsmFqNamePodCreated:
		if pod.Created != int64(desc.Value)*1000 {
			pod.UpdateTime = tools.GetMillis()
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Create Time. Before [ %s ] Currnet [ %s ]", pod.Name,
				pod.Namespace, tools.PrintTimeFromMilli(pod.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			pod.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNamePodRestartPolicy:
		policy := descLabels["type"]
		if pod.RestartPolicy != policy {
			pod.UpdateTime = tools.GetMillis()
			pod.RestartPolicy = policy
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Restart Policy[ %s ].", pod.Name, pod.Namespace, policy)
		}
	case ksm.KsmFqNamePodStatusscheduledTime:
		if pod.Status.ScheduledTime != int64(desc.Value)*1000 {
			pod.UpdateTime = tools.GetMillis()
			Log.Debugf("POD[ %s ] NS[ %s ] Changed Scheduled Time. Before [ %s ] Currnet [ %s ]", pod.Name,
				pod.Namespace, tools.PrintTimeFromMilli(pod.Status.ScheduledTime), tools.PrintTimeFromSec(int64(desc.Value)))
			pod.Status.ScheduledTime = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNamePodStatusunschedulable:
		// 스케쥴링되지 못했다는 것을 알 수는 있음.
		// 그러나 정상적으로 스케쥴링되었을 때는 이 메시지는
		// 수신되지 못하고, iris-cloud에서 처리가 어려워짐.
	case ksm.KsmFqNamePodStatusphase:
		if int(desc.Value) == 1 {
			phase := descLabels["phase"]
			if pod.Status.Phase != phase {
				pod.UpdateTime = tools.GetMillis()
				Log.Debugf("POD[ %s ] NS[ %s ] Changed Phase[ %s => %s ]",
					pod.Name, pod.Namespace, pod.Status.Phase, phase)
				pod.Status.Phase = phase
			}
		}
	case ksm.KsmFqNamePodStatusready:
		if int(desc.Value) == 1 {
			ready := descLabels["condition"]
			if pod.Status.Ready != ready {
				pod.UpdateTime = tools.GetMillis()
				Log.Debugf("POD[ %s ] NS[ %s ] Changed Ready[ %s => %s ]",
					pod.Name, pod.Namespace, pod.Status.Ready, ready)
				pod.Status.Ready = ready
			}
		}
	case ksm.KsmFqNamePodStatusScheduled:
		if int(desc.Value) == 1 {
			scheduled := descLabels["condition"]
			if pod.Status.Scheduled != scheduled {
				pod.UpdateTime = tools.GetMillis()
				Log.Debugf("POD[ %s ] NS[ %s ] Changed Scheduled[ %s => %s ]",
					pod.Name, pod.Namespace, pod.Status.Scheduled, scheduled)
				pod.Status.Scheduled = scheduled
			}
		}
	case ksm.KsmFqNamePodContainerInfo:
		err := proc.UpdateOrInsertContainerInfo(pod, descLabels, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodInitContainerInfo:
		err := proc.UpdateOrInsertInitContainerInfo(pod, descLabels, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodInitContainerStatusWaiting,
		ksm.KsmFqNamePodInitContainerStatusWaitingReason,
		ksm.KsmFqNamePodInitContainerStatusRunning,
		ksm.KsmFqNamePodInitContainerStatusTerminated,
		ksm.KsmFqNamePodInitContainerStatusTerminatedReason,
		ksm.KsmFqNamePodInitContainerStatusLastTerminatedReason,
		ksm.KsmFqNamePodInitContainerStatusReady,
		ksm.KsmFqNamePodInitContainerStatusRestartsTotal:
		err := proc.UpdatePodContainerStatus(pod, true, desc, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodContainerStatusWaiting,
		ksm.KsmFqNamePodContainerStatusWaitingReason,
		ksm.KsmFqNamePodContainerStatusRunning,
		ksm.KsmFqNamePodContainerStatusTerminated,
		ksm.KsmFqNamePodContainerStatusTerminatedReason,
		ksm.KsmFqNamePodContainerStatusLastTerminatedReason,
		ksm.KsmFqNamePodContainerStatusReady,
		ksm.KsmFqNamePodContainerStatusRestartsTotal:
		err := proc.UpdatePodContainerStatus(pod, false, desc, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodContainerResourceRequests:
		// ARGS : pod *ksm.Pod, isInitContainer bool, isRequest bool, desclabels map[string]string, value float64, Log
		err := proc.UpdatePodContainerResource(pod, false, true, descLabels, desc.Value, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodContainerResourceLimits:
		// ARGS : pod *ksm.Pod, isInitContainer bool, isRequest bool, desclabels map[string]string, value float64, Log
		err := proc.UpdatePodContainerResource(pod, false, false, descLabels, desc.Value, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodInitContainerResourceRequests:
		// ARGS : pod *ksm.Pod, isInitContainer bool, isRequest bool, desclabels map[string]string, value float64, Log
		err := proc.UpdatePodContainerResource(pod, true, true, descLabels, desc.Value, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodInitContainerResourceLimits:
		// ARGS : pod *ksm.Pod, isInitContainer bool, isRequest bool, desclabels map[string]string, value float64, Log
		err := proc.UpdatePodContainerResource(pod, false, true, descLabels, desc.Value, Log)
		if err != nil {
			pod.Put()
			return err
		}
	case ksm.KsmFqNamePodSpecVolumesPersistentvolumeclaimsInfo:
		if int(desc.Value) == 1 {
			err := proc.UpdateOrInsertPodVolumes(pod, descLabels, Log)
			if err != nil {
				pod.Put()
				return err
			}
		}
	case ksm.KsmFqNamePodSpecVolumesPersistentvolumeclaimsReadonly:
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	pod.Put()
	return nil
}

// KsmUpdateReplicaset replicaset
func (proc *Proc) KsmUpdateReplicaset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update RS State Metric")

	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Replicaset. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["replicaset"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Replicaset. Can Not Found Name.[ %+v ]", desc)
	}

	var rc *ksm.Replicaset
	rc, err := ksmManager.GetReplicaset(metricData.CID, namespace, name)
	if err != nil {
		rc, err = ksmManager.AllocReplicaset()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Replicaset[ %s-%s ]. [ %s ]",
				namespace, name, err.Error())
		}
		rc.CID = metricData.CID
		rc.Namespace = namespace
		rc.Name = name
		rc.Add()
	}
	// 갱신 시간 업데이트
	rc.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameReplicasetCreated:
		if rc.Created == 0 {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Created Time[ %s ]", rc.Name, rc.Namespace,
				tools.PrintTimeFromSec(int64(desc.Value)))
			rc.UpdateTime = tools.GetMillis()
			rc.Created = int64(desc.Value) * 1000
		} else if rc.Created != int64(desc.Value)*1000 {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]", rc.Name,
				rc.Namespace, tools.PrintTimeFromMilli(rc.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			rc.UpdateTime = tools.GetMillis()
			rc.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameReplicasetStatusReplicas:
		if rc.Status.Replicas != int(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Status Replicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicasetStatusFullyLabeledReplicas:
		if rc.Status.FullyLabeledReplicas != int(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Status FullyLabeledReplicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.FullyLabeledReplicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicasetStatusReadyReplicas:
		if rc.Status.ReadyReplicas != int(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Status ReadyReplicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.ReadyReplicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicasetStatusObservedGeneration:
		if rc.Status.ObservedGeneration != int(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Status ObservedGeneration.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.ObservedGeneration = int(desc.Value)
		}
	case ksm.KsmFqNameReplicasetSpecReplicas:
		if rc.Spec.Replicas != int(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Spec Replicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Spec.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicasetMetadataGeneration:
		if rc.Metadata.Generation != int64(desc.Value) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Meta Generation.[ %d ]",
				rc.Name, rc.Namespace, int64(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Metadata.Generation = int64(desc.Value)
		}
	case ksm.KsmFqNameReplicasetOwner:
		if !proc.CompareMapData(rc.Owner, descLabels, []string{"namespace", "replicaset"}) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Owner.",
				rc.Name, rc.Namespace)
			rc.UpdateTime = tools.GetMillis()
			rc.Owner = nil
			rc.Owner = make(map[string]string)
			proc.CopyMapData(rc.Owner, descLabels, []string{"namespace", "replicaset"})
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Update Owner.[ %+v ]",
				rc.Name, rc.Namespace, rc.Owner)
		}
	case ksm.KsmFqNameReplicasetLabels:
		if !proc.CompareMapData(rc.Labels, descLabels, []string{"namespace", "replicaset"}) {
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Changed Labels.",
				rc.Name, rc.Namespace)
			rc.UpdateTime = tools.GetMillis()
			rc.Labels = nil
			rc.Labels = make(map[string]string)
			proc.CopyMapData(rc.Labels, descLabels, []string{"namespace", "replicaset"})
			Log.Debugf("Replicaset[ %s ] NS[ %s ] Update Labels.[ %+v ]",
				rc.Name, rc.Namespace, rc.Labels)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	rc.Put()

	return nil
}

// KsmUpdateReplicationController replicationcontroller
func (proc *Proc) KsmUpdateReplicationController(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {

	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse ReplicationController. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["replicationcontroller"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse ReplicationController. Can Not Found Name.[ %+v ]", desc)
	}

	var rc *ksm.ReplicationController
	rc, err := ksmManager.GetReplicationController(metricData.CID, namespace, name)
	if err != nil {
		rc, err = ksmManager.AllocReplicationController()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update ReplicationController[ %s-%s ]. [ %s ]",
				namespace, name, err.Error())
		}
		rc.CID = metricData.CID
		rc.Namespace = namespace
		rc.Name = name
		rc.Add()
	}
	// 갱신 시간 업데이트
	rc.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameReplicationcontrollerCreated:
		if rc.Created == 0 {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Created Time[ %s ]", rc.Name, rc.Namespace,
				tools.PrintTimeFromSec(int64(desc.Value)))
			rc.UpdateTime = tools.GetMillis()
			rc.Created = int64(desc.Value) * 1000
		} else if rc.Created != int64(desc.Value)*1000 {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Created Time."+
				" Before [ %s ] Currnet [ %s ]", rc.Name, rc.Namespace,
				tools.PrintTimeFromMilli(rc.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			rc.UpdateTime = tools.GetMillis()
			rc.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameReplicationcontrollerStatusReplicas:
		if rc.Status.Replicas != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Status Replicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerStatusFullyLabeledReplicas:
		if rc.Status.FullyLabeledReplicas != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Status "+
				"FullyLabeledReplicas.[ %d ]", rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.FullyLabeledReplicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerStatusReadyReplicas:
		if rc.Status.ReadyReplicas != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Status ReadyReplicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.ReadyReplicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerStatusAvailableReplicas:
		if rc.Status.AvailableReplicas != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Status availableReplicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.AvailableReplicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerStatusObservedGeneration:
		if rc.Status.ObservedGeneration != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Status ObservedGeneration.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Status.ObservedGeneration = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerSpecReplicas:
		if rc.Spec.Replicas != int(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Spec Replicas.[ %d ]",
				rc.Name, rc.Namespace, int(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Spec.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameReplicationcontrollerMetadataGeneration:
		if rc.Metadata.Generation != int64(desc.Value) {
			Log.Debugf("ReplicationController[ %s ] NS[ %s ] Changed Meta Generation.[ %d ]",
				rc.Name, rc.Namespace, int64(desc.Value))
			rc.UpdateTime = tools.GetMillis()
			rc.Metadata.Generation = int64(desc.Value)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	rc.Put()
	return nil
}

// KsmUpdateSecret secret
func (proc *Proc) KsmUpdateSecret(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Secret. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["secret"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Secret. Can Not Found Name.[ %+v ]", desc)
	}

	var secret *ksm.Secret
	secret, err := ksmManager.GetSecret(metricData.CID, namespace, name)
	if err != nil {
		secret, err = ksmManager.AllocSecret()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Secret[ %s-%s ]. [ %s ]",
				namespace, name, err.Error())
		}
		secret.CID = metricData.CID
		secret.Namespace = namespace
		secret.Name = name
		secret.Add()
	}
	// 갱신 시간 업데이트
	secret.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameSecretInfo:
	case ksm.KsmFqNameSecretType:
		secretType := descLabels["type"]
		if secret.Type != secretType {
			Log.Debugf("Secret[ %s ] NS[ %s ] Changed Type[ %s ].", secret.Name,
				secret.Namespace, secretType)
			secret.UpdateTime = tools.GetMillis()
			secret.Type = secretType
		}
	case ksm.KsmFqNameSecretLabels:
		if !proc.CompareMapData(secret.Labels, descLabels, []string{"namespace", "secret"}) {
			Log.Debugf("Secret[ %s ] NS[ %s ] Changed Labels.", secret.Name, secret.Namespace)
			secret.UpdateTime = tools.GetMillis()
			secret.Labels = nil
			secret.Labels = make(map[string]string)
			proc.CopyMapData(secret.Labels, descLabels, []string{"namespace", "secret"})
			Log.Debugf("Secret[ %s ] NS[ %s ] Update Labels.[ %+v ]",
				secret.Name, secret.Namespace, secret.Labels)
		}
	case ksm.KsmFqNameSecretCreated:
		if secret.Created == 0 {
			Log.Debugf("Secret[ %s ] NS[ %s ] Created Time[ %s ]",
				secret.Name, secret.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			secret.UpdateTime = tools.GetMillis()
			secret.Created = int64(desc.Value) * 1000
		} else if secret.Created != int64(desc.Value)*1000 {
			Log.Debugf("Secret[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]",
				secret.Name, secret.Namespace, tools.PrintTimeFromMilli(secret.Created),
				tools.PrintTimeFromSec(int64(desc.Value)))
			secret.UpdateTime = tools.GetMillis()
			secret.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameSecretMetadataResourceVersion:
		if secret.Metadata.ResourceVersion != int64(desc.Value) {
			Log.Debugf("Secret[ %s ] NS[ %s ] Changed Meta ResourceVersion.[ %d ]",
				secret.Name, secret.Namespace, int64(desc.Value))
			secret.UpdateTime = tools.GetMillis()
			secret.Metadata.ResourceVersion = int64(desc.Value)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	secret.Put()
	return nil
}

// KsmUpdateService service
func (proc *Proc) KsmUpdateService(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Service. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["service"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Service. Can Not Found Name.[ %+v ]", desc)
	}

	var service *ksm.Service
	service, err := ksmManager.GetService(metricData.CID, namespace, name)
	if err != nil {
		service, err = ksmManager.AllocService()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Service[ %s-%s ]. [ %s ]",
				namespace, name, err.Error())
		}
		service.CID = metricData.CID
		service.Namespace = namespace
		service.Name = name
		service.Add()
	}
	// 갱신 시간 업데이트
	service.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameServiceInfo:
		if !proc.CompareMapData(service.Info, descLabels, []string{"namespace", "service"}) {
			Log.Debugf("Service[ %s ] NS[ %s ] Changed Infos.", service.Name, service.Namespace)
			service.UpdateTime = tools.GetMillis()
			service.Info = nil
			service.Info = make(map[string]string)
			proc.CopyMapData(service.Info, descLabels, []string{"namespace", "service"})
			Log.Debugf("Service[ %s ] NS[ %s ] Update Infos.[ %+v ]",
				service.Name, service.Namespace, service.Info)
		}
	case ksm.KsmFqNameServiceCreated:
		if service.Created == 0 {
			Log.Debugf("Service[ %s ] NS[ %s ] Created Time[ %s ]", service.Name, service.Namespace,
				tools.PrintTimeFromSec(int64(desc.Value)))
			service.UpdateTime = tools.GetMillis()
			service.Created = int64(desc.Value) * 1000
		} else if service.Created != int64(desc.Value)*1000 {
			Log.Debugf("Service[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]",
				service.Name, service.Namespace, tools.PrintTimeFromMilli(service.Created),
				tools.PrintTimeFromSec(int64(desc.Value)))
			service.UpdateTime = tools.GetMillis()
			service.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameServiceSpecType:
		serviceType := descLabels["type"]
		if service.Type != serviceType {
			Log.Debugf("service[ %s ] NS[ %s ] Changed Service Type.[ %s ]",
				service.Name, service.Namespace, serviceType)
			service.UpdateTime = tools.GetMillis()
			service.Type = serviceType
		}
	case ksm.KsmFqNameServiceLabels:
		if !proc.CompareMapData(service.Labels, descLabels, []string{"namespace", "service"}) {
			Log.Debugf("Service[ %s ] NS[ %s ] Changed Labelss.", service.Name, service.Namespace)
			service.UpdateTime = tools.GetMillis()
			service.Labels = nil
			service.Labels = make(map[string]string)
			proc.CopyMapData(service.Labels, descLabels, []string{"namespace", "service"})
			Log.Debugf("Service[ %s ] NS[ %s ] Update Labelss.[ %+v ]",
				service.Name, service.Namespace, service.Labels)
		}
	case ksm.KsmFqNameServiceSpecExternalIP:
		err := proc.UpdateServiceExternalIP(service, descLabels, Log)
		if err != nil {
			service.Put()
			return err
		}
	case ksm.KsmFqNameServiceStatusLoadBalancerIngress:
		// TODO ::
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	service.Put()
	return nil
}

// KsmUpdateStatefulset statefulset
func (proc *Proc) KsmUpdateStatefulset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Statefulset. Can Not Found Namespace.[ %+v ]", desc)
	}
	name, isExist := descLabels["statefulset"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Statefulset. Can Not Found Name.[ %+v ]", desc)
	}

	var statefulset *ksm.Statefulset
	statefulset, err := ksmManager.GetStatefulset(metricData.CID, namespace, name)
	if err != nil {
		statefulset, err = ksmManager.AllocStatefulset()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Statefulset[ %s-%s ]. [ %s ]",
				namespace, name, err.Error())
		}
		statefulset.CID = metricData.CID
		statefulset.Namespace = namespace
		statefulset.Name = name
		statefulset.Add()
	}
	// 갱신 시간 업데이트
	statefulset.RefreshTime = tools.GetMillis()
	switch desc.FqName {
	case ksm.KsmFqNameStatefulsetCreated:
		if statefulset.Created == 0 {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Created Time[ %s ]", statefulset.Name,
				statefulset.Namespace, tools.PrintTimeFromSec(int64(desc.Value)))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Created = int64(desc.Value) * 1000
		} else if statefulset.Created != int64(desc.Value)*1000 {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Created Time. Before [ %s ] Currnet [ %s ]",
				statefulset.Name, statefulset.Namespace,
				tools.PrintTimeFromMilli(statefulset.Created), tools.PrintTimeFromSec(int64(desc.Value)))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Created = int64(desc.Value) * 1000
		}
	case ksm.KsmFqNameStatefulsetStatusReplicas:
		if statefulset.Status.Replicas != int(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Replicas.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetStatusReplicasCurrent:
		if statefulset.Status.ReplicasCurrent != int(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Replicas Current.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.ReplicasCurrent = int(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetStatusReplicasReady:
		if statefulset.Status.ReplicasReady != int(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Replicas Current.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.ReplicasReady = int(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetStatusReplicasUpdated:
		if statefulset.Status.ReplicasUpdated != int(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Replicas Updated.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.ReplicasUpdated = int(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetStatusObservedGeneration:
		if statefulset.Status.ObservedGeneration != int64(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Observed Generation.[ %d ]",
				statefulset.Name, statefulset.Namespace, int64(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.ObservedGeneration = int64(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetReplicas:
		if statefulset.Replicas != int(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Replicas.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Status.Replicas = int(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetMetadataGeneration:
		if statefulset.Metadata.Generation != int64(desc.Value) {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Metadata Generation.[ %d ]",
				statefulset.Name, statefulset.Namespace, int(desc.Value))
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Metadata.Generation = int64(desc.Value)
		}
	case ksm.KsmFqNameStatefulsetLabels:
		if !proc.CompareMapData(statefulset.Labels, descLabels, []string{"namespace", "statefulset"}) {
			Log.Debugf("statefulset[ %s ] NS[ %s ] Changed Labels.", statefulset.Name, statefulset.Namespace)
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Labels = nil
			statefulset.Labels = make(map[string]string)
			proc.CopyMapData(statefulset.Labels, descLabels, []string{"namespace", "statefulset"})
			Log.Debugf("statefulset[ %s ] NS[ %s ] Update Labels.[ %+v ]",
				statefulset.Name, statefulset.Namespace, statefulset.Labels)
		}
	case ksm.KsmFqNameStatefulsetStatusCurrentRevision:
		revision := descLabels["revision"]
		if statefulset.Revision.Current != revision {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Current Revision.[ %s ]",
				statefulset.Name, statefulset.Namespace, revision)
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Revision.Current = revision
		}
	case ksm.KsmFqNameStatefulsetStatusUpdateRevision:
		revision := descLabels["revision"]
		if statefulset.Revision.Update != revision {
			Log.Debugf("Statefulset[ %s ] NS[ %s ] Changed Status Update Revision.[ %s ]",
				statefulset.Name, statefulset.Namespace, revision)
			statefulset.UpdateTime = tools.GetMillis()
			statefulset.Revision.Update = revision
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	statefulset.Put()
	return nil
}

// CompareMapData compare map[string]string
func (proc *Proc) CompareMapData(left map[string]string, right map[string]string, ignoreKey []string) bool {
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
		isIgnored := false
		for _, ignore := range ignoreKey {
			if rk == ignore {
				isIgnored = true
				break
			}
		}
		if isIgnored {
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
func (proc *Proc) CopyMapData(left map[string]string, right map[string]string, ignoreKey []string) {
	for rightKey, rightValue := range right {
		isIgnored := false
		for _, ignore := range ignoreKey {
			if rightKey == ignore {
				isIgnored = true
				break
			}
		}
		if isIgnored {
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
	key := descLabels["key"]
	value := descLabels["value"]
	effect := descLabels["effect"]
	// Key가 동일한 경우를 처리
	for idx, taint := range node.Spec.Taints {
		// map에서 키/밸류에 대한 오류 체크는 compare단계에서 이미 진행하여 생략한다.
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
		Key:         key,
		Value:       value,
		Effect:      effect,
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

// UpdateOrInsertIngressPath update or insert ingressPath
func (proc *Proc) UpdateOrInsertIngressPath(ingress *ksm.Ingress, descLabels map[string]string, Log logrus.FieldLogger) error {
	host, isExist := descLabels["host"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found host.[ %+v ]", descLabels)
	}
	inputPath, isExist := descLabels["path"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found path.[ %+v ]", descLabels)
	}
	service, isExist := descLabels["service_name"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found service_name.[ %+v ]", descLabels)
	}
	port, isExist := descLabels["service_port"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found service_port.[ %+v ]", descLabels)
	}

	for rIdx, rule := range ingress.Spec.Rules {
		if host == rule.Host {
			ingress.Spec.Rules[rIdx].RefreshTime = tools.GetMillis()
			// 동일한 path가 있다면 외부로 리턴
			for pIdx, path := range rule.Paths {
				if path.Path == inputPath && path.ServiceName == service && path.ServicePort == port {
					ingress.Spec.Rules[rIdx].Paths[pIdx].RefreshTime = tools.GetMillis()
					return nil
				}
			}
			// 동일한 path가 없으므로 추가
			tmp := ksm.IngressPath{
				RefreshTime: tools.GetMillis(),
				UpdateTime:  tools.GetMillis(),
				Path:        inputPath,
				ServiceName: service,
				ServicePort: port,
			}
			ingress.Spec.Rules[rIdx].Paths = append(ingress.Spec.Rules[rIdx].Paths, tmp)
			Log.Debugf("Ingress[ %s ] NS[ %s ] Add Path. Host[ %s ] Path[ %s ] Service[ %s ] Port[ %s ]",
				ingress.Name, ingress.Namespace, host, inputPath, service, port)
			return nil
		}
	}

	// Rule 없는 경우
	rule := ksm.IngressRule{
		RefreshTime: tools.GetMillis(),
		UpdateTime:  tools.GetMillis(),
		Host:        host,
		Paths: []ksm.IngressPath{
			{
				RefreshTime: tools.GetMillis(),
				UpdateTime:  tools.GetMillis(),
				Path:        inputPath,
				ServiceName: service,
				ServicePort: port,
			},
		},
	}
	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	Log.Debugf("Ingress[ %s ] NS[ %s ] Add Host. Host[ %s ] Path[ %s ] Service[ %s ] Port[ %s ]",
		ingress.Name, ingress.Namespace, host, inputPath, service, port)
	return nil
}

// UpdateOrInsertIngressTLS update or insert ingress TLS
func (proc *Proc) UpdateOrInsertIngressTLS(ingress *ksm.Ingress, descLabels map[string]string, Log logrus.FieldLogger) error {
	inputHost, isExist := descLabels["tls_host"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found host.[ %+v ]", descLabels)
	}
	inputSecret, isExist := descLabels["secret"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found path.[ %+v ]", descLabels)
	}

	for tlsIdx, tls := range ingress.Spec.TLS {
		if inputSecret == tls.Secretname {
			ingress.Spec.TLS[tlsIdx].RefreshTime = tools.GetMillis()
			// SecretName, Host가 동일한 경우
			for hostIdx, host := range tls.Hosts {
				if inputHost == host.Host {
					ingress.Spec.TLS[tlsIdx].Hosts[hostIdx].RefreshTime = tools.GetMillis()
					return nil
				}
			}
			// Host가 없는 경우
			host := ksm.IngressHost{
				RefreshTime: tools.GetMillis(),
				UpdateTime:  tools.GetMillis(),
				Host:        inputHost,
			}
			ingress.Spec.TLS[tlsIdx].Hosts = append(ingress.Spec.TLS[tlsIdx].Hosts, host)
			Log.Debugf("Ingress[ %s ] NS[ %s ] Add TLS Host. SecretName[ %s ] Host[ %s ]",
				ingress.Name, ingress.Namespace, inputHost, inputSecret)
			return nil
		}
	}
	// SecretName이 없는 경우
	tls := ksm.IngressTLS{
		RefreshTime: tools.GetMillis(),
		UpdateTime:  tools.GetMillis(),
		Secretname:  inputSecret,
		Hosts: []ksm.IngressHost{
			{
				RefreshTime: tools.GetMillis(),
				UpdateTime:  tools.GetMillis(),
				Host:        inputHost,
			},
		},
	}
	ingress.Spec.TLS = append(ingress.Spec.TLS, tls)
	Log.Debugf("Ingress[ %s ] NS[ %s ] Add TLS. SecretName[ %s ] Host[ %s ]",
		ingress.Name, ingress.Namespace, inputSecret, inputHost)
	return nil
}

// UpdateOrInsertPVCAccessModes update pvc access modes
func (proc *Proc) UpdateOrInsertPVCAccessModes(pvc *ksm.PersistentVolumeClaim,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	modeName, isExist := descLabels["access_mode"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Persistentvolumeclaim. Can Not Found access_mode.[ %+v ]", descLabels)
	}
	for idx, am := range pvc.AccessModes {
		if am.AccessMode == modeName {
			pvc.AccessModes[idx].RefreshTime = tools.GetMillis()
			return nil
		}
	}
	accessMode := ksm.PVCAccessMode{
		RefreshTime: tools.GetMillis(),
		AccessMode:  modeName,
	}
	pvc.AccessModes = append(pvc.AccessModes, accessMode)
	Log.Debugf("PVC[ %s ] NS[ %s ] Add AccessMode.[ %s ]",
		pvc.Name, pvc.Namespace, modeName)
	return nil
}

// UpdatePVCCondition update pvc condition
func (proc *Proc) UpdatePVCCondition(pvc *ksm.PersistentVolumeClaim,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	condition, isExist := descLabels["condition"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Persistentvolumeclaim. Can Not Found condition.[ %+v ]", descLabels)
	}
	inputStatus, isExist := descLabels["status"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Persistentvolumeclaim. Can Not Found status.[ %+v ]", descLabels)
	}
	status, isExist := pvc.Status.Conditions[condition]
	if !isExist {
		pvc.Status.UpdateTime = tools.GetMillis()
		pvc.Status.RefreshTime = pvc.Status.UpdateTime
		pvc.Status.Conditions[condition] = inputStatus
		Log.Debugf("PVC[ %s ] NS[ %s ] Add Condition.[ %s ][ %s ]",
			pvc.Name, pvc.Namespace, condition, pvc.Status.Conditions[condition])
		return nil
	}
	if status != inputStatus {
		pvc.Status.UpdateTime = tools.GetMillis()
		pvc.Status.RefreshTime = pvc.Status.UpdateTime
		pvc.Status.Conditions[condition] = inputStatus
		Log.Debugf("PVC[ %s ] NS[ %s ] Changed Condition.[ %s ][ %s ]",
			pvc.Name, pvc.Namespace, condition, pvc.Status.Conditions[condition])
		return nil
	}
	pvc.Status.RefreshTime = tools.GetMillis()
	return nil
}

// UpdateOrInsertContainerInfo update or insert container info
func (proc *Proc) UpdateOrInsertContainerInfo(pod *ksm.Pod,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	cID, isExist := descLabels["container_id"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Container Info. Can Not Found container_id.[ %+v ]", descLabels)
	}
	cName, isExist := descLabels["container"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Container Info. Can Not Found container.[ %+v ]", descLabels)
	}
	// 기존 데이터 확인
	if container, isExist := pod.Containers[cName]; isExist {
		if !proc.CompareMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"}) {
			Log.Debugf("Pod[ %s ] NS[ %s ] Changed Container[ %s ] Info.", pod.Name, pod.Namespace, cName)
			pod.UpdateTime = tools.GetMillis()
			if container.ID != cID {
				container.ID = cID
				Log.Debugf("Pod[ %s ] NS[ %s ] Update Container[ %s ] ID[ %s ]",
					pod.Name, pod.Namespace, cName, cID)
			}
			container.ContainerInfo = nil
			container.ContainerInfo = make(map[string]string)
			proc.CopyMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"})
			pod.Containers[cName] = container
			Log.Debugf("Pod[ %s ] NS[ %s ] Update Container[ %s ] Info.[ %+v ]", pod.Name, pod.Namespace, cName, container.ContainerInfo)
		}
		return nil
	}
	// 신규 정보
	container := ksm.PodContainer{
		ID:            cID,
		ContainerInfo: map[string]string{},
	}
	proc.CopyMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"})
	pod.UpdateTime = tools.GetMillis()
	pod.Containers[cName] = container
	Log.Debugf("Pod[ %s ] NS[ %s ] Add Container[ %s ] Info.[ %+v ]", pod.Name, pod.Namespace, cName, container.ContainerInfo)
	return nil
}

// UpdateOrInsertInitContainerInfo update or insert container info
func (proc *Proc) UpdateOrInsertInitContainerInfo(pod *ksm.Pod,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	cID, isExist := descLabels["container_id"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Init Container Info. Can Not Found container_id.[ %+v ]", descLabels)
	}
	cName, isExist := descLabels["container"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Init Container Info. Can Not Found container.[ %+v ]", descLabels)
	}
	if container, isExist := pod.InitContainers[cName]; isExist {
		if !proc.CompareMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"}) {
			Log.Debugf("Pod[ %s ] NS[ %s ] Changed Init Container[ %s ] Info.", pod.Name, pod.Namespace, cName)
			pod.UpdateTime = tools.GetMillis()
			if container.ID != cID {
				container.ID = cID
				Log.Debugf("Pod[ %s ] NS[ %s ] Update Init Container[ %s ] ID[ %s ]",
					pod.Name, pod.Namespace, cName, cID)
			}
			container.ContainerInfo = nil
			container.ContainerInfo = make(map[string]string)
			proc.CopyMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"})
			Log.Debugf("Pod[ %s ] NS[ %s ] Update Init Container[ %s ] Info.[ %+v ]",
				pod.Name, pod.Namespace, cName, container.ContainerInfo)
			pod.InitContainers[cName] = container
		}
		return nil
	}
	// 신규 정보
	container := ksm.PodContainer{
		ID:            cID,
		ContainerInfo: map[string]string{},
	}
	proc.CopyMapData(container.ContainerInfo, descLabels, []string{"namespace", "pod"})
	pod.UpdateTime = tools.GetMillis()
	pod.InitContainers[cName] = container
	Log.Debugf("Pod[ %s ] NS[ %s ] Add Init Container[ %s ] Info.[ %+v ]",
		pod.Name, pod.Namespace, cName, container.ContainerInfo)
	return nil
}

// UpdatePodContainerStatus update pod container status
func (proc *Proc) UpdatePodContainerStatus(pod *ksm.Pod, isInitContainer bool,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	descLabels := desc.Labels.(map[string]string)
	cName, isExist := descLabels["container"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Container Status. Can Not Found container.[ %+v ]", descLabels)
	}
	var container ksm.PodContainer
	if isInitContainer {
		if _, isExist := pod.InitContainers[cName]; !isExist {
			pod.InitContainers[cName] = ksm.PodContainer{}
		}
		container = pod.InitContainers[cName]
	} else {
		if _, isExist := pod.Containers[cName]; !isExist {
			pod.Containers[cName] = ksm.PodContainer{}
		}
		container = pod.Containers[cName]
	}

	switch desc.FqName {
	case ksm.KsmFqNamePodContainerStatusWaiting, ksm.KsmFqNamePodInitContainerStatusWaiting:
		if container.Status.Waiting != int(desc.Value) {
			if isInitContainer {
				Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Waiting Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			} else {
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Waiting Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			}
			container.Status.Waiting = int(desc.Value)
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusWaitingReason,
		ksm.KsmFqNamePodInitContainerStatusWaitingReason:
		if int(desc.Value) == 1 {
			reason := descLabels["reason"]
			if container.Status.WaitingReason != reason {
				if isInitContainer {
					Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Waiting Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				} else {
					Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Waiting Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				}
				container.Status.WaitingReason = reason
			} else {
				return nil
			}
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusRunning,
		ksm.KsmFqNamePodInitContainerStatusRunning:
		if pod.Containers[cName].Status.Running != int(desc.Value) {
			if isInitContainer {
				Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Running Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			} else {
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Running Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			}
			container.Status.Running = int(desc.Value)
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusTerminated,
		ksm.KsmFqNamePodInitContainerStatusTerminated:
		if container.Status.Terminated != int(desc.Value) {
			if isInitContainer {
				Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Terminated Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			} else {
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Terminated Status[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			}
			container.Status.Terminated = int(desc.Value)
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusTerminatedReason,
		ksm.KsmFqNamePodInitContainerStatusTerminatedReason:
		if int(desc.Value) == 1 {
			reason := descLabels["reason"]
			if container.Status.TerminatedReason != reason {
				if isInitContainer {
					Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Terminated Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				} else {
					Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Terminated Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				}
				container.Status.TerminatedReason = reason
			} else {
				return nil
			}
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusLastTerminatedReason,
		ksm.KsmFqNamePodInitContainerStatusLastTerminatedReason:
		if int(desc.Value) == 1 {
			reason := descLabels["reason"]
			if container.Status.LastTerminatedReason != reason {
				if isInitContainer {
					Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Last Terminated Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				} else {
					Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Last Terminated Reason[ %s ].",
						pod.Name, pod.Namespace, cName, reason)
				}
				container.Status.LastTerminatedReason = reason
			} else {
				return nil
			}
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusReady,
		ksm.KsmFqNamePodInitContainerStatusReady:
		if container.Status.Ready != int(desc.Value) {
			if isInitContainer {
				Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Status Ready[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			} else {
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Status Ready[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			}
			container.Status.Ready = int(desc.Value)
		} else {
			return nil
		}

	case ksm.KsmFqNamePodContainerStatusRestartsTotal,
		ksm.KsmFqNamePodInitContainerStatusRestartsTotal:
		if container.Status.RestartTotal != int(desc.Value) {
			if isInitContainer {
				Log.Debugf("POD[ %s ] NS[ %s ] Init Container[ %s ] Changed Status Restart[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			} else {
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Status Restart[ %d ].",
					pod.Name, pod.Namespace, cName, int(desc.Value))
			}
			container.Status.RestartTotal = int(desc.Value)
		} else {
			return nil
		}
	}

	pod.UpdateTime = tools.GetMillis()
	if isInitContainer {
		pod.InitContainers[cName] = container
	} else {
		pod.Containers[cName] = container
	}
	return nil
}

// UpdatePodContainerResource ARGS : pod *ksm.Pod, isInitContainer bool, isRequest bool,
// desclabels map[string]string, value float64, Log
func (proc *Proc) UpdatePodContainerResource(pod *ksm.Pod, isInitContainer, isRequest bool,
	descLabels map[string]string, value float64, Log logrus.FieldLogger) error {
	cName, isExist := descLabels["container"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Container Resource. "+
			"Can Not Found container.[ %+v ]", descLabels)
	}
	var container ksm.PodContainer
	if isInitContainer {
		if _, isExist := pod.InitContainers[cName]; !isExist {
			pod.InitContainers[cName] = ksm.PodContainer{}
		}
		container = pod.InitContainers[cName]
	} else {
		if _, isExist := pod.Containers[cName]; !isExist {
			pod.Containers[cName] = ksm.PodContainer{}
		}
		container = pod.Containers[cName]
	}
	resType := descLabels["resource"]
	strValue := fmt.Sprintf("%f", value)

	if isRequest {
		// resource request
		if resType == "cpu" {
			if container.RequestResource.CPU != strValue {
				container.RequestResource.CPU = strValue
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Resource Request CPU[ %s ].",
					pod.Name, pod.Namespace, cName, container.RequestResource.CPU)
			} else {
				return nil
			}
		} else if resType == "memory" {
			if container.RequestResource.Mem != strValue {
				container.RequestResource.Mem = strValue
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Resource Request MEM[ %s ].",
					pod.Name, pod.Namespace, cName, container.RequestResource.Mem)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("WARN. Not Supported Resource Type[ %s ]", resType)
		}
	} else {
		// resource limit
		if resType == "cpu" {
			if container.LimitResource.CPU != strValue {
				container.LimitResource.CPU = strValue
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Resource Limit CPU[ %s ].",
					pod.Name, pod.Namespace, cName, container.LimitResource.CPU)
			} else {
				return nil
			}
		} else if resType == "memory" {
			if container.LimitResource.Mem != strValue {
				container.LimitResource.Mem = strValue
				Log.Debugf("POD[ %s ] NS[ %s ] Container[ %s ] Changed Resource Limit MEM[ %s ].",
					pod.Name, pod.Namespace, cName, container.LimitResource.Mem)
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("WARN. Not Supported Resource Type[ %s ]", resType)
		}
	}

	// save
	Log.Debugf("Update Pod Updatetime")
	pod.UpdateTime = tools.GetMillis()
	if isInitContainer {
		pod.InitContainers[cName] = container
	} else {
		pod.Containers[cName] = container
	}
	return nil
}

// UpdateOrInsertPodVolumes update or insert pod - pvc info
func (proc *Proc) UpdateOrInsertPodVolumes(pod *ksm.Pod,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	pvcName, isExist := descLabels["persistentvolumeclaim"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Persistentvolumeclaims. "+
			"Can Not Found persistentvolumeclaim.[ %+v ]", descLabels)
	}
	volName, isExist := descLabels["volume"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Pod Persistentvolumeclaims. "+
			"Can Not Found volume.[ %+v ]", descLabels)
	}
	for _, pvc := range pod.PVCs {
		if pvc.Volume == volName && pvc.PVCName == pvcName {
			return nil
		}
	}
	pod.PVCs = append(pod.PVCs, ksm.PodPersistentVolumeClaim{
		Volume:  volName,
		PVCName: pvcName,
	})
	Log.Debugf("POD[ %s ] NS[ %s ] Insert Persistentvolumeclaim[ %+v ]",
		pod.Name, pod.Namespace, pod.PVCs)
	return nil
}

// UpdateServiceExternalIP update service spec external ips
func (proc *Proc) UpdateServiceExternalIP(service *ksm.Service,
	descLabels map[string]string, Log logrus.FieldLogger) error {
	externalIP, isExist := descLabels["external_ip"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Service External IP. "+
			"Can Not Found external_ip.[ %+v ]", descLabels)
	}
	for _, ip := range service.Spec.ExternalIP {
		if ip == externalIP {
			return nil
		}
	}
	service.UpdateTime = tools.GetMillis()
	service.Spec.ExternalIP = append(service.Spec.ExternalIP, externalIP)
	return nil
}
