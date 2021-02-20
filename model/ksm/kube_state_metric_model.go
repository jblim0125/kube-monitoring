package ksm

const (
	// KsmStatusTrue kube state metric status condition
	KsmStatusTrue = "true"
	// KsmStatusFalse kube state metric status condition
	KsmStatusFalse = "false"
	// KsmStatusUnknow kube state metric status condition
	KsmStatusUnknow = "unknow"
)

// FqName

// Node
const (
	KsmFqNameNodeInfo                    = "kube_node_info"
	KsmFqNameNodeCreated                 = "kube_node_created"
	KsmFqNameNodeLabels                  = "kube_node_labels"
	KsmFqNameNodeRole                    = "kube_node_role"
	KsmFqNameNodeSpecUnschedulable       = "kube_node_spec_unschedulable"
	KsmFqNameNodeSpecTaint               = "kube_node_spec_taint"
	KsmFqNameNodeStatusCondition         = "kube_node_status_condition"
	KsmFqNameNodeStatusPhase             = "kube_node_status_phase"
	KsmFqNameNodeStatusCapacityPods      = "kube_node_status_capacity_pods"
	KsmFqNameNodeStatusCapacityCPU       = "kube_node_status_capacity_cpu_cores"
	KsmFqNameNodeStatusCapacityMemory    = "kube_node_status_capacity_memory_bytes"
	KsmFqNameNodeStatusAllocatablePods   = "kube_node_status_allocatable_pods"
	KsmFqNameNodeStatusAllocatableCPU    = "kube_node_status_allocatable_cpu_cores"
	KsmFqNameNodeStatusAllocatableMemory = "kube_node_status_allocatable_memory_bytes"
	// KubeMetricNameNodeStatusCapacity          = "kube_node_status_capacity"
	// KubeMetricNameNodeStatusAllocatable       = "kube_node_status_allocatable"
)

// node status condition
const (
	KsmNodeStatusConditionNetworkUnavailable = "NetworkUnavailable"
	KsmNodeStatusConditionMemoryPressure     = "MemoryPressure"
	KsmNodeStatusConditionDiskPressure       = "DiskPressure"
	KsmNodeStatusConditionPIDPressure        = "PIDPressure"
	KsmNodeStatusConditionReady              = "Ready"
)

// Namespace
const (
	KsmFqNameNamespaceCreated     = "kube_namespace_created"
	KsmFqNameNamespaceLabels      = "kube_namespace_labels"
	KsmFqNameNamespaceStatusPhase = "kube_namespace_status_phase"
	// KsmFqNameNamespaceStatusCondition = "kube_namespace_status_condition"
)

// namespace phase
const (
	KsmNamespaceStatusPhaseActive      = "Active"
	KsmNamespaceStatusPhaseTerminating = "Terminating"
)

// Configmap
const (
	KsmFqNameConfigMapInfo    = "kube_configmap_info"
	KsmFqNameConfigMapCreated = "kube_configmap_created"
	KsmFqNameConfigMapVersion = "kube_configmap_metadata_resource_version"
)

// Cronjob
const (
	KsmFqNameCronjobLabels                      = "kube_cronjob_labels"
	KsmFqNameCronjobInfo                        = "kube_cronjob_info"
	KsmFqNameCronjobCreated                     = "kube_cronjob_created"
	KsmFqNameCronjobStatus                      = "kube_cronjob_status_active"
	KsmFqNameCronjobLastScheduleTime            = "kube_cronjob_status_last_schedule_time"
	KsmFqNameCronjobSpecSuspend                 = "kube_cronjob_spec_suspend"
	KsmFqNameCronjobSpecStartingDeadlineSeconds = "kube_cronjob_spec_starting_deadline_seconds"
	KsmFqNameCronjobNextSchedule                = "kube_cronjob_next_schedule_time"
)

