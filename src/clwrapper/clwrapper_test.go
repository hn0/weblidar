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

	valfnc := func(i int) float32 {
		return float32(i)
	}
	p := Program{"euclid_dist.cl", "euclid_dist", valfnc}

	if !RunProgram(&p, 59) {
		t.Error("Running matvec example failed")
		t.Fail()
	}
	t.Error("UNDER DEVELOPMENT!")
	t.FailNow()
}
