// računanje razlike elementov dveh vektorjev
// slaba rešitev: podpora samo za en blok niti

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorSubtract2(float *c, const float *a, const float *b, int len) {
	int gid = threadIdx.x;	
	// preprečimo pisanje v nerezervirani del pomnilnika
	if (gid < len)
		c[gid] = a[gid] - b[gid];
}

#ifdef __cplusplus
}
#endif
