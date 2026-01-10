#!/bin/bash
#SBATCH --job-name=razdalje
#SBATCH --partition=gpu
#SBATCH --gpus=1
#SBATCH --constraint=v100s
#SBATCH --output=urejanje.out

source ../cudago-init.sh

echo "osnovno"
cd osnovno
srun go run urejanje-o.go
srun go run urejanje-o.go
srun go run urejanje-o.go
srun go run urejanje-o.go
srun go run urejanje-o.go

echo "osnovno-vse-niti"
cd ../osnovno-vse-niti
srun go run urejanje-ov.go
srun go run urejanje-ov.go
srun go run urejanje-ov.go
srun go run urejanje-ov.go
srun go run urejanje-ov.go

echo "napredno"
cd ../napredno
srun go run urejanje-n.go
srun go run urejanje-n.go
srun go run urejanje-n.go
srun go run urejanje-n.go
srun go run urejanje-n.go   

echo "napredno-lok"
cd ../napredno-lok
srun go run urejanje-nl.go
srun go run urejanje-nl.go
srun go run urejanje-nl.go
srun go run urejanje-nl.go
srun go run urejanje-nl.go

