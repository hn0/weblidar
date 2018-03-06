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

	// var data [16]float32
	// for i := 0; i < 16; i++ {
	// 	data[i] = float32(i + 2)
	// }

	a := func(i int) float32 {
		return float32(i + 3)
	}
	p := Program{"square_test.cl", "square_test", a}

	if !RunProgram(&p, 16) {
		t.Error("Running matvec example failed")
		t.Fail()
	}
	t.Error("UNDER DEVELOPMENT!")
	t.FailNow()
}
