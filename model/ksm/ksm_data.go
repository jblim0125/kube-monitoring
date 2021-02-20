package ksm

// Cluster cluster info
type Cluster struct {
	used   int32
	cnt    int32
	ID     string `yaml:"cID"`
	Name   string `yaml:"cName"`
	KsmURL string `yaml:"ksmURL"`
}

// Clone cluster deep copy
func (cluster *Cluster) Clone() Cluster {
	ret := Cluster{
		ID:     cluster.ID,
		Name:   cluster.Name,
		KsmURL: cluster.KsmURL,
	}
	return ret
}

// Node kube state metric node info
type Node struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	NodeInfo    map[string]string `yaml:"nodeInfo"`
	Created     int64             `yaml:"created"`
	Labels      map[string]string `yaml:"labels"`
	Roles       map[string]string `yaml:"roles"`
	Spec        Nodespec          `yaml:"spec"`
	Status      NodeStatus        `yaml:"status"`
}

// Nodespec node spec
type Nodespec struct {
	Unschedulable bool        `yaml:"unschedulable"`
	reserved      [3]byte     `yaml:"reserved"`
	Taints        []NodeTaint `yaml:"taints"`
}

// NodeTaint node taint
type NodeTaint struct {
	RefreshTime int64  `yaml:"refreshTime"`
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Effect      string `yaml:"effect"`
}

// NodeStatus node status
type NodeStatus struct {
	Condition           NodeCondition   `yaml:"condition"`
	Capacity            NodeCapacity    `yaml:"capacity"`
	Allocatable         NodeAllocatable `yaml:"allocatable"`
	NumPods             int             `yaml:"Pods"`
	PodsRequestResource NodePodResource `yaml:"requestResource"`
	PodsLimitResource   NodePodResource `yaml:"limitResource"`
}

// NodeCondition node condition
type NodeCondition struct {
	NetworkUnavailable string `yaml:"networkUnavailable"`
	MemoryPressure     string `yaml:"memoryPressure"`
	DiskPressure       string `yaml:"diskPressure"`
	PIDPressure        string `yaml:"pidPressure"`
	Ready              string `yaml:"ready"`
}

// NodeCapacity node capacity ( 노드 전체의 자원 )
type NodeCapacity struct {
	Pods        int32 `yaml:"pods"`
	Cores       int32 `yaml:"cores"`
	MemoryBytes int64 `yaml:"memorybytes"`
}

// NodePodResource pod resource in node
type NodePodResource struct {
	Cores       int32 `yaml:"cores"`
	MemoryBytes int64 `yaml:"memorybytes"`
}

// NodeAllocatable node allocatable
// ( 시스템과 k8s daemon등으로 할당한 자원을 제외한
//   POD만을 위한 자원의 크기 )
type NodeAllocatable NodeCapacity

// Namespace namespace
type Namespace struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Created     int64             `yaml:"created"`
	Labels      map[string]string `yaml:"labels"`
	Phase       string            `yaml:"phase"`
}

// MetadataResVersion metadata_resource_version
type MetadataResVersion struct {
	ResourceVersion int64 `yaml:"resourceVersion"`
}

// MetadataGen metadata_generation
type MetadataGen struct {
	Generation int32 `yaml:"generation"`
}

// Configmap kube state metric info
type Configmap struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Created     int64             `yaml:"created"`
	Metadata    ConfigmapMetadata `yaml:"metadata"`
}

// ConfigmapMetadata metadata
type ConfigmapMetadata MetadataResVersion

// Cronjob cronjob info
type Cronjob struct {
	used             int32
	cnt              int32
	UpdateTime       int64             `yaml:"updateTime"`
	RefreshTime      int64             `yaml:"refreshTime"`
	CID              string            `yaml:"cID"`
	Name             string            `yaml:"name"`
	Namespace        string            `yaml:"namespace"`
	Labels           map[string]string `yaml:"labels"`
	Info             map[string]string `yaml:"schedule"`
	Created          int64             `yaml:"created"`
	Status           CronjobStatus     `yaml:"status"`
	Spec             CronjobSpec       `yaml:"spec"`
	NextScheduleTime int64             `yaml:"nextScheduleTime"`
}

