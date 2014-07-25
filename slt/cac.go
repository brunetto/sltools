package slt

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)

// Check And Continue
func CAC() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		err, mapErr error
		globName    string = "*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt"
		runMap      map[string]map[string][]string
		// for example runMap["08"]["err"][3]
		// will give ["err-....run08-rnd03.txt"]
		nProcs           int = 1
		inFileNameChan       = make(chan string, 1)
		cssInfo              = make(chan map[string]string, 1)
		done                 = make(chan struct{})
		runs             []string
		run              string
		lastErr, lastOut string
		errInfo, outInfo os.FileInfo
		toRemove         = []string{}
		machine          string
		machineDiscovery *exec.Cmd
		stdo             bytes.Buffer
	)

	log.Println("Try to discover machine name")
	machineDiscovery = exec.Command("hostname", "-A")
	if machineDiscovery.Stdout = &stdo; err != nil {
		log.Fatal("Error connecting STDOUT: ", err)
	}

	if err = machineDiscovery.Start(); err != nil {
		log.Fatal("Error starting machineDiscovery: ", err)
	}

	if err = machineDiscovery.Wait(); err != nil {
		log.Fatal("Error while waiting for machineDiscovery: ", err)
	}

	switch {
	case strings.Contains(stdo.String(), "eurora"):
		machine = "eurora"
	case strings.Contains(stdo.String(), "spritz"):
		machine = "spritz"
	case strings.Contains(stdo.String(), "longisland"):
		machine = "longisland"
	case strings.Contains(stdo.String(), "sfursat"):
		machine = "sfursat"
	case strings.Contains(stdo.String(), "plx"):
		machine = "plx"
	case strings.Contains(stdo.String(), "auriga"):
		machine = "auriga"
	case strings.Contains(stdo.String(), "tesla"):
		machine = "tesla"
	}

	log.Println("machine set to: ", machine)

	log.Println("Starting goroutines...")
	for idx := 0; idx < nProcs; idx++ {
		go Out2ICs(inFileNameChan, cssInfo)
		go CreateStartScripts(cssInfo, machine, done)
	}

	log.Println("Searching for files in the form: ", globName)
	// Find last round for each run in the folder
	// Runs are sorted
	runs, runMap, mapErr = FindLastRound(globName)
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

		outSize, outUnit := SizeUnit(outInfo.Size())
		errSize, errUnit := SizeUnit(errInfo.Size())

		fmt.Printf("%v\t%2.2f %v %v\n\t%2.2f %v %v\n",
			run, outSize, outUnit, lastOut,
			errSize, errUnit, lastErr)

		// Removing wrong files
		if outUnit == "bytes" || errUnit == "GB" {
			toRemove = append(toRemove, lastErr, lastOut)
			log.Printf("\tRemove because of suspicious dimensions (probably broken):\n\t%v\n\t%v\n ", lastOut, lastErr)
			for _, file := range []string{lastOut, lastErr} {
				if err = os.Remove(file); err != nil {
					log.Fatal("Error while removing ", file, ": ", err)
				}
			}
		} else {
			inFileNameChan <- runMap[run]["out"][len(runMap[run]["out"])-1]
			fmt.Println()
		}
		fmt.Println(".................................")

	}
	// Close the channel, if you forget it, goroutines
	// will wait forever
	close(inFileNameChan)

	// Wait the CreateStartScripts goroutines to finish
	for idx := 0; idx < nProcs; idx++ {
		<-done // wait the goroutine to finish
	}
	fmt.Println()
}
