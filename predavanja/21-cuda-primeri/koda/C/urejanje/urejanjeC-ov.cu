// bitonično urejanje tabele celih števil
// 		argumenta: število niti v bloku in velikost tabele
//		elementi tabele so inicializirani naključno
// delajo vse niti

#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>
#include "cuda.h"
#include "helper_cuda.h"

__global__ void bitonicSort(int *a, int len, int k, int j) {
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	while (gid < len/2) {
		int i1 = 2*j * (int)(gid / j) + (gid % j);	// prvi element
		int i2 = i1 ^ j;							// drugi element
		int dec = i1 & k;							// smer urejanja (padajoče: dec != 0)
		if ((dec == 0 && a[i1] > a[i2]) || (dec != 0 && a[i1] < a[i2])) {
			int temp = a[i1];
			a[i1] = a[i2];
			a[i2] = temp;
		}
		gid += gridDim.x * blockDim.x;
	}
}

int main(int argc, char **argv) {
	// preberemo argumente iz ukazne vrstice
	int numThreads = 0;
	int tableLength = 0;
	if (argc == 3) {
		numThreads = atoi(argv[1]);
		tableLength = atoi(argv[2]);
	}
	if (numThreads <= 0 || tableLength <= 0 || ceil(log2(tableLength)) != floor(log2(tableLength))) {
		printf("usage:\n\t%s <number of block threads> <table length (power of 2)>\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	// določimo potrebno število blokov niti (rabimo toliko niti, kot je parov elemntov)
	int numBlocks = (tableLength/2 - 1) / numThreads + 1;

	// rezerviramo pomnilnik na gostitelju
	int *a = (int *)malloc(tableLength * sizeof(int));
	int *ha = (int *)malloc(tableLength * sizeof(int));
	
	// rezerviramo pomnilnik na napravi
	int *da;
	checkCudaErrors(cudaMalloc((void**)&da, tableLength * sizeof(int)));

	// nastavimo vrednosti tabel a in ha na gostitelju
	srand(time(NULL));
	for (int i = 0; i < tableLength; i++) {
        a[i] = rand();
		ha[i] = a[i];
    }

	// merjenje časa na napravi - začetek
	struct timespec startDevice, stopDevice;
	clock_gettime(CLOCK_MONOTONIC, &startDevice);

	// prenesemo tabelo a iz gostitelja na napravo
	checkCudaErrors(cudaMemcpy(da, ha, tableLength * sizeof(int), cudaMemcpyHostToDevice));

	// zaženemo kodo na napravi
	dim3 gridSize(numBlocks, 1, 1);
	dim3 blockSize(numThreads, 1, 1);

    for (int k = 2; k <= tableLength; k <<= 1) 
        for (int j = k/2; j > 0; j >>= 1) {
        	bitonicSort<<<gridSize, blockSize>>>(da, tableLength, k, j);
	        checkCudaErrors(cudaGetLastError());
        }

	// počakamo, da vse niti na napravi zaključijo
	checkCudaErrors(cudaDeviceSynchronize());

	// tabelo a prekopiramo iz naprave na gostitelja
	checkCudaErrors(cudaMemcpy(ha, da, tableLength * sizeof(int), cudaMemcpyDeviceToHost));

	// merjenje časa na napravi - konec
	clock_gettime(CLOCK_MONOTONIC, &stopDevice);
	double timeDevice = (stopDevice.tv_sec - startDevice.tv_sec) * 1e3 + (stopDevice.tv_nsec - startDevice.tv_nsec) / 1e6;

	// urejanje na gostitelju
	struct timespec startHost, stopHost;
	clock_gettime(CLOCK_MONOTONIC, &startHost);

    int i2, dec, temp;
    for (int k = 2; k <= tableLength; k <<= 1) 
        for (int j = k/2; j > 0; j >>= 1)
            for (int i1 = 0; i1 < tableLength; i1++) {
                i2 = i1 ^ j;
                dec = i1 & k;
                if (i2 > i1)
                    if ((dec == 0 && a[i1] > a[i2]) || (dec != 0 && a[i1] < a[i2])) {
                        temp = a[i1];
                        a[i1] = a[i2];
                        a[i2] = temp;
                    }
            }

	clock_gettime(CLOCK_MONOTONIC, &stopHost);
	double timeHost = (stopHost.tv_sec - startHost.tv_sec) * 1e3 + (stopHost.tv_nsec - startHost.tv_nsec) / 1e6;

    // preverimo rešitev
    int okDevice = 1;
    int okHost = 1;
    for (int i = 1; i < tableLength; i++) {
        okDevice &= (ha[i-1] <= ha[i]);
        okHost &= (a[i-1] <= a[i]);
    }
    printf("Device: %s (%lf ms)\n", okDevice ? "correct" : "wrong", timeDevice);
    printf("Host  : %s (%lf ms)\n", okHost ? "correct" : "wrong", timeHost);

	// sprostimo pomnilnik na napravi
	checkCudaErrors(cudaFree(da));

	// sprostimo pomnilnik na gostitelju
	free(a);
	free(ha);

	return 0;
}
