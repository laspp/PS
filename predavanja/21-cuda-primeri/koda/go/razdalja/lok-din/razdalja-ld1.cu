// na napravi izračunamo vsote kvadratov za vsak blok:
//		uporabimo skupni pomnilnik, dinamična rezervacija
//		nit 0 v bloku sešteva zaporedno, uporabimo register

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorDistanceLD1(float *p, const float *a, const float *b, int len) {
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
	if (threadIdx.x == 0) {
		float sum = 0.0;
		for (int i = 0; i < blockDim.x; i++)
			sum += part[i];
		p[blockIdx.x] = sum;
    }
}

#ifdef __cplusplus
}
#endif
