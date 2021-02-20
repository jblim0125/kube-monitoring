package main

import (
	"fmt"
	"sync"
	"time"

	nm "mobigen.com/iris-cloud-monitoring/model/node"
	"mobigen.com/iris-cloud-monitoring/tools"
)

// StartSaver node, kube_state save to databases
func (m *Monitoring) StartSaver(id int) {
	Log := m.logger.WithFields(map[string]interface{}{
		"saver": id,
	})
	Log.Infof("Saver[ %d ] Started", id)

	var now, next, remain int64
	var period int64 = int64(m.config.ThreadConfig.Period)
	now = time.Now().UnixNano() / int64(time.Second)
	remain = period - (now % period)
	next = (now + remain + 1) * 1000

	nodeManager := nm.GetInstance()
	for {
		now = time.Now().UnixNano() / int64(time.Millisecond)
		if now >= next {
			Log.Infof("Wait All Finish.....")
			if waitTimeout(&m.syncWaitGroup, 2*time.Second) {
				Log.Infof("Incomplete Finish..")
			} else {
				Log.Infof("All Complete Finish..")
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
				Log.Infof("Node Update Info. Time[ %s ]", strTime)
				// PrintYaml(node)

				// Must
				node.Put()
			}

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
