package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func StdErrStich (stdErrs, run string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		fZip *gzip.Reader
		inFile *os.File
		inFiles []string
		outFileName string
		outFile *os.File
		err error
		nReader *bufio.Reader
		nWriter *bufio.Writer
		errSnapshot *ErrSnapshot
		timestep int64
		timesteps = make([]int64, 0)
		ext string
	)
	
	tGlob0 := time.Now()
	
	//
	// STDERR
	//
	
	log.Println("Stich STDERR")
	outFileName = "err" + conf.BaseName() + "-run" + run + "-all.txt"	
	log.Println("Output file will be ", outFileName)
	
	log.Println("Opening STDERR output file...")

	// Open output file
	if outFile, err = os.Create(outFileName); err != nil {log.Fatal(err)}
	defer outFile.Close()
	
	// Create reader and writerq
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	log.Println("Globbing and sorting STDERR input files")
	// Open infiles
	if inFiles, err = filepath.Glob(stdErrs); err != nil {
		log.Fatal("Error globbing STDERR files for output stiching: ", err)
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
			if errSnapshot, err = ReadErrSnapshot(nReader); err != nil {continue}
			if errSnapshot.Integrity == true {
				timestep, err = strconv.ParseInt(errSnapshot.Timestep, 10, 64)
				if timestep - timesteps[len(timesteps)-1] != 1 {
					log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
				}
				timesteps = append(timesteps, timestep)
				if err = errSnapshot.WriteSnapshot(nWriter); err != nil {
					log.Fatal("Error while writing snapshot to file: ", err)
				}
			} else {
				log.Println("Skipping incomplete snapshot at timestep", errSnapshot.Timestep)
			}
		}	
	}
	
	log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
	fmt.Println(timesteps)
	
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich STDERR output ", tGlob1.Sub(tGlob0))
}