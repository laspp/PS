#!/bin/bash
#SBATCH --job-name=razdalje
#SBATCH --partition=gpu
#SBATCH --gpus=1
#SBATCH --constraint=v100s
#SBATCH --output=razdalje.out

source ../cudago-init.sh

cd globalni
srun go run razdalja-g.go
srun go run razdalja-g.go
srun go run razdalja-g.go
srun go run razdalja-g.go   
srun go run razdalja-g.go   

cd ../lok-stat
srun go run razdalja-ls.go
srun go run razdalja-ls.go
srun go run razdalja-ls.go
srun go run razdalja-ls.go
srun go run razdalja-ls.go

cd ../lok-din
srun go run razdalja-ld.go -k 1
srun go run razdalja-ld.go -k 1
srun go run razdalja-ld.go -k 1
srun go run razdalja-ld.go -k 1
srun go run razdalja-ld.go -k 1

srun go run razdalja-ld.go -k 2
srun go run razdalja-ld.go -k 2
srun go run razdalja-ld.go -k 2
srun go run razdalja-ld.go -k 2
srun go run razdalja-ld.go -k 2

srun go run razdalja-ld.go -k 3
srun go run razdalja-ld.go -k 3
srun go run razdalja-ld.go -k 3
srun go run razdalja-ld.go -k 3
srun go run razdalja-ld.go -k 3

srun go run razdalja-ld.go -k 4
srun go run razdalja-ld.go -k 4
srun go run razdalja-ld.go -k 4
srun go run razdalja-ld.go -k 4
srun go run razdalja-ld.go -k 4

cd ../lok-atom
srun go run razdalja-la.go -k 1
srun go run razdalja-la.go -k 1
srun go run razdalja-la.go -k 1
srun go run razdalja-la.go -k 1
srun go run razdalja-la.go -k 1  

srun go run razdalja-la.go -k 2
srun go run razdalja-la.go -k 2
srun go run razdalja-la.go -k 2
srun go run razdalja-la.go -k 2
srun go run razdalja-la.go -k 2