// Daemonset
const (
	KsmFqNameDaemonsetCreated            = "kube_daemonset_created"
	KsmFqNameDaemonsetStatusCurrent      = "kube_daemonset_status_current_number_scheduled"
	KsmFqNameDaemonsetStatusDesired      = "kube_daemonset_status_desired_number_scheduled"
	KsmFqNameDaemonsetStatusAvailable    = "kube_daemonset_status_number_available"
	KsmFqNameDaemonsetStatusMisscheduled = "kube_daemonset_status_number_misscheduled"
	KsmFqNameDaemonsetStatusReady        = "kube_daemonset_status_number_ready"
	KsmFqNameDaemonsetStatusUnavailable  = "kube_daemonset_status_number_unavailable"
	KsmFqNameDaemonsetUpdated            = "kube_daemonset_updated_number_scheduled"
	KsmFqNameDaemonsetMetadataGeneration = "kube_daemonset_metadata_generation"
	KsmFqNameDaemonsetLabels             = "kube_daemonset_labels"
)

// Deployment
const (
	KsmFqNameDeploymentCreated                                 = "kube_deployment_created"
	KsmFqNameDeploymentStatusReplicas                          = "kube_deployment_status_replicas"
	KsmFqNameDeploymentStatusReplicasAvailable                 = "kube_deployment_status_replicas_available"
	KsmFqNameDeploymentStatusReplicasUnavailable               = "kube_deployment_status_replicas_unavailable"
	KsmFqNameDeploymentStatusReplicasUpdated                   = "kube_deployment_status_replicas_updated"
	KsmFqNameDeploymentStatusObserved                          = "kube_deployment_status_observed_generation"
	KsmFqNameDeploymentStatusCondition                         = "kube_deployment_status_condition"
	KsmFqNameDeploymentSpecReplicas                            = "kube_deployment_spec_replicas"
	KsmFqNameDeploymentSpecPaused                              = "kube_deployment_spec_paused"
	KsmFqNameDeploymentSpecStrategyRollingupdateMaxUnavailable = "kube_deployment_spec_strategy_rollingupdate_max_unavailable"
	KsmFqNameDeploymentSpecStrategyRollingupdateMaxSurge       = "kube_deployment_spec_strategy_rollingupdate_max_surge"
	KsmFqNameDeploymentMetadataGeneration                      = "kube_deployment_metadata_generation"
	KsmFqNameDeploymentLabels                                  = "kube_deployment_labels"
)

// Endpoint : 2021.02.09 시점에서는 무시

// HPA : 2021.02.09 시점에서는 무시

// Ingress
const (
	KsmFqNameIngressInfo                    = "kube_ingress_info"
	KsmFqNameIngressLabels                  = "kube_ingress_labels"
	KsmFqNameIngressCreated                 = "kube_ingress_created"
	KsmFqNameIngressMetadataResourceVersion = "kube_ingress_metadata_resource_version"
	KsmFqNameIngressPath                    = "kube_ingress_path"
	KsmFqNameIngressTLS                     = "kube_ingress_tls"
)

// Job
const (
	KsmFqNameJobLabels                    = "kube_job_labels"
	KsmFqNameJobInfo                      = "kube_job_info"
	KsmFqNameJobCreated                   = "kube_job_created"
	KsmFqNameJobSpecParallelism           = "kube_job_spec_parallelism"
	KsmFqNameJobSpecCompletions           = "kube_job_spec_completions"
	KsmFqNameJobSpecActiveDeadlineSeconds = "kube_job_spec_active_deadline_seconds"
	KsmFqNameJobStatusSucceeded           = "kube_job_status_succeeded"
	KsmFqNameJobStatusFailed              = "kube_job_status_failed"
	KsmFqNameJobStatusActive              = "kube_job_status_active"
	KsmFqNameJobComplete                  = "kube_job_complete"
	KsmFqNameJobFailed                    = "kube_job_failed"
	KsmFqNameJobStatusStartTime           = "kube_job_status_start_time"
	KsmFqNameJobStatusCompletionTime      = "kube_job_status_completion_time"
	KsmFqNameJobOwner                     = "kube_job_owner"
)

// Limitrange : 2021.02.09 시점에서는 무시
// kube_limitrange
// kube_limitrange_created

// MutatingWebhookConfiguration : 2021.02.09 시점에서는 무시
// kube_mutatingwebhookconfiguration_info
// kube_mutatingwebhookconfiguration_created
// kube_mutatingwebhookconfiguration_metadata_resource_version

// NetworkPolicy : 2021.02.09 시점에서는 무시
// kube_networkpolicy_created
// kube_networkpolicy_labels
// kube_networkpolicy_spec_ingress_rules
// kube_networkpolicy_spec_egress_rules

