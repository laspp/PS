// izpis podatkov o zagnanih nitih
//
// izvajanje:
//		source ../cudago-init.sh
//		CudaGo -precompile -package cudago pozdrav-gpe.cu
//      srun --partition=gpu --gpus=1 go run pozdrav-gpe.go

package main

import (
	"flag"
	"pozdrav-gpe/cudago"

	"github.com/InternatBlackhole/cudago/cuda"
)

func main() {

	// preberemo argumente iz ukazne vrstice
	bPtr := flag.Int("b", 1, "num blocks")
	tPtr := flag.Int("t", 1, "num threads")
	flag.Parse()

	var err error

	// inicializiramo napravo
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	// zaženemo kodo na napravi
	gridSize := cuda.Dim3{X: uint32(*bPtr), Y: 1, Z: 1}
	blockSize := cuda.Dim3{X: uint32(*tPtr), Y: 1, Z: 1}
	err = cudago.Hello(gridSize, blockSize)
	if err != nil {
		panic(err)
	}
}
