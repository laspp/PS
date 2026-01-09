// računanje razlike vektorjev
// 		argumenti: število blokov, število niti, dolžina vektorjev, oznaka ščepca
// uporabimo lokalni pomnilnik na napravi, ki ga rezerviramo statično
// delne vsote seštevamo z atomarnim ukazom
// dve rešitvi: VectorDistanceLA1Ex ... VectorDistanceLA2Ex
// prevajanje:
//		source ../../cudago-init.sh
//		CudaGo -precompile -package cudago *.cu
// izvajanje:
//		source ../../cudago-init.sh
// 		srun --partition=gpu --gpus=1 go run razdalja-la.go

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"razdalja/cudago"
	"time"
	"unsafe"

	"github.com/InternatBlackhole/cudago/cuda"
)

func main() {

	// preberemo argumente iz ukazne vrstice
	numBlocksPtr := flag.Int("b", 0, "num blocks")
	numThreadsPtr := flag.Int("t", 1024, "num threads")
	vectorSizePtr := flag.Int("s", 268435456, "vector size")
	kernelPtr := flag.Int("k", 0, "kernel")

	flag.Parse()
	if *numBlocksPtr < 0 || *numThreadsPtr <= 0 || *vectorSizePtr <= 0 {
		panic("Wrong arguments")
	}

	// izračunamo potrebno število blokov
	numBlocks := *numBlocksPtr
	if numBlocks == 0 {
		numBlocks = (*vectorSizePtr-1) / *numThreadsPtr + 1
	}

	var err error

	// inicializiramo napravo
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// rezerviramo pomnilnik na gostitelju
	hs := make([]float32, 1)
	ha := make([]float32, *vectorSizePtr)
	hb := make([]float32, *vectorSizePtr)

	// velikosti struktur v bajtih
	bytesFloat32 := int(unsafe.Sizeof(float32(0.0)))
	bytesVector := uint64(*vectorSizePtr * bytesFloat32)
	bytesPartialSum := uint64(1 * bytesFloat32)
	bytesLocalMemory := uint64(*numThreadsPtr * bytesFloat32)

	// rezerviramo pomnilnik na napravi
	ds, err := cuda.DeviceMemAlloc(bytesPartialSum)
	if err != nil {
		panic(err)
	}
	defer ds.Free()
	da, err := cuda.DeviceMemAlloc(bytesVector)
	if err != nil {
		panic(err)
	}
	defer da.Free()
	db, err := cuda.DeviceMemAlloc(bytesVector)
	if err != nil {
		panic(err)
	}
	defer db.Free()

	// nastavimo vrednosti vektorjev a in b na gostitelju
	for i := 0; i < *vectorSizePtr; i++ {
		ha[i] = rand.Float32()
		hb[i] = rand.Float32()
	}

	// merjenje časa na napravi - začetek
	startDevice := time.Now()

	// prenesemo vektorja a in b iz gostitelja na napravo
	err = da.MemcpyToDevice(unsafe.Pointer(&ha[0]), bytesVector)
	if err != nil {
		panic(err)
	}
	err = db.MemcpyToDevice(unsafe.Pointer(&hb[0]), bytesVector)
	if err != nil {
		panic(err)
	}

	// merjenje časa izvajanja ščepca na napravi - začetek
	startKernel := time.Now()

	// zaženemo kodo na napravi
	gridSize := cuda.Dim3{X: uint32(numBlocks), Y: 1, Z: 1}
	blockSize := cuda.Dim3{X: uint32(*numThreadsPtr), Y: 1, Z: 1}
	switch *kernelPtr {
	case 2:
		err = cudago.VectorDistanceLA2Ex(gridSize, blockSize, bytesLocalMemory, nil, ds.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	default:
		err = cudago.VectorDistanceLA1Ex(gridSize, blockSize, bytesLocalMemory, nil, ds.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	}
	if err != nil {
		panic(err)
	}

	// počakamo, da se procesiranje zahtev na napravi zaključi
	err = cuda.CurrentContextSynchronize()
	if err != nil {
		panic(err)
	}

	// merjenje časa izvajanja ščepca na napravi - konec
	timeKernel := time.Since(startKernel)

	// skalar s prekopiramo iz naprave na gostitelja
	err = ds.MemcpyFromDevice(unsafe.Pointer(&hs[0]), bytesPartialSum)

	// dokončamo izračun razdalje za napravo
	distDevice := math.Sqrt(float64(hs[0]))

	// merjenje časa na napravi - konec
	timeDevice := time.Since(startDevice)

	// preverimo rezultat
	startHost := time.Now()
	distHost := float64(0.0)
	var diff float32
	for i := 0; i < *vectorSizePtr; i++ {
		diff = ha[i] - hb[i]
		distHost += float64(diff * diff)
	}
	distHost = math.Sqrt(distHost)
	timeHost := time.Since(startHost)

	// izpišemo rezultate
	fmt.Printf("naprava:      %f (%v ms/%v us)\n", distDevice, timeDevice.Milliseconds(), timeKernel.Microseconds())
	fmt.Printf("gostitelj:    %f (%v ms)\n", distHost, timeHost.Milliseconds())
	fmt.Printf("napaka (rel): %e\n", math.Abs(distDevice/distHost-1))
}
