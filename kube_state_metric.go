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
		return fmt.Errorf("ERROR. Failed to Parse Configmap. Can Not Found Namespace.[ %+v ].", desc)
	}
	name, isExist := descLabels["configmap"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Configmap. Can Not Found Name.[ %+v ].", desc)
	}

	var configmap *ksm.Configmap
	configmap, err := ksmManager.GetConfigmap(metricData.CID, namespace, name)
	if err != nil {
		configmap, err = ksmManager.AllocConfigmap()
		if err != nil {
			return fmt.Errorf("ERROR. Failed to Update Configmap[ %s-%s ]. [ %s ]", namespace, name, err.Error())
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
		return fmt.Errorf("ERROR. Failed to Parse Cronjob. Can Not Found Namespace.[ %+v ].", desc)
	}
	name, isExist := descLabels["cronjob"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Cronjob. Can Not Found Name.[ %+v ].", desc)
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

// KsmUpdateDeamonset daemonset
func (proc *Proc) KsmUpdateDaemonset(ksmManager *ksm.Manager, metricData *model.MetricData,
	desc model.MetricDesc, Log logrus.FieldLogger) error {
	// Log.Debugf("Update Daemonset State Metric")

	descLabels := desc.Labels.(map[string]string)
	namespace, isExist := descLabels["namespace"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Daemonset. Can Not Found Namespace.[ %+v ].", desc)
	}
	name, isExist := descLabels["daemonset"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Daemonset. Can Not Found Name.[ %+v ].", desc)
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
		return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found Namespace.[ %+v ].", desc)
	}
	name, isExist := descLabels["deployment"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found Name.[ %+v ].", desc)
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
			return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found contition.[ %+v ].", desc)
		}
		status, isExist := descLabels["status"]
		if !isExist {
			deploy.Put()
			return fmt.Errorf("ERROR. Failed to Parse Deploy. Can Not Found status.[ %+v ].", desc)
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
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found Namespace.[ %+v ].", desc)
	}
	name, isExist := descLabels["ingress"]
	if !isExist {
		return fmt.Errorf("ERROR. Failed to Parse Ingressh. Can Not Found Name.[ %+v ].", desc)
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
		// tlsHost, isExist := descLabels["tls_host"]
		// secret, isExist := descLabels["secret"]
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
		Log.Debugf("Node[ %s ] Status Phase.[ %+v ][ %f ]", node.Name, descLabels, desc.Value)
		node.UpdateTime = tools.GetMillis()
	case ksm.KsmFqNameNodeStatusCapacityPods:
		if node.Status.Capacity.Pods != int32(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Capacity Pods.[ %d ]", node.Name, int32(desc.Value))
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.Pods = int32(desc.Value)
		}
	case ksm.KsmFqNameNodeStatusCapacityCPU:
		if node.Status.Capacity.Cores != int32(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Capacity Cores.[ %d ]", node.Name, int32(desc.Value))
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.Cores = int32(desc.Value)
		}
	case ksm.KsmFqNameNodeStatusCapacityMemory:
		if node.Status.Capacity.MemoryBytes != int64(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Capacity Mem Bytes.[ %d ]GiB", node.Name, int64(desc.Value)/1000/1000/1000)
			node.UpdateTime = tools.GetMillis()
			node.Status.Capacity.MemoryBytes = int64(desc.Value)
		}
	case ksm.KsmFqNameNodeStatusAllocatablePods:
		if node.Status.Allocatable.Pods != int32(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Allocatable Pods.[ %d ]", node.Name, int32(desc.Value))
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.Pods = int32(desc.Value)
		}
	case ksm.KsmFqNameNodeStatusAllocatableCPU:
		if node.Status.Allocatable.Cores != int32(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Allocatable Cores.[ %d ]", node.Name, int32(desc.Value))
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.Cores = int32(desc.Value)
		}
	case ksm.KsmFqNameNodeStatusAllocatableMemory:
		if node.Status.Allocatable.MemoryBytes != int64(desc.Value) {
			Log.Debugf("Node[ %s ] Changed Allocatable Mem Bytes.[ %d ]GiB", node.Name, int64(desc.Value)/1000/1000/1000)
			node.UpdateTime = tools.GetMillis()
			node.Status.Allocatable.MemoryBytes = int64(desc.Value)
		}
	default:
		Log.Errorf("ERROR. Not Supported FqName[ %s ][ %+v ]", desc.FqName, desc)
	}
	node.Put()
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
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found host.[ %+v ].", descLabels)
	}
	path, isExist := descLabels["path"]
	if !isExist {
		ingress.Put()
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found path.[ %+v ].", descLabels)
	}
	service, isExist := descLabels["service_name"]
	if !isExist {
		ingress.Put()
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found service_name.[ %+v ].", descLabels)
	}
	port, isExist := descLabels["service_port"]
	if !isExist {
		ingress.Put()
		return fmt.Errorf("ERROR. Failed to Parse Ingress. Can Not Found service_port.[ %+v ].", descLabels)
	}
	isExist = false
	// Key가 동일한 경우를 처리
	for idx, inPath := range ingress.Path {
		if host == inPath.Host && path == inPath.Path {
			ingress.Path[idx].RefreshTime = tools.GetMillis()
			if inPath.ServiceName != service || inPath.ServicePort != port {
				Log.Debugf("Ingress[ %s ] NS[ %s ] Changed Path. Before[ %+v ] After Host[ %s ] Path[ %s ] Service[ %s ] Port[ %s ]",
					ingress.Name, ingress.Namespace, inPath, host, path, service, port)
				ingress.Path[idx].ServiceName = service
				ingress.Path[idx].ServicePort = port
			}
			isExist = true
			break
		}
	}
	// 추가된 경우를 처리
	if !isExist {
		tmp := ksm.IngressPath{
			RefreshTime: tools.GetMillis(),
			Host:        host,
			Path:        path,
			ServiceName: service,
			ServicePort: port,
		}
		ingress.Path = append(ingress.Path, tmp)
		Log.Debugf("Ingress[ %s ] NS[ %s ] Add Path. Host[ %s ] Path[ %s ] Service[ %s ] Port[ %s ]",
			ingress.Name, ingress.Namespace, host, path, service, port)
	}
	return nil
}
