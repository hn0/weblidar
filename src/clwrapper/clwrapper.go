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
	Val     func(int) (float32, float32, float32)
}

type cldev struct {
	dev    []cl.CL_device_id
	numdev int
}

var v *cldev
var WORK_UNIT_SZ int = 16

func RunProgram(p *Program, datasz int) bool {

	// basic set of checks
	var chk string
	switch {
	case v.numdev < 1:
		chk = "Cannot run cl program without any device"
	case datasz < WORK_UNIT_SZ:
		chk = "Data size cannot be smaller than work unit size"
	}

	if len(chk) > 0 {
		fmt.Println(chk)
		return false
	}

	wunit := datasz / WORK_UNIT_SZ
	if datasz%WORK_UNIT_SZ != 0 {
		wunit++
	}
	datax := make([]float32, wunit*WORK_UNIT_SZ)
	datay := make([]float32, wunit*WORK_UNIT_SZ)
	dataz := make([]float32, wunit*WORK_UNIT_SZ)
	res := make([]float32, wunit*WORK_UNIT_SZ)

	// var data [64]float32
	// var res [64]float32

	for i := 0; i < datasz; i++ {
		datax[i], datay[i], dataz[i] = p.Val(i)
		// res[i] = 1
	}
	for i := datasz - 1; i < wunit*WORK_UNIT_SZ; i++ {
		datax[i] = 0
		datay[i] = 0
		dataz[i] = 0
		// res[i] = 1
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

		var matx_buff, maty_buff, matz_buff, res_buff cl.CL_mem
		matx_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
			cl.CL_size_t(int(unsafe.Sizeof(datax[0]))*len(datax)), unsafe.Pointer(&datax[0]), &err)
		if assert(err, "Cannot create input buffer") {
			return false
		}
		maty_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
			cl.CL_size_t(int(unsafe.Sizeof(datax[0]))*len(datax)), unsafe.Pointer(&datay[0]), nil)
		matz_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
			cl.CL_size_t(int(unsafe.Sizeof(datax[0]))*len(datax)), unsafe.Pointer(&dataz[0]), &err)
		res_buff = cl.CLCreateBuffer(ctx, cl.CL_MEM_WRITE_ONLY, cl.CL_size_t(int(unsafe.Sizeof(res[0]))*len(res)), nil, &err)
		if assert(err, "Cannot create output buffer") {
			return false
		}

		kernel := cl.CLCreateKernel(program, []byte(p.FncName), &err)
		if assert(err, "Cannot create kernel") {
			return false
		}
		defer cl.CLReleaseKernel(kernel)

		err = cl.CLSetKernelArg(kernel, 0, cl.CL_size_t(unsafe.Sizeof(matx_buff)), unsafe.Pointer(&matx_buff))
		if assert(err, "Cannot set kernel args") {
			return false
		}
		cl.CLSetKernelArg(kernel, 1, cl.CL_size_t(unsafe.Sizeof(maty_buff)), unsafe.Pointer(&maty_buff))
		cl.CLSetKernelArg(kernel, 2, cl.CL_size_t(unsafe.Sizeof(matz_buff)), unsafe.Pointer(&matz_buff))
		cl.CLSetKernelArg(kernel, 3, cl.CL_size_t(unsafe.Sizeof(res_buff)), unsafe.Pointer(&res_buff))

		// queue
		queue := cl.CLCreateCommandQueue(ctx, v.dev[0], 0, &err)
		if assert(err, "Failed at creating queue") {
			return false
		}

		// var work_unit_per_kernel = [1]cl.CL_size_t{cl.CL_size_t(WORK_UNIT_SZ)}
		// res is 64 length
		dim := cl.CL_uint(1)
		var global_size = [...]cl.CL_size_t{cl.CL_size_t(wunit * WORK_UNIT_SZ)} // total number of workitems product >= len(inp)
		var local_size = [...]cl.CL_size_t{cl.CL_size_t(wunit)}                 // for two dimm its half of first dimm!!!
		err = cl.CLEnqueueNDRangeKernel(queue, kernel, dim, nil, global_size[:], local_size[:], 0, nil, nil)
		if assert(err, "Cannot create kernel queue") {
			return false
		}
		cl.CLEnqueueReadBuffer(queue, res_buff, cl.CL_TRUE, 0, cl.CL_size_t(int(unsafe.Sizeof(res[0]))*len(res)), unsafe.Pointer(&res[0]), 0, nil, nil)

		cl.CLReleaseKernel(kernel)
		cl.CLReleaseCommandQueue(queue)
		cl.CLReleaseMemObject(matx_buff)
		cl.CLReleaseMemObject(maty_buff)
		cl.CLReleaseMemObject(matz_buff)
		cl.CLReleaseMemObject(res_buff)

		fmt.Println(datax)
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
