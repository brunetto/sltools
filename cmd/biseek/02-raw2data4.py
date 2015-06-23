#!/usr/bin/env python
# -*- coding: utf8 -*- 

from __future__ import division # no more "zero" integer division bugs!:P
import time, os, re, glob, sys, datetime, shutil
import numpy as np
import subprocess as sp #(Popen, PIPE, ...)
import multiprocessing as mp #(Process, Queue, ....)
from string import maketrans

""" Collect objects from the result of the stderr analysis and put the 
	compact object binaries in a single file and the compact triples in another file.
	Then put on record for each compact binary (that with the minimum sma) in one file
	and the history of each binary in one file for each one.
	
	changelog
	=========
	2013-01-07: check and fix for (i,j)=(j,i)
				ns-bh, bh-ns => mixed in the file name
"""

#=============================================================================
#   functions
#=============================================================================

def fish_orderer(infile, outpath, outfile, coolfile):  # outfile = all_the_fishes
	"""Collect objects from the result of the stderr analysis and put the 
	compact object binaries in a single file and the compact triples in another file.
	"""
	t = time.time()
	#regex = re.compile('fished_err-cineca-comb\d*-NCM(\d*)-fPB(\d*)-W(\d*)-Z(\d*)-run(\d*)-all\.\S*')
	regex = re.compile('fished_err-comb\d*-TF\S*-Rv\d*-NCM(\d*)-fPB(\d*)-W(\d*)-Z(\d*)-run(\d*)-all.\S*')
	params = regex.search(infile)
	print "Open file ", infile
	f = open(infile)
	out = open(os.path.join(outpath, outfile), 'a')
	cool = open(os.path.join(outpath, coolfile), 'a')
	n = params.group(5) # number of the simulation
	Z = params.group(4) # metallicity
	# Skip first line of comments (unsafe done in this way)
	line = f.readline()
	# Regex for the line:
	# ^ start of the line
	# (\d*\.*\S*) group 1 with one or more digits, 0 or more dots and chars 
	# different from whitespaces
	# \[(.*?)\] square brakets with chars inside
	# ...
	linereg = re.compile("^(\d*\.*\S*)\s+" + # GROUP 1: stars with system_time eg 10
						 "(\d*\.*\S*)\s+" + # GROUP 2: phys_time eg 2.62488
						 "\[(.*?)\]\s*"+ # GROUP 3: ids ['3460', '13460']
						 "([HS])\s+" + # GROUP 4: hardness flag H or S
						 "\[(.*?)\]\s*"+ # GROUP 5: objects ['bh++', '--']
						 "\[(\d*\.*\S*,\s\d*\.*\S*)\]\s+"+ # GROUP 6: masses [M_sun]  [28.622241655599996, 5.4053218129]
						 "(\d*\.*\S*)\s+"+ # GROUP 7: sma [pc] 0.000386024
						 "(\d*\.*\S*)\s+"+ # GROUP 8: period [Myr]  0.000121308304224
						 "(\d*\.*\S*)") # GROUP 9: ecc 0.562103
	# Loop through the file
	while True:
		line = f.readline()
		if len(line) == 0: # check for EOF
			print "EOF"
			break
		line = line.translate(None, "\n") # remove newline character
		#print line
		linesearch = linereg.search(line)
		if linesearch is None:
			print "Line search is void, exit..."
			print "Line is ", line
			sys.exit()
		system_time = linesearch.group(1)
		phys_time = linesearch.group(2)
		ids = linesearch.group(3).translate(None, "[ ']").split(",")
		hardflag = linesearch.group(4)
		n_obj = len(ids)
		types = linesearch.group(5).translate(None, "[ ']").split(",")
		if n_obj != len(types):
			print "Numer of ids differs from number of types, exit..."	
			print "Line is ", line
			sys.exit()
		masses = linesearch.group(6).split(",")
		sma = linesearch.group(7)
		period = linesearch.group(8)
		ecc = linesearch.group(9)
		if float(ecc) > 1.:
			continue # ignoring bugged systems with ecc > 1: they are not bound binaries 
		if n_obj > 2:
			if not "--" in types and not "ns++" in types and not "bh++" in types:
				print "Found cool multiple system!!!"
				print "ids ", ids, " len(ids) ", len(ids)
				cool.write(Z+"\t"+n+"\t"+system_time+"\t"+str(ids)+" "+hardflag+" "+str(types)+" "+str(masses)+" "+sma+" "+period+" "+ecc+"\n")
			print "Skipping system ", ids
			continue # ignoring multiple systems for now		
		#FIXME: inserito il flipping alfabetico qui
		binary_ids = "Z" + Z.ljust(3,"0") + "n" + n.rjust(3,"0") + "ids" + "a" +min(ids) + "b" +max(ids)#jointids
		out.write("{0:<5}".format(Z.ljust(3,"0"))+
			  "{0:<5}".format(n.rjust(3,"0"))+
			  "{0:<50}".format(binary_ids)+
			  "{0:<10}".format(str(int(float(system_time))))+
			  "{0:<20}".format(phys_time)+
			  "{0:<30}".format(ids[0]+"|"+ids[1])+
			  "{0:<14}".format(hardflag)+
			  "{0:<30}".format(types[0]+"|"+types[1])+
			  "{0:<30}".format(masses[0])+
			  "{0:<30}".format(masses[1])+
			  "{0:<20}".format(sma)+
			  "{0:<20}".format(period)+
			  "{0:<20}".format(ecc)+"\n"
			  )
	print "File ", infile, " analyzed in ", time.time()-t, " seconds."
		
