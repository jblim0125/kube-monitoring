package main

import (
	"sync"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	logrus "github.com/sirupsen/logrus"
	model "mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	nm "mobigen.com/iris-cloud-monitoring/model/node"
	"mobigen.com/iris-cloud-monitoring/tools"
)

// StartGatherer node-exporter, kube-state-metric data getter
func (m *Monitoring) StartGatherer(id int) {
	Log := m.logger.WithFields(map[string]interface{}{
		"gatherer": id,
	})

	Log.Infof("Gatherer[ %d ] Started", id)

	var now, next, remain int64
	var period int64 = int64(m.config.ThreadConfig.Period)
	now = time.Now().UnixNano() / int64(time.Second)
	remain = period - (now % period)
	next = (now + remain) * 1000

	nodeManager := nm.GetInstance()
	ksmManager := ksm.GetInstance()
	for {
		now = time.Now().UnixNano() / int64(time.Millisecond)
		if now >= next {
			// Node-Exporter
			for nodeIdx := 0; nodeIdx < m.config.HashConfig.NodeHashList.MaxSize; nodeIdx++ {
				node, err := nodeManager.GetByIdx(nodeIdx)
				if err != nil {
					continue
				}
				Log.Infof("Try Get Metrics. CID[ %s ] Node[ %s ] URL[ %s ]", node.CID, node.Name, node.NodeExporterURL)
				go GetMetricsNSend(model.NODE, node, m.queue[id], Log, &m.syncWaitGroup)
				node.Put()
			}
			// Kube State Metrics
			clusters := ksmManager.GetClusters()
			for _, cluster := range clusters {
				Log.Infof("Try Get Metrics. CID[ %s ] Kube State Metric URL[ %s ]", cluster.ID, cluster.KsmURL)
				go GetMetricsNSend(model.KUBESTATE, &cluster, m.queue[id], Log, &m.syncWaitGroup)
			}

			// Calc Next Runtime
			now = time.Now().UnixNano() / int64(time.Second)
			remain = period - (now % period)
			next = (now + remain) * 1000
		}
		// Sleep
		time.Sleep(20 * time.Millisecond)
	}
}

// GetMetricsNSend get and send
func GetMetricsNSend(metricType int32, data interface{}, queue *goconcurrentqueue.FIFO,
	Log logrus.FieldLogger, waitGroup *sync.WaitGroup) {
	var idx uint32
	var URL, cID, name string
	if metricType == model.NODE {
		node := data.(*nm.Node)
		idx = node.Idx
		cID = node.CID
		name = node.Name
		URL = node.NodeExporterURL
	} else if metricType == model.KUBESTATE {
		cluster := data.(*ksm.Cluster)
		cID = cluster.ID
		name = cluster.Name
		URL = cluster.KsmURL
	}

	res, err := tools.GetHTTPRequest(URL)
	if err != nil {
		Log.Errorf("ERROR. Get Http Request.[ %s ]", err.Error())
	} else {
		Log.Debugf("Success Get Metrics. URL[ %s ]", URL)

		// Add waitgroups
		Log.Debugf("Add sync.waitgroup")
		waitGroup.Add(1)

		// Send To Proc Th
		sendMsg := model.MetricData{
			Type: metricType,
			Idx:  idx,
			CID:  cID,
			Name: name,
			Time: tools.GetMillis(),
			Len:  uint32(len(res)),
			Data: res,
		}
		for {
			if queue.Enqueue(&sendMsg) == nil {
				if metricType == model.NODE {
					Log.Infof("[ GET >> PROC ] Send Metrics Data. Metric Data Type[ NODE ] Receive From Name[ %s ] URL[ %s ].", name, URL)
				} else {
					Log.Infof("[ GET >> PROC ] Send Metrics Data. Metric Data Type[ KUBE_STATE ] Receive From Name[ %s ] URL[ %s ].", name, URL)
				}
				break
			}
		}
	}
}
