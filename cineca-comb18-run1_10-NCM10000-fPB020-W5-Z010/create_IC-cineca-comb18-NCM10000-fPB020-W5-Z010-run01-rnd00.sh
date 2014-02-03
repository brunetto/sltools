#!/bin/bash
/home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/makeking -n 10000 -w 5 -i -u \
| /home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/makemass -f 8  -l 0.1 -u 150 \
| /home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/makesecondary -f 0.20 -q -l 0.1 \
| /home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/add_star -R 1 -Z 0.10 \
| /home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/scale -R 1 -M 1\
| /home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/makebinary -f 2 -o 1 -l 1 -u 107836.09 \
> ics-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.txt