package ksm

import (
	"fmt"
	"sync/atomic"

	logrus "github.com/sirupsen/logrus"
	"mobigen.com/iris-cloud-monitoring/model"
)

// Manager ksm manager
type Manager struct {
	Data *KSM
	Log  logrus.FieldLogger
}

var ksmMN *Manager

// Init init ksm data struct
func Init(_log logrus.FieldLogger) error {
	ksmMN = new(Manager)
	ksmMN.Data = &KSM{
		Clusters:               make(map[string]*Cluster),
		Nodes:                  make(map[DoubleKey]*Node),
		Namespaces:             make(map[DoubleKey]*Namespace),
		Configmaps:             make(map[TripleKey]*Configmap),
		Cronjobs:               make(map[TripleKey]*Cronjob),
		Daemonsets:             make(map[TripleKey]*Daemonset),
		Deployments:            make(map[TripleKey]*Deployment),
		Ingresses:              make(map[TripleKey]*Ingress),
		Jobs:                   make(map[TripleKey]*Job),
		PersistentVolumeClaims: make(map[TripleKey]*PersistentVolumeClaim),
		PersistentVolumes:      make(map[DoubleKey]*PersistentVolume),
		Pods:                   make(map[TripleKey]*Pod),
		Replicasets:            make(map[TripleKey]*Replicaset),
		ReplicationControllers: make(map[TripleKey]*ReplicationController),
		Secrets:                make(map[TripleKey]*Secret),
		Services:               make(map[TripleKey]*Service),
		Statefulsets:           make(map[TripleKey]*Statefulset),
	}
	ksmMN.Log = _log.WithFields(map[string]interface{}{
		"data": "ksm",
	})
	ksmMN.Log.Info("Kube State Metrics HashMap Init Success")
	return nil
}

// GetInstance singleton kube state metrics
func GetInstance() *Manager {
	return ksmMN
}

// AddCluster add new cluster
func (manager *Manager) AddCluster(cID, cName, ksmURL string) error {
	manager.Log.Infof("Add Cluster To Kube_State_Metrics. CID[ %s ] CName[ %s ] RequestURL[ %s ]",
		cID, cName, ksmURL)
	cluster := new(Cluster)
	cluster.ID = cID
	cluster.Name = cName
	cluster.KsmURL = ksmURL
	atomic.StoreInt32(&cluster.used, model.USED)
	atomic.StoreInt32(&cluster.cnt, 1)
	manager.Data.Clusters[cID] = cluster
	manager.Log.Infof("Success Add Cluster To Kube_State_Metrics. CID[ %s ] CName[ %s ] RequestURL[ %s ]",
		manager.Data.Clusters[cID].ID, manager.Data.Clusters[cID].Name, manager.Data.Clusters[cID].KsmURL)
	return nil
}

