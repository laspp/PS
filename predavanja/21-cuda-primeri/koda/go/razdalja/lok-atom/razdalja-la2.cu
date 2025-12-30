// na napravi izračunamo vsote kvadratov za vsak blok:
//		redukcija po drevesu, korak se zmanjšuje, v snopu ne potrebujemo sinhronizacije
//		atomarno seštevanje za vsako nit

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorDistanceLA2(float *sPtr, const float *a, const float *b, int len) {
	// skupni pomnilnik niti v bloku
	extern __shared__ float part[];
	part[threadIdx.x] = 0.0;

	// kvadriranje razlike
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	float diff;
	while (gid < len) {
		diff = a[gid] - b[gid];
		atomicAdd(&sPtr[0], diff * diff);
		gid += gridDim.x * blockDim.x;
	}
}

#ifdef __cplusplus
}
#endif
