package model

// Metric Data의 타입
const (
	// Node node
	NODE int32 = 0
	// KUBESTATE kube state
	KUBESTATE int32 = 1
)

// Data Management Value
const (
	UNUSED   int32 = 0
	CREATING int32 = 1
	USED     int32 = 2
	DELETE   int32 = 3
)
