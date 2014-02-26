package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sort"
	"strconv"
	"time"
)

// StichThemAll launch the stiching in parallel on all the simulation files in 
// the folder, accordingly to their names (run 01 is different from run 02 and so on).
func StichThemAll (conf *ConfigStruct) () {
	if Debug {Whoami(true)}
		
	var (
		err error
		wg sync.WaitGroup // to wait the end of the goroutines
		inFiles []string
		prefixes = []string{"out-", "err-"}
		outRegexp *regexp.Regexp = regexp.MustCompile(`\S` + conf.BaseName() + `-run(\d+)-rnd\d+.txt`)
		outRegResult []string
		runs StringSet // set = list of unique objects (run numbers)
		nRuns []int
	)
	
	tGlob0 := time.Now()
	
	nRuns = make([]int, 0)
	
	// Search for all the STDOUT and STDERR files in the folder
	for idx:=0; idx<2; idx++ {
		if inFiles, err = filepath.Glob(prefixes[idx]+conf.BaseName()+"*"); err != nil {
			log.Fatal("Error globbing for stiching all the run outputs in this folder: ", err)
		}
		// Sort file names
		sort.Strings(inFiles)
		
		runs = NewStringSet()
		
		// Find the numbers of the different runs
		for _, inFileName := range inFiles {
			outRegResult = outRegexp.FindStringSubmatch(inFileName) 
			if outRegResult == nil {
				log.Fatal("Can't find parameters in out name ", inFileName)
			}
			// Add the new number in the set
			runs.Add(outRegResult[1])
		}
		if Verb {
			log.Println("Found runs:")
			fmt.Println(runs.String())
		}
		nRuns = append(nRuns, len(runs))
	}
	
	// Check for missing run outputs
	if nRuns[0] != nRuns[1] {
			log.Println("WARNING, missing runs. Found ", nRuns[0], " STDOUT but ", nRuns[1], " STDERR.")
	}
	if nRuns[0] != conf.Runs {
			log.Println("WARNING, missing runs. Found ", nRuns[0], " STDOUT of ", conf.Runs, " planned in config file.")
	}
	
	log.Println("Found ", nRuns[0], " runs in this folder:")
	fmt.Println(runs.String())
	
	// Launch all the stiching
	for _, runIdx := range runs.Sorted() {
		name := "out-"+conf.BaseName()+"-run"+runIdx+"-rnd01.txt"
		if Verb {
			log.Println("Launching stich based on ", name)
		}
		// Count the goroutine as running
		wg.Add(1)
		// Launch one stiching
		go func(name string, conf *ConfigStruct) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			StichOutput (name, conf)
		}(name, conf)
	}
	
	// Wait for all the goroutine to finish
	wg.Wait()
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for parallel stich all ", tGlob1.Sub(tGlob0))
}

// FIXME: Workaround to call StichOutput not in parallel 
// because now StichOutput contain a call to wg.Done
// and I don't want to import "sync" in command.go
func StichOutputSingle (inFileName string, conf *ConfigStruct) () {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(name string, conf *ConfigStruct) {
		// Decrement the counter when the goroutine completes.
		defer wg.Done()
		StichOutput (name, conf)
		}(inFileName, conf)
	wg.Wait()
}
	
// StichOutput stiches the STDOUT and STDERR of a simulation.
func StichOutput (inFileName string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		outRegexp *regexp.Regexp = regexp.MustCompile(`\S` + conf.BaseName() + `-run(\d+)-rnd\d+.txt`)
		outRegResult []string
		run string
		stdOuts string
		stdErrs string
	)
	
	tGlob0 := time.Now()
	
	// Extract parameters from the name
	outRegResult = outRegexp.FindStringSubmatch(inFileName); 
	if outRegResult == nil {
		log.Fatal("Can't find parameters in out name ", inFileName)
	}
	
	run  = outRegResult[1]
	
	if inFileName == "" {
		log.Fatal("You need to specify an input file template with the -i flag!!!")
	}
	
	log.Println("Stiching *-" + conf.BaseName() + `-run` + run + `-rnd*.txt`)
	
	// Check if only have to run STDERR stich
	if !OnlyErr {
	//
	// STDOUT
	//
	stdOuts = "out-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdOuts, run, "out", conf)
	} else {
		log.Println("Only stich STDERRs")
	}
	
	// Check if only have to run STDOUT stich
	if !OnlyOut {
	//
	// STDERR
	//
	stdErrs = "err-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdErrs, run, "err", conf)
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich output ", tGlob1.Sub(tGlob0))
	} else {
		log.Println("Only stich STDOUTs")
	}
}

