package slt

import (
	"bitbucket.org/brunetto/goutils/readfile"
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Struct containing one snapshot
type ErrSnapshot struct {
	Timestep string
	Integrity bool
	Lines []string	
}

// Pick the snapshot line by line and write it to the
// output file
func (snap *ErrSnapshot) WriteSnapshot(nWriter *bufio.Writer) (err error) {
	if Debug {Whoami(true)}
	for _, line := range snap.Lines {
		_, err = nWriter.WriteString(line+"\n")
		if err = nWriter.Flush(); err != nil {log.Fatal(err)}
	}
	nWriter.Flush()
	return err
}

// This function read one and only one snapshot at a time
func ReadErrSnapshot(nReader *bufio.Reader) (*ErrSnapshot, error) {
	if Debug {Whoami(true)}
	var (
		snap *ErrSnapshot = new(ErrSnapshot)
		line string
		err error
		regSysTime = regexp.MustCompile(`^Time = (\d+)`)
		resSysTime []string
		endOfSnap string = "----------------------------------------"
		// This variables are the idxs to print the last or last 10 lines
		dataStartIdx int = 0
		dataEndIdx int
	)

	// Init snapshot container
	snap.Lines = make([]string, 0) //FIXME: check if 0 is ok!!!
	snap.Integrity = false
	snap.Timestep = "-1"
	
	for {
		// Read line by line
		if line, err = readfile.Readln(nReader); err != nil {
			if err.Error() == "EOF" {
				log.Println("File reading complete...")
				log.Println("Timestep not complete.")
				log.Println("Last ten lines:")
				dataEndIdx = len(snap.Lines)-1
				
				// Check that we have more than 10 lines
				if dataEndIdx > 10 {
					dataStartIdx = dataEndIdx - 10
				}
				for idx, row := range snap.Lines[dataStartIdx:dataEndIdx] {
					fmt.Println(idx, ": ", row)
				}
			} else {
				log.Fatal("Non EOF error while reading ", err)
			}
			// Mark snapshot as corrupted
			snap.Integrity = false
			return snap, err
		}
		
		// Add line to the snapshots in memory
		snap.Lines = append(snap.Lines, line)
		
		// Search for timestep number
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			snap.Timestep = resSysTime[1]
// 			log.Println("Reading timestep ", resSysTime[1])
		}
		
		// Check if entering or exiting a particle
		// and update the nesting level 
		if strings.Contains(line, endOfSnap) {
			snap.Integrity = true
			if Verb {
				log.Println("Timestep ", snap.Timestep, " integrity set to: ", snap.Integrity)
			} else {
				fmt.Print("\rTimestep ", snap.Timestep, " integrity set to: ", snap.Integrity)
			}
			fmt.Println() // To leave a space after the non verbose print
			return snap, err
		}
	}	
}
