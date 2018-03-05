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

	var data [16]float32
	for i := 0; i < 16; i++ {
		data[i] = float32(i)
	}

	p := Program{"square_test.cl", "square_test", data[:]}
	if !RunProgram(&p) {
		t.Error("Running matvec example failed")
		t.Fail()
	}
	t.Error("UNDER DEVELOPMENT!")
	t.FailNow()
}
