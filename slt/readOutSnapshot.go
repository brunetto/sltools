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
type OutSnapshot struct {
	Timestep string
	Integrity bool
	NestingLevel int
	Data []string	
}

// Pick the snapshot line by line and write it to the
// output file
func (snap *OutSnapshot) WriteSnapshot(nWriter *bufio.Writer) (err error) {
	for _, line := range snap.Data {
		_, err = nWriter.WriteString(line+"\n")
	}
	nWriter.Flush()
	return err
}

// This function read one and only one snapshot at a time
func ReadOutSnapshot(nReader *bufio.Reader) (*OutSnapshot, error) {
	var (
		snap *OutSnapshot = new(OutSnapshot)
		line string
		err error
		regSysTime = regexp.MustCompile(`system_time\s*=\s*(\d+)`)
		resSysTime []string
		// This variables are the idxs to print the last or last 10 lines
		dataStartIdx int = 0
		dataEndIdx int
	)

	// Init snapshot container
	snap.Data = make([]string, 1)
	snap.Integrity = false
	snap.NestingLevel = 0
	
	for {
		// Read line by line
		if line, err = readfile.Readln(nReader); err != nil {
			if err.Error() == "EOF" {
				log.Println("File reading complete...")
				log.Println("Timestep not complete.")
				log.Println("Last ten lines:")
				dataEndIdx = len(snap.Data)-1
				
				// Check that we have more than 10 lines
				if dataEndIdx > 10 {
					dataStartIdx = dataEndIdx - 10
				}
				for idx, row := range snap.Data[dataStartIdx:dataEndIdx] {
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
		snap.Data = append(snap.Data, line)
		
		// Search for timestep number
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			snap.Timestep = resSysTime[1]
// 			log.Println("Reading timestep ", resSysTime[1])
		}
		
		// Check if entering or exiting a particle
		// and update the nesting level 
		if strings.Contains(line, "(Particle") {
			snap.NestingLevel++
		} else if strings.Contains(line, ")Particle") {
			snap.NestingLevel--
		}
		
		// Check whether the whole snapshot is in memory
		// (root particle complete) and if true, return
		if snap.NestingLevel == 0 {
// 			log.Println("Timestep ended with nesting: ", snap.NestingLevel)
			snap.Integrity = true
			log.Println("Timestep ", snap.Timestep, " integrity set to: ", snap.Integrity)
			return snap, err
		}
	}	
}
