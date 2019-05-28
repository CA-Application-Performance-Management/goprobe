/*Package utils ...*/
/*
This package contains all the utility functions used by the SDK
This package is internal to the SDK
*/
package utils

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//IsValidPath finds the given path is a valid one or not
//Returns true if valid
//Retruns false if invalid
func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

//IsValidHost checks the given connection details are valid
//Returns true if valid hostName & hostPort
//Returns false if invalid hostName & hostPort
func IsValidHost(service string) bool {
	_, err := net.ResolveTCPAddr("tcp4", service)
	if err == nil {
		return true
	}
	return false
}

//GetSeqTag return the sequence tag based on the operation
//if the operation is "fnC" it returns "seq"
//if the operation is "fnR" it returns "cseq"
//if other operation return empty
func GetSeqTag(op string) string {
	if strings.Compare(op, "fnC") == 0 {
		return "seq"
	}

	if strings.Compare(op, "fnR") == 0 {
		return "cseq"
	}

	return ""
}

//ConvertMetricToARF converts the metric to ARF format
//it takes fields of metrics and convert to ARF format(JSON)
func ConvertMetricToARF(operation string, id string, mettype int, value string) string {
	var arfJSON string
	arfJSON += "{\"op\":\""
	arfJSON += operation
	arfJSON += "\",\"mid\":\""
	arfJSON += id
	arfJSON += "\",\"mtype\":"
	arfJSON += "\""
	arfJSON += strconv.Itoa(mettype)
	arfJSON += "\",\"val\":"
	arfJSON += "\""
	arfJSON += value
	arfJSON += "\""
	arfJSON += ",\"pid\":\""
	arfJSON += strconv.Itoa(os.Getpid())
	arfJSON += "\""
	arfJSON += "}\n"
	return arfJSON
}

//ConvertMessageToARF converts the message to ARF format
//it takes fields of message type and convert to ARF format(JSON)
func ConvertMessageToARF(tid string, op string, fn string, seq int, prms string) string {
	t := time.Now()
	ts := t.Format("20060102150405")

	var arfJSON string = ""

	arfJSON += "{\"op\":\""
	arfJSON += op
	arfJSON += "\",\"fn\":\""
	arfJSON += fn
	arfJSON += "\",\"ts\":"
	arfJSON += ts
	arfJSON += ",\"pid\":\""
	arfJSON += strconv.Itoa(os.Getpid())
	arfJSON += "\""
	arfJSON += ",\"tid\":\""
	arfJSON += tid
	arfJSON += "\",\""
	arfJSON += GetSeqTag(op)
	arfJSON += "\":"
	arfJSON += strconv.Itoa(seq)

	if prms != "" {
		arfJSON += prms
	}
	arfJSON += "}\n"

	return arfJSON
}
