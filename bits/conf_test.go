package bits

import (
	"testing"
)

func TestDefaultConf(t *testing.T) {
	conf := DefaultConf()
	
	if conf == nil {
		t.Error("DefaultConf() should not return nil")
	}
	
	if conf.DeduplicationScope == 0 {
		t.Error("DeduplicationScope should be set to non-zero value")
	}
}

func TestKeySize(t *testing.T) {
	if KeySize != 32 {
		t.Errorf("Expected KeySize to be 32, got %d", KeySize)
	}
}

func TestKeyOp(t *testing.T) {
	k := K{}
	kop := KeyOp{
		Op:      PushOp,
		K:       k,
		Skipped: false,
		CopyN:   100,
	}
	
	if kop.Op != PushOp {
		t.Errorf("Expected Op to be PushOp, got %v", kop.Op)
	}
	
	if kop.CopyN != 100 {
		t.Errorf("Expected CopyN to be 100, got %d", kop.CopyN)
	}
}

func TestOperationConstants(t *testing.T) {
	ops := []Op{PushOp, FetchOp, StageOp, IndexOp}
	
	for _, op := range ops {
		if string(op) == "" {
			t.Errorf("Operation %v should have non-empty string value", op)
		}
	}
}