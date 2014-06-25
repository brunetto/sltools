package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		err          error
		inFiles      []string
// 		idx int 
		fileName string
		globName string = "*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt"
		runMap = map[string]map[string][]string{}
// 		for example runMap["08"]["err"][3] will give ["err-....run08-rnd03.txt"]
		exists bool
		regRes                         map[string]string
		keys []string
		key string
		lastErr, lastOut string
		errInfo, outInfo os.FileInfo
		toRemove = []string{}
	)
	
	log.Println("Searching for files in the form: ", globName)
		
	if inFiles, err = filepath.Glob(globName); err != nil {
		log.Fatal("Error globbing files in this folder: ", err)
	}
	
	for _, fileName = range inFiles {
		// Try to detect file parameters (type, run, rnd) from fileName
		regRes, err = slt.Reg(fileName)
		// Not standard name
		if err != nil {log.Fatal("Can't find proper name to regex in ", fileName)}
		// Check if run is present, if not, create it in the map
		if _, exists = runMap[regRes["run"]]; !exists {
			runMap[regRes["run"]] = map[string][]string{
				"ics": []string{},
				"err": []string{},
				"out": []string{},
			}
		}
		// Fill the map entry with the fileName
		runMap[regRes["run"]][regRes["prefix"]] = append(runMap[regRes["run"]][regRes["prefix"]], fileName)
	}
	
	// Now runMap contains all the fileName 
	keys = make([]string, len(runMap))
	idx := 0 
	for key, value := range runMap {
        keys[idx] = key
        // Sort rounds
        sort.Strings(value["ics"])
		sort.Strings(value["err"])
		sort.Strings(value["out"])
        idx++
    }
    sort.Strings(keys)
	
	fmt.Println(".................................")
	for _, key = range keys {
		// In case we have only ics 
		if len(runMap[key]["err"]) == 0 || len(runMap[key]["out"]) == 0 {
			continue
		}
		lastErr = runMap[key]["err"][len(runMap[key]["err"])-1]
		lastOut = runMap[key]["out"][len(runMap[key]["out"])-1]
		
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
				   key, outSize, outUnit, lastOut, 
				   errSize, errUnit, lastErr)
		checkSnapshot(lastOut)
		fmt.Println()
		fmt.Println(".................................")	
	}
	// Suggest to delete broken rounds
	log.Println("Suggestion, run: ")
	fmt.Printf("rm ")
	for _, fileToDelete := range toRemove {
		fmt.Printf(" %v ", fileToDelete)
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
