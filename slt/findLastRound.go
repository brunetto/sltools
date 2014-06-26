package slt

import (
	"errors"
	"log"
	"path/filepath"
	"sort"
	"time"
	
	"github.com/brunetto/goutils/debug"
)


// FindLastRound gives you the last round ics, err and out 
// for each run in a folder 
func FindLastRound (globName string) (keys []string, runMap map[string]map[string][]string, err error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	var (
		fileName string
		inFiles      []string		
		exists bool
		regRes                         map[string]string
	)
	
	// Init runMap
	runMap = map[string]map[string][]string{}
	// for example runMap["08"]["err"][3] 
	// will give ["err-....run08-rnd03.txt"]

	
	if inFiles, err = filepath.Glob(globName); err != nil {
		log.Fatal("Error globbing files in this folder: ", err)
	}
	
	for _, fileName = range inFiles {
		// Try to detect file parameters (type, run, rnd) from fileName
		regRes, err = Reg(fileName)
		// Not standard name
		if err != nil {log.Fatal("Can't find proper name to regex in ", fileName)}
		// Check if run is present, if not, create it in the map
		if _, exists = runMap[regRes["run"]]; !exists {
			runMap[regRes["run"]] = map[string][]string{
				"ics": []string{},
				"err": []string{},
				"out": []string{},
			}
		}
		// Fill the map entry with the fileName
		runMap[regRes["run"]][regRes["prefix"]] = append(runMap[regRes["run"]][regRes["prefix"]], fileName)
	}
	
	// Now runMap contains all the fileName 
	keys = make([]string, len(runMap))
	idx := 0 
	err = nil
	for key, value := range runMap {
        keys[idx] = key
        // Sort rounds
        sort.Strings(value["ics"])
		sort.Strings(value["err"])
		sort.Strings(value["out"])
		if len(value["ics"]) != len(value["err"]) ||
			len(value["ics"]) != len(value["out"]) ||
			len(value["err"]) != len(value["out"]) {
				err = errors.New("Some ics/err/out missing, not all the slices of a run have the same size.")
			}
        idx++
    }
    sort.Strings(keys)
	
	return keys, runMap, err
}
	
	
