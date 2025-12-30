// računanje razlike vektorjev
// 		argumenti: število blokov, število niti, dolžina vektorjev in oznaka ščepca
//		število blokov lahko nastavimo ročno ali pa jih izračunamo (-b 0) glede na velikost vektorja in število niti
//		ločeni pomnilnik: sami poskrbimo za prenos podatkov
//		koda na napravi: izboljševanje od VectorSubtract1 do VectorSubtract4 (vrstica 97)
// prevajanje ščepcev:
//		source ../../cudago-init.sh
//		CudaGo -precompile -package cudago *.cu
// izvajanje:
//		source ../../cudago-init.sh
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 1 -b 1 -t 128 -s 128
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 1 -b 1 -t 128 -s 100
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 2 -b 1 -t 128 -s 100
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 2 -b 1 -t 128 -s 200
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 2 -b 2 -t 128 -s 200
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 3 -b 2 -t 128 -s 200
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 3 -b 1 -t 128 -s 200
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 4 -b 1 -t 128 -s 200
// 		srun --partition=gpu --gpus=1 go run razlika-l.go -k 4 -b 0 -t 128 -s 200

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

	// rezerviramo pomnilnik na gostitelju
	hc := make([]float32, *vectorSizePtr)
	ha := make([]float32, *vectorSizePtr)
	hb := make([]float32, *vectorSizePtr)

	// rezerviramo pomnilnik na napravi
	bytesFloat32 := int(unsafe.Sizeof(float32(0.0)))
	bytesVector := uint64(*vectorSizePtr * bytesFloat32)

	dc, err := cuda.DeviceMemAlloc(bytesVector)
	if err != nil {
		panic(err)
	}
	defer dc.Free()
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

	// prenesemo vektorja a in b iz gostitelja na napravo
	err = da.MemcpyToDevice(unsafe.Pointer(&ha[0]), bytesVector)
	if err != nil {
		panic(err)
	}
	err = db.MemcpyToDevice(unsafe.Pointer(&hb[0]), bytesVector)
	if err != nil {
		panic(err)
	}

	// zaženemo kodo na napravi
	gridSize := cuda.Dim3{X: uint32(numBlocks), Y: 1, Z: 1}
	blockSize := cuda.Dim3{X: uint32(*numThreadsPtr), Y: 1, Z: 1}
	switch *kernelPtr {
	case 1:
		err = cudago.VectorSubtract1(gridSize, blockSize, dc.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	case 2:
		err = cudago.VectorSubtract2(gridSize, blockSize, dc.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	case 3:
		err = cudago.VectorSubtract3(gridSize, blockSize, dc.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	default:
		err = cudago.VectorSubtract4(gridSize, blockSize, dc.Ptr, da.Ptr, db.Ptr, int32(*vectorSizePtr))
	}
	if err != nil {
		panic(err)
	}

	// vektor c prekopiramo iz naprave na gostitelja
	err = dc.MemcpyFromDevice(unsafe.Pointer(&hc[0]), bytesVector)

	// počakamo, da se procesiranje zahtev na napravi zaključi
	err = cuda.CurrentContextSynchronize()
	if err != nil {
		panic(err)
	}

	// preverimo rezultat
	ok := true
	for i := 0; i < *vectorSizePtr; i++ {
		if ha[i]-hb[i] != hc[i] {
			ok = false
		}
	}
	if ok {
		fmt.Println("Result is correct.")
	} else {
		fmt.Println("Result is wrong.")
	}
}
