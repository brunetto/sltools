package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
// 	"runtime"
// 	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)
/*
// StichThemAll launch the stiching in parallel on all the simulation files in
// the folder, accordingly to their names (run 01 is different from run 02 and so on).
func StichThemAll(sampleFile string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		err            error
		inFiles        []string
		prefixes       = []string{"out-", "err-"}
		run, baseName  string
		tmp            map[string]string
		runs           StringSet // set = list of unique objects (run numbers)
		nRuns          []int
		globName       string
		maxProcs       int = 1
		inFileNameChan     = make(chan string, maxProcs)
		done               = make(chan struct{})
	)

	runtime.GOMAXPROCS(maxProcs)

	nRuns = make([]int, 0)

	if tmp, err = Reg(sampleFile); err != nil {
		log.Fatal(err)
	}
	baseName = tmp["baseName"]

	// Search for all the STDOUT and STDERR files in the folder
	for idx := 0; idx < 2; idx++ {
		globName = "*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.*"
		if Verb {
			log.Println("Searching for: ", globName)
		}
		if inFiles, err = filepath.Glob(globName); err != nil {
			log.Fatal("Error globbing for stiching all the run outputs, ", globName, " ,  in this folder: ", err)
		}
		// Sort file names
		sort.Strings(inFiles)

		runs = NewStringSet()

		// Find the numbers of the different runs
		for _, inFileName := range inFiles {
			if tmp, err = Reg(inFileName); err != nil {
				log.Fatal(err)
			}
			run = tmp["run"]
			// Add the new number in the set
			runs.Add(run)
		}
		if Verb {
			log.Println("Found runs:")
			fmt.Println(runs.String())
		}
		nRuns = append(nRuns, len(runs))
	}

	// Launch maxProcs goroutines
	for idx := 0; idx < maxProcs; idx++ {
		go StichOutput(inFileNameChan, done)
	}

	// Launch all the stiching
	for _, runIdx := range runs.Sorted() {
		name := "out-" + baseName + "-run" + runIdx + "-rnd00.*"
		if Verb {
			log.Println("Launching stich based on ", name)
		}
		inFileNameChan <- name
	}
	close(inFileNameChan)
	for idx := 0; idx < maxProcs; idx++ {
		<-done // wait the goroutines to finish
	}
}

// FIXME: Workaround to call StichOutput not in parallel
// because now StichOutput contain a call to wg.Done
// and I don't want to import "sync" in command.go
func StichOutputSingle(inFileName string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		inFileNameChan = make(chan string, 1)
		done           = make(chan struct{})
	)
	go StichOutput(inFileNameChan, done)
	inFileNameChan <- inFileName
	close(inFileNameChan)
	<-done // wait the goroutine to finish
}

// StichOutput stiches the STDOUT and STDERR of a simulation.
func StichOutput(inFileNameChan chan string, done chan struct{}) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		inFileName string
		run        string
		stdOuts    string
		stdErrs    string
		baseName   string
		tmp        map[string]string
		err        error
	)

	for inFileName = range inFileNameChan {

		if inFileName == "" {
			log.Fatal("You need to specify an input file template with the -i flag!!!")
		}

		// Extract parameters from the name
		if tmp, err = Reg(inFileName); err != nil {
			log.Fatal(err)
		}
		run = tmp["run"]
		baseName = tmp["baseName"]

		log.Println("Stiching *-" + baseName + `-run` + run + `-rnd*.*`)

		// Check if only have to run STDERR stich
		if !OnlyErr {
			//
			// STDOUT
			//
			stdOuts = "out-" + baseName + `-run` + run + `-rnd*.*`
			StdStich(stdOuts, "out")
		} else {
			log.Println("Only stich STDERRs")
		}

		// Check if only have to run STDOUT stich
		if !OnlyOut {
			//
			// STDERR
			//
			stdErrs = "err-" + baseName + `-run` + run + `-rnd*.*`
			StdStich(stdErrs, "err")

		} else {
			log.Println("Only stich STDOUTs")
		}
	}
	done <- struct{}{}
}*/

