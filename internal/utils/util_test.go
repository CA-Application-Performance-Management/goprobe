package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetSeqTag(t *testing.T) {
	tag := GetSeqTag("fnC")
	if tag != "seq" {
		t.Errorf("Sequence tag for fnC, got : %s want : %s", tag, "seq")
	}

	tag = GetSeqTag("fnR")
	if tag != "cseq" {
		t.Errorf("Sequence tag for fnR, got : %s want : %s", tag, "cseq")
	}

	tag = GetSeqTag("fnRC")
	if len(tag) != 0 {
		t.Errorf("Sequence tag for Unknown tag should be empty, got : %s want : %s", tag, "")
	}
}

func TestIsValidPath(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Failed to retrieve the working directory")
	}

	result := IsValidPath(dir)
	if !result {
		t.Errorf("Is not a Valid directory, got : %t want : %t", result, true)
	}

	//Prepare a wrong directory
	dir = dir + "\test"
	result = IsValidPath(dir)
	if result {
		t.Errorf("Is not a Valid directory, got : %t want : %t", result, false)
	}
}

func TestIsValidHost(t *testing.T) {
	result := IsValidHost("localhost:80")
	if !result {
		t.Errorf("Is not a valid host, got : %t want : %t", result, true)
	}

	//Prepare wrong host name
	result = IsValidHost("localhostt:80")
	if result {
		t.Errorf("Is not a valid host, got : %t want : %t", result, false)
	}
}

func TestConvertMetricToARF(t *testing.T) {

	actual := ConvertMetricToARF("m", "GC Heap:Bytes In Use", 8194, "500")
	expected := fmt.Sprintf("{\"op\":\"m\",\"mid\":\"GC Heap:Bytes In Use\",\"mtype\":\"8194\",\"val\":\"500\",\"pid\":\"%d\"}\n", os.Getpid())

	if expected != actual {
		t.Errorf("Can not convert metric, got : %s want : %s", actual, expected)
	}

}
