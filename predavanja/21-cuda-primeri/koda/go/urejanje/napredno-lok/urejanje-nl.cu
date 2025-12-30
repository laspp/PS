// bitonično urejanje
// več ščepcev, en ščepec lahko izvaja več notranjih zank, uporabimo skupni pomnilnik

#ifdef __cplusplus
extern "C" {
#endif

__device__ void bitonicSort(int *a, int len, int k, int j) {
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	if (gid < len/2) {
		int i1 = 2*j * (int)(gid / j) + (gid % j);	// prvi element
		int i2 = i1 ^ j;							// drugi element
		int dec = i1 & k;							// smer urejanja (padajoče: dec != 0)
		if ((dec == 0 && a[i1] > a[i2]) || (dec != 0 && a[i1] < a[i2])) {
			int temp = a[i1];
			a[i1] = a[i2];
			a[i2] = temp;
		}
	}
}

__device__ void bitonicSortShared(int *as, int len, int k, int j) {
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	if (gid < len/2) {
		int i1 = 2*j * (int)(gid / j) + (gid % j);	// prvi element
		int i2 = i1 ^ j;							// drugi element
		int dec = i1 & k;							// smer urejanja (padajoče: dec != 0)
		int i1s = i1 % (2 * blockDim.x);
		int i2s = i2 % (2 * blockDim.x);
		if ((dec == 0 && as[i1s] > as[i2s]) || (dec != 0 && as[i1s] < as[i2s])) {
			int temp = as[i1s];
			as[i1s] = as[i2s];
			as[i2s] = temp;
		}
	}
}

__device__ void copyToShared(int *as, int *a) {
	int i1Start = (2 * blockDim.x) * blockIdx.x;
	as[threadIdx.x] = a[i1Start + threadIdx.x];	
	as[blockDim.x + threadIdx.x] = a[i1Start + blockDim.x + threadIdx.x];
}

__device__ void copyFromShared(int *a, int *as) {
	int i1Start = (2 * blockDim.x) * blockIdx.x;
	a[i1Start + threadIdx.x] = as[threadIdx.x];	
	a[i1Start + blockDim.x + threadIdx.x] = as[blockDim.x + threadIdx.x];
}

__global__ void bitonicSortStart(int *a, int len) {
	extern __shared__ int as[];
	copyToShared(as, a);
	__syncthreads();
	for (int k = 2; k <= 2 * blockDim.x; k <<= 1) 
		for (int j = k/2; j > 0; j >>= 1) {
			bitonicSortShared(as, len, k, j);
			__syncthreads();
	}
	copyFromShared(a, as);
	__syncthreads();
}

__global__ void bitonicSortMiddle(int *a, int len, int k, int j) {
	bitonicSort(a, len, k, j);
}

__global__ void bitonicSortFinish(int *a, int len, int k) {
	extern __shared__ int as[];
	copyToShared(as, a);
	__syncthreads();
	for (int j = blockDim.x; j > 0; j >>= 1) {
		bitonicSortShared(as, len, k, j);
		__syncthreads();
	}
	copyFromShared(a, as);
	__syncthreads();
}

#ifdef __cplusplus
}
#endif