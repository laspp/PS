// na napravi izračunamo vsote kvadratov za vsak blok:
//		uporabimo skupni pomnilnik, dinamična rezervacija
//		redukcija po drevesu, korak se zmanjšuje, v snopu ne potrebujemo sinhronizacije

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorDistanceLD4(float *p, const float *a, const float *b, int len) {
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
		p[blockIdx.x] = part[0];
}

#ifdef __cplusplus
}
#endif
