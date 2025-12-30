// informacije o napravi
// prevajanje:
//      module load CUDA
//      srun --partition=gpu --gpus=1 nvcc -o napravaC napravaC.c
// izvajanje:
//      srun --partition=gpu --gpus=1 ./napravaC

#include <stdio.h>
#include "cuda.h"
#include "cuda_runtime.h"
#include "helper_cuda.h"

int main(int argc, char **argv) 
{
    
    // Get number of GPUs
    int deviceCount = 0;
    cudaError_t error = cudaGetDeviceCount(&deviceCount);
    if (error != cudaSuccess) 
    {
        printf("cudaGetDeviceCount error %d\n-> %s\n", error, cudaGetErrorString(error));
        exit(EXIT_FAILURE);
    }

    // Get propreties of each device
    for (int dev = 0; dev < deviceCount; dev++) 
    {

        struct cudaDeviceProp prop;
        int value;

        cudaGetDeviceProperties(&prop, dev);
        
        printf("\nDevice %d: %s\n", dev, prop.name);
        printf("\ncudaDeviceGetProperties:\n");
        printf("  CUDA Architecture:                             %s, %d.%d\n", _ConvertSMVer2ArchName(prop.major, prop.minor), prop.major, prop.minor);
        printf("\n");
        printf("  GPU clock rate (MHz):                          %d\n", prop.clockRate/1000);
        printf("  Memory clock rate (MHz):                       %d\n", prop.memoryClockRate/1000);
        printf("  Memory bus width (bits):                       %d\n", prop.memoryBusWidth);
        printf("  Peak memory bandwidth (GB/s):                  %.0f\n", 2.0*prop.memoryClockRate*(prop.memoryBusWidth/8)/1.0e6);
        printf("\n");
        printf("  Number of MPs:                                 %d\n", prop.multiProcessorCount);
        printf("  Number of cores per MP:                        %d\n", _ConvertSMVer2Cores(prop.major, prop.minor));
        printf("  Total number of cores:                         %d\n", _ConvertSMVer2Cores(prop.major, prop.minor) * prop.multiProcessorCount);
        printf("\n");
        printf("  Total amount of global memory (GB):            %.0f\n", prop.totalGlobalMem / 1073741824.0f);
        printf("  Total amount of shared memory per MP (kB):     %d\n", prop.sharedMemPerMultiprocessor/1024);
        printf("  Total amount of shared memory per block (kB):  %zu\n", prop.sharedMemPerBlock/1024);
        cudaDeviceGetAttribute (&value, cudaDevAttrL2CacheSize, dev);
        printf("  Size of L2 cache in MB:                        %.0f\n", value/1048576.0f);
        printf("  Maximum number of registers per MP:            %d\n", prop.regsPerMultiprocessor);
        printf("  Total number of registers per block:           %d\n", prop.regsPerBlock);
        printf("\n");
        printf("  Maximum number of threads per MP:              %d\n", prop.maxThreadsPerMultiProcessor);
        printf("  Maximum number of threads per block:           %d\n", prop.maxThreadsPerBlock);
        cudaDeviceGetAttribute (&value, cudaDevAttrMaxBlocksPerMultiprocessor, dev);  // works in newer CUDA versions
        printf("  Maximum number of blocks per MP:                %d\n", value);  
        printf("  Warp size:                                     %d\n", prop.warpSize);
        printf("\n");
        printf("  Max dimension of a thread block (x,y,z):      (%d, %d, %d)\n", prop.maxThreadsDim[0], prop.maxThreadsDim[1], prop.maxThreadsDim[2]);
        printf("  Max dimension of a grid (x,y,z):              (%d, %d, %d)\n", prop.maxGridSize[0], prop.maxGridSize[1], prop.maxGridSize[2]);

        printf("\ncudaDeviceGetAttribute:\n");
    }
}
