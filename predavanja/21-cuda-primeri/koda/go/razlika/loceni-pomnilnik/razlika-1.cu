// računanje razlike elementov dveh vektorjev
// slaba rešitev: indeks niti je lahko večji od velikosti tabele

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorSubtract1(float *c, const float *a, const float *b, int len) {
	int gid = threadIdx.x;	
	c[gid] = a[gid] - b[gid];
}

#ifdef __cplusplus
}
#endif
