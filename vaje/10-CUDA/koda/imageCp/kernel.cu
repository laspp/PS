#ifdef __cplusplus
extern "C" {
#endif

// Copy image from input to output
__global__ void process(unsigned char *img_in, unsigned char *img_out, int width, int height) {
    // row
    int j = threadIdx.x + blockIdx.x * blockDim.x;
    // col
    int i = threadIdx.y + blockIdx.y * blockDim.y;
    int ipx = i * width + j;
    while (ipx < width * height) {
        img_out[ipx] = img_in[ipx];
        ipx += blockDim.x * gridDim.x * blockDim.y * gridDim.y;
    }
}
#ifdef __cplusplus
}
#endif