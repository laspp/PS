// bitonično urejanje
// več ščepcev, en ščepec lahko izvaja več notranjih zank

#ifdef __cplusplus
extern "C" {
#endif

__device__ void bitonicSort(int *a, int len, int k, int j) {
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

__global__ void bitonicSortStart(int *a, int len) {
	for (int k = 2; k <= 2 * blockDim.x; k <<= 1) 
		for (int j = k/2; j > 0; j >>= 1) {
			bitonicSort(a, len, k, j);
			__syncthreads();
	}
}

__global__ void bitonicSortMiddle(int *a, int len, int k, int j) {
	bitonicSort(a, len, k, j);
}

__global__ void bitonicSortFinish(int *a, int len, int k) {
	for (int j = blockDim.x; j > 0; j >>= 1) {
		bitonicSort(a, len, k, j);
		__syncthreads();
	}
}


#ifdef __cplusplus
}
#endif