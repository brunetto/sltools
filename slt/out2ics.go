package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Out2ICs (inFileName string) (string, string, string) {
	if Debug {Whoami(true)}
	
	var (
		outFileName string
		inFile *os.File
		fZip *gzip.Reader
		outFile *os.File
		err error
		nReader *bufio.Reader
		nWriter *bufio.Writer
		snapshots = make([]*OutSnapshot, 2)
		snpN int
		simulationStop int64 = 500
		thisTimestep int64 = 0
		randomSeed string
		remainingTime int64
		ext string
	)
	
	tGlob0 := time.Now()
	
	if inFileName == "" {
		log.Fatal("You need to specify an input file with the -i flag!!!")
	}
	
	// Open infile, both text or gzip
	log.Println("Opening input and output files...")
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
			log.Fatal("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
		}
		default: {
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extention ", ext)
		}
	}

	// outFile name
	outFileName = OutName2ICName (inFileName)	
	log.Println("Output file will be ", outFileName)
	
	// Open outFile and outWriter
	if outFile, err = os.Create(outFileName); err != nil {log.Fatal(err)}
	defer outFile.Close()
	
	// Create writer
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
	fmt.Println() // To leave a space after the non verbose print
	log.Println("Done reading, last complete timestep is ", snapshots[snpN].Timestep)
	thisTimestep, _ = strconv.ParseInt(snapshots[snpN].Timestep, 10, 64)
	remainingTime = simulationStop - thisTimestep
	log.Println("Set -t flag to ", remainingTime)
	
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
	log.Println("Wall time for creating ICs from STDOUT ", tGlob1.Sub(tGlob0))
	
	return strconv.Itoa(int(remainingTime)), randomSeed, outFileName
}