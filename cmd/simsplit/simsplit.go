package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/sltools/slt"
	"github.com/codegangsta/cli"
)

func main() {
	defer debug.TimeMe(time.Now())
	app := cli.NewApp()
	app.Name = "simsplit"
	app.Usage = "Split a Starlab simulation STDOUT or STDERR in the single snapshots."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Value: "",
			Usage: "File type (out/err)",
		},
		cli.StringFlag{
			Name:  "infile, i",
			Value: "",
			Usage: "File to split",
		},
		cli.StringFlag{
			Name:  "extract, e",
			Value: "",
			Usage: "Extract only the specify snapshot",
		},
	}

	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		if c.String("type") != "out" && c.String("type") != "err" {
			log.Fatal("Please provide the type of the file to split, either 'out' for STDOUT or 'err' for STDERR.")
		}
		if c.String("infile") == "" {
			log.Fatal("Please provide the name of the file to split.")
		}
		cutStd(c.String("infile"), c.String("type"), c.String("extract"))
	}

	app.Run(os.Args)

}

func cutStd(inFileName, stdWhat, extract string) {
	// 	defer debug.TimeMe(time.Now())
	var (
		baseName, ext  string
		inFile         *os.File
		nReader        *bufio.Reader
		err            error
		snapshot       *slt.DumbSnapshot
		timestep       int64
		timesteps          = make([]int64, 0)
		wroteTimesteps     = make([]string, 0)
		snapChan           = make(chan *slt.DumbSnapshot)
		done               = make(chan struct{}, 1)
		nProcs         int = 4
	)

	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}

	ext = filepath.Ext(inFileName)
	baseName = strings.TrimSuffix(inFileName, ext)

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
	
	// Start gorutines
	for idx := 0; idx < nProcs; idx++ {
		go write(baseName, snapChan, done)
	}

SnapLoop: // label
	for {
		if stdWhat == "out" {
			snapshot, err = slt.ReadOutSnapshot(nReader, false)
		} else if stdWhat == "err" {
			snapshot, err = slt.ReadErrSnapshot(nReader)
		} else {
			log.Fatal("Unrecognized stdWhat: ", stdWhat)
		}
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal("Incomplete snapshot.", err)
			}
			break
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
				if slt.AbsInt(timestep-timesteps[len(timesteps)-1]) > 1 {
					if true {
						log.Println("Read timestep: ")
						for _, ts := range timesteps {
							fmt.Print(ts, " ")
						}
						fmt.Println()
					}
					log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
				} else if slt.AbsInt(timestep-timesteps[len(timesteps)-1]) < 1 {
					log.Println("Duplicated timestep ", timestep, ", continue.")
					continue SnapLoop /*to the next timestep*/
				}
			}
			timesteps = append(timesteps, timestep) // Write the snapshot
			if extract == "" || extract == snapshot.Timestep {
				snapChan <- snapshot
				wroteTimesteps = append(wroteTimesteps, snapshot.Timestep)
			}

		} else {
			// This shouldn't happend because of the break in reading the snapshots
			// This shoud be a redundant check
			// TODO: check if it is true!!!
			log.Fatal("Bad snapshot.")
		}
	} // end reading snapshot from a single file loop
	
	// Close channel and shutdown gorutines
	close(snapChan)
	for idx := 0; idx < nProcs; idx++ {
		<-done 
	}

	fmt.Println("\n")
	fmt.Println("Analyzed timesteps, \n", timesteps)
	fmt.Println("\n\n")
	fmt.Println("Wrote timesteps, \n", wroteTimesteps)
	fmt.Println("\n\n")
}

func write(baseName string, snapChan chan *slt.DumbSnapshot, done chan struct{}) {
	var (
		outFileName string
		outFile     *os.File
		nWriter     *bufio.Writer
		err         error
	)
	
	for snapshot := range snapChan {
		outFileName = baseName + "-part_" + snapshot.Timestep + ".txt"
		if outFile, err = os.Create(outFileName); err != nil {
			log.Fatal(err)
		}

		// Create reader and writer
		nWriter = bufio.NewWriter(outFile)
		if err = snapshot.WriteSnapshot(nWriter); err != nil {
			log.Fatal("Error while writing snapshot to file: ", err)
		}
		nWriter.Flush()
		outFile.Close()
	}
	// Halt signal
	done <- struct{}{}
}