// PersistentVolumeClaim
const (
	KsmFqNamePersistentvolumeclaimLabels                = "kube_persistentvolumeclaim_labels"
	KsmFqNamePersistentvolumeclaimInfo                  = "kube_persistentvolumeclaim_info"
	KsmFqNamePersistentvolumeclaimStatusPhase           = "kube_persistentvolumeclaim_status_phase"
	KsmFqNamePersistentvolumeclaimResourceRequestsBytes = "kube_persistentvolumeclaim_resource_requests_storage_bytes"
	KsmFqNamePersistentvolumeclaimAccessMode            = "kube_persistentvolumeclaim_access_mode"
	KsmFqNamePersistentvolumeclaimStatusCondition       = "kube_persistentvolumeclaim_status_condition"
)

// PersistentVolume
const (
	KsmFqNamePersistentvolumeLabels        = "kube_persistentvolume_labels"
	KsmFqNamePersistentvolumeStatusPhase   = "kube_persistentvolume_status_phase"
	KsmFqNamePersistentvolumeInfo          = "kube_persistentvolume_info"
	KsmFqNamePersistentvolumeCapacityBytes = "kube_persistentvolume_capacity_bytes"
)

// pv status
const (
	KsmPvStatusPhasePending   = "Pending"
	KsmPvStatusPhaseAvailable = "Available"
	KsmPvStatusPhaseBound     = "Bound"
	KsmPvStatusPhaseReleased  = "Released"
	KsmPvStatusPhaseFailed    = "Failed"
)

// Pod Disruption Budget : 2021.02.09 시점에서는 무시
// kube_poddisruptionbudget_created
// kube_poddisruptionbudget_status_current_healthy
// kube_poddisruptionbudget_status_desired_healthy
// kube_poddisruptionbudget_status_pod_disruptions_allowed
// kube_poddisruptionbudget_status_expected_pods
// kube_poddisruptionbudget_status_observed_generatio

// Pod
const (
	KsmFqNamePodInfo                                      = "kube_pod_info"
	KsmFqNamePodStartTime                                 = "kube_pod_start_time"
	KsmFqNamePodCompletionTime                            = "kube_pod_completion_time"
	KsmFqNamePodOwner                                     = "kube_pod_owner"
	KsmFqNamePodLabels                                    = "kube_pod_labels"
	KsmFqNamePodCreated                                   = "kube_pod_created"
	KsmFqNamePodRestartPolicy                             = "kube_pod_restart_policy"
	KsmFqNamePodStatusscheduledTime                       = "kube_pod_status_scheduled_time"
	KsmFqNamePodStatusunschedulable                       = "kube_pod_status_unschedulable"
	KsmFqNamePodStatusphase                               = "kube_pod_status_phase"
	KsmFqNamePodStatusready                               = "kube_pod_status_ready"
	KsmFqNamePodStatusScheduled                           = "kube_pod_status_scheduled"
	KsmFqNamePodContainerInfo                             = "kube_pod_container_info"
	KsmFqNamePodInitContainerInfo                         = "kube_pod_init_container_info"
	KsmFqNamePodInitContainerStatusWaiting                = "kube_pod_init_container_status_waiting"
	KsmFqNamePodInitContainerStatusWaitingReason          = "kube_pod_init_container_status_waiting_reason"
	KsmFqNamePodInitContainerStatusRunning                = "kube_pod_init_container_status_running"
	KsmFqNamePodInitContainerStatusTerminated             = "kube_pod_init_container_status_terminated"
	KsmFqNamePodInitContainerStatusTerminatedReason       = "kube_pod_init_container_status_terminated_reason"
	KsmFqNamePodInitContainerStatusLastTerminatedReason   = "kube_pod_init_container_status_last_terminated_reason"
	KsmFqNamePodInitContainerStatusReady                  = "kube_pod_init_container_status_ready"
	KsmFqNamePodInitContainerStatusRestartsTotal          = "kube_pod_init_container_status_restarts_total"
	KsmFqNamePodInitContainerResourceLimits               = "kube_pod_init_container_resource_limits"
	KsmFqNamePodContainerStatusWaiting                    = "kube_pod_container_status_waiting"
	KsmFqNamePodContainerStatusWaitingReason              = "kube_pod_container_status_waiting_reason"
	KsmFqNamePodContainerStatusRunning                    = "kube_pod_container_status_running"
	KsmFqNamePodContainerStatusTerminated                 = "kube_pod_container_status_terminated"
	KsmFqNamePodContainerStatusTerminatedReason           = "kube_pod_container_status_terminated_reason"
	KsmFqNamePodContainerStatusLastTerminatedReason       = "kube_pod_container_status_last_terminated_reason"
	KsmFqNamePodContainerStatusReady                      = "kube_pod_container_status_ready"
	KsmFqNamePodContainerStatusRestartsTotal              = "kube_pod_container_status_restarts_total"
	KsmFqNamePodContainerResourceRequests                 = "kube_pod_container_resource_requests"
	KsmFqNamePodContainerResourceLimits                   = "kube_pod_container_resource_limits"
	KsmFqNamePodContainerResourceRequestsCPU              = "kube_pod_container_resource_requests_cpu_cores"
	KsmFqNamePodContainerResourceRequestsMemory           = "kube_pod_container_resource_requests_memory_bytes"
	KsmFqNamePodContainerResourceLimitsCPU                = "kube_pod_container_resource_limits_cpu_cores"
	KsmFqNamePodContainerResourceLimitsMemory             = "kube_pod_container_resource_limits_memory_bytes"
	KsmFqNamePodSpecVolumesPersistentvolumeclaimsInfo     = "kube_pod_spec_volumes_persistentvolumeclaims_info"
	KsmFqNamePodSpecVolumesPersistentvolumeclaimsReadonly = "kube_pod_spec_volumes_persistentvolumeclaims_readonly"
)

