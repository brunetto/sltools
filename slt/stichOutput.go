package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)

// StdStich stiches a given STD??? according to the type passed with stdWhat.
// func StdStich(stdFiles, stdWhat string) {
func StdStich(inFilesList chan []string, done chan struct{}) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		nReader *bufio.Reader
		inFile                                *os.File
		inFiles []string
		snapshot/*s = make([]*/ *DumbSnapshot /*, 2)*/
		outFileName                           string
		outFile                               *os.File
		err                                   error
		nWriter                               *bufio.Writer
		gzipWriter                            *gzip.Writer
		timestep                              int64
		timesteps                             = make([]int64, 0)
		ext                                   string
		stdWhat string
	)
	
	for inFiles = range inFilesList {
		// Remove suffix to create the output name
		r := regexp.MustCompile(`-rnd\S+.\S+`)
		tmp := r.ReplaceAllString(inFiles[0], "")
		log.Println("Stich " + tmp)

		outFileName = tmp + "-all.txt.gz"
		log.Println("Output file will be ", outFileName)

		log.Println("Opening STDOUT output file...")

		// Open output file
		if outFile, err = os.Create(outFileName); err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		gzipWriter = gzip.NewWriter(outFile)
		defer gzipWriter.Close()
		defer gzipWriter.Flush()
	
		// Create reader and writer
		nWriter = bufio.NewWriter(gzipWriter)
		defer nWriter.Flush()

		if strings.HasPrefix(inFiles[0], "out") {
			stdWhat = "out"
		} else if strings.HasPrefix(inFiles[0], "err") {
			stdWhat = "err"
		} else {
			log.Fatal("Wrong prefix in ", inFiles[0])
		}
		
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
					var fZip *gzip.Reader
					if fZip, err = gzip.NewReader(inFile); err != nil {
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
					snapshot, err = ReadOutSnapshot(nReader, true)
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
		timesteps = nil
	}
	// Send end signal
	done <- struct{}{}
}
