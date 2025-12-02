#!/bin/bash
#SBATCH --nodes=1
#SBATCH --array=0-4
#SBATCH --reservation=fri
#SBATCH --output=telefon-%a.txt

#Zazeni s sbatch run_telefon.sh na gruci Arnes
path=./telefonUDP

module load Go
go build $path/telefon.go
srun telefon -p 9000 -id $SLURM_ARRAY_TASK_ID -n $SLURM_ARRAY_TASK_COUNT