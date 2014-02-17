package slt

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

func StichOutput (inFileName string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		outRegString string = `\S` + conf.BaseName() + `-run(\d+)-rnd\d+.txt`
		outRegexp *regexp.Regexp = regexp.MustCompile(outRegString)
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
	
	//
	// STDOUT
	//
	stdOuts = "out-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdOuts, run, "out")
	
	//
	// STDERR
	//
	stdErrs = "err-" + conf.BaseName() + `-run` + run + `-rnd*.txt`
	StdStich (stdErrs, run, "err")
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich output ", tGlob1.Sub(tGlob0))
}

func StdStich (stdFiles, run, stdWhat string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		fZip *gzip.Reader
		inFile *os.File
		var snapshot *DumbSnapshot
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
	outFile, err = os.Create(outFileName); err != nil {log.Fatal(err)}
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
	
	for _, inFileName := range inFiles {
		if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
		defer inFile.Close()
		ext = filepath.Ext(inFileName)
		switch ext {
			case "txt":{
				nReader = bufio.NewReader(inFile)
			}
			case "gz": {
				fZip, err = gzip.NewReader(inFile)
				if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFile, err)
				os.Exit(1)
				}
				nReader = bufio.NewReader(fZip)
			}
			default: {
				log.Fatal("Unrecognized file type", inFile)
			}
		}
		
		for {
			if stdWhat == "out" {
				if snapshot, err = ReadOutSnapshot(nReader); err != nil {continue}
			} else if stdWhat == "err" {
				if snapshot, err = ReadErrSnapshot(nReader); err != nil {continue}
			} else {
				log.Fatal("Unrecognized stdWhat: ", stdWhat)
			}
			if snapshot.Integrity == true {
				timestep, err = strconv.ParseInt(snapshot.Timestep, 10, 64)
				if timestep - timesteps[len(timesteps)-1] != 1 {
					log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
				}
				timesteps = append(timesteps, timestep)
				if err = snapshot.WriteSnapshot(nWriter); err != nil {
					log.Fatal("Error while writing snapshot to file: ", err)
				}
			} else {
				log.Println("Skipping incomplete snapshot at timestep", snapshot.Timestep)
			}
		}	
	}
	
	log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
	fmt.Println(timesteps)
		
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich STDOUT output ", tGlob1.Sub(tGlob0))
}