package clwrapper

import (
	"fmt"
	"github.com/rainliu/gocl/cl"
	"os"
	"unsafe"
)

type clwrapper struct {
	// plat  cl.CL_platform_id
	dev     []cl.CL_device_id
	ctx     cl.CL_context
	program cl.CL_program
	lock    bool
	valid   bool
}

var v *clwrapper

// simple test matvec!
func MatVec() bool {
	// create context
	if get_context() && read_program("dotvect.cl") {
		// read program

		run_program([]byte("matvec_mult"))
		// fmt.Println("ready to continue")

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

func run_program(name []byte) {
	v.lock = true
	var err cl.CL_int
	kernel := cl.CLCreateKernel(v.program, name, &err)
	if err < 0 {
		fmt.Println("ERROR!", kernel, err)
		return
	}

	// now the memory struff
	var mat [16]float64
	for i := 0; i < 16; i++ {
		mat[i] = float64(i)
	}

	var vec, res [4]float32
	var vec_buff, res_buff cl.CL_mem

	mat_buff := cl.CLCreateBuffer(v.ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		cl.CL_size_t(unsafe.Sizeof(mat)), unsafe.Pointer(&mat[0]), &err)

	if err < 0 {
		fmt.Println("Error creating buffer")
		return
	}

	vec_buff = cl.CLCreateBuffer(v.ctx, cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		cl.CL_size_t(unsafe.Sizeof(vec)), unsafe.Pointer(&vec[0]), &err)
	res_buff = cl.CLCreateBuffer(v.ctx, cl.CL_MEM_WRITE_ONLY, cl.CL_size_t(unsafe.Sizeof(res)), unsafe.Pointer(&res[0]), &err)

	fmt.Println(mat_buff, vec_buff, res_buff)
	err = cl.CLSetKernelArg(kernel, 0, cl.CL_size_t(unsafe.Sizeof(mat_buff)), unsafe.Pointer(&mat_buff))
	if err < 0 {
		fmt.Println("Failed to set kernel args")
	}
	cl.CLSetKernelArg(kernel, 1, cl.CL_size_t(unsafe.Sizeof(vec_buff)), unsafe.Pointer(&vec_buff))
	cl.CLSetKernelArg(kernel, 2, cl.CL_size_t(unsafe.Sizeof(res_buff)), unsafe.Pointer(&res_buff))

	// at last, create command queue
	queue := cl.CLCreateCommandQueue(v.ctx, v.dev[0], 0, &err)
	if err < 0 {
		fmt.Println("Failed at creating queue")
	}

	var work_unit_per_kernel = [1]cl.CL_size_t{4}
	err = cl.CLEnqueueNDRangeKernel(queue, kernel, 1, nil, work_unit_per_kernel[:], nil, 0, nil, nil)
	if err < 0 {
		fmt.Println("Assertion needed!", err)
	}
	err = cl.CLEnqueueReadBuffer(queue, res_buff, cl.CL_TRUE, 0, cl.CL_size_t(unsafe.Sizeof(res)), unsafe.Pointer(&res[0]), 0, nil, nil)
	if err < 0 {
		fmt.Println("Assertion needed2!", err)
	}

	fmt.Println(res)

	// cleanup the stuff
	cl.CLReleaseMemObject(mat_buff)
	cl.CLReleaseMemObject(vec_buff)
	cl.CLReleaseMemObject(res_buff)
	cl.CLReleaseKernel(kernel)
	cl.CLReleaseCommandQueue(queue)
	cl.CLReleaseProgram(v.program)
	cl.CLReleaseContext(v.ctx)
}

func read_program(filename string) bool {
	if v.lock {
		return false
	}
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
			err = cl.CLBuildProgram(program, 1, v.dev[:], nil, nil, nil)
			if err < 0 {
				fmt.Println("Failed to compile provided program!")
				var log_size cl.CL_size_t
				var err_msg interface{}
				cl.CLGetProgramBuildInfo(program, v.dev[0], cl.CL_PROGRAM_BUILD_LOG, 0, nil, &log_size)
				cl.CLGetProgramBuildInfo(program, v.dev[0], cl.CL_PROGRAM_BUILD_LOG, log_size, &err_msg, nil)
				fmt.Printf("\tfail msg:\n%s\n", err_msg)
				return false
			}

			v.program = program
			return true
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