// StdStich stiches a given STD??? according to the type passed with stdWhat.
func StdStich (stdFiles, run, stdWhat string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		fZip *gzip.Reader
		inFile *os.File
		snapshot/*s = make([]*/*DumbSnapshot/*, 2)*/
		inFiles []string
		outFileName string
		outFile *os.File
		err error
		nReader *bufio.Reader
		nWriter *bufio.Writer
		timestep int64
		timesteps = make([]int64, 0)
		ext string
	)
	
	tGlob0 := time.Now()
	
	log.Println("Stich std"+stdWhat)
	outFileName = stdWhat + "-" + conf.BaseName() + "-run" + run + "-all.txt"	
	log.Println("Output file will be ", outFileName)
	
	log.Println("Opening STDOUT output file...")

	// Open output file
	if outFile, err = os.Create(outFileName); err != nil {log.Fatal(err)}
	defer outFile.Close()
	
	// Create reader and writerq
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	log.Println("Globbing and sorting " + stdWhat + " input files")
	// Open infiles
	if inFiles, err = filepath.Glob(stdFiles); err != nil {
		log.Fatal("Error globbing " + stdWhat + " files for output stiching: ", err)
	}

	sort.Strings(inFiles)
	
	if Verb {
		log.Println("Found:")
		for idx, file := range inFiles {
			fmt.Println(idx, ": ", file)
		}
	}
	
	for _, inFileName := range inFiles {
		if Verb {
			log.Println("Working on ", inFileName)
		}
		if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
		defer inFile.Close()
		ext = filepath.Ext(inFileName)
		// Try to open the file if it is a gzipped one or a simple txt
		switch ext {
			case ".txt":{
				nReader = bufio.NewReader(inFile)
			}
			case ".gz": {
				fZip, err = gzip.NewReader(inFile)
				if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFileName, err)
				os.Exit(1)
				}
				nReader = bufio.NewReader(fZip)
			}
			default: {
				log.Fatal("Unrecognized file type", inFileName)
			}
		}
		
		//Read snapshots and write them if everything is OK
		SnapLoop: for {
			if stdWhat == "out" {
				snapshot, err = ReadOutSnapshot(nReader)
			} else if stdWhat == "err" {
				snapshot, err = ReadErrSnapshot(nReader)
			} else {
				log.Fatal("Unrecognized stdWhat: ", stdWhat)
			}
			if err != nil {
				if Verb {
					log.Println("Incomplete snapshot, moving to the next file")
				}
				break SnapLoop
			}
			// -1 is the "ICs to 0" timestep, skipping
			// I will skip this also because it creates problems of duplication 
			// and timestep check
			if snapshot.Timestep == "-1" {continue SnapLoop /*to the next timestep*/}
			
			// I will loose the last timestep on STDERR because it is probably not complete
			// TODO: find out how to manage this
			// BUG: I can't find a univoque way to define the last snapshot complete
			if snapshot.Integrity == true {
				timestep, err = strconv.ParseInt(snapshot.Timestep, 10, 64)
				// Skip the first loop (=first timestep) with len = 0
				if len(timesteps) > 0 {
					if AbsInt(timestep - timesteps[len(timesteps)-1]) > 1 {
						if Verb {
							log.Println("Read timestep: ")
							for _, ts := range timesteps {
								fmt.Print(ts, " ")
							}
							fmt.Println()
						}
						log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
					} else if AbsInt(timestep - timesteps[len(timesteps)-1]) < 1 {
						log.Println("Duplicated timestep ", timestep, ", continue.")
						continue SnapLoop /*to the next timestep*/
					}
				}
				timesteps = append(timesteps, timestep)// Write the snapshot
				if err = snapshot.WriteSnapshot(nWriter); err != nil {
					log.Fatal("Error while writing snapshot to file: ", err)
				}
			} else { 
				// This shouldn't happend because of the break in reading the snapshots
				// This shoud be a redundant check
				// TODO: check if it is true!!!
				fmt.Println("************************ ATTENTION *************************")
				fmt.Println("************************************************************")
				log.Println("Skipping incomplete snapshot at timestep", snapshot.Timestep)
				fmt.Println("************************************************************")
				fmt.Println("************************************************************")
			}
		} // end reading snapshot from a single file loop
	} // end reading file loop
	fmt.Println("\n")
	log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
	fmt.Println(timesteps)
		
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich STDOUT output ", tGlob1.Sub(tGlob0))
}
