/*Package metric ...*/
/*
This package contains all the runtime metrics collected by the SDK
This is used to collect the runtime statistics of the process and publish it to IA
This runs in a separate go process based on time interval configured in config
*/
package metric

import (
	"fmt"
	"os"
	"runtime"
)

//The Metric type has the information about the metrics collected like
// op - Type of Operation
// Id - Display ID of the metric in IA
// Type - Type of the metrics like counter,timepahse
// Value - Value of the metric
type Metric struct {
	Op    string
	Id    string
	Type  int
	Value string
}

const (
	kintervalCounter = 8194
	kstringType      = 21
	kstringIndEvents = 4101
)

//CollectRuntimeMetrics is used to retrieve the runtime metrics of the process
//The following metrics are collected at present and this can be extended to collect more metrics
// heapObjects -- Cumulative count of heap objects allocation(Total Heap Allocated)
// heapObjectsFreed -- Cumulative count of heap object freed(Total Heap Free)
// allocation -- Bytes of allocated heap objects (Bytes In Use)
// totalAllocation -- Cumulative bytes allocated for heap objects (Bytes Total)
// goRoutines -- Represents the number of currently active goRoutines for the process(Routines Total)
func CollectRuntimeMetrics() []Metric {
	var metrics []Metric
	var rm runtime.MemStats
	runtime.ReadMemStats(&rm)

	heapObjects := Metric{
		Op:    "m",
		Id:    "GC Heap:Total Heap Allocated",
		Type:  kintervalCounter,
		Value: fmt.Sprintf("%v", rm.Mallocs),
	}
	metrics = append(metrics, heapObjects)

	heapObjectsFreed := Metric{
		Op:    "m",
		Id:    "GC Heap:Total Heap Free",
		Type:  kintervalCounter,
		Value: fmt.Sprintf("%v", rm.Frees),
	}
	metrics = append(metrics, heapObjectsFreed)

	allocation := Metric{
		Op:    "m",
		Id:    "GC Heap:Bytes In Use",
		Type:  kintervalCounter,
		Value: fmt.Sprintf("%v", rm.Alloc),
	}
	metrics = append(metrics, allocation)

	totalAllocation := Metric{
		Op:    "m",
		Id:    "GC Heap:Bytes Total",
		Type:  kintervalCounter,
		Value: fmt.Sprintf("%v", rm.TotalAlloc),
	}
	metrics = append(metrics, totalAllocation)

	goRoutines := Metric{
		Op:    "m",
		Id:    "GC Heap:Routines Total",
		Type:  kintervalCounter,
		Value: fmt.Sprintf("%v", runtime.NumGoroutine()),
	}
	metrics = append(metrics, goRoutines)

	return metrics
}

//CollectStaticMetrics is used to retrieve the metrics of the process
//The following metrics are collected at only once
// processID -- Current ID of the running process
// version -- GO Version by the running process
// cpu -- Logical CPU count of the process
// HostName -- Name of the host the probe is running
func CollectStaticMetrics() []Metric {
	var metrics []Metric

	processID := Metric{
		Op:    "m",
		Id:    "ProcessID",
		Type:  kstringIndEvents,
		Value: fmt.Sprintf("%v", os.Getpid()),
	}
	metrics = append(metrics, processID)

	version := Metric{
		Op:    "m",
		Id:    "GO Version",
		Type:  kstringIndEvents,
		Value: fmt.Sprintf("%v", runtime.Version()),
	}
	metrics = append(metrics, version)

	cpu := Metric{
		Op:    "m",
		Id:    "Logical CPU count",
		Type:  kstringIndEvents,
		Value: fmt.Sprintf("%v", runtime.NumCPU()),
	}
	metrics = append(metrics, cpu)

	hname, err := os.Hostname()
	if err == nil {
		name := Metric{
			Op:    "m",
			Id:    "Host Name",
			Type:  kstringIndEvents,
			Value: fmt.Sprintf("%v", hname),
		}
		metrics = append(metrics, name)
	}
	return metrics
}
