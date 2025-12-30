// računanje razlike vektorjev
// 		argumenti: število blokov, število niti in dolžina vektorjev
// 		elementi vektorjev so inicializirani naključno
// uporaba enotnega pomnilnika

#include <stdio.h>
#include <stdlib.h>
#include "cuda.h"
#include "helper_cuda.h"

__global__ void vectorSubtract(float *c, const float *a, const float *b, int len) {
	// določimo globalni indeks elementov
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	// če je niti manj kot je dolžina vektorjev, morajo nekatere narediti več elementov
	while (gid < len) {
		c[gid] = a[gid] - b[gid];
		gid += gridDim.x * blockDim.x;
	}
}

int main(int argc, char **argv) {
	// preberemo argumente iz ukazne vrstice
	int numBlocks = 0;
	int numThreads = 0;
	int vectorLength = 0;
	if (argc == 4) {
		numBlocks = atoi(argv[1]);
		numThreads = atoi(argv[2]);
		vectorLength = atoi(argv[3]);
	}
	if (numBlocks < 0 || numThreads <= 0 || vectorLength <= 0) {
		printf("usage:\n\t%s <number of blocks> <number of threads> <vector length>\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	// določimo potrebno število blokov niti
	if (numBlocks == 0)
		numBlocks = (vectorLength - 1) / numThreads + 1;

	// rezerviramo pomnilnik
	float *c, *a, *b;
	checkCudaErrors(cudaMallocManaged((void**)&c, vectorLength * sizeof(float)));
	checkCudaErrors(cudaMallocManaged((void**)&a, vectorLength * sizeof(float)));
	checkCudaErrors(cudaMallocManaged((void**)&b, vectorLength * sizeof(float)));

	// nastavimo vrednosti vektorjev a in b na gostitelju
	srand(time(NULL));
	for (int i = 0; i < vectorLength; i++) {
		a[i] = (float)rand()/RAND_MAX;
		b[i] = (float)rand()/RAND_MAX;
	}

	// zaženemo kodo na napravi
	dim3 gridSize(numBlocks, 1, 1);
	dim3 blockSize(numThreads, 1, 1);
	vectorSubtract<<<gridSize, blockSize>>>(c, a, b, vectorLength);
	checkCudaErrors(cudaGetLastError());

	// počakamo, da vse niti na napravi zaključijo
	checkCudaErrors(cudaDeviceSynchronize());

	// preverimo rezultat
	int ok = 1;
	for (int i = 0; i < vectorLength; i++)
		ok &= (a[i] - b[i]) == c[i];
	printf("Result is %s.\n", ok == 1 ? "correct": "wrong");

	// sprostimo pomnilnik
	checkCudaErrors(cudaFree(c));
	checkCudaErrors(cudaFree(a));
	checkCudaErrors(cudaFree(b));

	return 0;
}
