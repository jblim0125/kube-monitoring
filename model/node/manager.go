package node

import (
	"fmt"
	"hash/fnv"
	"sync"
	atomic "sync/atomic"

	logrus "github.com/sirupsen/logrus"
	"mobigen.com/iris-cloud-monitoring/model"
)

// Manager node manager
type Manager struct {
	hashList *Nodes
	maxSize  uint32
	hashSize uint32
	Log      logrus.FieldLogger
}

// var singletonLock *sync.Mutex
var manager *Manager

// Init init node hash list
func Init(_maxSize, _hashSize int, _log logrus.FieldLogger) error {
	// singletonLock = new(sync.Mutex)
	manager = new(Manager)

	manager.hashList = new(Nodes)
	manager.hashList.List = make(map[uint32][]*Node)
	if manager.hashList.List == nil {
		return fmt.Errorf("failed to alloc node map list")
	}
	manager.hashSize = uint32(_hashSize)
	manager.maxSize = uint32(_maxSize)
	// Lock
	for i := 0; uint32(i) < manager.hashSize; i++ {
		lock := new(sync.Mutex)
		manager.hashList.Lock = append(manager.hashList.Lock, lock)
	}
	// Node
	for i := 0; uint32(i) < manager.maxSize; i++ {
		node := new(Node)
		node.Idx = uint32(i)
		atomic.StoreInt32(&node.used, model.UNUSED)
		manager.hashList.NodePool = append(manager.hashList.NodePool, node)
	}
	manager.Log = _log.WithFields(map[string]interface{}{
		"data": "node",
	})
	manager.Log.Info("Node HashMap Init Success")
	return nil
}

// GetInstance singleton nodemanager
func GetInstance() *Manager {
	return manager
	// if nodeManager == nil {
	// 	singletonLock.Lock()
	// 	defer singletonLock.Unlock()
	// 	if nodeManager== nil {
	// 		nodeManager = &foo{1}
	// 	}
	// }
	// return instance
}

// Alloc alloc node
func (manager *Manager) Alloc() (*Node, error) {
	// Log.Debugf("Node Alloc")
	// node := new(Node)
	// if node == nil {
	// 	return nil, fmt.Errorf("error. failed to alloc node")
	// }
	// atomic.StoreInt32(&node.used, USED)
	// atomic.StoreInt32(&node.cnt, 2)

	for idx := 0; uint32(idx) < manager.maxSize; idx++ {
		if atomic.LoadInt32(&manager.hashList.NodePool[idx].used) != model.UNUSED {
			continue
		}
		atomic.StoreInt32(&manager.hashList.NodePool[idx].used, model.CREATING)
		atomic.StoreInt32(&manager.hashList.NodePool[idx].cnt, 2)
		return manager.hashList.NodePool[idx], nil
	}
	return nil, fmt.Errorf("error. can not alloc node. node full[ %d ]", manager.maxSize)
}

// GetByIdx get node by idx
func (manager *Manager) GetByIdx(idx int) (*Node, error) {
	if idx < 0 || uint32(idx) > manager.maxSize {
		return nil, fmt.Errorf("error. input index out of bound[ %d ][ %d - %d ]", idx, 0, manager.maxSize)
	}
	if atomic.LoadInt32(&manager.hashList.NodePool[idx].used) != model.USED {
		return nil, fmt.Errorf("error. not used node idx[ %d ]", idx)
	}
	atomic.AddInt32(&manager.hashList.NodePool[idx].cnt, 1)
	return manager.hashList.NodePool[idx], nil
}

func (manager *Manager) getHashValue(name string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(name))
	return h.Sum32() % manager.hashSize
}

// Get get node from cid(clusterid), name
func (manager *Manager) Get(cid, name string) (*Node, error) {
	key := manager.getHashValue(name)
	manager.hashList.Lock[key].Lock()
	for _, node := range manager.hashList.List[key] {
		if atomic.LoadInt32(&node.used) != model.USED ||
			node.CID != cid || node.Name != name {
			continue
		}
		atomic.AddInt32(&node.cnt, 1)
		manager.hashList.Lock[key].Unlock()
		return node, nil
	}
	manager.hashList.Lock[key].Unlock()
	return nil, fmt.Errorf("error. can not found node cid[ %s ] name[ %s ]", cid, name)
}

