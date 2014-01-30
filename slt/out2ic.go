package slt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func Out2IC (inFileName string, fileN string) () {
	// FIXME: generate ICs with templates
	// http://golang.org/pkg/text/template/
	// filepath.Glob(pattern string) (matches []string, err error)
	
	var (
		outFileName string
		inFile *os.File
		outFile *os.File
		err error
		nReader *bufio.Reader
		nWriter *bufio.Writer
		snapshots = make([]*OutSnapshot, 2)
		snpN int
		simulationStop int64 = 500
		thisTimestep int64 = 0
		randomSeed string
	)
	
	tGlob0 := time.Now()
	
	if inFileName == "" {
		log.Fatal("You need to specify an input file with the -i flag!!!")
	}
	
	if fileN == "" {
		log.Fatal("You need to specify a number for the new ICs with the -n flag!!!")
	}
	
	outFileName = OutName2ICName (inFileName, fileN)	
	log.Println("Output file will be ", outFileName)
	
	log.Println("Opening input and output files...")
	// FIXME: take into account .gz files
	
	// Open files
	if inFile, err = os.Open(inFileName); err != nil {panic(err)}
	defer inFile.Close()
	
	if outFile, err = os.Create(outFileName); err != nil {panic(err)}
	defer outFile.Close()
	
	// Create reader and writerq
	nReader = bufio.NewReader(inFile)
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	log.Println("Start reading...")
	// Read two snapshot each loop to ensure at least one of them is complete
	// (= I keep the previous read in memory in case the last is corrupted)
	for {
		if snapshots[0], err = ReadOutSnapshot(nReader); err != nil {break}
		if snapshots[1], err = ReadOutSnapshot(nReader); err != nil {break}
	}
	
	// Check integrity once the file reading is ended
	// First the last read, then the previous one
	if snapshots[1].Integrity == true {
		snpN = 1
	} else if snapshots[0].Integrity == true {
		snpN = 0
	} else {
		log.Println("Both last two snapshots corrupted on file ", inFileName)
		fmt.Println("Snapshot ", snapshots[1].Timestep, " is ", snapshots[1].Integrity)
		fmt.Println("Snapshot ", snapshots[0].Timestep, " is ", snapshots[0].Integrity)
		log.Fatal("Reading exit with error ", err)
	}
	// Info
	log.Println("Done reading, last complete timestep is ", snapshots[snpN].Timestep)
	thisTimestep, _ = strconv.ParseInt(snapshots[snpN].Timestep, 10, 64)
	log.Println("Set -t flag to ", simulationStop - thisTimestep)
	
	// Write last complete snapshot to file
	log.Println("Writing snapshot to ", outFileName)
	if err = snapshots[snpN].WriteSnapshot(nWriter); err != nil {
		log.Fatal("Error while writing snapshot to file: ", err)
	}
	
	log.Println("Search for random seed...")
	randomSeed = DetectRandomSeed(inFileName)
	log.Println("Set -s flag to ", randomSeed)
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for continue ", tGlob1.Sub(tGlob0))
}