def read_big_fish(infile="all_the_fishes.txt"):
	"""Load the table of all the compact objects from the file
	"""
	data = np.genfromtxt(infile, 
						dtype=([("Z", "|S5"),
								("n", int),
								("bin_id", "|S60"), 
								("sys_time", int), 
								("phys_time", float), 
								("ids", "|S40"), 
								("hardflag", "|S1"), 
								("types", "|S20"), 
								("mass_0", float), 
								("mass_1", float), 
								("sma", float), 
								("period", float), 
								("ecc", float)]))
	return data
	
def unique_seeker(outpath, data, selection="bh|bh"):
	"""Retrieve the unique id of the bhs and save their data
	"""
	# Select data of type "selection"
	if selection == "bh|bh" or selection == "ns|ns":
		cond = data["types"] == selection
		seeked_data = data[cond]
	elif selection == "mixed":
		print "Retrieving bh|ns..."
		cond = data["types"] == "bh|ns"
		seeked_data = data[cond]
		print "Stacking bh|ns with ns|bh retrieved on the fly..."
		cond = data["types"] == "ns|bh"
		seeked_data = np.hstack((seeked_data, data[cond]))
	# check and fix for (i,j)=(j,i)
	# FIXME: it is really sub-optimal because is substitute all the occurrences
	# each time it encounter the ids, but the first time only it's enough
	# but it takes not too long so for now it's ok...
	t_s = time.time()
	# FIXME: inserito il flipping sopra, lo disattivo qui
	#print "Check and fix for (i,j)=(j,i)..."
	#print "Really sub-optimal version..."
	#for i in xrange(seeked_data.size):
		#if i % 1000 == 0:
			#print "Loop ", i
		#corrected_bin_id = seeked_data["bin_id"][i]
		#oid = seeked_data["ids"][i].split("|")
		#flipped_temp = corrected_bin_id.split("idsa")
		#flipped_bin_id = flipped_temp[0]+"idsa"+oid[1]+"b"+oid[0]
		##reversed_oid = str(oid[1])+"|"+str(oid[0])
		#print "flipped_bin_id ", flipped_bin_id
		#cond = seeked_data["bin_id"] == flipped_bin_id
		#print "to be modified ", seeked_data["bin_id"][cond]
		#print "new value ", corrected_bin_id
		#seeked_data["bin_id"][cond] = corrected_bin_id
		#print "modified ", seeked_data[cond]["bin_id"]
	# All the data found
	fff = open("seeked.txt", 'a')
	for i in range(seeked_data.size):
		fff.write(str(seeked_data[i]["bin_id"])+"\n")
	fff.flush()
	fff.close()
	# Find unique id of the selected data
	uniques_id = np.unique(seeked_data["bin_id"])
	# Unique id from the (all) data found
	fff = open("uniques.txt", 'a')
	for i in range(uniques_id.size):
		fff.write(str(uniques_id[i])+"\n")
	fff.flush()
	fff.close()
	print "Found ", uniques_id.size, " uniques objects like ", selection
	# Find type and translate it into the file prefix
	intt = "|"
	outt = "-"
	ttt = maketrans(intt, outt)
	pref = selection.translate(ttt)
	if pref == "bh-ns" or pref == "ns-bh":
		pref = "mixed"
	# If file already exists don't rewrite the header
	# the header contains the datatype
	if not os.path.exists(outpath):
		print "File doesn't exists, header to be written..."
		header = True
	else:
		print "File exists, no header..."
		header = False
	ftot = open(os.path.join(outpath, pref+"_total.txt"), 'a')
	if header:
		ftot.write("#" + str(data.dtype)+"\n")
	for id_ in uniques_id:
		# Save obj history
		print "Saving history for ", id_
		udata = seeked_data[seeked_data["bin_id"] == id_]
		f = open(os.path.join(outpath, pref+"_"+id_+".txt"), 'w')
		f.write("#" + str(udata.dtype)+"\n")
		#np.savetxt(f, udata)
		for obj in udata:
			f.write(str(obj).translate(None, "(')")+"\n")
		f.flush()
		f.close()
		# Save the line with the minimum sma in the total file
		print "Saving line with minimum sma for ", id_
		idd = np.argmin(udata["sma"])
		#np.savetxt(ftot, udata[idd])
		ftot.write(str(udata[idd]).translate(None, "(')")+"\n")
		ftot.flush()
	ftot.flush()
	ftot.close()
	