// CronjobStatus status
type CronjobStatus struct {
	Active           bool    `yaml:"active"`
	reserved         [3]byte `yaml:"reserved"`
	LastScheduleTime int64   `yaml:"lastScheduleTime"`
}

// CronjobSpec spec
type CronjobSpec struct {
	Suspend                 int32 `yaml:"suspend"`
	StartingDeadlineSeconds int64 `yaml:"startingDeadlineSeconds"`
	NextScheduletime        int64 `yaml:"nextScheduleTime"`
}

// Daemonset daemonset
type Daemonset struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Labels      map[string]string `yaml:"labels"`
	Created     int64             `yaml:"created"`
	Status      DaemonsetStatus   `yaml:"status"`
	Metadata    DaemonsetMetadata `yaml:"metadata"`
}

// DaemonsetStatus daemonset status
type DaemonsetStatus struct {
	Desired      int16 `yaml:"desired"`
	Current      int16 `yaml:"current"`
	Available    int16 `yaml:"available"`
	Misscheduled int16 `yaml:"misscheduled"`
	Ready        int16 `yaml:"ready"`
	Unavailable  int16 `yaml:"unavailable"`
	Scheduled    int16 `yaml:"scheduled"`
	Res          int16 `yaml:"reserved"` // byte aligned
}

// DaemonsetMetadata daemonset metadata
type DaemonsetMetadata MetadataGen

// Deployment deployment
type Deployment struct {
	used        int32
	cnt         int32
	UpdateTime  int64              `yaml:"updateTime"`
	RefreshTime int64              `yaml:"refreshTime"`
	CID         string             `yaml:"cID"`
	Name        string             `yaml:"name"`
	Namespace   string             `yaml:"namespace"`
	Labels      map[string]string  `yaml:"labels"`
	Created     int64              `yaml:"created"`
	Status      DeploymentStatus   `yaml:"status"`
	Spec        DeploymentSpec     `yaml:"spec"`
	Metadata    DeploymentMetadata `yaml:"metadata"`
}

// DeploymentStatus deployment status
type DeploymentStatus struct {
	Replicas            int16               `yaml:"replicas"`
	ReplicasAvailable   int16               `yaml:"replicasAvailable"`
	ReplicasUnavailable int16               `yaml:"replicasUnavailable"`
	ReplicasUpdated     int16               `yaml:"replicasUpdated"`
	ObservedGeneration  int16               `yaml:"observedGeneration"`
	reserved            int16               `yaml:"reserved"`
	Condition           DeploymentCondition `yaml:"condition"`
}

// DeploymentCondition deployment condition
type DeploymentCondition struct {
	Available      string `yaml:"Available"`
	Progressing    string `yaml:"Progressing"`
	ReplicaFailure string `yaml:"Progressing"`
}

// DeploymentSpec deployment spec
type DeploymentSpec struct {
	Replicas int16              `yaml:"replicas"`
	Paused   int16              `yaml:"paused"`
	Strategy DeploymentStrategy `yaml:"strategy"`
}

// DeploymentStrategy Deployment strategy
type DeploymentStrategy struct {
	MaxUnavailable int16 `yaml:"maxUnavailable"`
	MaxSurge       int16 `yaml:"maxSurge"`
}

// DeploymentMetadata Deployment metadata
type DeploymentMetadata MetadataGen

// Endpoint : 2021-02-15 1차 모니터링에서는 무시

// Ingress ingress
type Ingress struct {
	used        int32
	cnt         int32
	UpdateTime  int64                  `yaml:"updateTime"`
	RefreshTime int64                  `yaml:"refreshTime"`
	CID         string                 `yaml:"cID"`
	Name        string                 `yaml:"name"`
	Namespace   string                 `yaml:"namespace"`
	Labels      map[string]string      `yaml:"labels"`
	Created     int64                  `yaml:"created"`
	Metadata    IngressMetadata        `yaml:"metadata"`
	Path        map[string]IngressPath `yaml:"path"`
}

