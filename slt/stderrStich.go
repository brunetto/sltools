package slt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"strconv"
	"time"
)

func StdErrStich (inFileTmpl string) () {
	if Debug {Whoami(true)}
	// FIXME: generate ICs with templates
	// http://golang.org/pkg/text/template/
	// filepath.Glob(pattern string) (matches []string, err error)
	
	var (
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
	)
	
	tGlob0 := time.Now()
	
	if inFileTmpl == "" {
		log.Fatal("You need to specify a STDOUT input file template with the -i flag!!!")
	}
		
	//
	// STDERR
	//
	
	log.Println("Stich STDERR")
	outFileName = strings.TrimPrefix(inFileTmpl, "n") + "-all.txt"	
	log.Println("Output file will be ", outFileName)
	
	log.Println("Opening STDERR output file...")

	// Open output file
	if outFile, err = os.Create(outFileName); err != nil {panic(err)}
	defer outFile.Close()
	
	// Create reader and writerq
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	log.Println("Globbing and sorting STDERR input files")
	// Open infiles
	if inFiles, err = filepath.Glob(strings.TrimPrefix(inFileTmpl, "n")); err != nil {
		log.Fatal("Error globbing STDERR files for output stiching: ", err)
	}
	// FIXME: take into account .gz files
	sort.Strings(inFiles)
	
	for _, inFileName := range inFiles {
		if inFile, err = os.Open(inFileName); err != nil {panic(err)}
		defer inFile.Close()
		nReader = bufio.NewReader(inFile)
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