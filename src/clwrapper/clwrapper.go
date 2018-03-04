package clwrapper

import (
	"fmt"
	"github.com/rainliu/gocl/cl"
	"os"
)

type clwrapper struct {
	// plat  cl.CL_platform_id
	dev   []cl.CL_device_id
	ctx   cl.CL_context
	lock  bool
	valid bool
}

var v *clwrapper

// simple test matvec!
func MatVec() bool {
	// create context
	if get_context() {
		// read program
		if read_program("dotvect.cl") {

			fmt.Println("ready to continue")
		}
	}
	return false
}

func HasSupport() bool {

	v = new(clwrapper)
	var errNum cl.CL_int
	var numPlatforms cl.CL_uint

	errNum = cl.CLGetPlatformIDs(0, nil, &numPlatforms)
	if errNum == cl.CL_SUCCESS && numPlatforms > 0 {

		plids := make([]cl.CL_platform_id, numPlatforms)
		if err := cl.CLGetPlatformIDs(1, plids, nil); err == cl.CL_SUCCESS && len(plids) > 0 {
			// access the device!
			devices := make([]cl.CL_device_id, 1)
			if err := cl.CLGetDeviceIDs(plids[0], cl.CL_DEVICE_TYPE_GPU, 1, devices, nil); err == cl.CL_SUCCESS && len(devices) > 0 {
				// name of device!
				var name string
				var paramValueSize cl.CL_size_t
				if err := cl.CLGetDeviceInfo(devices[0], cl.CL_DEVICE_NAME, 0, nil, &paramValueSize); err == cl.CL_SUCCESS {
					var info interface{}
					if err := cl.CLGetDeviceInfo(devices[0], cl.CL_DEVICE_NAME, paramValueSize, &info, nil); err == cl.CL_SUCCESS {
						name = info.(string)
					}
				}
				fmt.Printf("Found %d open cl platforms, using first one: %s\n", numPlatforms, name)
				// v.plat = plids[0]
				v.dev = devices
				v.valid = true
				return true
			}
		}
	}

	return false
}

func read_program(filename string) bool {
	if fp, err := os.Open(filename); err == nil {
		defer fp.Close()

		var program_size [1]cl.CL_size_t
		var program_buffer [1][]byte

		fi, _ := fp.Stat()
		program_size[0] = cl.CL_size_t(fi.Size())
		program_buffer[0] = make([]byte, program_size[0])
		fp.Read(program_buffer[0])

		var err cl.CL_int
		program := cl.CLCreateProgramWithSource(v.ctx, 1, program_buffer[:], program_size[:], &err)
		if err >= 0 {
			fmt.Println(program)
		}
	}
	return false
}

func get_context() bool {
	if v.valid {
		var err cl.CL_int
		v.ctx = cl.CLCreateContext(nil, 1, v.dev[:], nil, nil, &err)
		return !(err < 0)
	}
	return false
}
