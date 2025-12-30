// računanje razdalje med vektorjema
// 		argumenti: število blokov, število niti in dolžina vektorjev
// 		elementi vektorjev so inicializirani naključno
// kvadarate izračunamo na napravi

#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>
#include "cuda.h"
#include "helper_cuda.h"

__global__ void vectorDistance(float *c, const float *a, const float *b, int len) {
	// kvadriranje razlike
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	float diff;
	while (gid < len) {
		diff = a[gid] - b[gid];
		c[gid] = diff * diff;
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

	// rezerviramo pomnilnik na gostitelju
	float *hc = (float *)malloc(vectorLength * sizeof(float));
	float *ha = (float *)malloc(vectorLength * sizeof(float));
	float *hb = (float *)malloc(vectorLength * sizeof(float));

	// rezerviramo pomnilnik na napravi
	float *dc, *da, *db;
	checkCudaErrors(cudaMalloc((void**)&dc, vectorLength * sizeof(float)));
	checkCudaErrors(cudaMalloc((void**)&da, vectorLength * sizeof(float)));
	checkCudaErrors(cudaMalloc((void**)&db, vectorLength * sizeof(float)));

	// nastavimo vrednosti vektorjev a in b na gostitelju
	srand(time(NULL));
	for (int i = 0; i < vectorLength; i++) {
		ha[i] = (float)rand()/RAND_MAX;
		hb[i] = (float)rand()/RAND_MAX;
	}

	// merjenje časa na napravi - začetek
	struct timespec startDevice, stopDevice;
	clock_gettime(CLOCK_MONOTONIC, &startDevice);

	// prenesemo vektorja a in b iz gostitelja na napravo
	checkCudaErrors(cudaMemcpy(da, ha, vectorLength * sizeof(float), cudaMemcpyHostToDevice));
	checkCudaErrors(cudaMemcpy(db, hb, vectorLength * sizeof(float), cudaMemcpyHostToDevice));

	// zaženemo kodo na napravi
	dim3 gridSize(numBlocks, 1, 1);
	dim3 blockSize(numThreads, 1, 1);
	vectorDistance<<<gridSize, blockSize>>>(dc, da, db, vectorLength);
	checkCudaErrors(cudaGetLastError());

	// počakamo, da vse niti na napravi zaključijo
	checkCudaErrors(cudaDeviceSynchronize());

	// vektor c prekopiramo iz naprave na gostitelja
	checkCudaErrors(cudaMemcpy(hc, dc, vectorLength * sizeof(float), cudaMemcpyDeviceToHost));

	// dokončamo izračun razdalje za napravo
	double distDevice = 0.0;
	for (int i = 0; i < vectorLength; i++)
		distDevice += hc[i];
	distDevice = sqrt(distDevice);

	// merjenje časa na napravi - konec
	clock_gettime(CLOCK_MONOTONIC, &stopDevice);
	double timeDevice = (stopDevice.tv_sec - startDevice.tv_sec) * 1e3 + (stopDevice.tv_nsec - startDevice.tv_nsec) / 1e6;

	// izračunamo razdaljo in izmerimo čas na gostitelju
	struct timespec startHost, stopHost;
	clock_gettime(CLOCK_MONOTONIC, &startHost);
	double distHost = 0.0;
	float diff;
	for (int i = 0; i < vectorLength; i++) {
		diff = ha[i] - hb[i];
		distHost += diff * diff;
	}
	distHost = sqrt(distHost);
	clock_gettime(CLOCK_MONOTONIC, &stopHost);
	double timeHost = (stopHost.tv_sec - startHost.tv_sec) * 1e3 + (stopHost.tv_nsec - startHost.tv_nsec) / 1e6;

	// rezultata izpišemo
	printf("naprava:      %lf (%lf ms)\ngostitelj:    %lf (%lf ms)\nnapaka (rel): %e\n", distDevice, timeDevice, distHost, timeHost, fabs(distDevice/distHost-1));

	// sprostimo pomnilnik na napravi
	checkCudaErrors(cudaFree(dc));
	checkCudaErrors(cudaFree(da));
	checkCudaErrors(cudaFree(db));

	// sprostimo pomnilnik na gostitelju
	free(hc);
	free(ha);
	free(hb);

	return 0;
}
