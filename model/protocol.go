package model

// MetricData Queue Protocol ( Get => Proc )
type MetricData struct {
	Type int32
	Idx  uint32
	CID  string
	Name string
	Time int64
	Len  uint32
	Data []byte
}