// Container Waiting Reason
const (
	KsmContainerCreating          = "ContainerCreating"
	KsmCrashLoopBackOff           = "CrashLoopBackOff"
	KsmCreateContainerConfigError = "CreateContainerConfigError"
	KsmErrImagePull               = "ErrImagePull"
	KsmImagePullBackOff           = "ImagePullBackOff"
	KsmCreateContainerError       = "CreateContainerError"
	KsmInvalidImageName           = "InvalidImageName"
)

// Container Terminated Reason
const (
	KsmOOMKilled          = "OOMKilled"
	KsmCompleted          = "Completed"
	KsmError              = "Error"
	KsmContainerCannotRun = "ContainerCannotRun"
	KsmDeadlineExceeded   = "DeadlineExceeded"
	KsmEvicted            = "Evicted"
)

// Replicaset
const (
	KsmFqNameReplicasetCreated                    = "kube_replicaset_created"
	KsmFqNameReplicasetStatusReplicas             = "kube_replicaset_status_replicas"
	KsmFqNameReplicasetStatusFullyLabeledReplicas = "kube_replicaset_status_fully_labeled_replicas"
	KsmFqNameReplicasetStatusReadyReplicas        = "kube_replicaset_status_ready_replicas"
	KsmFqNameReplicasetStatusObservedGeneration   = "kube_replicaset_status_observed_generation"
	KsmFqNameReplicasetSpecReplicas               = "kube_replicaset_spec_replicas"
	KsmFqNameReplicasetMetadataGeneration         = "kube_replicaset_metadata_generation"
	KsmFqNameReplicasetOwner                      = "kube_replicaset_owner"
	KsmFqNameReplicasetLabels                     = "kube_replicaset_labels"
)

// Replicationcontroller
const (
	KsmFqNameReplicationcontrollerCreated                    = "kube_replicationcontroller_created"
	KsmFqNameReplicationcontrollerStatusReplicas             = "kube_replicationcontroller_status_replicas"
	KsmFqNameReplicationcontrollerStatusFullyLabeledReplicas = "kube_replicationcontroller_status_fully_labeled_replicas"
	KsmFqNameReplicationcontrollerStatusReadyReplicas        = "kube_replicationcontroller_status_ready_replicas"
	KsmFqNameReplicationcontrollerStatusAvailableReplicas    = "kube_replicationcontroller_status_available_replicas"
	KsmFqNameReplicationcontrollerStatusObservedGeneration   = "kube_replicationcontroller_status_observed_generation"
	KsmFqNameReplicationcontrollerSpecReplicas               = "kube_replicationcontroller_spec_replicas"
	KsmFqNameReplicationcontrollerMetadataGeneration         = "kube_replicationcontroller_metadata_generation"
)

