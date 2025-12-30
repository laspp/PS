// informacije o napravi
//
// izvajanje:
//		source cudago-init.sh
//      srun --partition=gpu --gpus=1 go run naprava.go

package main

import (
	"fmt"
	"math"

	"github.com/InternatBlackhole/cudago/cuda"
)

func main() {

	var err error

	// inicializiramo napravo
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// poizvedbe
	name, err := dev.Device.Name()
	fmt.Printf("\nDevice %d: %s\n", dev.DeviceIndex, name)

	version_major, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MAJOR)
	version_minor, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MINOR)
	fmt.Printf("  CUDA architecture:                            %v, %v.%v\n", ConvertSMVer2ArchName(version_major, version_minor), version_major, version_minor)

	fmt.Println()

	gpu_clock_rate, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_CLOCK_RATE)
	fmt.Println("  GPU clock rate (MHz):                        ", gpu_clock_rate/1e3)

	mem_clock_rate, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MEMORY_CLOCK_RATE)
	fmt.Println("  Memory clock rate (MHz):                     ", mem_clock_rate/1e3)

	mem_bus_width, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_GLOBAL_MEMORY_BUS_WIDTH)
	fmt.Println("  Memory bus width (bits):                     ", mem_bus_width)
	fmt.Println("  Peak memory bandwidth (GB/s):                ", math.Round(2.0*float64(mem_clock_rate)*(float64(mem_bus_width)/8)/1.0e6))

	fmt.Println()

	mp_count, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MULTIPROCESSOR_COUNT)
	fmt.Println("  Number of MPs:                               ", mp_count)
	fmt.Println("  Number of cores per MP:                      ", ConvertSMVer2Cores(version_major, version_minor))
	fmt.Println("  Total number of cores:                       ", mp_count*ConvertSMVer2Cores(version_major, version_minor))

	fmt.Println()

	mem_global, err := dev.Device.TotalMem()
	fmt.Println("  Total amount of global memory (GB):          ", math.Round(float64(mem_global)/1073741824.0))
	l2cache, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_L2_CACHE_SIZE)
	fmt.Println("  Size of L2 cache (MB):                       ", l2cache/1024/1024)
	mem_shared_mp, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_SHARED_MEMORY_PER_MULTIPROCESSOR)
	fmt.Println("  Total amount of shared memory per MP (kB):   ", mem_shared_mp/1024)
	mem_shared_block, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_SHARED_MEMORY_PER_BLOCK)
	fmt.Println("  Total amount of shared memory per block (kB):", mem_shared_block/1024)
	reg_mp, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_REGISTERS_PER_MULTIPROCESSOR)
	fmt.Println("  Maximum number of registers per MP:          ", reg_mp)
	reg_block, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_REGISTERS_PER_BLOCK)
	fmt.Println("  Maximum number of registers per block:       ", reg_block)

	fmt.Println()

	threads_mp, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_THREADS_PER_MULTIPROCESSOR)
	fmt.Println("  Maximum number of threads per MP:            ", threads_mp)
	threads_block, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_THREADS_PER_BLOCK)
	fmt.Println("  Maximum number of threads per block:         ", threads_block)
	blocks_mp, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_BLOCKS_PER_MULTIPROCESSOR)
	fmt.Println("  Max number of blocks per MP:                 ", blocks_mp)
	threads_warp, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_WARP_SIZE)
	fmt.Println("  Warp size:                                   ", threads_warp)

	fmt.Println()

	max_block_dim_x, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_BLOCK_DIM_X)
	max_block_dim_y, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_BLOCK_DIM_Y)
	max_block_dim_z, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_BLOCK_DIM_Z)
	fmt.Printf("  Max dimension size of a thread block (x,y,z): (%v,%v,%v)\n", max_block_dim_x, max_block_dim_y, max_block_dim_z)
	max_grid_dim_x, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_GRID_DIM_X)
	max_grid_dim_y, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_GRID_DIM_Y)
	max_grid_dim_z, err := dev.Device.GetAttribute(cuda.CU_DEVICE_ATTRIBUTE_MAX_GRID_DIM_Z)
	fmt.Printf("  Max dimension size of a grid (x,y,z):         (%v,%v,%v)\n", max_grid_dim_x, max_grid_dim_y, max_grid_dim_z)
}