// GetCluster return cluster
func (manager *Manager) GetCluster(cID string) (*Cluster, error) {
	if cluster, ok := manager.Data.Clusters[cID]; ok {
		if atomic.LoadInt32(&cluster.used) != model.USED {
			return nil, fmt.Errorf("error. cluster not used.. id[ %s ]", cID)
		}
		atomic.AddInt32(&cluster.cnt, 1)
		return cluster, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Cluster ID[ %s ]", cID)
}

// GetClusters return cluster all
func (manager *Manager) GetClusters() []Cluster {
	ret := []Cluster{}
	for _, v := range manager.Data.Clusters {
		if atomic.LoadInt32(&v.used) != model.USED {
			continue
		}
		ret = append(ret, v.Clone())
	}
	return ret
}

// Del return delete cluster
func (cluster *Cluster) Del() error {
	if atomic.CompareAndSwapInt32(&cluster.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Cluster CID[ %s ] Name[ %s ]", cluster.ID, cluster.Name)
		cluster.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Cluster CID[ %s ] Name[ %s ]", cluster.ID, cluster.Name)
}

// Put cluster put
func (cluster *Cluster) Put() {
	if atomic.LoadInt32(&cluster.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&cluster.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Cluster CID[ %s ] Name[ %s ]", cluster.ID, cluster.Name)
		delete(manager.Data.Clusters, cluster.ID)
	}
}

// AllocNode alloc node
func (manager *Manager) AllocNode() (*Node, error) {
	node := new(Node)
	if node != nil {
		atomic.StoreInt32(&node.used, model.CREATING)
		atomic.StoreInt32(&node.cnt, 2)
		return node, nil
	}
	return nil, fmt.Errorf("error. failed to alloc node")
}

// Add Add Node
func (node *Node) Add() error {
	if len(node.CID) <= 0 || len(node.Name) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Node[ %s ]", node.CID, node.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Node To Kube_State_Metrics. CID[ %s ] Node[ %s ]", node.CID, node.Name)
	key := DoubleKey{
		KeyA: node.CID,
		KeyB: node.Name,
	}
	manager.Data.Nodes[key] = node
	atomic.StoreInt32(&node.used, model.USED)
	return nil
}

// GetNode get node by cid(clusterid) and node name
func (manager *Manager) GetNode(cid, name string) (*Node, error) {
	key := DoubleKey{
		KeyA: cid,
		KeyB: name,
	}
	if node, ok := manager.Data.Nodes[key]; ok {
		if atomic.LoadInt32(&node.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Node Not Used.. CID[ %s ] Name[ %s ]", cid, name)
		}
		atomic.AddInt32(&node.cnt, 1)
		return node, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Node CID[ %s ] Name[ %s ]", cid, name)
}

// Put put node
func (node *Node) Put() {
	if atomic.LoadInt32(&node.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&node.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Node CID[ %s ] Name[ %s ]", node.CID, node.Name)
		key := DoubleKey{
			KeyA: node.CID,
			KeyB: node.Name,
		}
		delete(manager.Data.Nodes, key)
	}
}

// Del del node
func (node *Node) Del() error {
	if atomic.CompareAndSwapInt32(&node.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Node CID[ %s ] Name[ %s ]", node.CID, node.Name)
		node.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Node. CID[ %s ] Name[ %s ]", node.CID, node.Name)
}

// AllocNamespace alloc Namespace
func (manager *Manager) AllocNamespace() (*Namespace, error) {
	ret := new(Namespace)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Namespace")
}

// Add Add Namespace
func (namespace *Namespace) Add() error {
	if len(namespace.CID) <= 0 || len(namespace.Name) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Name[ %s ]", namespace.CID, namespace.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Namespace To Kube_State_Metrics. CID[ %s ] Name[ %s ]", namespace.CID, namespace.Name)
	key := DoubleKey{
		KeyA: namespace.CID,
		KeyB: namespace.Name,
	}
	manager.Data.Namespaces[key] = namespace
	atomic.StoreInt32(&namespace.used, model.USED)
	return nil
}

// GetNamespace get Namespace by cid(clusterid) and name
func (manager *Manager) GetNamespace(cid, name string) (*Namespace, error) {
	key := DoubleKey{
		KeyA: cid,
		KeyB: name,
	}
	if namespace, ok := manager.Data.Namespaces[key]; ok {
		if atomic.LoadInt32(&namespace.used) != model.USED {
			return nil, fmt.Errorf("error. namespace not used.. CID[ %s ] Name[ %s ]", cid, name)
		}
		atomic.AddInt32(&namespace.cnt, 1)
		return namespace, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Namespace CID[ %s ] Name[ %s ]", cid, name)
}

// Put put Namespace
func (namespace *Namespace) Put() {
	if atomic.LoadInt32(&namespace.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&namespace.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Namespace CID[ %s ] Name[ %s ]", namespace.CID, namespace.Name)
		key := DoubleKey{
			KeyA: namespace.CID,
			KeyB: namespace.Name,
		}
		delete(manager.Data.Namespaces, key)
	}
}

// Del del Namespace
func (namespace *Namespace) Del() error {
	if atomic.CompareAndSwapInt32(&namespace.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Namespace CID[ %s ] Name[ %s ]", namespace.CID, namespace.Name)
		namespace.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Namespace. CID[ %s ] Name[ %s ]", namespace.CID, namespace.Name)
}

// AllocConfigmap alloc Configmap
func (manager *Manager) AllocConfigmap() (*Configmap, error) {
	ret := new(Configmap)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Configmap")
}

// Add Add Configmap
func (configmap *Configmap) Add() error {
	if len(configmap.CID) <= 0 || len(configmap.Name) <= 0 || len(configmap.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			configmap.CID, configmap.Namespace, configmap.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Configmap To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		configmap.CID, configmap.Namespace, configmap.Name)
	key := TripleKey{
		KeyA: configmap.CID,
		KeyB: configmap.Namespace,
		KeyC: configmap.Name,
	}
	manager.Data.Configmaps[key] = configmap
	atomic.StoreInt32(&configmap.used, model.USED)
	return nil
}

// GetConfigmap get Configmap by cid(clusterid), namespace, name
func (manager *Manager) GetConfigmap(cid, namespace, name string) (*Configmap, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if configmap, ok := manager.Data.Configmaps[key]; ok {
		if atomic.LoadInt32(&configmap.used) != model.USED {
			return nil, fmt.Errorf("error. configmap not used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&configmap.cnt, 1)
		return configmap, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Configmap CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Configmap
func (configmap *Configmap) Put() {
	if atomic.LoadInt32(&configmap.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&configmap.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Configmap CID[ %s ] Namespace[ %s ] Name[ %s ]", configmap.CID, configmap.Namespace, configmap.Name)
		key := TripleKey{
			KeyA: configmap.CID,
			KeyB: configmap.Namespace,
			KeyC: configmap.Name,
		}
		delete(manager.Data.Configmaps, key)
	}
}

// Del del Configmap
func (configmap *Configmap) Del() error {
	if atomic.CompareAndSwapInt32(&configmap.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Configmap CID[ %s ] Namespace[ %s ] Name[ %s ]", configmap.CID, configmap.Namespace, configmap.Name)
		configmap.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Configmap. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		configmap.CID, configmap.Namespace, configmap.Name)
}

// AllocDaemonset alloc Daemonset
func (manager *Manager) AllocDaemonset() (*Daemonset, error) {
	ret := new(Daemonset)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Daemonset")
}

// Add Add Daemonset
func (daemonset *Daemonset) Add() error {
	if len(daemonset.CID) <= 0 || len(daemonset.Name) <= 0 || len(daemonset.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			daemonset.CID, daemonset.Namespace, daemonset.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Daemonset To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		daemonset.CID, daemonset.Namespace, daemonset.Name)
	key := TripleKey{
		KeyA: daemonset.CID,
		KeyB: daemonset.Namespace,
		KeyC: daemonset.Name,
	}
	manager.Data.Daemonsets[key] = daemonset
	atomic.StoreInt32(&daemonset.used, model.USED)
	return nil
}

// GetDaemonset get Daemonset by cid(clusterid), namespace, name
func (manager *Manager) GetDaemonset(cid, namespace, name string) (*Daemonset, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if daemonset, ok := manager.Data.Daemonsets[key]; ok {
		if atomic.LoadInt32(&daemonset.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Daemonset Not Used. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&daemonset.cnt, 1)
		return daemonset, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Daemonset CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Daemonset
func (daemonset *Daemonset) Put() {
	if atomic.LoadInt32(&daemonset.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&daemonset.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Daemonset CID[ %s ] Namespace[ %s ] Name[ %s ]", daemonset.CID, daemonset.Namespace, daemonset.Name)
		key := TripleKey{
			KeyA: daemonset.CID,
			KeyB: daemonset.Namespace,
			KeyC: daemonset.Name,
		}
		delete(manager.Data.Daemonsets, key)
	}
}

// Del del Daemonset
func (daemonset *Daemonset) Del() error {
	if atomic.CompareAndSwapInt32(&daemonset.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Daemonset CID[ %s ] Namespace[ %s ] Name[ %s ]", daemonset.CID, daemonset.Namespace, daemonset.Name)
		daemonset.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Daemonset. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		daemonset.CID, daemonset.Namespace, daemonset.Name)
}

// AllocDeployment alloc Deployment
func (manager *Manager) AllocDeployment() (*Deployment, error) {
	ret := new(Deployment)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Deployment")
}

// Add Add Deployment
func (deployment *Deployment) Add() error {
	if len(deployment.CID) <= 0 || len(deployment.Name) <= 0 || len(deployment.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			deployment.CID, deployment.Namespace, deployment.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Deployment To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		deployment.CID, deployment.Namespace, deployment.Name)
	key := TripleKey{
		KeyA: deployment.CID,
		KeyB: deployment.Namespace,
		KeyC: deployment.Name,
	}
	manager.Data.Deployments[key] = deployment
	atomic.StoreInt32(&deployment.used, model.USED)
	return nil
}

// GetDeployment get Deployment by cid(clusterid), namespace, name
func (manager *Manager) GetDeployment(cid, namespace, name string) (*Deployment, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if deployment, ok := manager.Data.Deployments[key]; ok {
		if atomic.LoadInt32(&deployment.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Deployment Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&deployment.cnt, 1)
		return deployment, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Deployment CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Deployment
func (deployment *Deployment) Put() {
	if atomic.LoadInt32(&deployment.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&deployment.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Deployment CID[ %s ] Namespace[ %s ] Name[ %s ]", deployment.CID, deployment.Namespace, deployment.Name)
		key := TripleKey{
			KeyA: deployment.CID,
			KeyB: deployment.Namespace,
			KeyC: deployment.Name,
		}
		delete(manager.Data.Deployments, key)
	}
}

// Del del Deployment
func (deployment *Deployment) Del() error {
	if atomic.CompareAndSwapInt32(&deployment.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Deployment CID[ %s ] Namespace[ %s ] Name[ %s ]", deployment.CID, deployment.Namespace, deployment.Name)
		deployment.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Deployment. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		deployment.CID, deployment.Namespace, deployment.Name)
}

// AllocIngress alloc Ingress
func (manager *Manager) AllocIngress() (*Ingress, error) {
	ret := new(Ingress)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Ingress")
}

// Add Add Ingress
func (ingress *Ingress) Add() error {
	if len(ingress.CID) <= 0 || len(ingress.Name) <= 0 || len(ingress.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			ingress.CID, ingress.Namespace, ingress.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Ingress To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		ingress.CID, ingress.Namespace, ingress.Name)
	key := TripleKey{
		KeyA: ingress.CID,
		KeyB: ingress.Namespace,
		KeyC: ingress.Name,
	}
	manager.Data.Ingresses[key] = ingress
	atomic.StoreInt32(&ingress.used, model.USED)
	return nil
}

// GetIngress get Ingress by cid(clusterid), namespace, name
func (manager *Manager) GetIngress(cid, namespace, name string) (*Ingress, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if ingress, ok := manager.Data.Ingresses[key]; ok {
		if atomic.LoadInt32(&ingress.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Ingress Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&ingress.cnt, 1)
		return ingress, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Ingress CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Ingress
func (ingress *Ingress) Put() {
	if atomic.LoadInt32(&ingress.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&ingress.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Ingress CID[ %s ] Namespace[ %s ] Name[ %s ]", ingress.CID, ingress.Namespace, ingress.Name)
		key := TripleKey{
			KeyA: ingress.CID,
			KeyB: ingress.Namespace,
			KeyC: ingress.Name,
		}
		delete(manager.Data.Ingresses, key)
	}
}

// Del del Ingress
func (ingress *Ingress) Del() error {
	if atomic.CompareAndSwapInt32(&ingress.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Ingress CID[ %s ] Namespace[ %s ] Name[ %s ]", ingress.CID, ingress.Namespace, ingress.Name)
		ingress.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Ingress. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		ingress.CID, ingress.Namespace, ingress.Name)
}

// AllocCronjob alloc Cronjob
func (manager *Manager) AllocCronjob() (*Cronjob, error) {
	ret := new(Cronjob)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Cronjob")
}

// Add Add Cronjob
func (cronjob *Cronjob) Add() error {
	if len(cronjob.CID) <= 0 || len(cronjob.Name) <= 0 || len(cronjob.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			cronjob.CID, cronjob.Namespace, cronjob.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Cronjob To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cronjob.CID, cronjob.Namespace, cronjob.Name)
	key := TripleKey{
		KeyA: cronjob.CID,
		KeyB: cronjob.Namespace,
		KeyC: cronjob.Name,
	}
	manager.Data.Cronjobs[key] = cronjob
	atomic.StoreInt32(&cronjob.used, model.USED)
	return nil
}

// GetCronjob get Cronjob by cid(clusterid), namespace, name
func (manager *Manager) GetCronjob(cid, namespace, name string) (*Cronjob, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if cronjob, ok := manager.Data.Cronjobs[key]; ok {
		if atomic.LoadInt32(&cronjob.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Cronjob Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&cronjob.cnt, 1)
		return cronjob, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Cronjob CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Cronjob
func (cronjob *Cronjob) Put() {
	if atomic.LoadInt32(&cronjob.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&cronjob.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Cronjob CID[ %s ] Namespace[ %s ] Name[ %s ]", cronjob.CID, cronjob.Namespace, cronjob.Name)
		key := TripleKey{
			KeyA: cronjob.CID,
			KeyB: cronjob.Namespace,
			KeyC: cronjob.Name,
		}
		delete(manager.Data.Cronjobs, key)
	}
}

// Del del Cronjob
func (cronjob *Cronjob) Del() error {
	if atomic.CompareAndSwapInt32(&cronjob.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Cronjob CID[ %s ] Namespace[ %s ] Name[ %s ]", cronjob.CID, cronjob.Namespace, cronjob.Name)
		cronjob.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Cronjob. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cronjob.CID, cronjob.Namespace, cronjob.Name)
}

// AllocJob alloc Job
func (manager *Manager) AllocJob() (*Job, error) {
	ret := new(Job)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Job")
}

// Add Add Job
func (job *Job) Add() error {
	if len(job.CID) <= 0 || len(job.Name) <= 0 || len(job.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			job.CID, job.Namespace, job.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Job To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		job.CID, job.Namespace, job.Name)
	key := TripleKey{
		KeyA: job.CID,
		KeyB: job.Namespace,
		KeyC: job.Name,
	}
	manager.Data.Jobs[key] = job
	atomic.StoreInt32(&job.used, model.USED)
	return nil
}

// GetJob get Job by cid(clusterid), namespace, name
func (manager *Manager) GetJob(cid, namespace, name string) (*Job, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if job, ok := manager.Data.Jobs[key]; ok {
		if atomic.LoadInt32(&job.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Job Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&job.cnt, 1)
		return job, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Job CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Job
func (job *Job) Put() {
	if atomic.LoadInt32(&job.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&job.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Job CID[ %s ] Namespace[ %s ] Name[ %s ]", job.CID, job.Namespace, job.Name)
		key := TripleKey{
			KeyA: job.CID,
			KeyB: job.Namespace,
			KeyC: job.Name,
		}
		delete(manager.Data.Jobs, key)
	}
}

// Del del Job
func (job *Job) Del() error {
	if atomic.CompareAndSwapInt32(&job.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Job CID[ %s ] Namespace[ %s ] Name[ %s ]", job.CID, job.Namespace, job.Name)
		job.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Job. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		job.CID, job.Namespace, job.Name)
}

// AllocPersistentVolumeClaim alloc PersistentVolumeClaim
func (manager *Manager) AllocPersistentVolumeClaim() (*PersistentVolumeClaim, error) {
	ret := new(PersistentVolumeClaim)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc PersistentVolumeClaim")
}

// Add Add PersistentVolumeClaim
func (pvc *PersistentVolumeClaim) Add() error {
	if len(pvc.CID) <= 0 || len(pvc.Name) <= 0 || len(pvc.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			pvc.CID, pvc.Namespace, pvc.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add PersistentVolumeClaim To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		pvc.CID, pvc.Namespace, pvc.Name)
	key := TripleKey{
		KeyA: pvc.CID,
		KeyB: pvc.Namespace,
		KeyC: pvc.Name,
	}
	manager.Data.PersistentVolumeClaims[key] = pvc
	atomic.StoreInt32(&pvc.used, model.USED)
	return nil
}

// GetPersistentVolumeClaim get PersistentVolumeClaim by cid(clusterid), namespace, name
func (manager *Manager) GetPersistentVolumeClaim(cid, namespace, name string) (*PersistentVolumeClaim, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if pvc, ok := manager.Data.PersistentVolumeClaims[key]; ok {
		if atomic.LoadInt32(&pvc.used) != model.USED {
			return nil, fmt.Errorf("ERROR. PersistentVolumeClaim Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&pvc.cnt, 1)
		return pvc, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found PersistentVolumeClaim CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put PersistentVolumeClaim
func (pvc *PersistentVolumeClaim) Put() {
	if atomic.LoadInt32(&pvc.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&pvc.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. PersistentVolumeClaim CID[ %s ] Namespace[ %s ] Name[ %s ]", pvc.CID, pvc.Namespace, pvc.Name)
		key := TripleKey{
			KeyA: pvc.CID,
			KeyB: pvc.Namespace,
			KeyC: pvc.Name,
		}
		delete(manager.Data.PersistentVolumeClaims, key)
	}
}

// Del del PersistentVolumeClaim
func (pvc *PersistentVolumeClaim) Del() error {
	if atomic.CompareAndSwapInt32(&pvc.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. PersistentVolumeClaim CID[ %s ] Namespace[ %s ] Name[ %s ]", pvc.CID, pvc.Namespace, pvc.Name)
		pvc.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete PersistentVolumeClaim. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		pvc.CID, pvc.Namespace, pvc.Name)
}

// AllocPersistentVolume alloc PersistentVolume
func (manager *Manager) AllocPersistentVolume() (*PersistentVolume, error) {
	ret := new(PersistentVolume)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc PersistentVolume")
}

// Add Add PersistentVolume
func (pv *PersistentVolume) Add() error {
	if len(pv.CID) <= 0 || len(pv.Name) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Name[ %s ]",
			pv.CID, pv.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add PersistentVolume To Kube_State_Metrics. CID[ %s ] Name[ %s ]",
		pv.CID, pv.Name)
	key := DoubleKey{
		KeyA: pv.CID,
		KeyB: pv.Name,
	}
	manager.Data.PersistentVolumes[key] = pv
	atomic.StoreInt32(&pv.used, model.USED)
	return nil
}

// GetPersistentVolume get PersistentVolume by cid(clusterid), namespace, name
func (manager *Manager) GetPersistentVolume(cid, name string) (*PersistentVolume, error) {
	key := DoubleKey{
		KeyA: cid,
		KeyB: name,
	}
	if pv, ok := manager.Data.PersistentVolumes[key]; ok {
		if atomic.LoadInt32(&pv.used) != model.USED {
			return nil, fmt.Errorf("ERROR. PersistentVolume Not Used.. CID[ %s ] Name[ %s ]",
				cid, name)
		}
		atomic.AddInt32(&pv.cnt, 1)
		return pv, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found PersistentVolume CID[ %s ] Name[ %s ]",
		cid, name)
}

// Put put PersistentVolume
func (pv *PersistentVolume) Put() {
	if atomic.LoadInt32(&pv.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&pv.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. PersistentVolume CID[ %s ] Name[ %s ]", pv.CID, pv.Name)
		key := DoubleKey{
			KeyA: pv.CID,
			KeyB: pv.Name,
		}
		delete(manager.Data.PersistentVolumes, key)
	}
}

// Del del PersistentVolume
func (pv *PersistentVolume) Del() error {
	if atomic.CompareAndSwapInt32(&pv.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. PersistentVolume CID[ %s ] Name[ %s ]", pv.CID, pv.Name)
		pv.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete PersistentVolume. CID[ %s ] Name[ %s ]",
		pv.CID, pv.Name)
}

// AllocPod alloc Pod
func (manager *Manager) AllocPod() (*Pod, error) {
	ret := new(Pod)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Pod")
}

// Add Add Pod
func (pod *Pod) Add() error {
	if len(pod.CID) <= 0 || len(pod.Name) <= 0 || len(pod.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			pod.CID, pod.Namespace, pod.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Pod To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		pod.CID, pod.Namespace, pod.Name)
	key := TripleKey{
		KeyA: pod.CID,
		KeyB: pod.Namespace,
		KeyC: pod.Name,
	}
	manager.Data.Pods[key] = pod
	atomic.StoreInt32(&pod.used, model.USED)
	return nil
}

// GetPod get Pod by cid(clusterid), namespace, name
func (manager *Manager) GetPod(cid, namespace, name string) (*Pod, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if pod, ok := manager.Data.Pods[key]; ok {
		if atomic.LoadInt32(&pod.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Pod Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&pod.cnt, 1)
		return pod, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Pod CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Pod
func (pod *Pod) Put() {
	if atomic.LoadInt32(&pod.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&pod.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Pod CID[ %s ] Namespace[ %s ] Name[ %s ]", pod.CID, pod.Namespace, pod.Name)
		key := TripleKey{
			KeyA: pod.CID,
			KeyB: pod.Namespace,
			KeyC: pod.Name,
		}
		delete(manager.Data.Pods, key)
	}
}

// Del del Pod
func (pod *Pod) Del() error {
	if atomic.CompareAndSwapInt32(&pod.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Pod CID[ %s ] Namespace[ %s ] Name[ %s ]", pod.CID, pod.Namespace, pod.Name)
		pod.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Pod. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		pod.CID, pod.Namespace, pod.Name)
}

// AllocReplicaset alloc Replicaset
func (manager *Manager) AllocReplicaset() (*Replicaset, error) {
	ret := new(Replicaset)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Replicaset")
}

// Add Add Replicaset
func (replicaset *Replicaset) Add() error {
	if len(replicaset.CID) <= 0 || len(replicaset.Name) <= 0 || len(replicaset.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			replicaset.CID, replicaset.Namespace, replicaset.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Replicaset To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		replicaset.CID, replicaset.Namespace, replicaset.Name)
	key := TripleKey{
		KeyA: replicaset.CID,
		KeyB: replicaset.Namespace,
		KeyC: replicaset.Name,
	}
	manager.Data.Replicasets[key] = replicaset
	atomic.StoreInt32(&replicaset.used, model.USED)
	return nil
}

// GetReplicaset get Replicaset by cid(clusterid), namespace, name
func (manager *Manager) GetReplicaset(cid, namespace, name string) (*Replicaset, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if replicaset, ok := manager.Data.Replicasets[key]; ok {
		if atomic.LoadInt32(&replicaset.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Replicaset Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&replicaset.cnt, 1)
		return replicaset, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Replicaset CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Replicaset
func (replicaset *Replicaset) Put() {
	if atomic.LoadInt32(&replicaset.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&replicaset.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Replicaset CID[ %s ] Namespace[ %s ] Name[ %s ]", replicaset.CID, replicaset.Namespace, replicaset.Name)
		key := TripleKey{
			KeyA: replicaset.CID,
			KeyB: replicaset.Namespace,
			KeyC: replicaset.Name,
		}
		delete(manager.Data.Replicasets, key)
	}
}

// Del del Replicaset
func (replicaset *Replicaset) Del() error {
	if atomic.CompareAndSwapInt32(&replicaset.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Replicaset CID[ %s ] Namespace[ %s ] Name[ %s ]", replicaset.CID, replicaset.Namespace, replicaset.Name)
		replicaset.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Replicaset. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		replicaset.CID, replicaset.Namespace, replicaset.Name)
}

// AllocReplicationController alloc ReplicationController
func (manager *Manager) AllocReplicationController() (*ReplicationController, error) {
	ret := new(ReplicationController)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc ReplicationController")
}

// Add Add ReplicationController
func (rc *ReplicationController) Add() error {
	if len(rc.CID) <= 0 || len(rc.Name) <= 0 || len(rc.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			rc.CID, rc.Namespace, rc.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add ReplicationController To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		rc.CID, rc.Namespace, rc.Name)
	key := TripleKey{
		KeyA: rc.CID,
		KeyB: rc.Namespace,
		KeyC: rc.Name,
	}
	manager.Data.ReplicationControllers[key] = rc
	atomic.StoreInt32(&rc.used, model.USED)
	return nil
}

// GetReplicationController get ReplicationController by cid(clusterid), namespace, name
func (manager *Manager) GetReplicationController(cid, namespace, name string) (*ReplicationController, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if rc, ok := manager.Data.ReplicationControllers[key]; ok {
		if atomic.LoadInt32(&rc.used) != model.USED {
			return nil, fmt.Errorf("ERROR. ReplicationController Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&rc.cnt, 1)
		return rc, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found ReplicationController CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put ReplicationController
func (rc *ReplicationController) Put() {
	if atomic.LoadInt32(&rc.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&rc.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. ReplicationController CID[ %s ] Namespace[ %s ] Name[ %s ]", rc.CID, rc.Namespace, rc.Name)
		key := TripleKey{
			KeyA: rc.CID,
			KeyB: rc.Namespace,
			KeyC: rc.Name,
		}
		delete(manager.Data.ReplicationControllers, key)
	}
}

// Del del ReplicationController
func (rc *ReplicationController) Del() error {
	if atomic.CompareAndSwapInt32(&rc.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. ReplicationController CID[ %s ] Namespace[ %s ] Name[ %s ]", rc.CID, rc.Namespace, rc.Name)
		rc.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete ReplicationController. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		rc.CID, rc.Namespace, rc.Name)
}

// AllocSecret alloc Secret
func (manager *Manager) AllocSecret() (*Secret, error) {
	ret := new(Secret)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Secret")
}

// Add Add Secret
func (secret *Secret) Add() error {
	if len(secret.CID) <= 0 || len(secret.Name) <= 0 || len(secret.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			secret.CID, secret.Namespace, secret.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Secret To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		secret.CID, secret.Namespace, secret.Name)
	key := TripleKey{
		KeyA: secret.CID,
		KeyB: secret.Namespace,
		KeyC: secret.Name,
	}
	manager.Data.Secrets[key] = secret
	atomic.StoreInt32(&secret.used, model.USED)
	return nil
}

// GetSecret get Secret by cid(clusterid), namespace, name
func (manager *Manager) GetSecret(cid, namespace, name string) (*Secret, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if secret, ok := manager.Data.Secrets[key]; ok {
		if atomic.LoadInt32(&secret.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Secret Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&secret.cnt, 1)
		return secret, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Secret CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Secret
func (secret *Secret) Put() {
	if atomic.LoadInt32(&secret.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&secret.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Secret CID[ %s ] Namespace[ %s ] Name[ %s ]", secret.CID, secret.Namespace, secret.Name)
		key := TripleKey{
			KeyA: secret.CID,
			KeyB: secret.Namespace,
			KeyC: secret.Name,
		}
		delete(manager.Data.Secrets, key)
	}
}

// Del del Secret
func (secret *Secret) Del() error {
	if atomic.CompareAndSwapInt32(&secret.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Secret CID[ %s ] Namespace[ %s ] Name[ %s ]", secret.CID, secret.Namespace, secret.Name)
		secret.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Secret. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		secret.CID, secret.Namespace, secret.Name)
}

// AllocService alloc Service
func (manager *Manager) AllocService() (*Service, error) {
	ret := new(Service)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Service")
}

// Add Add Service
func (service *Service) Add() error {
	if len(service.CID) <= 0 || len(service.Name) <= 0 || len(service.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			service.CID, service.Namespace, service.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Service To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		service.CID, service.Namespace, service.Name)
	key := TripleKey{
		KeyA: service.CID,
		KeyB: service.Namespace,
		KeyC: service.Name,
	}
	manager.Data.Services[key] = service
	atomic.StoreInt32(&service.used, model.USED)
	return nil
}

// GetService get Service by cid(clusterid), namespace, name
func (manager *Manager) GetService(cid, namespace, name string) (*Service, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if service, ok := manager.Data.Services[key]; ok {
		if atomic.LoadInt32(&service.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Service Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&service.cnt, 1)
		return service, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Service CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Service
func (service *Service) Put() {
	if atomic.LoadInt32(&service.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&service.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Service CID[ %s ] Namespace[ %s ] Name[ %s ]", service.CID, service.Namespace, service.Name)
		key := TripleKey{
			KeyA: service.CID,
			KeyB: service.Namespace,
			KeyC: service.Name,
		}
		delete(manager.Data.Services, key)
	}
}

// Del del Service
func (service *Service) Del() error {
	if atomic.CompareAndSwapInt32(&service.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Service CID[ %s ] Namespace[ %s ] Name[ %s ]", service.CID, service.Namespace, service.Name)
		service.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Service. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		service.CID, service.Namespace, service.Name)
}

// AllocStatefulset alloc Statefulset
func (manager *Manager) AllocStatefulset() (*Statefulset, error) {
	ret := new(Statefulset)
	if ret != nil {
		atomic.StoreInt32(&ret.used, model.CREATING)
		atomic.StoreInt32(&ret.cnt, 2)
		return ret, nil
	}
	return nil, fmt.Errorf("ERROR. Failed To Alloc Statefulset")
}

// Add Add Statefulset
func (sts *Statefulset) Add() error {
	if len(sts.CID) <= 0 || len(sts.Name) <= 0 || len(sts.Namespace) <= 0 {
		return fmt.Errorf("ERROR. Not Found Key Value CID[ %s ], Namespace[ %s ], Name[ %s ]",
			sts.CID, sts.Namespace, sts.Name)
	}
	manager := GetInstance()
	manager.Log.Infof("Add Statefulset To Kube_State_Metrics. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		sts.CID, sts.Namespace, sts.Name)
	key := TripleKey{
		KeyA: sts.CID,
		KeyB: sts.Namespace,
		KeyC: sts.Name,
	}
	manager.Data.Statefulsets[key] = sts
	atomic.StoreInt32(&sts.used, model.USED)
	return nil
}

// GetStatefulset get Statefulset by cid(clusterid), namespace, name
func (manager *Manager) GetStatefulset(cid, namespace, name string) (*Statefulset, error) {
	key := TripleKey{
		KeyA: cid,
		KeyB: namespace,
		KeyC: name,
	}
	if sts, ok := manager.Data.Statefulsets[key]; ok {
		if atomic.LoadInt32(&sts.used) != model.USED {
			return nil, fmt.Errorf("ERROR. Statefulset Not Used.. CID[ %s ] Namespace[ %s ] Name[ %s ]",
				cid, namespace, name)
		}
		atomic.AddInt32(&sts.cnt, 1)
		return sts, nil
	}
	return nil, fmt.Errorf("ERROR. Not Found Statefulset CID[ %s ] Namespace[ %s ] Name[ %s ]",
		cid, namespace, name)
}

// Put put Statefulset
func (sts *Statefulset) Put() {
	if atomic.LoadInt32(&sts.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&sts.cnt, -1) == 0 {
		manager := GetInstance()
		manager.Log.Infof("Put. Statefulset CID[ %s ] Namespace[ %s ] Name[ %s ]", sts.CID, sts.Namespace, sts.Name)
		key := TripleKey{
			KeyA: sts.CID,
			KeyB: sts.Namespace,
			KeyC: sts.Name,
		}
		delete(manager.Data.Statefulsets, key)
	}
}

// Del del Statefulset
func (sts *Statefulset) Del() error {
	if atomic.CompareAndSwapInt32(&sts.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Statefulset CID[ %s ] Namespace[ %s ] Name[ %s ]", sts.CID, sts.Namespace, sts.Name)
		sts.Put()
		return nil
	}
	return fmt.Errorf("Del. Already Delete Statefulset. CID[ %s ] Namespace[ %s ] Name[ %s ]",
		sts.CID, sts.Namespace, sts.Name)
}

// DebugForeach debug
func (manager *Manager) DebugForeach() {
	manager.Log.Debugf("Kube State Metrics Datas")
	manager.Log.Debugf("Clusters")
	for k, v := range manager.Data.Clusters {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Nodes")
	for k, v := range manager.Data.Nodes {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Namespace")
	for k, v := range manager.Data.Namespaces {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("PV")
	for k, v := range manager.Data.PersistentVolumes {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Configmaps")
	for k, v := range manager.Data.Configmaps {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Cronjob")
	for k, v := range manager.Data.Cronjobs {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Daemonsets")
	for k, v := range manager.Data.Daemonsets {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Deployments")
	for k, v := range manager.Data.Deployments {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Ingress")
	for k, v := range manager.Data.Ingresses {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Jobs")
	for k, v := range manager.Data.Jobs {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("PVCs")
	for k, v := range manager.Data.PersistentVolumeClaims {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Pods")
	for k, v := range manager.Data.Pods {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Replicaset")
	for k, v := range manager.Data.Replicasets {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("ReplicationControllers")
	for k, v := range manager.Data.ReplicationControllers {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Secret")
	for k, v := range manager.Data.Secrets {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Services")
	for k, v := range manager.Data.Services {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
	manager.Log.Debugf("Statefulsets")
	for k, v := range manager.Data.Statefulsets {
		manager.Log.Debugf("Key[ %+v ] Value[ %+v ]", k, *v)
	}
}
