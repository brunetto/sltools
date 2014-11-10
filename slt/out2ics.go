package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/brunetto/goutils/debug"
)

var mute bool = false

func Out2ICsEmbed(inFileNameChan chan string, cssInfo chan map[string]string) {
	mute = true
	Out2ICs(inFileNameChan, cssInfo)
}

// Out2ICs read the STDOUT and write the new ICs with the last snapshot.
func Out2ICs(inFileNameChan chan string, cssInfo chan map[string]string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		inFileName                     string
		err                            error    // errora container
		newICsFileName                 string   // new ICs file names
		inFile, newICsFile             *os.File // last STDOUT and new ICs file
		nReader                        *bufio.Reader
		nWriter                        *bufio.Writer
		fileNameBody, newRnd, ext      string                                              // newRnd is the number of the new run round
		snapshots                      = []*DumbSnapshot{&DumbSnapshot{}, &DumbSnapshot{}} // slice for two snapshots
		snpN                           int                                                 // number of the snapshot
		simulationStop                 int64                                               // when to stop the simulation
		thisTimestep, remainingTime    int64                                               // current timestep number and remaining timesteps to reach simulationStop
		randomSeed                     string                                              // random seed from STDERR
		runString                      string                                              // string to run the next round from terminal
		newErrFileName, newOutFileName string                                              // new names from STDERR and STDOUT
		regRes                         map[string]string
		rnd                            string
		fZip                           *gzip.Reader
		lengthUnit, length, nStars, fPB, ncm       float64
	)

	// 	simulationStop = 500

	fmt.Printf("\tSimulation stop set to (slightly more than) 100 Myr and calculated in the code\n")

	// Retrieve infile from channel and use it
	for inFileName = range inFileNameChan {

		// Extract fileNameBody, round and ext
		if regRes, err = DeepReg(inFileName); err == nil {
			if regRes["prefix"] != "out" {
				log.Fatalf("Please specify a STDOUT file, found %v prefix", regRes["prefix"])
			}
			
			fileNameBody = regRes["baseName"]
			rnd = regRes["rnd"]
			ext = regRes["ext"]
			temp, _ := strconv.ParseInt(rnd, 10, 64)
			newRnd = strconv.Itoa(int(temp + 1))
			
			// Creating new filenames
			newICsFileName = "ics-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext
			newErrFileName = "err-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext
			newOutFileName = "out-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext
			
			if length, err = strconv.ParseFloat(regRes["Rv"], 64); err != nil {
				log.Fatalf("Can't convert cluster length %v to float64: %v\n", regRes["Rv"], err)
			}
			if fPB, err = strconv.ParseFloat(regRes["fPB"][:1]+"."+regRes["fPB"][1:], 64); err != nil {
				log.Fatalf("Can't convert cluster fPB %v to float64: %v\n", regRes["fPB"], err)
			}
			if ncm, err = strconv.ParseFloat(regRes["NCM"], 64); err != nil {
				log.Fatalf("Can't convert cluster NCM %v to float64: %v\n", regRes["NCM"], err)
			}
			nStars = ncm * (1 + fPB)
			
		} else if regRes, err = Reg(inFileName); err == nil {			
			if regRes["prefix"] != "out" {
				log.Fatalf("Please specify a STDOUT file, found %v prefix", regRes["prefix"])
			}
			
			log.Println("Can't derive deep info from name, only take the names and go standard for the rest")
			
			fileNameBody = regRes["baseName"]
			rnd = regRes["rnd"]
			ext = regRes["ext"]
			temp, _ := strconv.ParseInt(rnd, 10, 64)
			newRnd = strconv.Itoa(int(temp + 1))
			
			// Creating new filenames
			newICsFileName = "ics-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext
			newErrFileName = "err-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext
			newOutFileName = "out-" + fileNameBody + "-run" + regRes["run"] + "-rnd" + LeftPad(newRnd, "0", 2) + ext

			// Default
			length = float64(1)
			nStars = float64(5500)
			
			fmt.Printf("Set default parameters for cluster: \n")
			fmt.Printf("Radius: %v\n", length)
			fmt.Printf("Number of stars: %v\n", nStars)
		} else {
			log.Println("Can't derive standard names or deep info from STDOUT => wrap it!!")
			ext = filepath.Ext(inFileName)
			newICsFileName = "ics-" + inFileName + ext
			newErrFileName = "err-" + inFileName + ext
			newOutFileName = "out-" + inFileName + ext
			
			// Default
			length = float64(1)
			nStars = float64(5500)
			
			fmt.Printf("Set default parameters for cluster: \n")
			fmt.Printf("Radius: %v\n", length)
			fmt.Printf("Number of stars: %v\n", nStars)
		}

		// Now simulation stop is calculated scaling that of the clusters of the first simlations
		// with the formula
		//
		// maxTimeStep2 = ( maxTime / sqrt(timeUnit1**2 * (length2/length1)**3 * (m1 / m2)))
		//
		// with
		//
		// maxTime ~ 100 Myr
		// timeUnit1 ~ 0.25 Myr
		// length1 = 1 pc
		// m1 / m2 approximated with the number of stars, so m2 = NCM * (1 + fPB) and m1 = 5500
		// FIXME maybe leng should be timeUnit or something similar????

		lengthUnit = math.Sqrt(0.25*0.25*math.Pow(length, 3)*(5500./nStars))
		simulationStop = 1 + int64(math.Floor(110./lengthUnit))
		fmt.Printf("\tApprox length unit: %2.2f || simulationStop: %v\n", lengthUnit, simulationStop)

		// Open infile, both text or gzip and create the reader
		if !mute {
			log.Println("Opening STDOUT file: ", inFileName)
		}
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
				log.Fatal("with extension ", ext)
			}
		}

		if !mute {
			log.Println("Start reading...")
		}
		// Read two snapshot each loop to ensure at least one of them is complete
		// (= I keep the previous read in memory in case the last is corrupted)
		for {
			if snapshots[0], err = ReadOutSnapshot(nReader); err != nil {
				break
			}
			if snapshots[1], err = ReadOutSnapshot(nReader); err != nil {
				break
			}
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
			fmt.Println("Maybe your output file is empty")
			log.Fatal("Reading exit with error ", err)
		}
		// Info
		fmt.Println() // To leave a space after the non verbose print
		if !mute {
			log.Println("Done reading, last complete timestep is ", snapshots[snpN].Timestep)
		}
		thisTimestep, _ = strconv.ParseInt(snapshots[snpN].Timestep, 10, 64)
		remainingTime = simulationStop - thisTimestep

		// Write last complete snapshot to file
		if !force && remainingTime < 1 {
			fmt.Println("\tNo need to create a new ICs, simulation complete.")
			cssInfo <- map[string]string{} // empty map if no need to create css scripts
			continue
		} else {
			// Create the new ICs file
			if !mute {
				fmt.Println("Creating new ICs file ", newICsFileName)
			}
			if newICsFile, err = os.Create(newICsFileName); err != nil {
				log.Fatal(err)
			}
			defer newICsFile.Close()
			nWriter = bufio.NewWriter(newICsFile)
			defer nWriter.Flush()

			fmt.Println("\tWriting snapshot to ", newICsFileName)
			if err = snapshots[snpN].WriteSnapshot(nWriter); err != nil {
				log.Fatal("Error while writing snapshot to file: ", err)
			}
			fmt.Println("\tSet -t flag to ", remainingTime)
		}

		if !mute {
			fmt.Fprint(os.Stderr, "\n")
		}
		if !mute {
			log.Println("Search for random seed...")
		}
		randomSeed = DetectRandomSeed(inFileName)
		fmt.Println("\tSet -s flag to ", randomSeed)

		cssInfo <- map[string]string{
			"remainingTime":  strconv.Itoa(int(remainingTime)),
			"randomSeed":     randomSeed,
			"newICsFileName": newICsFileName,
		}

		runString = "\nYou can run the new round from the terminal with:\n" +
			"----------------------\n" +
			"(" + os.Getenv("HOME") + "/bin/kira -F -t " +
			strconv.Itoa(int(remainingTime)) +
			" -d 1 -D 1 -b 1 -f 0 " +
			"-n 10 -e 0.000 -B -s " + randomSeed +
			" < " + newICsFileName + " >  " + newOutFileName + " 2> " + newErrFileName + ")& \n" +
			"\nor\n\n" +
			"($HOME/bin/kiraWrap " + "-i " + newICsFileName + " -t " +
			strconv.Itoa(int(remainingTime)) + " -s " +
			randomSeed + ")\n\n" +
			"----------------------\n\n" +
			"You can watch the status of the simulation by running: \n" +
			"----------------------\n" +
			"watch stat " + newErrFileName + "\n" +
			"----------\n" +
			"cat " + newErrFileName + ` | grep "Time = " | tail -n 1` + "\n" +
			"----------------------\n"

		if !mute {
			fmt.Println(runString)
		}
		fmt.Println()
	}
	close(cssInfo)
	// 	done <- struct{}{}
}
