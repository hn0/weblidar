package clwrapper

import (
	"fmt"
	"testing"
)

func TestSupport(t *testing.T) {
	if !HasSupport() {
		t.Error("No supporting opencl hardware detected")
		t.FailNow()
	}
}

func TestMatvec(t *testing.T) {

	npts := 10

	// max distance 1.73
	valfnc := func(i int) (float32, float32, float32) {
		ret := -(float32(i) + 1) / float32(npts)
		return ret, ret, ret
	}
	resfnc := func(i int, xyz [3]float32, cat float32) {
		fmt.Printf("%d: x:%f y:%f z:%f  <-> c:%f\n", i, xyz[0], xyz[1], xyz[2], cat)
	}

	p := Program{"euclid_dist.cl", "euclid_dist", valfnc, resfnc}

	// this is weird but why invalid work item error is thrown (8192 is max sz, yep limitation of gpu! 256 * worksizeitem!)
	if !RunProgram(&p, npts) {
		t.Error("Running matvec example failed")
		t.Fail()
	}
	t.Error("UNDER DEVELOPMENT!")
	t.FailNow()
}
