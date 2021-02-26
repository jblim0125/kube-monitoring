package main

import (
	"fmt"
	"sync"
	"time"

	logrus "github.com/sirupsen/logrus"
	model "mobigen.com/iris-cloud-monitoring/model"
	ksm "mobigen.com/iris-cloud-monitoring/model/ksm"
	nm "mobigen.com/iris-cloud-monitoring/model/node"
	"mobigen.com/iris-cloud-monitoring/tools"
)

// Saver saver data struct
type Saver struct {
	ID     int
	Log    logrus.FieldLogger
	config model.MonitoringConfig
}

// StartSaver node, kube_state save to databases
func (m *Monitoring) StartSaver(id int) {
	saver := Saver{
		ID:     id,
		Log:    m.logger.WithFields(map[string]interface{}{"saver": id}),
		config: m.config,
	}
	saver.Log.Infof("Saver[ %d ] Started", id)

	var now, next, remain int64
	var period int64 = int64(m.config.ThreadConfig.Period)
	now = time.Now().UnixNano() / int64(time.Second)
	remain = period - (now % period)
	next = (now + remain + 1) * 1000

	nodeManager := nm.GetInstance()
	ksmManager := ksm.GetInstance()
	for {
		now = time.Now().UnixNano() / int64(time.Millisecond)
		if now >= next {
			saver.Log.Infof("[ PROC >> SAVE ] Waiting.....")
			if waitTimeout(&m.syncWaitGroup, 2*time.Second) {
				saver.Log.Infof("[ PROC >> SAVE ] Incomplete Finish...")
			} else {
				saver.Log.Infof("[ PROC >> SAVE ] All Complete Finish.")
			}

			// Get Node
			for nodeIdx := 0; nodeIdx < m.config.HashConfig.NodeHashList.MaxSize; nodeIdx++ {
				node, err := nodeManager.GetByIdx(nodeIdx)
				if err != nil {
					continue
				}

				t, _ := tools.GetLocalTime(node.Time, "Korea")
				strTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
					t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
				saver.Log.Infof("Node Update Info. Time[ %s ]", strTime)
				// PrintYaml(node)

				// Must
				node.Put()
			}

			// Kube State Metrics Check
			saver.SaveKubeStateMetrics(ksmManager)

			// Calc Next Runtime
			now = time.Now().UnixNano() / int64(time.Second)
			remain = period - (now % period)
			next = (now + remain + 1) * 1000
		}
		// Sleep
		time.Sleep(100 * time.Millisecond)
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// SaveKubeStateMetrics save or data clean kube state metrics
func (saver *Saver) SaveKubeStateMetrics(ksmManager *ksm.Manager) error {
	// POD
	ksmManager.CheckKSMData()

	// ksmManager.GetPod()
	return nil
}
