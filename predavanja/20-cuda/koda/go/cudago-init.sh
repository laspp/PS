#!/usr/bin/env bash

# run as: 
#   source cuda-init.sh 
# or
#   . cuda-init.sh

# check CUDA version with 
#   nvcc --version
# and replace cudart version in the below environment variables appropriately

module load Go
module load CUDA
export CGO_CFLAGS=$(pkg-config --cflags cudart-12.8)
export CGO_LDFLAGS=$(pkg-config --libs cudart-12.8)
export PATH="~/go/bin/:$PATH"
