package clwrapper

import (
	"fmt"
	"github.com/rainliu/gocl/cl"
)

func HasSupport() bool {

	var errNum cl.CL_int
	var numPlatforms cl.CL_uint

	errNum = cl.CLGetPlatformIDs(0, nil, &numPlatforms)
	if errNum == cl.CL_SUCCESS && numPlatforms > 0 {

		// while at it, setup the target platform
		// OK, this remains something TODO?
		// plids := make([]cl.CL_platform_id, numPlatforms)
		// if err := cl.CLGetPlatformIDs(numPlatforms, plids, nil); err == cl.CL_SUCCESS {

		//  // SHOW basic platform device id
		//  var paramValueSize cl.CL_size_t
		//  if errNum := cl.CLGetDeviceInfo(plids[0], cl.CL_PLATFORM_PROFILE, 0, nil, &paramValueSize); errNum == cl.CL_SUCCESS {
		//      var info interface{}
		//      cl.GetDeviceInfo(plids[0], cl.CL_PLATFORM_PROFILE, paramValueSize, &info, nil)
		//      fmt.Println(info)
		//  }

		//  DisplayPlatformInfo(plids[0], cl.CL_PLATFORM_VENDOR, "CL_PLATFORM_VENDOR")
		// }

		fmt.Printf("Found %d open cl platforms, using first one\n", numPlatforms)
		return true
	}

	return false
}