// IngressMetadata Ingress metadata
type IngressMetadata MetadataResVersion

// IngressPath ingress host - service
type IngressPath struct {
	Host        string `yaml:"host"`
	Path        string `yaml:"path"`
	ServiceName string `yaml:"serviceName"`
	ServicePort string `yaml:"servicePort"`
	TLSHost     string `yaml:"tlsHost"`
	TLSSecret   string `yaml:"tlsSecret"`
}

// Job job struct
type Job struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Labels      map[string]string `yaml:"labels"`
	Created     int64             `yaml:"created"`
	Spec        JobSpec           `yaml:"spec"`
	Status      JobStatus         `yaml:"status"`
	Complete    int               `yaml:"complete"`
	Failed      int               `yaml:"failed"`
	Owner       map[string]string `yaml:"ownera"`
}

// JobSpec job spec
type JobSpec struct {
	Parallelism           int `yaml:"parallelism"`
	Completions           int `yaml:"completions"`
	ActiveDeadlineSeconds int `yaml:"activeDeadlineSeconds"`
}

// JobStatus job status
type JobStatus struct {
	Succeeded      int   `yaml:"succeeded"`
	Failed         int   `yaml:"failed"`
	Active         int   `yaml:"active"`
	StartTime      int64 `yaml:"startTime"`
	CompletionTime int64 `yaml:"completiontime"`
}

// LimitRange

// PersistentVolumeClaim persistent volume claim
type PersistentVolumeClaim struct {
	used         int32
	cnt          int32
	UpdateTime   int64                       `yaml:"updateTime"`
	RefreshTime  int64                       `yaml:"refreshTime"`
	CID          string                      `yaml:"cID"`
	Name         string                      `yaml:"name"`
	Namespace    string                      `yaml:"namespace"`
	Labels       map[string]string           `yaml:"labels"`
	Info         map[string]string           `yaml:"info"`
	Phase        string                      `yaml:"phase"`
	RequestBytes int64                       `yaml:"requestBytes"`
	AccessModes  []string                    `yaml:"accessModes"`
	Status       PersistentVolumeClaimStatus `yaml:"status"`
}

// PersistentVolumeClaimStatus persistent volume claim status
type PersistentVolumeClaimStatus struct {
	Condition string `yaml:"condition"` // Resizing, FileSystemResizePending
}

// PersistentVolume persistent volume
type PersistentVolume struct {
	used          int32
	cnt           int32
	UpdateTime    int64             `yaml:"updateTime"`
	RefreshTime   int64             `yaml:"refreshTime"`
	CID           string            `yaml:"cID"`
	Name          string            `yaml:"name"`
	Labels        map[string]string `yaml:"labels"`
	Phase         string            `yaml:"phase"`
	CapacityBytes int64             `yaml:"capacityBytes"`
}

// Pod pod
type Pod struct {
	used           int32
	cnt            int32
	UpdateTime     int64                    `yaml:"updateTime"`
	RefreshTime    int64                    `yaml:"refreshTime"`
	CID            string                   `yaml:"cID"`
	Name           string                   `yaml:"name"`
	Namespace      string                   `yaml:"namespace"`
	PodInfo        map[string]string        `yaml:"info"`
	StartTime      int64                    `yaml:"startTime"`
	Owner          map[string]string        `yaml:"owner"`
	Labels         map[string]string        `yaml:"labels"`
	Created        int64                    `yaml:"created"`
	RestartPolicy  string                   `yaml:"restartPolicy"`
	Status         PodStatus                `yaml:"status"`
	InitContainers map[string]PodContainer  `yaml:"initContainers"`
	Containers     map[string]PodContainer  `yaml:"containers"`
	PVC            PodPersistentVolumeClaim `yaml:"PVC"`
}

// PodStatus pod status
type PodStatus struct {
	ScheduledTime int64
	Phase         string
	Ready         string
	Scheduled     string
}

