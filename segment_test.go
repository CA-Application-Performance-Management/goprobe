package goprobe

import (
	"context"
	"testing"
)

func TestStartSegment(t *testing.T) {
	//data prep
	apmTidData["123456"] = &seqGen{1, "1234"}
	es := Segment{2, "123456", "TestFunc1", "fs"}
	var ctx context.Context
	ctx = context.WithValue(ctx, "tid", "123456")
	ctx = context.WithValue(ctx, "corid", "1234")
	ctx = context.WithValue(ctx, "seqid", 1)
	InitGoProbe()

	as := StartSegment(ctx, "TestFunc1")
	if es != as {
		t.Errorf("Failed to Start Segment, got : %v want : %v", as, es)
	}

	seqgen := apmTidData["123456"]
	if seqgen.seqNo != 2 {
		t.Errorf("Invalid sequence in Start Segment, got : %d want : %d", seqgen.seqNo, 2)
	}
}

func TestStartSegmentSQL(t *testing.T) {
	//data prep
	apmTidData["123456"] = &seqGen{1, "1234"}
	es := Segment{2, "123456", "TestSegment1", "mysql"}
	var ctx context.Context
	ctx = context.WithValue(ctx, "tid", "123456")
	ctx = context.WithValue(ctx, "corid", "1234")
	ctx = context.WithValue(ctx, "seqid", 1)
	InitGoProbe()

	as := StartSegment(ctx, "TestSegment1", "models.booking", "select * from booking")
	if es != as {
		t.Errorf("Failed to Start Segment, got : %v want : %v", as, es)
	}

	seqgen := apmTidData["123456"]
	if seqgen.seqNo != 2 {
		t.Errorf("Invalid sequence in Start Segment, got : %d want : %d", seqgen.seqNo, 2)
	}
}

func TestStartSegmentNC(t *testing.T) {
	InitGoProbe()
	es := Segment{1, "123456", "TestFunc1", "fs"}

	as := StartSegment(nil, "TestFunc1")

	if as.seqNo != es.seqNo {
		t.Errorf("Invalid sequence in Start Segment, got : %d want : %d", as.seqNo, es.seqNo)
	}

	if as.segType != es.segType {
		t.Errorf("Invalid sequence type in Start Segment, got : %s want : %s", as.segType, es.segType)
	}

	if as.funcName != es.funcName {
		t.Errorf("Invalid function name in Start Segment, got : %s want : %s", as.funcName, es.funcName)
	}
}

func TestEndSegment(t *testing.T) {
	apmTidData["123456"] = &seqGen{1, "1234"}
	apmTidData["654321"] = &seqGen{2, "4321"}
	es := Segment{1, "123456", "TestSegment1", "fs"}
	InitGoProbe()

	es.EndSegment()
	if len(apmTidData) != 2 {
		t.Errorf("Failed to End Segment, got : %d want : %d", len(apmTidData), 2)
	}
	seqgen := apmTidData["123456"]
	if seqgen.seqNo != 0 {
		t.Errorf("Failed to End Segment wrong seqNo, got : %d want : %d", seqgen.seqNo, 0)
	}

	es = Segment{2, "654321", "TestSegment2", "fs"}
	es.EndSegment()
	if len(apmTidData) != 2 {
		t.Errorf("Failed to End Segment, got : %d want : %d", len(apmTidData), 2)
	}

	seqgen = apmTidData["654321"]
	if seqgen.seqNo != 1 {
		t.Errorf("Failed to End Segment wrong seqNo, got : %d want : %d", seqgen.seqNo, 1)
	}
}
