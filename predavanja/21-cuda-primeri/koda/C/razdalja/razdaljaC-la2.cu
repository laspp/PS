// računanje razdalje med vektorjema
// 		argumenti: število blokov, število niti in dolžina vektorjev
// 		elementi vektorjev so inicializirani naključno
// na napravi izračunamo vsote kvadratov za vsak blok:
//		atomarno seštevanje - prekomerna uporaba in enotni pomnilnik

#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>
#include "cuda.h"
#include "helper_cuda.h"

__managed__ float sum;

__global__ void vectorDistance(const float *a, const float *b, int len) {
	// skupni pomnilnik niti v bloku
	extern __shared__ float part[];
	part[threadIdx.x] = 0.0;

	// kvadriranje razlike
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	float diff;
	while (gid < len) {
		diff = a[gid] - b[gid];
		atomicAdd(&sum, diff * diff);
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
	if (numBlocks < 0 || numThreads <= 0 || ceil(log2(numThreads)) != floor(log2(numThreads)) || vectorLength <= 0) {
		printf("usage:\n\t%s <number of blocks> <number of threads (power of 2)> <vector length>\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	// določimo potrebno število blokov niti
	if (numBlocks == 0)
		numBlocks = (vectorLength - 1) / numThreads + 1;

	// rezerviramo pomnilnik
	float *a, *b;
	checkCudaErrors(cudaMallocManaged((void**)&a, vectorLength * sizeof(float)));
	checkCudaErrors(cudaMallocManaged((void**)&b, vectorLength * sizeof(float)));

	// iniciliziramo vsoto
	sum = 0.0;
	// nastavimo vrednosti vektorjev a in b na gostitelju
	srand(time(NULL));
	for (int i = 0; i < vectorLength; i++) {
		a[i] = (float)rand()/RAND_MAX;
		b[i] = (float)rand()/RAND_MAX;
	}

	// merjenje časa na napravi - začetek
	struct timespec startDevice, stopDevice;
	clock_gettime(CLOCK_MONOTONIC, &startDevice);

	// zaženemo kodo na napravi
	dim3 gridSize(numBlocks, 1, 1);
	dim3 blockSize(numThreads, 1, 1);
	vectorDistance<<<gridSize, blockSize, numThreads*sizeof(float)>>>(a, b, vectorLength);
	checkCudaErrors(cudaGetLastError());

	// počakamo, da vse niti na napravi zaključijo
	checkCudaErrors(cudaDeviceSynchronize());

	// dokončamo izračun razdalje za napravo
	double distDevice = 0.0;
	distDevice = sqrt(sum);

	// merjenje časa na napravi - konec
	clock_gettime(CLOCK_MONOTONIC, &stopDevice);
	double timeDevice = (stopDevice.tv_sec - startDevice.tv_sec) * 1e3 + (stopDevice.tv_nsec - startDevice.tv_nsec) / 1e6;

	// izračunamo razdaljo in izmerimo čas na gostitelju
	struct timespec startHost, stopHost;
	clock_gettime(CLOCK_MONOTONIC, &startHost);
	double distHost = 0.0;
	float diff;
	for (int i = 0; i < vectorLength; i++) {
		diff = a[i] - b[i];
		distHost += diff * diff;
	}
	distHost = sqrt(distHost);
	clock_gettime(CLOCK_MONOTONIC, &stopHost);
	double timeHost = (stopHost.tv_sec - startHost.tv_sec) * 1e3 + (stopHost.tv_nsec - startHost.tv_nsec) / 1e6;

	// rezultata izpišemo
	printf("naprava:      %lf (%lf ms)\ngostitelj:    %lf (%lf ms)\nnapaka (rel): %e\n", distDevice, timeDevice, distHost, timeHost, fabs(distDevice/distHost-1));

	// sprostimo pomnilnik
	checkCudaErrors(cudaFree(a));
	checkCudaErrors(cudaFree(b));

	return 0;
}
