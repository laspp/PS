// bitonično urejanje
// delajo vse niti

#ifdef __cplusplus
extern "C" {
#endif

__global__ void bitonicSortOV(int *a, int len, int k, int j) {
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

#ifdef __cplusplus
}
#endif