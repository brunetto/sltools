package slt

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)

// Continue provide a lazy function to prepare a simulation for he next run.
// It will convert the last snapshot from a StarLab stdout file in an ICs file
// by calling Out2ICs and then it will create the needed scripts for launching
// the simulation (kiraLaunch and PBSlaunch) calling CreateStartScripts.
// It needs a valid configuration file.





func Continue(inFileName, machine string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		nProcs int = 1
		inFileNameChan = make(chan string, 1)
		cssInfo = make(chan map[string]string, 1)
		done = make(chan struct{})
	)
	
	for idx:=0; idx<nProcs; idx++ {
		go Out2ICs(inFileNameChan, cssInfo)
		go CreateStartScripts(cssInfo, machine, done)
	}
	
	// Check if we have to run on all the files in the folder 
	// and not only on a selected one 
	if inFileName == "all" || inFileName == "*" || 
		inFileName == "" || strings.Contains(inFileName, "*") {
		runs, runMap, mapErr := FindLastRound("*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt")
		log.Println("Selected to continue round for all the runs in the folder")
		log.Println("Found: ")
		for _, run := range runs {
			if mapErr != nil && (len(runMap[run]["err"]) == 0 || len(runMap[run]["out"]) == 0) {
				continue
			}
			fmt.Printf("%v\n", runMap[run]["out"][len(runMap[run]["out"])-1])
		}
		fmt.Println()
		// Fill the channel with the last round of each run
		for _, run := range runs {
			if mapErr != nil && (len(runMap[run]["err"]) == 0 || len(runMap[run]["out"]) == 0) {
				continue
			}
			inFileNameChan <- runMap[run]["out"][len(runMap[run]["out"])-1]
		}
	} else {
		// Only continue the selected file
		inFileNameChan <- inFileName
	}
	
	// Close the channel, if you forget it, goroutines 
	// will wait forever
	close(inFileNameChan)
	
	// Wait the goroutines to finish
	for idx:=0; idx<nProcs; idx++ {
		<-done // wait the goroutine to finish
	}
}

