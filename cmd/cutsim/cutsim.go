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
	restartFromHereCmd.Execute()
	
}

var (
	inFileName string
	selectedSnapshot string
)

var restartFromHereCmd = &cobra.Command {
	Use:   "restartFromHere",
	Short: "Prepare a pp3-stalled simulation to be restarted",
	Long: `Too often StarLab stalled while integrating a binary,
	this tool let you easily restart a stalled simulation.
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	restartFromHere out --inFile <STDOUT file> --cut <snapshot where to cut>
	restartFromHere err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var stdOutCutCmd = &cobra.Command {
	Use:   "out",
	Short: "Prepare a pp3-stalled stdout to restart the simulation",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cutStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrCutCmd = &cobra.Command {
	Use:   "err",
	Short: "Prepare a pp3-stalled stderr so that it is synced with the stdout",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cutStdErr(inFileName, selectedSnapshot)
	},	
}	

func initCommands () {
	restartFromHereCmd.AddCommand(stdOutCutCmd)
	restartFromHereCmd.AddCommand(stdErrCutCmd)
	
	restartFromHereCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	restartFromHereCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
}

func cutStdOut(inFileName, selectedSnapshot string) {
	defer debug.TimeMe(time.Now())

	
	var (
		err                            error    // errora container
		newICsFileName                 string   // new ICs file names
		inFile, newICsFile             *os.File // last STDOUT and new ICs file
		nReader                        *bufio.Reader
		nWriter                        *bufio.Writer
		fileNameBody, newRnd, ext      string                           // newRnd is the number of the new run round
		snapshots                      = make([]*slt.DumbSnapshot, 2)       // slice for two snapshots
		snpN                           int                              // number of the snapshot
		simulationStop                 int64                      = 500 // when to stop the simulation
		thisTimestep, remainingTime    int64                            // current timestep number and remaining timesteps to reach simulationStop
		randomSeed                     string                           // random seed from STDERR
		runString                      string                           // string to run the next round from terminal
		newErrFileName, newOutFileName string                           // new names from STDERR and STDOUT
		regRes                         map[string]string
		rnd                            string
		fZip                           *gzip.Reader
	)
	
	// Extract fileNameBody, round and ext
	ext = filepath.Ext(inFileName)
	regRes, err = slt.Reg(inFileName)
	if err != nil {
		log.Println("Can't derive standard names from STDOUT => wrap it!!")
		newICsFileName = "ics-" + inFileName + ext
		newErrFileName = "err-" + inFileName + ext
		newOutFileName = "out-" + inFileName + ext
	} else {
		if regRes["prefix"] != "out" {
			log.Fatalf("Please specify a STDOUT file, found %v prefix", regRes["prefix"])
		}

		fileNameBody = regRes["baseName"]
		rnd = regRes["rnd"]
		temp, _ := strconv.ParseInt(rnd, 10, 64)
		newRnd = strconv.Itoa(int(temp + 1))

		// Creating new filenames
		newICsFileName = "ics-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + slt.LeftPad(newRnd, "0", 2) + ext
		newErrFileName = "err-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + slt.LeftPad(newRnd, "0", 2) + ext
		newOutFileName = "out-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + slt.LeftPad(newRnd, "0", 2) + ext
	}

	log.Println("New ICs file will be ", newICsFileName)

	// Open infile, both text or gzip and create the reader
	log.Println("Opening input and output files...")
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	switch ext {
	case ".txt":
		{
			nReader = bufio.NewReader(inFile)
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
		}
	default:
		{
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extention ", ext)
		}
	}

	// Create the new ICs file
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
		if snapshots[0].Timestep == selectedSnapshot {break}
		if snapshots[1], err = slt.ReadOutSnapshot(nReader); err != nil {
			break
		}
		if snapshots[1].Timestep == selectedSnapshot {break}
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

	runString = "\nYou can run the new round from the terminal with:\n" +
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
		"----------\n" +
		"cat " + newErrFileName + ` | grep "Time = " | tail -n 1` + "\n" +
		"----------------------\n\n"

	fmt.Println(runString)
	fmt.Println()
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
	nWriter = bufio.NewWriter(outFile)

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
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFileName, err)
				os.Exit(1)
			}
			nReader = bufio.NewReader(fZip)
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
			if len(timesteps) > 0 {
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
