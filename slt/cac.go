package slt

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils"
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
		cssInfo0              = make(chan map[string]string, 1)
		cssInfo1              = make(chan map[string]string, 1)
		pbsLaunchChannel0      = make(chan string, 1)
		pbsLaunchChannel1      = make(chan string, 1)
		done                 = make(chan struct{}, nProcs)
		done1                 = make(chan struct{}, 1)
		runs             []string
		run              string
		lastErr, lastOut string
		errInfo, outInfo os.FileInfo
		toRemove         = []string{}
		machine          string
		machineDiscovery *exec.Cmd
		stdo             bytes.Buffer
		removedFileName string = "Removed.txt"
		removedFile *os.File
		tmp map[string]string
		tmp1 string
		toContinue = []map[string]string{}
		completeFile *os.File
		fileName string
		fInfo os.FileInfo
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
	case strings.Contains(stdo.String(), "swin"):
		machine = "g2swin"
	}

	log.Println("machine set to: ", machine)

	log.Println("Starting goroutines...")
	
	for idx := 0; idx < nProcs; idx++ {
		go Out2ICsEmbed(inFileNameChan, cssInfo0)
		go CreateStartScripts(cssInfo1, machine, pbsLaunchChannel0, done)
		go PbsLaunchOnTheFly(pbsLaunchChannel1, done1)
	}
	
	log.Println("Searching for files in the form: ", globName)
	// Find last round for each run in the folder
	// Runs are sorted
	runs, runMap, mapErr = FindLastRound(globName)
	// Some round are present because of the ics ma don't have errs or outSize,
	// probably they were run somewhere else (Spritz?)
	if mapErr != nil {
		log.Println(mapErr)
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
			fmt.Printf("\tRemove because of suspicious dimensions (probably broken):\n\t%v\n\t%v\n ", lastOut, lastErr)
			for _, file := range []string{lastOut, lastErr} {
				if err = os.Remove(file); err != nil {
					log.Fatal("Error while removing ", file, ": ", err)
				}
			}
			if len(runMap[run]["out"])-2 > 0 {
				inFileNameChan <- runMap[run]["out"][len(runMap[run]["out"])-2] // rerun previous run
				tmp = <- cssInfo0
			} else { // Only ics is still here: need to only recreate start script, no new ics
				tmp = map[string]string{
					"remainingTime": "500",
					"randomSeed": "",
					"newICsFileName": runMap[run]["ics"][0],
				}				
			}
			
		} else {
			inFileNameChan <- runMap[run]["out"][len(runMap[run]["out"])-1]
			tmp = <- cssInfo0
		}
		if len(tmp) != 0 {
			toContinue = append(toContinue, tmp)
		} else {
			// Check "Rounds" folder exists, in case create it
			if fInfo, err = os.Stat("Rounds"); err != nil {
				if os.IsNotExist(err) {
					if err = os.Mkdir("Rounds", 0700); err != nil {
						log.Fatal("Can't create folder ", err)
					}
					fmt.Println("\tCreated Rounds folder")
				} else {
					log.Fatal("Can't check Rounds folder existance: ", err)
				}
			} else {
				if fInfo.IsDir() {
					fmt.Println("\tRounds already exists and is a folder")
				} else {
					log.Fatal("Rounds already exists but is not a folder")
				}
			}
			// Move all the rounds of this run to the "Rounds" folder
			fmt.Printf("\tMove all the run %v files to Rounds", run)
			for _, kindOfFile := range []string{"ics", "err", "out"} {
				for _, fileName = range runMap[run][kindOfFile]{
					if err = os.Rename(fileName, filepath.Join("Rounds", fileName)); err != nil {
						log.Fatalf("Can't rename %v because %v\n", fileName, err)
					}
				}
			}
			
		}
		// Create start scripts
		cssInfo1 <- tmp
		tmp1 = <- pbsLaunchChannel0
		pbsLaunchChannel1 <- tmp1
		<- done1
		
		fmt.Println()
		fmt.Println(".................................")
	}
	
	if len(toContinue) == 0 {
		log.Println("It seems that all the runs are complete, creating the 'complete' file")
		
		if completeFile, err = os.Create("complete"); err != nil {
			log.Fatal("Can't create complete file with error: ", err)
		}
		completeFile.Close()
	}
	
	// Close the channel, if you forget it, goroutines
	// will wait forever
	close(inFileNameChan)
	close(cssInfo1)
	close(pbsLaunchChannel0)
	close(pbsLaunchChannel1)
	
	// Wait the CreateStartScripts goroutines to finish
	for idx := 0; idx < nProcs; idx++ {
		<-done // wait the goroutine to finish
	}
	fmt.Println()
	
	log.Println("Write removed files to file")
	if !goutils.Exists(removedFileName) {
		removedFile, err = os.Create(removedFileName)
	} else {
		removedFile, err = os.OpenFile(removedFileName, os.O_APPEND|os.O_WRONLY, 0600)
	}
	if err != nil {
		log.Fatal("Error while opening removed files file: ", err)
	}
	defer removedFile.Close()
	if _, err = removedFile.WriteString(fmt.Sprintf("%v %v\n", time.Now().Format(time.RFC850), toRemove)); err != nil {
		log.Fatal("Error while writing removed files to file: ", err)
	}
	
	

	
}
