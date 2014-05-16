package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"bufio"

	"bitbucket.org/brunetto/sltools/slt"

	"github.com/brunetto/goutils/debug"
)

func main() {
	defer debug.TimeMe(time.Now())

	var (
		err                            error
		inFileName, newICsFileName     string
		inFile, newICsFile             *os.File
		nReader                        *bufio.Reader
		nWriter                        *bufio.Writer
		regString                      string         = `(\w{3})-(\S*-rnd)(\d*)(\.\S*)`
		regExp                         *regexp.Regexp = regexp.MustCompile(regString)
		regResult                      []string
		fileNameBody, newRnd, ext      string
		snapshots                      = make([]*slt.DumbSnapshot, 2)
		snpN                           int
		simulationStop                 int64 = 500
		thisTimestep, remainingTime    int64
		randomSeed                     string
		runString                      string
		newErrFileName, newOutFileName string
	)

	if len(os.Args) < 2 {
		log.Fatal("You need to provide a STDOUT file to parse!!")
	}

	inFileName = os.Args[1]

	if regResult = regExp.FindStringSubmatch(inFileName); regResult == nil {
		log.Fatal("Can't find fileNameBody in ", inFileName)
	}
	if regResult[1] != "out" {
		log.Fatal("Please specify a STDOUT file, found ", regResult[1])
	}

	fileNameBody = regResult[2]
	ext = regResult[4]

	temp, _ := strconv.ParseInt(regResult[3], 10, 64)
	newRnd = strconv.Itoa(int(temp + 1))

	newICsFileName = "ics-" + fileNameBody + slt.LeftPad(newRnd, "0", 2) + ext
	newErrFileName = "err-" + fileNameBody + slt.LeftPad(newRnd, "0", 2) + ext
	newOutFileName = "out-" + fileNameBody + slt.LeftPad(newRnd, "0", 2) + ext

	log.Println("Creating files")
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	nReader = bufio.NewReader(inFile)

	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	if newICsFile, err = os.Create(newICsFileName); err != nil {
		log.Fatal(err)
	}
	defer newICsFile.Close()
	nWriter = bufio.NewWriter(newICsFile)
	defer nWriter.Flush()

	log.Println("Start reading...")
	// Read two snapshot each loop to ensure at least one of them is complete
	// (= I keep the previous read in memory in case the last is corrupted)
	for {
		if snapshots[0], err = slt.ReadOutSnapshot(nReader); err != nil {
			break
		}
		if snapshots[1], err = slt.ReadOutSnapshot(nReader); err != nil {
			break
		}
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
	log.Println("Writing snapshot to ", newICsFileName)
	if err = snapshots[snpN].WriteSnapshot(nWriter); err != nil {
		log.Fatal("Error while writing snapshot to file: ", err)
	}

	fmt.Fprint(os.Stderr, "\n")
	log.Println("Search for random seed...")
	randomSeed = slt.DetectRandomSeed(inFileName)
	log.Println("Set -s flag to ", randomSeed)

	runString = "\nRun with:\n" +
		"----------------------\n" +
		"(/home/ziosi/Code/Mapelli/slpack/starlab/usr/bin/kira -F -t " +
		strconv.Itoa(int(remainingTime)) +
		" -d 1 -D 1 -b 1 -f 0 " +
		"-n 10 -e 0.000 -B -s " + randomSeed +
		" < " + newICsFileName + " >  " + newOutFileName + " 2> " + newErrFileName + ")& \n" +
		"----------------------\n\n" +
		"You can watch the status of the simulation by running: \n" + 
		"----------------------\n" +
		"watch stat " + newErrFileName + "\n" +
		"cat " + newErrFileName + ` | grep "Time = " | tail -n 1` + "\n" +
		"----------------------\n\n"
	
	
		
	fmt.Println(runString)

	fmt.Println()
}
