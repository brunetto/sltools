package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		err, mapErr          error
		globName string = "*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt"
		runMap map[string]map[string][]string
		// for example runMap["08"]["err"][3] 
		// will give ["err-....run08-rnd03.txt"]
		
		runs []string
		run string
		lastErr, lastOut string
		errInfo, outInfo os.FileInfo
		toRemove = []string{}
	)
	
	log.Println("Searching for files in the form: ", globName)
	
	// Find last round for each run in the folder
	// Runs are sorted
	runs, runMap, mapErr = slt.FindLastRound(globName)
	// Some round are present because of the ics ma don't have errs or outSize,
	// probably they were run somewhere else (Spritz?)
	if err != nil {
		log.Println(err)
	}
	
	// Loop over the last rounds found and print infos
	fmt.Println(".................................")
	for _, run = range runs {
		// In case we have only ics 
		if mapErr != nil && (len(runMap[run]["err"]) == 0 || len(runMap[run]["out"]) == 0) {
			continue
		}
		lastErr = runMap[run]["err"][len(runMap[run]["err"])-1]
		lastOut = runMap[run]["out"][len(runMap[run]["out"])-1]
		
		// Check files dimension
		if errInfo, err = os.Stat(lastErr); err != nil {
			log.Fatal("Error checking STDERR file size, err")
		}
		if outInfo, err = os.Stat(lastOut); err != nil {
			log.Fatal("Error checking STDOUT file size, err")
		}
		
		outSize, outUnit := sizeUnit(outInfo.Size())
		errSize, errUnit := sizeUnit(errInfo.Size())
		
		if outUnit == "bytes" || errUnit == "GB" {
			toRemove = append(toRemove, lastErr, lastOut)
		}
		
		fmt.Printf("%v\t%2.2f %v %v\n\t%2.2f %v %v\n", 
				   run, outSize, outUnit, lastOut, 
				   errSize, errUnit, lastErr)
		checkSnapshot(lastOut)
		fmt.Println()
		fmt.Println(".................................")	
	}
	// Suggest to delete broken rounds if it is the case
	if len(toRemove) > 0 {
		log.Println("Suggestion, run: ")
		fmt.Printf("rm ")
		for _, fileToDelete := range toRemove {
			fmt.Printf(" %v ", fileToDelete)
		}
	}
	fmt.Println()
	
	fmt.Print("\x07") // Beep when finish!!:D
}

func checkSnapshot (inFileName string) () {
	var (
		err error
		inFile *os.File
		nReader  *bufio.Reader
	)
	
	if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
	defer inFile.Close()
	
	nReader = bufio.NewReader(inFile)
	
	for {
		if _, err = slt.ReadOutSnapshot(nReader); err != nil {break}
	}
}

func sizeUnit(size int64) (outSize float64, unit string) {
	const tokB = float64(1. / (1024))
	const toMB = float64(1. / (1024*1024))
	const toGB = float64(1. / (1024*1024*1024))
	
	switch {
		case size > (1024*1024*1024):
			return float64(size)*toGB, "GB"
		case size > (1024*1024):
			return float64(size)*toMB, "MB"
		case size > 1024:
			return float64(size)*tokB, "kB"
		default:
			return float64(size), "bytes"
	}
}

// printf "\n"; pwd; printf "\n"; for (( c=0; c<=9; c++ )); do printf "$c "; ls -lah out-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; prStintf "  "; ls -lah err-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; printf "  "; cat $(ls err-*-run0$c-rnd0* | tail -n 1) | grep "Time = " | tail -n 1; done