// PodContainer pod container
type PodContainer struct {
	ID              string            `yaml:"id"`
	ContainerInfo   map[string]string `yaml:"info"`
	Status          ContainerStatus   `yaml:"status"`
	RequestResource ContainerResource `yaml:"requestResource"`
	LimitResource   ContainerResource `yaml:"limitResource"`
}

// ContainerStatus container status
type ContainerStatus struct {
	Waiting              int    `yaml:"waiting"`
	WaitingReason        string `yaml:"waitingReason"`
	Running              int    `yaml:"running"`
	Terminated           int    `yaml:"terminated"`
	TerminatedReason     string `yaml:"terminatedReason"`
	LastTerminatedReason string `yaml:"lastTerminatedReason"`
	Ready                int    `yaml:"ready"`
	RestartTotal         int    `yaml:"restartTotal"`
}

// ContainerResource container resource
type ContainerResource struct {
	Node string `yaml:"node"`
	CPU  string `yaml:"cpu"`
	Mem  string `yaml:"mem"`
}

// PodPersistentVolumeClaim pod - persistent volume claim info
type PodPersistentVolumeClaim struct {
	ReadOnly bool   `yaml:"readyOnly"`
	PVCName  string `yaml:"pvcName"`
}

// Replicaset replicaset
type Replicaset struct {
	used        int32
	cnt         int32
	UpdateTime  int64              `yaml:"updateTime"`
	RefreshTime int64              `yaml:"refreshTime"`
	CID         string             `yaml:"cID"`
	Name        string             `yaml:"name"`
	Namespace   string             `yaml:"namespace"`
	Created     int64              `yaml:"created"`
	Status      ReplicasetStatus   `yaml:"status"`
	Spec        ReplicasetSpec     `yaml:"spec"`
	Metadata    ReplicasetMetadata `yaml:"metadata"`
	Owner       map[string]string  `yaml:"owner"`
	Labels      map[string]string  `yaml:"labels"`
}

// ReplicasetStatus replicaset status
type ReplicasetStatus struct {
	Replicas             int `yaml:"replicas"`
	FullyLabeledReplicas int `yaml:"fullyLabeledReplicas"` // 이 리플리카셋으로 실행되는 복제본(pod?)의 수
	ReadyReplicas        int `yaml:"readyReplicas"`
	ObservedGeneration   int `yaml:"observedGeneration"`
}

// ReplicasetSpec replicaset spec
type ReplicasetSpec struct {
	Replicas int `yaml:"replicas"`
}

// ReplicasetMetadata replicaset metadata
type ReplicasetMetadata MetadataGen

// ReplicationController replication controller
type ReplicationController struct {
	used        int32
	cnt         int32
	UpdateTime  int64                         `yaml:"updateTime"`
	RefreshTime int64                         `yaml:"refreshTime"`
	CID         string                        `yaml:"cID"`
	Name        string                        `yaml:"name"`
	Namespace   string                        `yaml:"namespace"`
	Created     int64                         `yaml:"created"`
	Status      ReplicationControllerStatus   `yaml:"status"`
	Spec        ReplicationControllerSpec     `yaml:"spec"`
	Metadata    ReplicationControllerMetadata `yaml:"metadata"`
}

// ReplicationControllerStatus replication controller status
type ReplicationControllerStatus struct {
	Replicas             int `yaml:"replicas"`
	FullyLabeledReplicas int `yaml:"fullyLabeledReplicas"` // 이 리플리카셋으로 실행되는 복제본(pod?)의 수
	ReadyReplicas        int `yaml:"readyReplicas"`
	AvailableReplicas    int `yaml:"availableReplicas"`
	ObservedGeneration   int `yaml:"observedGeneration"`
}

// ReplicationControllerSpec replication controller spec
type ReplicationControllerSpec ReplicasetSpec

// ReplicationControllerMetadata replication controller metadata
type ReplicationControllerMetadata MetadataGen

// Secret secret
type Secret struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Type        string            `yaml:"type"`
	Labels      map[string]string `yaml:"labels"`
	Created     int64             `yaml:"created"`
	Metadata    SecretMetadata    `yaml:"metadata"`
}

