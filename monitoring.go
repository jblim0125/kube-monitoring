package main

import (
	"fmt"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"
	logrus "github.com/sirupsen/logrus"
	model "mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	nm "mobigen.com/iris-cloud-monitoring/model/node"
)

// MonitoringStore database interface
type MonitoringStore interface {
	// GetClusters(clusterFilter *model.ClusterFilter) ([]*model.Cluster, error)
}

// MonitoringProvisioner k8s cluster interface
type MonitoringProvisioner interface {
	// GetNodeInfo(cluster *model.Cluster) (*v1.NodeList, error)
}

// Monitoring monitoring object
type Monitoring struct {
	instanceID    string
	store         MonitoringStore
	queue         []*goconcurrentqueue.FIFO
	config        model.MonitoringConfig
	syncWaitGroup sync.WaitGroup
	logger        logrus.FieldLogger
}

// NewMonitoring new concreate
func NewMonitoring(store MonitoringStore, instanceID string, logger logrus.FieldLogger) *Monitoring {
	monitoring := &Monitoring{
		store:      store,
		instanceID: instanceID,
		logger:     logger,
	}
	return monitoring
}

// InitMonitoring init monitoring arg = (proc + save) count
func (monitoring *Monitoring) InitMonitoring(config model.MonitoringConfig) error {
	monitoring.logger.Infof("Init Monitoring")
	if config.ThreadConfig.NumProc <= 0 || config.ThreadConfig.NumSave <= 0 || config.ThreadConfig.Period <= 0 {
		return fmt.Errorf("Must be Proc[ %d ] > 0, Save[ %d ] > 0, Period[ %d ] >= 5",
			config.ThreadConfig.NumProc, config.ThreadConfig.NumSave, config.ThreadConfig.Period)
	}

	monitoring.config = config
	// Receive 기준으로 Queue를 계산
	qSize := monitoring.config.ThreadConfig.NumProc
	for i := 0; i < qSize; i++ {
		queue := goconcurrentqueue.NewFIFO()
		monitoring.queue = append(monitoring.queue, queue)
	}

	err := nm.Init(monitoring.config.HashConfig.NodeHashList.MaxSize,
		monitoring.config.HashConfig.NodeHashList.HashSize, monitoring.logger)
	if err != nil {
		return err
	}

	err = ksm.Init(monitoring.logger)
	if err != nil {
		return err
	}

	return nil
}

// StartMonitoring start sub thread(goroutin)
func (monitoring *Monitoring) StartMonitoring() error {
	monitoring.logger.Infof("Start Monitoring")

	// Start Get
	monitoring.logger.Infof("Start Get Period[ %d ]", monitoring.config.ThreadConfig.Period)
	go monitoring.StartGatherer(0)

	// Start Proc
	for i := 0; i < monitoring.config.ThreadConfig.NumProc; i++ {
		monitoring.logger.Infof("Start Proc ID[ %d ]", i)
		go monitoring.StartProc(i)
	}

	// Start Save
	for i := 0; i < monitoring.config.ThreadConfig.NumSave; i++ {
		monitoring.logger.Infof("Start Save ID[ %d ]", i)
		go monitoring.StartSaver(i)
	}

	return nil
}

// TestMonitoring TestMonitoring
func (monitoring *Monitoring) TestMonitoring() {
	var node *nm.Node
	nodeManager := nm.GetInstance()

	node, _ = nodeManager.Alloc()
	node.CID = "test1"
	node.Name = "test1"
	node.Time = 1234566789
	node.Add()
	node.Put()

	node, _ = nodeManager.Alloc()
	node.CID = "test2"
	node.Name = "test2"
	node.Time = 2234566789
	node.Add()
	node.Put()

	nodeManager.DebugForeach()

	node, err := nodeManager.Get("test1", "test1")
	if err != nil {
		monitoring.logger.Debugf("%s", err.Error())
	} else {
		monitoring.logger.Debugf("SUCCESS : %+v", node)

		monitoring.logger.Debugf("Delete test1")
		err = node.Del()
		if err != nil {
			monitoring.logger.Debugf("error.[ %s ]", err.Error())
		}
		node.Put()
	}

	node, err = nodeManager.Get("test1", "test1")
	if err != nil {
		monitoring.logger.Debugf("%s", err.Error())
	} else {
		monitoring.logger.Debugf("SUCCESS : %+v", node)
		node.Put()
	}

	nodeManager.DebugForeach()

	node, err = nodeManager.Get("test2", "test2")
	if err != nil {
		monitoring.logger.Debugf("%s", err.Error())
	} else {
		monitoring.logger.Debugf("SUCCESS : %+v", node)
		node.Put()
	}

	for i := 0; i < 20; i++ {
		node, _ = nodeManager.Alloc()
		node.CID = fmt.Sprintf("worker%d", i)
		node.Name = fmt.Sprintf("worker%d", i)
		node.Time = 1234566789
		node.Add()
		node.Put()
	}
	nodeManager.DebugForeach()
}

// AddTestNodeExporterInfo for node exporter test
func (monitoring *Monitoring) AddTestNodeExporterInfo() {
	var node *nm.Node
	nodeManager := nm.GetInstance()
	node, _ = nodeManager.Alloc()
	node.CID = "external-cluster"
	node.Name = "k8s-master"
	node.NodeExporterURL = "http://192.168.50.67:9100/metrics"
	node.Add()
	node.Put()

	node, _ = nodeManager.Alloc()
	node.CID = "external-cluster"
	node.Name = "k8s-node01"
	node.NodeExporterURL = "http://192.168.50.56:9100/metrics"
	node.Add()
	node.Put()

	node, _ = nodeManager.Alloc()
	node.CID = "external-cluster"
	node.Name = "k8s-node02"
	node.NodeExporterURL = "http://192.168.50.65:9100/metrics"
	node.Add()
	node.Put()

	node, _ = nodeManager.Alloc()
	node.CID = "external-cluster"
	node.Name = "k8s-node03"
	node.NodeExporterURL = "http://192.168.50.53:9100/metrics"
	node.Add()
	node.Put()

	node, _ = nodeManager.Alloc()
	node.CID = "external-cluster"
	node.Name = "k8s-node04.novalocal"
	node.NodeExporterURL = "http://192.168.50.68:9100/metrics"
	node.Add()
	node.Put()
}

// AddTestKubeStateMetrics add kube state metric for test
func (monitoring *Monitoring) AddTestKubeStateMetrics() {
	manager := ksm.GetInstance()
	manager.AddCluster("0123456789abcde", "external-cluster", "http://192.168.50.67:30445/metrics")

	manager.DebugForeach()
}