def file_joiner(path, outfile="bh-bh_all.txt", pref='bh-bh_Z*'):
	if not os.path.exists(os.path.join(path, outfile)):
		outfile=open(os.path.join(path, outfile),"w")
		outfile.write("#[('Z', '|S5'), ('n', int), ('bin_id', '|S60'), "+
						"('sys_time', '<i8'), ('phys_time', '<f8'), "+
						"('ids', '|S40'), ('hardflag', '|S1'), ('types', '|S20'), "+
						"('mass_0', '<f8'), ('mass_1', '<f8'), ('sma', '<f8'), "+
						"('period', '<f8'), ('ecc', '<f8')]\n")
		for infile in glob.glob(os.path.join(path, pref)):
			for line in open(infile):
				if not line.startswith("#"):
					outfile.write(line)    
		outfile.flush()
		outfile.close()
	
	
#=============================================================================
# Parameters
#=============================================================================
tt_glob = time.time()

# Initialize queues
in_queue = mp.Queue()

#inpath = "/home/ziosi/data_mapelli/2013-04-19-results_sim_1"
#outpath = "/home/ziosi/Dropbox/Uni/PhD_Mapelli/1-BHs_binaries_from_StarLab&SL_fisher/data/only_1/03-final"

inpath = "./Analysis/02-intermediate"
outpath = "./Analysis/03-final"


pathfile = open("conf.txt", "w")
pathfile.write("coocked_path = "+outpath+"\n")
pathfile.flush()
pathfile.close()

if os.path.exists(outpath):
	print "Output folder already existing: deleting..."
	shutil.rmtree(outpath)
print "Creating output folder..."
os.mkdir(outpath)


outfile = "all_the_fishes.txt"
	
#=============================================================================
#   MAIN
#=============================================================================
if os.path.exists("seeked.txt"):
	os.remove("seeked.txt")
if os.path.exists("uniques.txt"):
	os.remove("uniques.txt")

if True:
	outfile = "all_the_fishes.txt"
	coolfile = "cool_triples.txt"
	print "Collecting all interesting binaries into a single file: ", os.path.join(outpath, outfile)
	print "and collecting interesting triples into : ", os.path.join(outpath, coolfile)
	out = open(os.path.join(outpath, outfile), 'w')
	out.write("{0:<5}".format("# Z")+
			  "{0:<5}".format("n")+
			  "{0:<50}".format("binary_ids")+
			  "{0:<10}".format("sys_time")+
			  "{0:<20}".format("phys_time [Myr]")+
			  "{0:<30}".format("objects ids")+
			  "{0:<14}".format("hardflag")+
			  "{0:<30}".format("types")+
			  "{0:<30}".format("masses[0]")+
			  "{0:<30}".format("masses[1]")+
			  "{0:<20}".format("sma")+
			  "{0:<20}".format("period")+
			  "{0:<20}".format("ecc")+"\n"
			  )
	out.flush()
	out.close()

	cool = open(coolfile, 'w')
	cool.write("# Interesting objects\n")
	cool.write("#Z"+"\t"+"n"+"\t"+"system_time"+"\t"+"ids"+" "+"hardflag"+" "+"types"+" "+"masses"+" "+"sma"+" "+"period"+" "+"ecc"+"\n")
	cool.flush()
	cool.close()
	
	print "Start joining the fishes..."
	#for inf in glob.glob(os.path.join(inpath, 'fished_err-cineca-comb*-NCM*-fPB*-W*-Z*-run*-all.txt.gz.txt')):
	for inf in glob.glob(os.path.join(inpath, 'fished_err-comb*-TF*-Rv*-NCM*-fPB*-W*-Z*-run*-*.*')):
		print "Considering ", inf
		if "uninteresting" in inf:
			pass
		else:
			print "Enqueue ", inf
			in_queue.put(inf)
	print "Number of files to process: ", in_queue.qsize()

	print "Queue size ", in_queue.qsize()
	while in_queue.qsize() != 0:
		infile = in_queue.get()
		fish_orderer(infile, outpath, outfile, coolfile)

if True:
	print "Reading data from the total file and separate objects..."
	data = read_big_fish(os.path.join(outpath, outfile))
	print "Searching for bh-bh pairs..."
	unique_seeker(outpath, data, selection="bh|bh")
	print "Searching for ns-ns pairs..."
	unique_seeker(outpath, data, selection="ns|ns")
	print "Searching for mixed pairs..."
	unique_seeker(outpath, data, selection="mixed")
	
if True:
	print "Joining histories..."
	file_joiner(path = outpath, outfile="bh-bh_all.txt", pref='bh-bh_Z*')
	file_joiner(path = outpath, outfile="ns-ns_all.txt", pref='ns-ns_Z*')
	file_joiner(path = outpath, outfile="mixed_all.txt", pref='mixed_Z*')
		
print "Wall time for all: ", time.time()-tt_glob, " seconds."

print "That's All Folks and thanks for all the fish!"

