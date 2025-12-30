// računanje razlike elementov dveh vektorjev
// slaba rešitev: ne deluje, če je elementov več kot niti

#ifdef __cplusplus
extern "C" {
#endif

__global__ void vectorSubtract3(float *c, const float *a, const float *b, int len) {
	// določimo globalni indeks elementov
	int gid = blockDim.x * blockIdx.x + threadIdx.x;	
	// preprečimo pisanje v nerezervirani del pomnilnika
	if (gid < len)
		c[gid] = a[gid] - b[gid];
}

#ifdef __cplusplus
}
#endif
