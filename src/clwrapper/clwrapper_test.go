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

	npts := 12

	// max distance 1.73, given abs value of dmax is equal to a*sqrt(3) where a=1
	valfnc := func(i int) (float32, float32, float32) {
		pos := npts - i
		ret := float32(2*pos) / float32(2*npts)
		return ret, ret, ret
	}
	resfnc := func(i int, xyz [3]float32, cat [2]float32) {
		fmt.Printf("%d: x:%f y:%f z:%f  <-> d:%f a:%f\n", i, xyz[0], xyz[1], xyz[2], cat[0], cat[1])
	}

	p := Program{"euclid_dist.cl", "euclid_dist", valfnc, resfnc}

	// this is weird but why invalid work item error is thrown (8192 is max sz, yep limitation of gpu! 256 * worksizeitem!)
	if !RunProgram(&p, 2*npts) {
		t.Error("Running matvec example failed")
		t.Fail()
	}
	t.Error("UNDER DEVELOPMENT!")
	t.FailNow()
}
