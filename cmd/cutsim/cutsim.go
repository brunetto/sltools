package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	
	"github.com/spf13/cobra"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	initCommands()
	cutsimCmd.Execute()
	
}

var (
	inFileName string
	selectedSnapshot string
)

var cutsimCmd = &cobra.Command {
	Use:   "cutsim",
	Short: `Shorten a give snapshot to a certain timestep
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	cutsim out --inFile <STDOUT file> --cut <snapshot where to cut>
	cutsim err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var bothCutCmd = &cobra.Command {
	Use:   "cut",
	Short: "cut <STDOUT or STDERR>",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cutStdBoth(inFileName, selectedSnapshot)
	},	
}

var stdOutCutCmd = &cobra.Command {
	Use:   "out",
	Short: "cut STDOUT",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cutStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrCutCmd = &cobra.Command {
	Use:   "err",
	Short: "cut STDERR",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cutStdErr(inFileName, selectedSnapshot)
	},	
}	

func initCommands () {
	cutsimCmd.AddCommand(stdOutCutCmd)
	cutsimCmd.AddCommand(stdErrCutCmd)
	
	cutsimCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	cutsimCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
}


func cutStdBoth(inFileName, selectedSnapshot string) {
}

func cutStdOut(inFileName, selectedSnapshot string) {
	defer debug.TimeMe(time.Now())

	
	var (
		err                            error    // errora container
		inFile, outFile             *os.File // last STDOUT and new ICs file
		nReader                        *bufio.Reader
		nOutWriter                        *bufio.Writer
		ext      string                           // newRnd is the number of the new run round
		snapshots                      = make([]*slt.DumbSnapshot, 2)       // slice for two snapshots
		fZip                           *gzip.Reader
		wZip *gzip.Writer
	)
	
	// Backup old STDOUT
	if err = os.Rename(inFileName, inFileName+".bck"); err != nil {
		log.Fatalf("Error renaming %v: %v\n", inFileName, err)
	}
	
	// Extract fileNameBody, round and ext
	ext = filepath.Ext(inFileName)
		
	// Open infile, both text or gzip and create the reader
	log.Println("Opening input and output files...")
	if inFile, err = os.Open(inFileName+".bck"); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	// Create the new old out file
	if outFile, err = os.Create(inFileName); err != nil {
		log.Fatal(err)
	}

	switch ext {
	case ".txt":
		{
			nReader = bufio.NewReader(inFile)
			nOutWriter = bufio.NewWriter(outFile)
			defer nOutWriter.Flush()
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatalf("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
			
			wZip = gzip.NewWriter(outFile)
			defer wZip.Close()
			defer wZip.Flush()
			nOutWriter = bufio.NewWriter(wZip)
			defer nOutWriter.Flush()
		}
	case ".txt.gz":
	{
		fZip, err = gzip.NewReader(inFile)
		if err != nil {
			log.Fatalf("Can't open %s: error: %s\n", inFile, err)
		}
		nReader = bufio.NewReader(fZip)
		
		wZip = gzip.NewWriter(outFile)
		defer wZip.Close()
		defer wZip.Flush()
		nOutWriter = bufio.NewWriter(wZip)
		defer nOutWriter.Flush()
	}
	default:
		{
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extension ", ext)
		}
	}

	log.Println("Start reading...")
	// Read two snapshot each loop to ensure at least one of them is complete
	// (= I keep the previous read in memory in case the last is corrupted)
	for {
		if snapshots[0], err = slt.ReadOutSnapshot(nReader); err != nil {
			break
		}
		
		// Write to the old cutted file the integer snapshots
		if snapshots[0].Integrity == true {
			if err = snapshots[0].WriteSnapshot(nOutWriter); err != nil {
				log.Fatal("Error while writing snapshot to file: ", err)
			} 
		}
		if snapshots[0].Timestep == selectedSnapshot {break}
		if snapshots[1], err = slt.ReadOutSnapshot(nReader); err != nil {
			break
		}
		
		// Write to the old cutted file the integer snapshots
		if snapshots[1].Integrity == true {
			if err = snapshots[1].WriteSnapshot(nOutWriter); err != nil {
				log.Fatal("Error while writing snapshot to file: ", err)
			}
		}
		if snapshots[1].Timestep == selectedSnapshot {break}
	}
	fmt.Println() // To leave a space after the non verbose print
}
	



func cutStdErr(inFileName, selectedSnapshot string) () {
	defer debug.TimeMe(time.Now())

	var (
		fZip                                  *gzip.Reader
		inFile                                *os.File
		snapshot/*s = make([]*/ *slt.DumbSnapshot /*, 2)*/
		outFile                               *os.File
		err                                   error
		nReader                               *bufio.Reader
		nWriter                               *bufio.Writer
		timestep                              int64
		timesteps                             = make([]int64, 0)
		ext                                   string
		wZip *gzip.Writer
	)

	// Backup old STDERR
	if err = os.Rename(inFileName, inFileName+".bck"); err != nil {
		log.Fatalf("Error renaming %v: %v\n", inFileName, err)
	}
		
	// Open output file
	if outFile, err = os.Create(inFileName); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if inFile, err = os.Open(inFileName+".bck"); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	ext = filepath.Ext(inFileName)
	
	// Try to open the file if it is a gzipped one or a simple txt
	switch ext {
	case ".txt":
		{
			nReader = bufio.NewReader(inFile)
			nWriter = bufio.NewWriter(outFile)
			defer nWriter.Flush()
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatalf("Can't open %s: error: %s\n", inFileName, err)
				os.Exit(1)
			}
			nReader = bufio.NewReader(fZip)
			
			wZip = gzip.NewWriter(outFile)
			defer wZip.Close()
			defer wZip.Flush()
			nWriter = bufio.NewWriter(wZip)
			defer nWriter.Flush()
		}
	default:
		{
			log.Fatal("Unrecognized file type", inFileName)
		}
	}

	//Read snapshots and write them if everything is OK
	SnapLoop:
	for {
		snapshot, err = slt.ReadErrSnapshot(nReader)
		if err != nil {
			log.Printf("Incomplete snapshot %v\n", snapshot.Timestep)
			break SnapLoop
		}
// 		// -1 is the "ICs to 0" timestep, skipping
// 		// I will skip this also because it creates problems of duplication
// 		// and timestep check
// 		if snapshot.Timestep == "-1" && len(timesteps) > 0 {
// 			continue SnapLoop /*to the next timestep*/
// 		}

		// I will loose the last timestep on STDERR because it is probably not complete
		// TODO: find out how to manage this
		// BUG: I can't find a univoque way to define the last snapshot complete
		if snapshot.Integrity == true {
			timestep, err = strconv.ParseInt(snapshot.Timestep, 10, 64)
			// Skip the first loop (=first timestep) with len = 0
			if len(timesteps) > 1 {
				if slt.AbsInt(timestep-timesteps[len(timesteps)-1]) > 1 {
					log.Fatal("More that one timestep of distance between ", timesteps[len(timesteps)-1], " and ", timestep)
				} else if slt.AbsInt(timestep-timesteps[len(timesteps)-1]) < 1 {
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
		if snapshot.Timestep == selectedSnapshot {break}
	} // end reading snapshot from a single file loop
	fmt.Println("\n")
	log.Println("Wrote ", len(timesteps), "snapshots to ", inFileName)
	fmt.Println(timesteps)
}