func ConvertSMVer2ArchName(major int, minor int) string {
	// Defines for GPU Architecture types (using the SM version to determine
	// the GPU Arch name)

	type sSMtoArchName struct {
		SM   int // 0xMm (hexadecimal notation), M = SM Major version, and m = SM minor version
		name string
	}

	var nGpuArchNameSM = []sSMtoArchName{
		sSMtoArchName{SM: 0x30, name: "Kepler"},
		sSMtoArchName{SM: 0x32, name: "Kepler"},
		sSMtoArchName{SM: 0x35, name: "Kepler"},
		sSMtoArchName{SM: 0x37, name: "Kepler"},
		sSMtoArchName{SM: 0x50, name: "Maxwell"},
		sSMtoArchName{SM: 0x52, name: "Maxwell"},
		sSMtoArchName{SM: 0x53, name: "Maxwell"},
		sSMtoArchName{SM: 0x60, name: "Pascal"},
		sSMtoArchName{SM: 0x61, name: "Pascal"},
		sSMtoArchName{SM: 0x62, name: "Pascal"},
		sSMtoArchName{SM: 0x70, name: "Volta"},
		sSMtoArchName{SM: 0x72, name: "Xavier"},
		sSMtoArchName{SM: 0x75, name: "Turing"},
		sSMtoArchName{SM: 0x80, name: "Ampere"},
		sSMtoArchName{SM: 0x86, name: "Ampere"},
		sSMtoArchName{SM: 0x89, name: "Ada"},
		sSMtoArchName{SM: 0x90, name: "Hopper"},
		sSMtoArchName{SM: -1, name: "Graphics Device"},
	}

	index := 0
	for nGpuArchNameSM[index].SM != -1 {
		if nGpuArchNameSM[index].SM == ((major << 4) + minor) {
			return nGpuArchNameSM[index].name
		}
		index++
	}
	// If we don't find the values, we default use the previous one to run properly
	fmt.Printf("MapSMtoArchName for SM %v.%v is undefined. Default to use %s\n", major, minor, nGpuArchNameSM[index-1].name)
	return nGpuArchNameSM[index-1].name
}

// Beginning of GPU Architecture definitions
func ConvertSMVer2Cores(major int, minor int) int {
	// Defines for GPU Architecture types (using the SM version to determine
	// the # of cores per SM
	type sSMtoCores struct {
		SM    int // 0xMm (hexadecimal notation), M = SM Major version, and m = SM minor version
		Cores int
	}

	var nGpuArchCoresPerSM = []sSMtoCores{
		sSMtoCores{SM: 0x30, Cores: 192},
		sSMtoCores{SM: 0x32, Cores: 192},
		sSMtoCores{SM: 0x35, Cores: 192},
		sSMtoCores{SM: 0x37, Cores: 192},
		sSMtoCores{SM: 0x50, Cores: 128},
		sSMtoCores{SM: 0x52, Cores: 128},
		sSMtoCores{SM: 0x53, Cores: 128},
		sSMtoCores{SM: 0x60, Cores: 64},
		sSMtoCores{SM: 0x61, Cores: 128},
		sSMtoCores{SM: 0x62, Cores: 128},
		sSMtoCores{SM: 0x70, Cores: 64},
		sSMtoCores{SM: 0x72, Cores: 64},
		sSMtoCores{SM: 0x75, Cores: 64},
		sSMtoCores{SM: 0x80, Cores: 64},
		sSMtoCores{SM: 0x86, Cores: 128},
		sSMtoCores{SM: 0x87, Cores: 128},
		sSMtoCores{SM: 0x89, Cores: 128},
		sSMtoCores{SM: 0x90, Cores: 128},
		sSMtoCores{SM: -1, Cores: -1},
	}

	index := 0
	for nGpuArchCoresPerSM[index].SM != -1 {
		if nGpuArchCoresPerSM[index].SM == ((major << 4) + minor) {
			return nGpuArchCoresPerSM[index].Cores
		}
		index++
	}
	// If we don't find the values, we default use the previous one to run properly
	fmt.Printf("  MapSMtoCores for SM %v.%v is undefined. Default to use %v Cores/SM\n", major, minor, nGpuArchCoresPerSM[index-1].Cores)
	return nGpuArchCoresPerSM[index-1].Cores
}