// SecretMetadata secret metadata
type SecretMetadata MetadataResVersion

// Service service data struct
type Service struct {
	used        int32
	cnt         int32
	UpdateTime  int64             `yaml:"updateTime"`
	RefreshTime int64             `yaml:"refreshTime"`
	CID         string            `yaml:"cID"`
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Info        map[string]string `yaml:"info"`
	Created     int64             `yaml:"created"`
	Type        string            `yaml:"type"`
	Labels      map[string]string `yaml:"labels"`
	Spec        ServiceSpec       `yaml:"spec"`
}

// ServiceSpec service spec
type ServiceSpec struct {
	ExternalIP []string `yaml:"externalIP"`
}

// Statefulset statefulset
type Statefulset struct {
	used        int32
	cnt         int32
	UpdateTime  int64               `yaml:"updateTime"`
	RefreshTime int64               `yaml:"refreshTime"`
	CID         string              `yaml:"cID"`
	Name        string              `yaml:"name"`
	Namespace   string              `yaml:"namespace"`
	Created     int64               `yaml:"created"`
	Status      StatefulsetStatus   `yaml:"status"`
	Replicas    int                 `yaml:"replicas"`
	Metadata    StatefulsetMetadata `yaml:"metadata"`
	Labels      map[string]string   `yaml:"labels"`
	Revision    StatefulsetRevision `yaml:"revision"`
}

// StatefulsetStatus statefulset status
type StatefulsetStatus struct {
	Replicas           int `yaml:"replicas"`
	ReplicasCurrent    int `yaml:"replicasCurrent"`
	ReplicasReady      int `yaml:"replicasReady"`
	ReplicasUpdated    int `yaml:"replicasUpdated"`
	ObservedGeneration int `yaml:"observedGeneration"`
}

// StatefulsetMetadata statefulset metadata
type StatefulsetMetadata MetadataGen

// StatefulsetRevision statefulset revision
type StatefulsetRevision struct {
	Current string `yaml:"current"`
	Update  string `yaml:"update"`
}

// DoubleKey double key
type DoubleKey struct {
	KeyA string `yaml:"keyA"`
	KeyB string `yaml:"keyB"`
}

// TripleKey tiple key
type TripleKey struct {
	KeyA string `yaml:"keyA"`
	KeyB string `yaml:"keyB"`
	KeyC string `yaml:"keyC"`
}

// KSM kube state metrics info
type KSM struct {
	Clusters               map[string]*Cluster                  `yaml:"clusters"`
	Nodes                  map[DoubleKey]*Node                  `yaml:"nodes"`
	Namespaces             map[DoubleKey]*Namespace             `yaml:"namespaces"`
	Cronjobs               map[TripleKey]*Cronjob               `yaml:"cronjobs"`
	Configmaps             map[TripleKey]*Configmap             `yaml:"configmaps"`
	Daemonsets             map[TripleKey]*Daemonset             `yaml:"daemonsets"`
	Deployments            map[TripleKey]*Deployment            `yaml:"deployments"`
	Ingresses              map[TripleKey]*Ingress               `yaml:"ingresses"`
	Jobs                   map[TripleKey]*Job                   `yaml:"jobs"`
	PersistentVolumeClaims map[TripleKey]*PersistentVolumeClaim `yaml:"persistentVolumeClaims"`
	PersistentVolumes      map[DoubleKey]*PersistentVolume      `yaml:"persistentVolumes"`
	Pods                   map[TripleKey]*Pod                   `yaml:"pods"`
	Replicasets            map[TripleKey]*Replicaset            `yaml:"replicasets"`
	ReplicationControllers map[TripleKey]*ReplicationController `yaml:"replicationControllers"`
	Secrets                map[TripleKey]*Secret                `yaml:"secrets"`
	Services               map[TripleKey]*Service               `yaml:"services"`
	Statefulsets           map[TripleKey]*Statefulset           `yaml:"statefulsets"`
}
