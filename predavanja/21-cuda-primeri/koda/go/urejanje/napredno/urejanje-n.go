// bitonično urejanje
// 		argumenti: število niti, velikost tabele
// več ščepcev, en ščepec lahko izvaja več notranjih zank
// prevajanje
//		source ../../cudago-init.sh
//		CudaGo -precompile -package cudago *.cu
// izvajanje
//		source ../../cudago-init.sh
// 		srun --partition=gpu --gpus=1 go run urejanje-n.go

package main

import (
	"flag"
	"fmt"
	"math"
	"time"
	"unsafe"
	"urejanje/cudago"

	"github.com/InternatBlackhole/cudago/cuda"
	"golang.org/x/exp/rand"
)

func main() {

	// preberemo argumente iz ukazne vrstice
	numThreadsPtr := flag.Int("t", 256, "num threads")
	tableSizePtr := flag.Int("s", 16777216, "table size")
	flag.Parse()
	if *numThreadsPtr <= 0 || *tableSizePtr <= 0 ||
		*tableSizePtr < *numThreadsPtr ||
		math.Ceil(math.Log2(float64(*tableSizePtr))) != math.Floor(math.Log2(float64(*tableSizePtr))) {
		panic("Wrong arguments")
	}

	// izračunamo potrebno število blokov
	numBlocks := (*tableSizePtr/2-1) / *numThreadsPtr + 1

	var err error

	// inicializiramo napravo
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// rezerviramo pomnilnik na gostitelju
	a := make([]int32, *tableSizePtr)
	ha := make([]int32, *tableSizePtr)

	// velikosti struktur v bajtih
	bytesInt := uint64(unsafe.Sizeof(int32(0)))
	bytesTable := uint64(*tableSizePtr) * bytesInt

	// rezerviramo pomnilnik na napravi
	da, err := cuda.DeviceMemAlloc(bytesTable)
	if err != nil {
		panic(err)
	}
	defer da.Free()

	// nastavimo vrednosti tabel a in ha na gostitelju
	for i := 0; i < *tableSizePtr; i++ {
		a[i] = rand.Int31()
		ha[i] = a[i]
	}

	// merjenje časa na napravi - začetek
	startDevice := time.Now()

	// prenesemo tabelo a iz gostitelja na napravo
	err = da.MemcpyToDevice(unsafe.Pointer(&ha[0]), bytesTable)
	if err != nil {
		panic(err)
	}

	// merjenje časa izvajanja ščepca na napravi - začetek
	startKernel := time.Now()

	// zaženemo kodo na napravi
	gridSize := cuda.Dim3{X: uint32(numBlocks), Y: 1, Z: 1}
	blockSize := cuda.Dim3{X: uint32(*numThreadsPtr), Y: 1, Z: 1}

	cudago.BitonicSortStart(gridSize, blockSize, da.Ptr, int32(*tableSizePtr)) // k = 2 ... 2 * blockSize.x
	for k := 4 * int32(blockSize.X); k <= int32(*tableSizePtr); k <<= 1 {      // k = 4 * blockSize ... tableLength
		for j := k / 2; j >= 2*int32(blockSize.X); j >>= 1 { //   j = k/2 ... 2 * blockSize.x
			err = cudago.BitonicSortMiddle(gridSize, blockSize, da.Ptr, int32(*tableSizePtr), k, j)
			if err != nil {
				panic(err)
			}
		}
		cudago.BitonicSortFinish(gridSize, blockSize, da.Ptr, int32(*tableSizePtr), k) //   j = 2 * blockSize.x ... 1
	}

	// počakamo, da se procesiranje zahtev na napravi zaključi
	err = cuda.CurrentContextSynchronize()
	if err != nil {
		panic(err)
	}

	// merjenje časa izvajanja ščepca na napravi - konec
	timeKernel := time.Since(startKernel)

	// prenesemo tabelo a iz naprave na gostitelja
	err = da.MemcpyFromDevice(unsafe.Pointer(&ha[0]), bytesTable)
	if err != nil {
		panic(err)
	}

	// merjenje časa na napravi - konec
	timeDevice := time.Since(startDevice)

	// bitonično rurejanje na gostitelju
	startHost := time.Now()
	for k := 2; k <= *tableSizePtr; k <<= 1 {
		for j := k / 2; j > 0; j >>= 1 {
			for i1 := 0; i1 < *tableSizePtr; i1++ {
				i2 := i1 ^ j
				dec := i1 & k
				if i2 > i1 {
					if (dec == 0 && a[i1] > a[i2]) || (dec != 0 && a[i1] < a[i2]) {
						temp := a[i1]
						a[i1] = a[i2]
						a[i2] = temp
					}
				}
			}
		}
	}
	timeHost := time.Since(startHost)

	// preverimo rezultat
	okDevice, okHost := true, true
	for i := 1; i < *tableSizePtr; i++ {
		okDevice = okDevice && (ha[i-1] <= ha[i])
		okHost = okHost && (a[i-1] <= a[i])
	}

	// izpisi
	fmt.Printf("Device: correct: %v, time: %v ms/%v us\n", okDevice, timeDevice.Milliseconds(), timeKernel.Microseconds())
	fmt.Printf("Host  : correct: %v, time: %v ms\n", okHost, timeHost.Milliseconds())
}
