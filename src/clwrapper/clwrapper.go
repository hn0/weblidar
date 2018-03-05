package clwrapper

import (
	"fmt"
	"github.com/rainliu/gocl/cl"
	"os"
	"unsafe"
)

type Program struct {
	Source  string
	FncName string
}

type cldev struct {
	dev    []cl.CL_device_id
	numdev int
}

var v *cldev

func RunProgram(p *Program, data [16]float32) bool {

	if v.numdev < 1 {
		fmt.Println("Cannot run cl program without any device")
		return false
	}

	var err cl.CL_int
	ctx := cl.CLCreateContext(nil, 1, v.dev, nil, nil, &err)
	if assert(err, "Cannot create ctx") {
		return false
	}
	defer cl.CLReleaseContext(ctx)

	psz, pbuff := read_program(p.Source)
	if psz[0] > 0 {
		program := cl.CLCreateProgramWithSource(ctx, 1, pbuff[:], psz[:], &err)
		if assert(err, "Cannot create program from source") {
			return false
		}
		defer cl.CLReleaseProgram(program)

		// build program
		err = cl.CLBuildProgram(program, 1, v.dev[:], nil, nil, nil)
		if assert(err, "Cannot build the program") {
			var log_size cl.CL_size_t
			var err_msg interface{}
			cl.CLGetProgramBuildInfo(program, v.dev[0], cl.CL_PROGRAM_BUILD_LOG, 0, nil, &log_size)
			cl.CLGetProgramBuildInfo(program, v.dev[0], cl.CL_PROGRAM_BUILD_LOG, log_size, &err_msg, nil)
			fmt.Printf("\tfail msg:\n%s\n", err_msg)
			return false
		}

		var mat_buff, res_buff cl.CL_mem
		mat_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
			cl.CL_size_t(unsafe.Sizeof(data)), unsafe.Pointer(&data[0]), &err)
		if assert(err, "Cannot create input buffer") {
			return false
		}
		var res [16]float32
		res_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_WRITE_ONLY, cl.CL_size_t(unsafe.Sizeof(res)), nil, nil)

		// returns -45
		kernel := cl.CLCreateKernel(program, []byte(p.FncName), &err)
		if assert(err, "Cannot create kernel") {
			return false
		}
		defer cl.CLReleaseKernel(kernel)

		err = cl.CLSetKernelArg(kernel, 0, cl.CL_size_t(unsafe.Sizeof(mat_buff)), unsafe.Pointer(&mat_buff))
		if assert(err, "Cannot set kernel args") {
			return false
		}
		cl.CLSetKernelArg(kernel, 1, cl.CL_size_t(unsafe.Sizeof(res_buff)), unsafe.Pointer(&res_buff))

		// queue
		queue := cl.CLCreateCommandQueue(ctx, v.dev[0], 0, &err)
		if assert(err, "Failed at creating queue") {
			return false
		}

		var work_unit_per_kernel = [2]cl.CL_size_t{8} // split!?
		err = cl.CLEnqueueNDRangeKernel(queue, kernel, 1, nil, work_unit_per_kernel[:], nil, 0, nil, nil)
		if assert(err, "cannot create kernel queue") {
			return false
		}
		cl.CLEnqueueReadBuffer(queue, res_buff, cl.CL_TRUE, 0, cl.CL_size_t(unsafe.Sizeof(res)), unsafe.Pointer(&res[0]), 0, nil, nil)

		cl.CLReleaseKernel(kernel)
		cl.CLReleaseCommandQueue(queue)
		cl.CLReleaseMemObject(mat_buff)
		cl.CLReleaseMemObject(res_buff)

		fmt.Println(data)
		fmt.Println(res)

	}
	return false
}

func read_program(fname string) ([1]cl.CL_size_t, [1][]byte) {
	var program_size [1]cl.CL_size_t
	var program_buffer [1][]byte

	fp, err := os.Open(fname)
	if err == nil {
		defer fp.Close()

		fi, err := fp.Stat()
		if err == nil {
			program_size[0] = cl.CL_size_t(fi.Size())
			program_buffer[0] = make([]byte, program_size[0])
			fp.Read(program_buffer[0])
		} else {
			fmt.Println("Cannot read cl program")
		}
	} else {
		fmt.Println("Cannot open cl program source")
	}

	return program_size, program_buffer
}

func assert(err cl.CL_int, msg string) bool {
	if err < 0 {
		fmt.Printf("Assertion error: no:%d msg:%s\n", err, msg)
		return true
	}
	return false
}

func HasSupport() bool {

	v = new(cldev)
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
				v.numdev = len(devices)
				return true
			}
		}
	}

	return false
}
