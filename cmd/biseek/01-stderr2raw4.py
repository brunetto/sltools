#!/usr/bin/env python
# -*- coding: utf8 -*- 

from __future__ import division # no more "zero" integer division bugs!:P
import gzip, sys, re, time, glob, os
import multiprocessing as mp #(Process, Queue, ....)
import subprocess as sp
from string import maketrans
import SL_seeker_module as SLsm

"""This script retrieve BH, NS and mergers from the StarLab STDERR. This was chosen because
   it is smaller than the STDOUT and because has unique signatures. The STDOUT uses
   different label for neutron stars.
   This is intended as a replacement for Mapelli's script to fin BH and NS.
   
   For binary mergers we search for 
   
   binary_evolution:  merger within (2170,12170) triggered by 12170 at time 0.203125
   
   maybe we can include (FIXME) ids in the form \d{1,5}+\d{1,5} in case one or both the mergers
   are result of a previous merger.
   
   For BH and neutrons stars we search for 
   
   3606 hyper_giant_to_black_hole_at_time =   3.28 Myr (old mass = 88.8764)
   2976 super_giant_to_neutron_star_at_time =   9.27 Myr (old mass = 15.8217)
"""
		
#=============================================================================
#   MAIN
#=============================================================================

tt_glob = time.time()

separator = "==============================="

t_max = 5000

paths = ["."]
#filebase = "err-*-comb*-NCM*-fPB*-W*-Z*-run*-all.txt.gz"
filebase = "err-comb*-TF*-Rv*-NCM*-fPB*-W*-Z*-run*-all.*"
#filebase = "*ew_cineca1_bin_N5000_frac01_W5_Z001.txtot.gz" # TEST
outpath = "./Analysis/02-intermediate"

loop = 0
for path in paths:
	# Load one after the other the STDERR files
	for file_ in glob.glob(os.path.join(path, filebase)):
		print "Loop ", loop, " file ", file_
		fname = os.path.split(file_)[-1]
		print fname
		# Create object and init common variables
		obj = SLsm.CO_and_merger_seeker_obj(fname, path, outpath, t_max)
		# Init CO and merger variables
		obj.CO_and_merger_seeker_var()
		# Init binaries variables
		obj.binaries_var()
		
		print "Start reading..."
		n_lines = 0
		while True:
			# Read line by line
			line = obj.infile.readline()
			# Check for EOF
			if len(line) == 0: 
				print "EOF"
				break
			# Time checks
			line = obj.time_check(line)
			if line is None:
				break
			# Search for compact objects and mergers
			line = obj.CO_and_merger_seek(line)
			if line is None:
				break
			# Search for interesting binaries
			obj.search_for_binaries(line)
		print "Store and print..."
		obj.seek_print()
		obj.store_interesting_binaries()
		loop += 1

print separator
print "Done in ", time.time()-tt_glob, " seconds."
print "That's All Folks!"


