// na napravi izračunamo vsote kvadratov za vsak blok:
// 		uporabimo skupni pomnilnik, rezerviramo ga statično
//		nit 0 v bloku sešteva zaporedno

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorDistanceLS(float *p, const float *a, const float *b, int len) {
	// skupni pomnilnik niti v bloku
	__shared__ float part[1024];
	part[threadIdx.x] = 0.0;

	// kvadriranje razlike
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	float diff;
	while (gid < len) {
		diff = a[gid] - b[gid];
		part[threadIdx.x] += diff * diff;
		gid += gridDim.x * blockDim.x;
	}

	// počakamo, da vse niti zaključijo s kvadriranjem razlike
	__syncthreads();

	// izračunamo delno vsoto za blok niti
	if (threadIdx.x == 0) {
        p[blockIdx.x] = 0.0;
        for (int i = 0; i < blockDim.x; i++)
			p[blockIdx.x] += part[i];
    }
}


#ifdef __cplusplus
}
#endif