// ResourceQuota
const (
	KsmFqNameResourcequotaCreated = "kube_resourcequota_created"
	KsmFqNameResourcequota        = "kube_resourcequota"
)

// Secret
const (
	KsmFqNameSecretInfo                    = "kube_secret_info"
	KsmFqNameSecretType                    = "kube_secret_type"
	KsmFqNameSecretLabels                  = "kube_secret_labels"
	KsmFqNameSecretCreated                 = "kube_secret_created"
	KsmFqNameSecretMetadataResourceVersion = "kube_secret_metadata_resource_version"
)

// Service
const (
	KsmFqNameServiceInfo                      = "kube_service_info"
	KsmFqNameServiceCreated                   = "kube_service_created"
	KsmFqNameServiceSpecType                  = "kube_service_spec_type"
	KsmFqNameServiceLabels                    = "kube_service_labels"
	KsmFqNameServiceSpecExternalIP            = "kube_service_spec_external_ip"
	KsmFqNameServiceStatusLoadBalancerIngress = "kube_service_status_load_balancer_ingress"
)

// Statefulset
const (
	KsmFqNameStatefulsetCreated                  = "kube_statefulset_created"
	KsmFqNameStatefulsetStatusReplicas           = "kube_statefulset_status_replicas"
	KsmFqNameStatefulsetStatusReplicasCurrent    = "kube_statefulset_status_replicas_current"
	KsmFqNameStatefulsetStatusReplicasReady      = "kube_statefulset_status_replicas_ready"
	KsmFqNameStatefulsetStatusReplicasUpdated    = "kube_statefulset_status_replicas_updated"
	KsmFqNameStatefulsetStatusObservedGeneration = "kube_statefulset_status_observed_generation"
	KsmFqNameStatefulsetReplicas                 = "kube_statefulset_replicas"
	KsmFqNameStatefulsetMetadataGeneration       = "kube_statefulset_metadata_generation"
	KsmFqNameStatefulsetLabels                   = "kube_statefulset_labels"
	KsmFqNameStatefulsetStatusCurrentRevision    = "kube_statefulset_status_current_revision"
	KsmFqNameStatefulsetStatusUpdateRevision     = "kube_statefulset_status_update_revision"
)

// Storageclass
const (
	KsmFqNameStorageclassInfo    = "kube_storageclass_info"
	KsmFqNameStorageclassCreated = "kube_storageclass_created"
	KsmFqNameStorageclassLabels  = "kube_storageclass_labels"
)

// Validating Webhook Configuration : 2021.02.09 시점에서는 무시
// kube_validatingwebhookconfiguration_info
// kube_validatingwebhookconfiguration_created
// kube_validatingwebhookconfiguration_metadata_resource_version

// Volume Attachment  : 2021.02.09 시점에서는 무시
// kube_volumeattachment_labels
// kube_volumeattachment_info
// kube_volumeattachment_created
// kube_volumeattachment_spec_source_persistentvolume
// kube_volumeattachment_status_attached
// kube_volumeattachment_status_attachment_metadata

// Labels Key
// const (
// 	KubeMetricKey_Namespace = "namespace"
// 	KubeMetricKey_Configmap = "configmap"
// 	KubeMetricKey_Daemonset = "daemonset"
// 	KubeMetricKey_Deployment= "deployment"
// 	KubeMetricKey_Condition  = "condition"
// 	KubeMetricKey_Status = "status"
// 	KubeMetricKey_Endpoint = "endpoint"
// 	KubeMetricKey_Ingress= "ingress"
// 	KubeMetricKey_Host= "host"
// 	KubeMetricKey_Path= "path"
// 	KubeMetricKey_ServiceName= "service_name"
// 	KubeMetricKey_ServicePort= "service_port"
// 	KubeMetricKey_Phase= "phase"
// 	KubeMetricKey_Node = "node"
// 	KubeMetricKey_Role= "role"
// 	KubeMetricKey_Resource= "resource"
// 	KubeMetricKey_Unit= "unit"
// 	KubeMetricKey_persistentvolumeclaim= "persistentvolumeclaim"
// 	KubeMetricKey_storageclass= "storageclass"
// 	KubeMetricKey_volumename= "volumename"
// 	KubeMetricKey_persistentvolume= "persistentvolume"
// )
