package clwrapper

import (
	"fmt"
	"github.com/samuel/go-opencl/cl"
)

func HasSupport() bool {

	platforms, _ := cl.GetPlatforms()
	fmt.Println(platforms)

	return false
}
