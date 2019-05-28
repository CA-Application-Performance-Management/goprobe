package metric

import (
	"testing"
)

func TestCollectRuntimeMetrics(t *testing.T) {
	var expectedMetrics []Metric
	metric1 := Metric{"m", "GC Heap:Total Heap Allocated", 8194, "test"}
	expectedMetrics = append(expectedMetrics, metric1)
	metric2 := Metric{"m", "GC Heap:Total Heap Free", 8194, "test"}
	expectedMetrics = append(expectedMetrics, metric2)
	metric3 := Metric{"m", "GC Heap:Bytes In Use", 8194, "test"}
	expectedMetrics = append(expectedMetrics, metric3)
	metric4 := Metric{"m", "GC Heap:Bytes Total", 8194, "test"}
	expectedMetrics = append(expectedMetrics, metric4)
	metric5 := Metric{"m", "GC Heap:Routines Total", 8194, "test"}
	expectedMetrics = append(expectedMetrics, metric5)

	actualMetrics := CollectRuntimeMetrics()

	if len(actualMetrics) != len(expectedMetrics) {
		t.Errorf("Invalid count of runtime metrics, got : %d want : %d", len(actualMetrics), 5)
	}

	for _, am := range actualMetrics {
		for _, em := range expectedMetrics {
			if am.Id == em.Id {
				//runtime value can not be asserted
				if am.Op != em.Op {
					t.Errorf("Failed to collect Runtime Metrics, got : %s want : %s", am.Op, em.Op)
				}
				if am.Type != em.Type {
					t.Errorf("Failed to collect Runtime Metrics, got : %d want : %d", am.Type, em.Type)
				}
				break
			}
		}
	}
}

func TestCollectStaticMetrics(t *testing.T) {
	var expectedMetrics []Metric
	metric1 := Metric{"m", "ProcessID", 4101, "test"}
	expectedMetrics = append(expectedMetrics, metric1)
	metric2 := Metric{"m", "GO Version", 4101, "test"}
	expectedMetrics = append(expectedMetrics, metric2)
	metric3 := Metric{"m", "Logical CPU count", 4101, "test"}
	expectedMetrics = append(expectedMetrics, metric3)
	metric4 := Metric{"m", "Host Name", 4101, "test"}
	expectedMetrics = append(expectedMetrics, metric4)

	actualMetrics := CollectStaticMetrics()

	if len(actualMetrics) != len(expectedMetrics) {
		t.Errorf("Invalid count of static metrics, got : %d want : %d", len(actualMetrics), 3)
	}

	for _, am := range actualMetrics {
		for _, em := range expectedMetrics {
			if am.Id == em.Id {
				//runtime value can not be asserted
				if am.Op != em.Op {
					t.Errorf("Failed to collect Static Metrics, got : %s want : %s", am.Op, em.Op)
				}
				if am.Type != em.Type {
					t.Errorf("Failed to collect Static Metrics, got : %d want : %d", am.Type, em.Type)
				}
				break
			}
		}
	}
}
