#!/bin/bash
#PBS -N r18-01-00
#PBS -A IscrC_VMStars
#PBS -q longpar
#PBS -l walltime=24:00:00
#PBS -l select=1:ncpus=1:ngpus=2

module purge
module load gnu/4.1.2
module load profile/advanced
module load boost/1.41.0--intel--11.1--binary
module load cuda/4.0

LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/4.0/none/lib64:/cineca/prod/compilers/cuda/4.0/none/lib:/cineca/prod/libraries/boost/1.41.0/intel--11.1--binary/lib:/cineca/prod/compilers/intel/11.1/binary/lib/intel64
export LD_LIBRARY_PATH

sh /gpfs/scratch/userexternal/bziosi00/plx-parameterSpace/cineca-comb18-run1_10-NCM10000-fPB020-W5-Z010/kiraLaunch-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.sh