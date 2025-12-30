// bitonično urejanje
// dela samo polovica niti

#ifdef __cplusplus
extern "C" {
#endif

__global__ void bitonicSortO(int *a, int len, int k, int j) {
	int gid = blockIdx.x * blockDim.x + threadIdx.x;
	while (gid < len) {
		int i1 = gid;								// prvi element
		int i2 = i1 ^ j;							// drugi element
		int dec = i1 & k;							// smer urejanja (padajoče: dec != 0)
		if (i2 > i1)
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