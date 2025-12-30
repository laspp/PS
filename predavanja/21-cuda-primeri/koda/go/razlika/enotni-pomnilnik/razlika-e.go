// računanje razlike vektorjev
// 		argumenti: število blokov, število niti, dolžina vektorjev in oznaka ščepca
//		število blokov lahko nastavimo ali pa jih izračunamo (-b 0) glede na velikost vektorja in število niti
//		enotni pomnilnik: ogrodje CUDA poskrbi za prenos podatkov
//		koda na napravi: VectorSubtract4 je enaka kot pri rešitvi z ločenim pomnilnikom
// prevajanje ščepcev:
//		source ../../cudago-init.sh
//		CudaGo -precompile -package cudago *.cu
// izvajanje:
//		source ../../cudago-init.sh
// 		VectorSubtract4: srun --partition=gpu --gpus=1 go run razlika-e.go -b 0 -t 128 -s 200

package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"razlika/cudago"
	"unsafe"

	"github.com/InternatBlackhole/cudago/cuda"
)

func main() {

	// preberemo argumente iz ukazne vrstice
	numBlocksPtr := flag.Int("b", 1, "num blocks")
	numThreadsPtr := flag.Int("t", 1, "num threads")
	vectorSizePtr := flag.Int("s", 1, "vector size")
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

	// rezerviramo pomnilnik
	bytesFloat32 := uint64(unsafe.Sizeof(float32(0.0)))
	c, err := cuda.ManagedMemAlloc[float32](uint64(*vectorSizePtr), bytesFloat32)
	if err != nil {
		panic(err)
	}
	defer c.Free()
	a, err := cuda.ManagedMemAlloc[float32](uint64(*vectorSizePtr), bytesFloat32)
	if err != nil {
		panic(err)
	}
	defer a.Free()
	b, err := cuda.ManagedMemAlloc[float32](uint64(*vectorSizePtr), bytesFloat32)
	if err != nil {
		panic(err)
	}
	defer b.Free()

	// nastavimo vrednosti vektorjev a in b na gostitelju
	for i := 0; i < *vectorSizePtr; i++ {
		a.Arr[i] = rand.Float32()
		b.Arr[i] = rand.Float32()
	}

	// zaženemo kodo na napravi
	gridSize := cuda.Dim3{X: uint32(numBlocks), Y: 1, Z: 1}
	blockSize := cuda.Dim3{X: uint32(*numThreadsPtr), Y: 1, Z: 1}
	switch *kernelPtr {
	default:
		err = cudago.VectorSubtract4(gridSize, blockSize, c.Ptr, a.Ptr, b.Ptr, int32(*vectorSizePtr))
	}
	if err != nil {
		panic(err)
	}

	// počakamo, da se procesiranje zahtev na napravi zaključi
	err = cuda.CurrentContextSynchronize()
	if err != nil {
		panic(err)
	}

	// preverimo rezultat
	ok := true
	for i := 0; i < *vectorSizePtr; i++ {
		if a.Arr[i]-b.Arr[i] != c.Arr[i] {
			ok = false
		}
	}
	if ok {
		fmt.Println("Result is correct.")
	} else {
		fmt.Println("Result is wrong.")
	}
}