// Add node attach to hashlist
func (n *Node) Add() {
	manager := GetInstance()
	key := manager.getHashValue(n.Name)
	n.key = key
	manager.hashList.Lock[key].Lock()
	atomic.StoreInt32(&n.used, model.USED)
	n.listIdx = int32(len(manager.hashList.List[key]))
	manager.hashList.List[key] = append(manager.hashList.List[key], n)
	manager.Log.Infof("Add. Key[ %d ] Idx[ %d ] CID[ %s ] Name[ %s ]", n.key, n.listIdx, n.CID, n.Name)
	manager.hashList.Lock[key].Unlock()
}

// Put put node
func (n *Node) Put() {
	if atomic.LoadInt32(&n.cnt) <= 0 {
		return
	}
	if atomic.AddInt32(&n.cnt, -1) < 1 {
		manager := GetInstance()
		hash := n.key
		manager.hashList.Lock[hash].Lock()
		manager.Log.Infof("Put. Key[ %d ] Idx[ %d ] CID[ %s ] Name[ %s ]", n.key, n.listIdx, n.CID, n.Name)
		// Delete Node From HashMap List
		manager.hashList.List[hash] = append(manager.hashList.List[hash][:n.listIdx], manager.hashList.List[hash][n.listIdx+1:]...)
		// Clear Node Data
		n.Clear()
		manager.hashList.Lock[hash].Unlock()
	}
}

// Del del node
func (n *Node) Del() error {
	if atomic.CompareAndSwapInt32(&n.used, model.USED, model.DELETE) {
		manager := GetInstance()
		manager.Log.Infof("Del. Key[ %d ] Idx[ %d ] CID[ %s ] Name[ %s ]", n.key, n.listIdx, n.CID, n.Name)
		n.Put()
		return nil
	}
	return fmt.Errorf("Del. Already delete Key[ %d ] Idx[ %d ] CID[ %s ] Name[ %s ]", n.key, n.listIdx, n.CID, n.Name)
}

// Clear golang no have memset
func (n *Node) Clear() {
	manager := GetInstance()
	manager.Log.Infof("Clear. Key[ %d ] Idx[ %d ] CID[ %s ] Name[ %s ]", n.key, n.listIdx, n.CID, n.Name)

	n.key = 0
	n.listIdx = 0
	n.CID = ""
	n.Time = 0
	n.Name = ""
	n.NodeExporterURL = ""
	n.CPU = CPU{Cores: make(map[string]Cores)}
	n.Mem = Memory{}
	n.Net = Network{Devices: []NetDevice{}}
	n.File = FileSystem{Devices: []FileDevice{}}
	n.LoadAverage = LoadAverage{}
	n.Disk = Disks{Devices: []DiskDevice{}}

	atomic.StoreInt32(&n.cnt, 0)
	atomic.StoreInt32(&n.used, model.UNUSED)
}

// NodeGetPart get node list
// func NodeGetPart(start, end uint32) ([]*Node, error) {
// 	if start > end {
// 		return nil, fmt.Errorf("error. wrong arguments. start[ %d ] end[ %d ]", start, end)
// 	}
// 	var ret []*Node
// 	for key, list := range hashList.List {
// 		hashList.Lock[key].Lock()
// 		if key < start || end <= key {
// 			hashList.Lock[key].Unlock()
// 			continue
// 		}
// 		for _, v := range list {
// 			// Clone ( deep copy )
// 			clone := v.Clone()
// 			// Append
// 			ret = append(ret, clone)
// 		}
// 		hashList.Lock[key].Unlock()
// 	}
// 	return ret, nil
// }

// DebugForeach debug
func (manager *Manager) DebugForeach() {
	manager.Log.Debugf("HashMap List foreach")
	for key, list := range manager.hashList.List {
		for _, v := range list {
			manager.Log.Debugf("HashKey[ %d ] %+v", key, *v)
		}
	}
	manager.Log.Debugf("Node List foreach")
	for _, node := range manager.hashList.NodePool {
		manager.Log.Debugf("Idx[ %d ] %+v", node.Idx, node)
	}
}
