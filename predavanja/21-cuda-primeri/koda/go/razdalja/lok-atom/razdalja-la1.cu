// na napravi izračunamo vsote kvadratov za vsak blok:
//		redukcija po drevesu, korak se zmanjšuje, v snopu ne potrebujemo sinhronizacije
//		atomarno seštevanje delnih vsot (eno na blok niti)

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorDistanceLA1(float *sPtr, const float *a, const float *b, int len) {
	// skupni pomnilnik niti v bloku
	extern __shared__ float part[];
	part[threadIdx.x] = 0.0;

	// kvadriranje razlike
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	float diff;
	while (gid < len) {
		diff = a[gid] - b[gid];
		part[threadIdx.x] += diff * diff;
		gid += gridDim.x * blockDim.x;
	}

	// počakamo, da vse niti zaključijo
	__syncthreads();

	// izračunamo delno vsoto za blok niti
	int idxStep;
	for(idxStep = blockDim.x >> 1; idxStep > 32 ; idxStep >>= 1) {
		if (threadIdx.x < idxStep)
			part[threadIdx.x] += part[threadIdx.x+idxStep];
		__syncthreads();
	}
	for( ; idxStep > 0 ; idxStep >>= 1 ) {
		if (threadIdx.x < idxStep)
			part[threadIdx.x] += part[threadIdx.x+idxStep];
		__syncwarp();
	}

	if (threadIdx.x == 0)
		atomicAdd(sPtr, part[0]);
}

#ifdef __cplusplus
}
#endif