// StdStich stiches a given STD??? according to the type passed with stdWhat.
// func StdStich(stdFiles, stdWhat string) {
func StdStich(inFilesList chan []string, done chan struct{}) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		fZip                                  *gzip.Reader
		inFile                                *os.File
		inFiles []string
		snapshot/*s = make([]*/ *DumbSnapshot /*, 2)*/
		outFileName                           string
		outFile                               *os.File
		err                                   error
		nReader                               *bufio.Reader
		nWriter                               *bufio.Writer
		timestep                              int64
		timesteps                             = make([]int64, 0)
		ext                                   string
		stdWhat string
	)
	
	for inFiles = range inFilesList {
		
		tmp := strings.TrimSuffix(inFiles[0], "-rnd*.*")
		log.Println("Stich std" + tmp)

		outFileName = tmp + "-all.txt"
		log.Println("Output file will be ", outFileName)

		log.Println("Opening STDOUT output file...")

		// Open output file
		if outFile, err = os.Create(outFileName); err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		// Create reader and writer
		nWriter = bufio.NewWriter(outFile)
		defer nWriter.Flush()

		if strings.HasPrefix(inFiles[0], "out") {
			stdWhat = "out"
		} else if strings.HasPrefix(inFiles[0], "err") {
			stdWhat = "err"
		} else {
			log.Fatal("Wrong prefix in ", inFiles[0])
		}
		
// 		log.Println("Globbing and sorting " + stdWhat + " input files")
// 		// Open infiles
// 		if inFiles, err = filepath.Glob(stdFiles); err != nil {
// 			log.Fatal("Error globbing "+stdWhat+" files for output stiching: ", err)
// 		}

// 		sort.Strings(inFiles)

		if Verb {
			log.Println("Found:")
			for idx, file := range inFiles {
				fmt.Println(idx, ": ", file)
			}
		}

		for _, inFileName := range inFiles {
			if Verb {
				log.Println("Working on ", inFileName)
			}
			if inFile, err = os.Open(inFileName); err != nil {
				log.Fatal(err)
			}
			defer inFile.Close()
			ext = filepath.Ext(inFileName)
			// Try to open the file if it is a gzipped one or a simple txt
			switch ext {
			case ".txt":
				{
					nReader = bufio.NewReader(inFile)
				}
			case ".gz":
				{
					fZip, err = gzip.NewReader(inFile)
					if err != nil {
						log.Fatal("Can't open %s: error: %s\n", inFileName, err)
					}
					nReader = bufio.NewReader(fZip)
				}
			default:
				{
					log.Fatal("Unrecognized file type", inFileName)
				}
			}

			//Read snapshots and write them if everything is OK
		SnapLoop: // label
			for {
				if stdWhat == "out" {
					snapshot, err = ReadOutSnapshot(nReader)
				} else if stdWhat == "err" {
					snapshot, err = ReadErrSnapshot(nReader)
				} else {
					log.Fatal("Unrecognized stdWhat: ", stdWhat)
				}
				if err != nil {
					if Verb {
						log.Println("Incomplete snapshot, moving to the next file")
					}
					break SnapLoop
				}
				// -1 is the "ICs to 0" timestep, skipping
				// I will skip this also because it creates problems of duplication
				// and timestep check
				if snapshot.Timestep == "-1" && len(timesteps) > 0 {
					continue SnapLoop /*to the next timestep*/
				}

				// I will loose the last timestep on STDERR because it is probably not complete
				// TODO: find out how to manage this
				// BUG: I can't find a univoque way to define the last snapshot complete
				if snapshot.Integrity == true {
					timestep, err = strconv.ParseInt(snapshot.Timestep, 10, 64)
					// Skip the first loop (=first timestep) with len = 0
					if len(timesteps) > 0 {
						if AbsInt(timestep-timesteps[len(timesteps)-1]) > 1 {
							if Verb {
								log.Println("Read timestep: ")
								for _, ts := range timesteps {
									fmt.Print(ts, " ")
								}
								fmt.Println()
							}
							log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
						} else if AbsInt(timestep-timesteps[len(timesteps)-1]) < 1 {
							log.Println("Duplicated timestep ", timestep, ", continue.")
							continue SnapLoop /*to the next timestep*/
						}
					}
					timesteps = append(timesteps, timestep) // Write the snapshot
					if err = snapshot.WriteSnapshot(nWriter); err != nil {
						log.Fatal("Error while writing snapshot to file: ", err)
					}
				} else {
					// This shouldn't happend because of the break in reading the snapshots
					// This shoud be a redundant check
					// TODO: check if it is true!!!
					fmt.Println("************************ ATTENTION *************************")
					fmt.Println("************************************************************")
					log.Println("Skipping incomplete snapshot at timestep", snapshot.Timestep)
					fmt.Println("************************************************************")
					fmt.Println("************************************************************")
				}
			} // end reading snapshot from a single file loop
		} // end reading file loop
		fmt.Println("\n")
		log.Println("Wrote ", len(timesteps), "snapshots to ", outFileName)
		fmt.Println(timesteps)
	}
	// Send end signal
	done <- struct{}{}
}
