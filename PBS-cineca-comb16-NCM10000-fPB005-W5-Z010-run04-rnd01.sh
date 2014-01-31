#!/bin/bash#PBS -N r16 04 01
#PBS -A IscrC_VMStars
#PBS -q longpar
#PBS -l walltime=24:00:00
#PBS -l select=1:ncpus=1:ngpus=2
module purge
module load gnu/4.1.2
module load profile/advanced
module load boost/1.41.0--intel--11.1--binary
module load cuda/4.0

LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/5.0.35/none/lib64:/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib
export LD_LIBRARY_PATH

sh /gpfs/scratch/userexternal/bziosi00/plx-parameterSpace/cineca-comb16-run1_10-NCM10000-fPB005-W5-Z010/kiraLaunch-cineca-comb16-NCM10000-fPB005-W5-Z010-run04-rnd01.sh.sh