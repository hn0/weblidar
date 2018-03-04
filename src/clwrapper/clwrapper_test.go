package clwrapper

import (
	"testing"
)

func TestSupport(t *testing.T) {
	if !HasSupport() {
		t.Error("No supporting opencl hardware detected")
		t.FailNow()
	}
}

func TestMatvec(t *testing.T) {
	if !MatVec() {
		t.Error("UNDER DEVELOPMENT!")
		t.FailNow()
	}
}
