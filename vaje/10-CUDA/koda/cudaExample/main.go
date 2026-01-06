package main

import (
	"cudaExample/helloCuda"
	"fmt"
	"unsafe"

	"github.com/InternatBlackhole/cudago/cuda"
)

func main() {
	//Initialize CUDA API for this OS thread
	var err error
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	fmt.Println("Cuda initialized")

	//Prepare message on host
	//String in go are not null terminated
	helloStrHost := []byte("Hello from the GPU!" + string(0))
	strLen := uint64(len(helloStrHost))

	//Allocate memory on the device
	helloStrDev, err := cuda.DeviceMemAlloc(strLen)
	defer helloStrDev.Free()

	//Copy message to the device
	//Use a pointer to the first element of the slice, and not the the slice itself
	err = helloStrDev.MemcpyToDevice(unsafe.Pointer(&helloStrHost[0]), strLen)
	if err != nil {
		panic(err)
	}

	//Specify grid and block size
	grid, block := cuda.Dim3{X: 1, Y: 1, Z: 1}, cuda.Dim3{X: 1, Y: 1, Z: 1}

	//Call the kernel to execute on the device
	err = helloCuda.Hello(grid, block, helloStrDev.Ptr)
	if err != nil {
		panic(err)
	}
}
