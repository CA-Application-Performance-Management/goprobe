package goprobe

import (
	"strings"
	"testing"
)

func TestGetCommandConnection(t *testing.T) {
	ecmd := "{\"op\":\"commandConn\",\"probe\":\"go\",\"ver\":1,\"instid\":\"10752\",\"pgm\":\"goProbe\",\"prms\":{\"appName\":\"Go APP\",\"hostName\":\"localhost\"}}"
	cmdConn := commandconnection{}
	cmdStr := cmdConn.getCmdConnection()

	if !strings.Contains(cmdStr, "\"op\":\"commandConn\"") {
		t.Errorf("Command connection string is wrong for op, got : %s want : %s", cmdStr, ecmd)
	}

	if !strings.Contains(cmdStr, "\"probe\":\"go\"") {
		t.Errorf("Command connection string is wrong for probe, got : %s want : %s", cmdStr, ecmd)
	}

	if !strings.Contains(cmdStr, "\"ver\":1") {
		t.Errorf("Command connection string is wrong for version, got : %s want : %s", cmdStr, ecmd)
	}

	if !strings.Contains(cmdStr, "\"pgm\":\"goProbe\"") {
		t.Errorf("Command connection string is wrong for pgm, got : %s want : %s", cmdStr, ecmd)
	}

	if !strings.Contains(cmdStr, "\"prms\":{\"appName\":\"Go APP\",\"hostName\":\"localhost\"}") {
		t.Errorf("Command connection string is wrong for prms, got : %s want : %s", cmdStr, ecmd)
	}
}

func TestGetDataConnection(t *testing.T) {
	edata := "{\"op\":\"dataConn\",\"probe\":\"go\",\"ver\":1,\"instid\":\"4644\",\"tid\":\"4644\",\"pid\":\"4644\"}"
	dataConn := dataconnection{}
	adata := dataConn.getDataConnection()

	if !strings.Contains(adata, "\"op\":\"dataConn\"") {
		t.Errorf("Data connection string is wrong for op, got : %s want : %s", adata, edata)
	}

	if !strings.Contains(adata, "\"probe\":\"go\"") {
		t.Errorf("Data connection string is wrong for probe, got : %s want : %s", adata, edata)
	}

	if !strings.Contains(adata, "\"ver\":1") {
		t.Errorf("Data connection string is wrong for ver, got : %s want : %s", adata, edata)
	}
}
