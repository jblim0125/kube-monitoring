package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"mobigen.com/iris-cloud-monitoring/model"
	"mobigen.com/iris-cloud-monitoring/tools"
)

func main() {
	logger.SetLevel(logrus.DebugLevel)
	monitoring := NewMonitoring(nil, tools.NewID(), logger)
	config := model.MonitoringConfig{
		ThreadConfig: model.MonitoringThInfo{
			Period:  5,
			NumProc: 2,
			NumSave: 1,
		},
		HashConfig: model.HashConfig{
			NodeHashList: model.HashListConfig{
				MaxSize: 50, HashSize: 10,
			},
			ContainerHashList: model.HashListConfig{
				MaxSize: 200, HashSize: 20,
			},
		},
		CollectionConfig: model.CollectionConfig{
			FileSystems: map[string]string{
				"/dev/mapper/centos-root": "/",
			},
			NetDevices: []string{
				"eth0",
			},
			DiskDevices: []string{
				"vda",
			},
		},
		Additional: model.AdditionalConfig{},
	}
	monitoring.InitMonitoring(config)
	PrintYaml(config)

	// Add kube state metrics Cluster Info
	monitoring.AddTestKubeStateMetrics()
	// Add Node Info
	// monitoring.AddTestNodeExporterInfo()

	monitoring.StartMonitoring()

	// monitoring.TestMonitoring()
	for {
		time.Sleep(1 * time.Second)
	}
}

// PrintYaml print yaml format
func PrintYaml(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	// encoder.SetIndent("", "    ")
	return encoder.Encode(data)
}
