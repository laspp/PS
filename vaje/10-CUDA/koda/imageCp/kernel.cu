#ifdef __cplusplus
extern "C" {
#endif

// Copy image from input to output
__global__ void process(unsigned char *img_in, unsigned char *img_out, int width, int height) {

    for (int i = threadIdx.y + blockIdx.y * blockDim.y; i < height; i += blockDim.y * gridDim.y) {
        for (int j = threadIdx.x + blockIdx.x * blockDim.x; j < width; j += blockDim.x * gridDim.x) {
            int ipx = i * width + j;
            img_out[ipx] = img_in[ipx];
        }
    }
}
#ifdef __cplusplus
}
#endif