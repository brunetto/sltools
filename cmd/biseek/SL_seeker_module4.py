#!/usr/bin/env python
# -*- coding: utf8 -*- 

from __future__ import division # no more "zero" integer division bugs!:P
import gzip
import time, os, re, glob, sys
import subprocess as sp #(Popen, PIPE, ...)
import multiprocessing as mp #(Process, Queue, ....)
from string import maketrans

""" This module contains some tools to manipulate Starlab outputs.
"""

#=============================================================================
#   functions
#=============================================================================
#FIXME separare i misti
class CO_and_merger_seeker_obj(object):
	"""Slice object to handle all the information about a slice.
	"""
	def __init__(self, fname, fpath, outpath, t_max):
		"""Construct slice object."""
		self.fpath = fpath
		self.fname = fname
		self.t_max = t_max
		
		#binaries_var(infile, Z, N, W, frac, n, filebase_path, t_max)
		
		self.separator = "==============================="
		self.filename = os.path.join(self.fpath, self.fname)
		# Find parameters from infile name	
		name_param_reg = re.compile('err-comb\d*-TF\S*-Rv\d*-NCM(\d*)-fPB(\d*)-W(\d*)-Z(\d*)-run(\d*)-\S*\.\S*')
		name_param_res = name_param_reg.search(fname)
		if name_param_res is None:
			print "Error retrieving data from file name, exit without creating object!"
			sys.exit(1)
		self.n = name_param_res.group(5)
		self.N = name_param_res.group(1)
		self.frac = name_param_res.group(2)
		self.W = name_param_res.group(3)
		self.Z = name_param_res.group(4)
		self.n_lines = 0
		
		# Open STDERR
		print "Open file ", self.filename
		if self.filename.endswith(".gz"):
			print "Gzipped file"
			import gzip
			self.infile = gzip.open(self.filename, 'r')
		else:
			print "Uncompressed file"
			self.infile = open(self.filename, 'r')
		
		# Outputs
		#========
		self.outpath = outpath
		
		
	def CO_and_merger_seeker_var(self):
		# NS outfile
		self.ns_outfile_name = os.path.join(self.outpath, "ns_list_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		# BH outfile
		self.bh_outfile_name = os.path.join(self.outpath, "bh_list_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		# Normal mergers outfile
		self.nCO_mergers_outfile_name = os.path.join(self.outpath, "normal_mergers_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		# NS mergers outfile
		self.ns_mergers_outfile_name = os.path.join(self.outpath, "ns_mergers_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		# BH mergers outfile
		self.bh_mergers_outfile_name = os.path.join(self.outpath, "bh_mergers_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		# MIXED mergers outfile
		self.mixed_mergers_outfile_name = os.path.join(self.outpath, "mixed_mergers_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt")
		
		print "Initialize containers"
		# Initialize containers list
		self.ns_list = []
		self.bh_list = []
		self.nCO_mergers_list = []
		self.ns_mergers_list = []
		self.bh_mergers_list = []
		self.mixed_mergers_list = []
	
		print "Initialize regexp"
		# Compile regex for time limit, if (sys-)time is greater than 400 we can stop
		# group(1) is the time
		self.time_step_reg = re.compile("Time\s*=\s*(\d+)\s*")
		# Initialize time_limit
		self.time_step = 0
		# Compile regex for merger
		# group(1) are the two ids
		# group(2) is the merger time in units of FIXME
		self.merger_reg = re.compile("binary_evolution:\s*merger within\s*\((\d{1,5}\+*\d*\+*\d*,\d{1,5}\+*\d*\+*\d*)\)\s*triggered by \d{1,5}\+*\d*\s*at time\s*(\d+\.*\d*)")
		self.merger_time_reg = re.compile("Collision at time =\s*\d+\.*\d*\s*\((\d+\.*\d*)\s*\[Myr\]\)\s*between")
		self.collision_result_reg = re.compile("\d+\s*\((\S+);\s*M\s*=\s*(\d+\.*\S*)\s*\[Msun\]")
		# Compile regex for NS
		# group(1) is the NS id
		# group(2) is the NS formation time in units of FIXME
		# group(3) is the time unit
		# group(4) is the progenitor mass is the progenitor mass
		#bh_reg = re.compile("(\d{1,5}\+*\d*\+*\d*)\s+\S+_to_black_hole_at_time\s*=\s*(\d+\.*\d*)\s*(\S+)\s*\(old mass\s*=\s*(\d+\.*\d*)\)")
		self.bh_reg = re.compile("(\d{1,5}\+*\d*\+*\d*)\s+\S+_to_black_hole_at_time\s*=\s*(\d+\.*\d*)\s*(\S+)")
		# Compile regex for BH
		# group(1) is the BH id
		# group(2) is the HB formation time in units of FIXME
		# group(3) is the time unit
		# group(4) is the progenitor mass
		#ns_reg = re.compile("(\d{1,5}\+*\d*\+*\d*)\s+\S+_to_neutron_star_at_time\s*=\s*(\d+\.*\d*)\s*(\S+)\s*\(old mass\s*=\s*(\d+\.*\d*)\)")
		self.ns_reg = re.compile("(\d{1,5}\+*\d*\+*\d*)\s+\S+_to_neutron_star_at_time\s*=\s*(\d+\.*\d*)\s*(\S+)")
		
		print "Init broken log flag"
		# Initialize some flags
		self.tt_res_check_1 = 0
		self.tt_res_check_2 = 0
		
		print "Init timesteps"
		self.time_step = None
		self.time_back = None
		# In case of simulation crash may be 
		#that some timesteps are partially duplicated
		self.duplicated_times = []
		
		print "Create translation table"
		# Create table for translation from , to + in the binary id
		instr = ","
		outstr = "+"
		self.transtable = maketrans(instr, outstr)
		self.container_list = []	

	def time_check(self, line):
		#FIXME: Move somewhere more general
		self.n_lines += 1
		#======================================================
		#    TIME CHECKS
		#======================================================
		# Check for time greater than 400 to exit
		if "Time = " in line:
			self.time_back = self.time_step
			self.time_step = self.time_step_reg.search(line).group(1)
			print "Coming from timestep ", self.time_back
			print "Going into timestep ", self.time_step
			if self.time_back == self.time_step:
				print "Found duplicated timestep ", self.time_back
				self.duplicated_times.append(self.time_step)
			# Zero the broken log check counters
			self.tt_res_check_1 = 0 
			self.tt_res_check_2 = 0
			if int(self.time_step) > self.t_max:
				print self.separator
				print "Reached time-step max: ", line[:12]
				print "Exit main loop!"
				print "Parsed ", self.n_lines, " lines"
				print self.separator
				return #FIXME??? -> it seems ok
		# Check for resonance corrupted logs
		if "pp3 output" in line: 
			self.tt_res_check_1 += 1
		if "time step warning" in line:
			self.tt_res_check_2 += 1
		if self.tt_res_check_1 > 10 or self.tt_res_check_2 > 10:
			print "Probably broken log, save and exit..."
			return #FIXME??? -> it seems ok
		return line

	def CO_and_merger_seek(self, line):
		#======================================================
		#    MERGERS
		#======================================================
		# Search for mergers
		merger_res = self.merger_reg.search(line)
		# Search for NS formation
		ns_res = self.ns_reg.search(line)
		# Search for BH formation
		bh_res = self.bh_reg.search(line)
		# Found merger???
		if merger_res is not None:
			# Find time of merging searching for the Collision time line
			while True:
				line = self.infile.readline()
				# Check for EOF
				if len(line) == 0: 
					print "EOF"
					break
				self.n_lines += 1
				# If line is collapsing time store the time in Myr
				if "Collision at time" in line:
					merger_time_res = self.merger_time_reg.search(line)
					if merger_time_res is not None:
						merger_time = merger_time_res.group(1)
					else:
						print "Error searching for merger time, exit"
						sys.exit(1)
				if "merge_nodes: merger product:" in line:
					# Read the line after the line containing "merger product" string
					line = self.infile.readline()
					self.n_lines += 1
					# Search for the collision result
					collision_result_res = self.collision_result_reg.search(line)
					if collision_result_res is not None:
						collision_result = collision_result_res.group(1)
						collision_mass = collision_result_res.group(2)
						break
					else:
						print "Error searching for collision result, exit"
						sys.exit(1)
			# Find merger id
			merger_id = merger_res.group(1).translate(self.transtable, " ")
			# Check if one or both merging object are COs
			# ns_mergers and bh_mergers are redefined each loop so they don't need to be cleared
			localNSlist = [item['id'] for item in self.ns_list]
			localBHlist = [item['id'] for item in self.bh_list]
			ns_mergers = [element for element in merger_id.split("+") if element in localNSlist]  
			bh_mergers = [element for element in merger_id.split("+") if element in localBHlist]
			merger_data = {"id": merger_id, 
						   "time": merger_time,
						   "mass": collision_mass,
						   "merger_result": collision_result}
			# BH+BH merger
			if len(ns_mergers) == 0 and len(bh_mergers) == 2:
				if "bh" in collision_result:
					# FIXME: controllare non sia già presente
					print "Found BH merger ", merger_data
					self.bh_list.append(merger_data)
					self.bh_mergers_list.append(merger_data)
				else:
					print "Error matching merging objects and collision result"
			# NS+NS merger
			elif len(ns_mergers) == 2 and len(bh_mergers) == 0:
				if "ns" in collision_result:
					print "Found NS merger with NS result"
					self.ns_list.append(merger_data)
					self.ns_mergers_list.append(merger_data)
				elif "bh" in collision_result:
					print "Found NS merger with BH result"
					self.bh_list.append(merger_data)
					self.ns_mergers_list.append(merger_data)
				else:
					print "Error matching merging objects and collision result"
			# BH+NS (mixed) merger
			elif len(ns_mergers) == 1 and len(bh_mergers) == 1:
				if "bh" in collision_result:
					print "Found BH+NS merger ", merger_data
					self.bh_list.append(merger_data)
					self.mixed_mergers_list.append(merger_data)
				else:
					print "Error matching merging objects and collision result"
			# BH+star merger
			elif len(ns_mergers) == 0 and len(bh_mergers) == 1:
				if "bh" in collision_result:
					print "Found BH+star merger ", merger_data
					self.bh_list.append(merger_data)
					self.nCO_mergers_list.append(merger_data)
				else:
					print "Error matching merging objects and collision result"
			# NS+star or star+star
			elif len(ns_mergers) == 0 or len(ns_mergers) == 1 and len(bh_mergers) == 0:
				if "bh" in collision_result:
					print "Found merger with BH result ", merger_data
					self.bh_list.append(merger_data)
					self.nCO_mergers_list.append(merger_data)
				elif "ns" in collision_result:
					print "Found merger with NS result ", merger_data
					self.ns_list.append(merger_data)
					self.nCO_mergers_list.append(merger_data)
				else:
					print "Found merger with no compacts", merger_data
					self.nCO_mergers_list.append(merger_data)
			else:
				print "No matching merger, exit!"
				print merger_data
				sys,exit(1)
			return line # FIXME???
		#======================================================
		#    COMPACT OBJECTS
		#======================================================
		# Found NS
		elif ns_res is not None:
			# Store NS data if not already stored
			ns_data = {"id": ns_res.group(1),
				  "time": ns_res.group(2),
				  "mass": "--", #ns_res.group(4),
				  "merger_result": "no"}
			if ns_data not in self.ns_list:
				self.ns_list.append(ns_data)
				print "Found NS ", self.ns_list[-1]
			# Skip to the next line
			return line # FIXME???
		# Found BH
		elif bh_res is not None:
			# Store NS data if not already stored
			bh_data = {"id": bh_res.group(1),
				  "time": bh_res.group(2),
				  "mass": "--", #bh_res.group(4),
				  "merger_result": "no"}
			if bh_data not in self.bh_list:
				self.bh_list.append(bh_data)
				print "Found BH ", self.bh_list[-1]
			# Skip to the next line
			return line # FIXME???
		# Found nothing!:P
		else:    
			pass
		return line
		
	def list_print(self, list_):		
		for item in list_:
			print str(item)#+"\n"

	def seek_print(self):
		# Print some info
		self.ns2bh_list = [item for item in self.ns_list if item in self.bh_list]
		if len(self.ns2bh_list) == 0:   
			self.ns2bh_check = "ok"
		else:
			self.ns2bh_check = "something went wrong"
		
		print self.separator
		print "Merger_list"
		self.list_print(self.nCO_mergers_list)
		print self.separator
		print "ns_mergers_list"
		self.list_print(self.ns_mergers_list)
		print self.separator
		print "bh_mergers_list"
		self.list_print(self.bh_mergers_list)
		print self.separator
		print "ns_list"
		self.list_print(self.ns_list)
		print self.separator
		print "ns2bh_list"
		self.list_print(self.ns2bh_list)
		print self.separator
		print "bh_list"
		self.list_print(self.bh_list)
		print self.separator
		print "Found"
		print len(self.ns_mergers_list), " NS mergers"
		print len(self.bh_mergers_list), " BH mergers"
		print len(self.mixed_mergers_list), " NS+BH mergers"
		print len(self.nCO_mergers_list), " non_CO_mergers"
		print len(self.ns_list), " neutron stars"
		print len(self.ns2bh_list), " neutron stars-> BHs => ", self.ns2bh_check
		print len(self.bh_list), " black holes"
		
		print self.separator
		print "Writing data to disk..."
		
		print "Writing found NSs"
		#self.ns_list = set(self.ns_list)
		ns_file = open(self.ns_outfile_name, 'w')
		ns_file.write("# read with eval(f.readline().split('#')[1])\n")
		ns_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.ns_list:
			ns_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		ns_file.close()
		
		print self.separator
		print "Writing found BHs"
		#self.bh_list = set(self.bh_list)
		bh_file = open(self.bh_outfile_name, 'w')
		bh_file.write("# read with eval(f.readline().split('#')[1])\n")
		bh_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.bh_list:
			bh_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		bh_file.close()
		
		print self.separator
		print "Writing found non compact mergers"
		#self.nCO_mergers_list = set(self.nCO_mergers_list)
		nCO_mergers_file = open(self.nCO_mergers_outfile_name, 'w')
		nCO_mergers_file.write("# read with eval(f.readline().split('#')[1])\n")
		nCO_mergers_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.nCO_mergers_list:
			nCO_mergers_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		nCO_mergers_file.close()
		print self.separator
		print "Writing found NS mergers"
		#self.ns_mergers_list = set(self.ns_mergers_list)
		ns_mergers_file = open(self.ns_mergers_outfile_name, 'w')
		ns_mergers_file.write("# read with eval(f.readline().split('#')[1])\n")
		ns_mergers_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.ns_mergers_list:
			ns_mergers_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		ns_mergers_file.close()
		
		print self.separator
		print "Writing found BH mergers"
		#self.bh_mergers_list = set(self.bh_mergers_list)
		bh_mergers_file = open(self.bh_mergers_outfile_name, 'w')
		bh_mergers_file.write("# read with eval(f.readline().split('#')[1])\n")
		bh_mergers_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.bh_mergers_list:
			bh_mergers_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		bh_mergers_file.close()
		
		print self.separator
		print "Writing found mixed mergers"
		#self.mixed_mergers_list = set(self.mixed_mergers_list)
		mixed_mergers_file = open(self.mixed_mergers_outfile_name, 'w')
		mixed_mergers_file.write("# read with eval(f.readline().split('#')[1])\n")
		mixed_mergers_file.write("# [('id','|S30'), ('mass', float), ('phys_time',float), ('from_merger', '|S5')\n")
		for i in self.mixed_mergers_list:
			mixed_mergers_file.write(i["id"]+" "+i["mass"]+" "+i["time"]+" "+i["merger_result"]+"\n")
		mixed_mergers_file.close()
		print self.separator

	def flush_err_data(self, out, container):
		""" Flush data to file and reset the container.
		"""
		out.write("{0:<15}".format(container["system_time"])+
				"{0:<20}".format(container["phys_time"])+
				"{0:<50}".format(container["ids"])+
				"{0:<14}".format(container["hardflag"])+
				"{0:<35}".format(container["objects"])+
				"{0:<50}".format(container["masses"])+
				"{0:<20}".format(container["sma"])+
				"{0:<20}".format(container["period"])+
				"{0:<15}".format(container["ecc"])
				)
		out.write("\n")
		out.flush()
		return

	def binaries_var(self):
		"""Extract objects and their properties for the given snapshots
		from the Starlab infile
		"""
		t = time.time()
		# Open files
		# FIXME: I think they are not useful anymore
		#filebase = "stripped_cineca"+self.n+"_bin_N"+self.N+"_frac"+self.frac+"_W"+self.W+"_Z"+self.Z+".txtot.gz" 
		#baselist = "_list_n"+self.n+"_W"+self.W+"_N"+self.N+"_Z"+self.Z+".txt2"
		print "Done in ", time.time()-t, "secs."
		t = time.time()
		print "Start fishing"
		# Find units
		print "Search for units"
		mass_re = re.compile('\[m\]:\s+(\d*\.*\d*)\s*(\S*)')
		r_re    = re.compile('\[R\]:\s+(\d*\.*\d*)\s*(\S*)')
		time_re = re.compile('\[T\]:\s+(\d*\.*\d*)\s*(\S*)')
		self.mass_unit = self.r_unit = self.time_unit = None
		# Units loop
		while True: # search for units
			line = self.infile.readline()
			if len(line) == 0: # check for EOF
				print "EOF searching for units, ERROR!!!"
				sys.exit()
			line = line.translate(None, "\n") # remove newline character
			mass_search = mass_re.search(line)
			if mass_search:
				self.mass_unit = [float(mass_search.group(1)), mass_search.group(2)]
			r_search = r_re.search(line)
			if r_search:
				self.r_unit = [float(r_search.group(1)), r_search.group(2)]
			time_search = time_re.search(line)
			if time_search:
				self.time_unit = [float(time_search.group(1)), time_search.group(2)]
			if self.mass_unit and self.r_unit and self.time_unit:
				print "Found units"
				break
		print "Units:"
		print "Mass: ", self.mass_unit[0], self.mass_unit[1]
		print "Radius: ", self.r_unit[0], self.r_unit[1]
		print "Time: ", self.time_unit[0], self.time_unit[1]
		# Regex: ^\s+ line starts with one or more spaces
		#        U* maybe you will find 0 or more U characters
		#        \s* and 0 or more space
		#        () put the things inside in group(1), group(0) is the whole match
		#        \(\S+, \S+\) one or more no whitespaces characters then "," then again, inside brackets
		#        :* zero or more ":"
		#        \s+ one or more whitespaces
		#        a\s= a followed by one whitespace
		#...
		# This is to find the id line of the binaries/pairs
		self.bin_ids_search = re.compile('^\s*(U*)\s*(\(\S+,\S+\)):*\s+a\s=\s(\d+\.*\S*)\s+e\s=\s(\d+\.*\S*)\s+P\s=\s(\d+\.*\S*)')
		self.container_list = []
		self.hardflag = "--"
		return
	
	
	def search_for_binaries(self, line):	
		# Initialize the container
		container = {
			"system_time": None,
			"phys_time": None,
			"ids": "--",
			"hardflag": "--",
			"objects": [],
			"masses": [],
			"sma": "--",
			"period": "--",
			"ecc": "--"#,
			}
		if True:	# to avoid re-indent all the function
			line = line.translate(None, "\n") # remove newline character
			#n_lines = 0
			if ("Binaries/multiples:" in line or "Bound nn pairs:" in line): 
				# Separate the two categories
				if "Binaries/multiples:" in line:
					self.hardflag = "H"
					print "H"
				if "Bound nn pairs:" in line:
					self.hardflag = "S"
					print "S"
				print "Found binary or bound section at time_step ", self.time_step
				while True: # inside interesting section loop	
					# Read line with EOF check
					line = self.infile.readline()
					if len(line) == 0: # check for EOF
						print "EOF"
						break
					self.n_lines += 1
					line = line.translate(None, "\n") # remove newline character
					# End of interesting sections
					if "Total binary energy" in line or "user_diag:" in line: 
						break# exit the inner loop because the interesting section is finished	
					# Search for objects IDs
					bin_ids = self.bin_ids_search.search(line)
					if bin_ids: 
						# Find timestep
						#print "Found binary ids ", bin_ids.group(2), " at timestep ", time_step
						container["system_time"] = int(self.time_step)
						container["phys_time"] = container["system_time"] * float(self.time_unit[0])
						#print "Found binary in binaries section ", bin_ids.group(2)
						container["ids"] = bin_ids.group(2).translate(None, "()").split(",")
						container["hardflag"] = self.hardflag
						container["sma"] = float(bin_ids.group(3)) * self.r_unit[0]
						container["ecc"] = bin_ids.group(4)
						container["period"] = float(bin_ids.group(5)) * self.time_unit[0]
						bin_ids = None# reset ids variable
					if "masses" in line:
						mass_line = line.split()[1:-3]
						for mass in mass_line:
							container["masses"].append(float(mass)*self.mass_unit[0])	
						self.container_list.append(container)
						# Reinit container
						container = {
						"system_time": None,
						"phys_time": None,
						"ids": "--",
						"objects": [],
						"masses": [],
						"sma": "--",
						"period": "--",
						"ecc": "--"#,
						}
			else:
				pass
		return 

	def store_interesting_binaries(self):
		"""Store interesting binaries, that are that found in the NS or BH lists
		"""
		print "Inside function to store interesting binaries"
		print "Init containers"
		containers_out = []
		#containers_out_nn_int = []
		print "Init lists"
		bh_list_id = [item["id"] for item in self.bh_list]
		ns_list_id = [item["id"] for item in self.ns_list]
		bh_list_form_time = [item["time"] for item in self.bh_list]
		ns_list_form_time = [item["time"] for item in self.ns_list]
		# If we found ids we now search for types
		#print "Main loop on the container list "
		for container in self.container_list:
			interesting = False
			container["objects"] = [] # be sure we don't have any residual
			#print "Loop on the container "
			for i in container["ids"]:
				if i in bh_list_id:
					if float(bh_list_form_time[bh_list_id.index(i)]) < container["phys_time"]:
						container["objects"].append("bh") # if obj is a bh, append
						interesting = True
					else:
						container["objects"].append("bh++")
						interesting = True
				elif i in ns_list_id:
					if float(ns_list_form_time[ns_list_id.index(i)]) < container["phys_time"]:
						container["objects"].append("ns") # if obj is a ns, append 
						interesting = True
					else:
						container["objects"].append("ns++")
						interesting = True
				else:
					# nothing interesting:P
					container["objects"].append("--") # obj not a bh nor a ns
					#interesting = False # Here it is wrong, it may set False if the second id is not
					# interesting but the first is
			
			if interesting == True:
				#print "True"
				containers_out.append(container)
			elif interesting == False:
				#print "False"
				#containers_out_nn_int.append(container)
				pass
			else:
				print "ERROR with the interesting value ", interesting
				print "Exit!!!"
				sys.exit()

		tt1 = time.time()
		# Avoid appending duplicated items (in case of duplicated timesteps due to crashes)
		print "Check iteresting binaries for duplication"
		containers_out_single = []
		for item in containers_out:
			if item not in containers_out_single:
				containers_out_single.append(item)
		print "Done in ", time.time()-tt1
		
		#tt2 = time.time()
		#print "Check uniteresting binaries for duplication"
		#containers_out_nn_int_single = []
		#for item in containers_out_nn_int:
			#if item not in containers_out_nn_int_single:
				#containers_out_nn_int_single.append(item)
		#print "Done in ", time.time()-tt2
		
		out_list = os.path.join(self.outpath, "fished_"+self.fname+".txt")
		out = open(out_list, "w")
		print "Write interesting binaries in ", out_list
		out.write("#") # initialize file
		out.write("{0:<14}".format("system_time")+
				"{0:<20}".format("phys_time ["+self.time_unit[1]+"]")+
				"{0:<50}".format("ids")+
				"{0:<14}".format("hardflag")+
				"{0:<35}".format("objects")+
				"{0:<50}".format("masses ["+self.mass_unit[1]+"]")+
				"{0:<20}".format("sma ["+self.r_unit[1]+"]")+
				"{0:<20}".format("period ["+self.time_unit[1]+"]")+
				"{0:<15}".format("ecc")#+
				)
		out.write("\n")
		out.flush()	
		for container in containers_out_single:
			self.flush_err_data(out, container) # write last data of current section before exit and reset the container!!!
			#print container
		out.flush()
		out.close()	
		# Non interesting binaries output
		#out_list = os.path.join(self.outpath, "uninteresting_fished_"+self.fname+".txt")
		#out = open(out_list, "w")
		#out.write("#") # initialize file
		#out.write("{0:<14}".format("system_time")+
				#"{0:<20}".format("phys_time ["+self.time_unit[1]+"]")+
				#"{0:<50}".format("ids")+
				#"{0:<35}".format("objects")+
				#"{0:<50}".format("masses ["+self.mass_unit[1]+"]")+
				#"{0:<20}".format("sma ["+self.r_unit[1]+"]")+
				#"{0:<20}".format("period ["+self.time_unit[1]+"]")+
				#"{0:<15}".format("ecc")#+
				#)
		#out.write("\n")
		#out.flush()	
		#for container in containers_out_nn_int:
			#self.flush_err_data(out, container) # write last data of current section before exit and reset the container!!!
		#out.flush()
		#out.close()	
		return

	def final_messages(self):
		print self.separator
		print "Total number of lines parsed ", self.n_lines
		print "Duplicated time_steps ", self.duplicated_times
