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

func StdOutStich (inFileTmpl string) () {
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
		outSnapshot *OutSnapshot
		timestep int64
		timesteps = make([]int64, 0)
		ext string
	)
	
	tGlob0 := time.Now()
	
	if inFileTmpl == "" {
		log.Fatal("You need to specify a STDOUT input file template with the -i flag!!!")
	}
	
	log.Println("Stich STDOUT")
	outFileName = inFileTmpl + "-all.txt"	
	log.Println("Output file will be ", outFileName)
	
	log.Println("Opening STDOUT output file...")

	// Open output file
	if outFile, err = os.Create(outFileName); err != nil {panic(err)}
	defer outFile.Close()
	
	// Create reader and writerq
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	log.Println("Globbing and sorting STDOUT input files")
	// Open infiles
	if inFiles, err = filepath.Glob(inFileTmpl); err != nil {
		log.Fatal("Error globbing STDOUT files for output stiching: ", err)
	}
	// FIXME: take into account .gz files
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
			if outSnapshot, err = ReadOutSnapshot(nReader); err != nil {continue}
			if outSnapshot.Integrity == true {
				timestep, err = strconv.ParseInt(outSnapshot.Timestep, 10, 64)
				if timestep - timesteps[len(timesteps)-1] != 1 {
					log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
				}
				timesteps = append(timesteps, timestep)
				if err = outSnapshot.WriteSnapshot(nWriter); err != nil {
					log.Fatal("Error while writing snapshot to file: ", err)
				}
			} else {
				log.Println("Skipping incomplete snapshot at timestep", outSnapshot.Timestep)
			}
		}	
	}
	
	log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
	fmt.Println(timesteps)
		
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich STDOUT output ", tGlob1.Sub(tGlob0))
}