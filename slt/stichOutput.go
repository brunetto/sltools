package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

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
	
	outRegResult = outRegexp.FindStringSubmatch(inFileName); 
	if outRegResult == nil {
		log.Fatal("Can't find parameters in out name ", inFileName)
	}
	
	run  = outRegResult[1]
	
	if inFileName == "" {
		log.Fatal("You need to specify an input file template with the -i flag!!!")
	}
	
	log.Println("Stiching *-" + conf.BaseName() + `-run` + run + `-rnd*.txt`)
	
	//
	// STDOUT
	//
	stdOuts = "out-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdOuts, run, "out", conf)
	
	//
	// STDERR
	//
	stdErrs = "err-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdErrs, run, "err", conf)
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich output ", tGlob1.Sub(tGlob0))
}

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
			fmt.Println(idx, file)
		}
	}
	
	for _, inFileName := range inFiles {
		if Verb {
			log.Println("Working on ", inFileName)
		}
		if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
		defer inFile.Close()
		ext = filepath.Ext(inFileName)
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
		
		for {
			if stdWhat == "out" {
				// Why continue??? 
				// FIXME, TODO: However, reading 2 snapshots at a time avoid 
				// problems in comparing the timesteps the first time
// 				if snapshot, err = ReadOutSnapshot(nReader); err != nil {continue}
				if snapshot/*s[0]*/, err = ReadOutSnapshot(nReader); err != nil {break}
				/*if snapshots[1], err = ReadOutSnapshot(nReader); err != nil {break}*/
			} else if stdWhat == "err" {
// 				if snapshot, err = ReadErrSnapshot(nReader); err != nil {continue}
				if snapshot/*s[0]*/, err = ReadErrSnapshot(nReader); err != nil {break}
// 				if snapshots[1], err = ReadErrSnapshot(nReader); err != nil {break}
			} else {
				log.Fatal("Unrecognized stdWhat: ", stdWhat)
			}
			if snapshot.Integrity == true {
				timestep, err = strconv.ParseInt(snapshot.Timestep, 10, 64)
				// Skip the first loop (=first timestep) with len = 0
				if len(timesteps) > 0 {
					if AbsInt(timestep - timesteps[len(timesteps)-1]) != 1 {
						log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
					}
				}
				timesteps = append(timesteps, timestep)
				if err = snapshot.WriteSnapshot(nWriter); err != nil {
					log.Fatal("Error while writing snapshot to file: ", err)
				}
			} else {
				fmt.Println("************************ ATTENTION *************************")
				fmt.Println("************************************************************")
				log.Println("Skipping incomplete snapshot at timestep", snapshot.Timestep)
				fmt.Println("************************************************************")
				fmt.Println("************************************************************")
			}
		}	
	}
	
	log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
	fmt.Println(timesteps)
		
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich STDOUT output ", tGlob1.Sub(tGlob0))
}
