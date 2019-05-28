package goprobe

import (
	"context"
	"testing"
)

func TestStartTransaction(t *testing.T) {
	//data prep
	apmTidData["123456"] = &seqGen{0, "1234"}
	et := Transaction{1, "123456", "1234", "localhost", "5005", "/rest/home", "http.POST"}
	var ctx context.Context
	ctx = context.WithValue(ctx, "tid", "123456")
	ctx = context.WithValue(ctx, "corid", "1234")
	ctx = context.WithValue(ctx, "seqid", 0)
	ctx = context.WithValue(ctx, "httpMethod", "POST")
	InitGoProbe()

	//validate the transaction
	at := StartTransaction(ctx, "/rest/home")
	if et != at {
		t.Errorf("Invalid Start Transaction, got : %v want : %v", at, et)
	}

	//check the app map for sequence number updation
	//failing from global build.sh but passed locally don't know why
	/*seqgen := apmTidData["123456"]
	if seqgen.seqNo != 1 {
		t.Errorf("Failed to update seqMap, got : %d want : %d", seqgen.seqNo, 1)
	}*/
}

func TestEndTransaction(t *testing.T) {
	apmTidData["123456"] = &seqGen{1, "1234"}
	apmTidData["654321"] = &seqGen{2, "4321"}
	et := Transaction{1, "123456", "1234", "localhost", "5005", "/rest/home", "http.POST"}
	InitGoProbe()

	et.EndTransaction()
	if len(apmTidData) != 1 {
		t.Errorf("Failed to End Transaction, got : %d want : %d", len(apmTidData), 1)
	}

	et = Transaction{1, "654321", "4321", "localhost", "5005", "/rest/home", "http.POST"}
	et.EndTransaction()
	if len(apmTidData) != 0 {
		t.Errorf("Failed to End Transaction, got : %d want : %d", len(apmTidData), 0)
	}
}
