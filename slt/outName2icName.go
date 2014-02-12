package slt

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)


// Create the output file name that will be the new IC for the restart
func OutName2ICName (inFileName string, conf *ConfigStruct) (outFileName string) {
	if Debug {Whoami(true)}
	var (
		extension string
		baseName string
		file string
		dir string
		outRegString string = `out-` + conf.BaseName() + `-run(\d+)-rnd(\d+).txt`
		outRegexp *regexp.Regexp = regexp.MustCompile(outRegString)
		outRegResult []string
		run string
		rnd string
		runString string
	)
	
	outRegResult = outRegexp.FindStringSubmatch(inFileName); 
	if outRegResult == nil {
		log.Fatal("Can't find parameters in out name ", inFileName)
	}
	
	rnd  = outRegResult[1]
	rnd  = outRegResult[2]
	
	// Retrieve the round number and increment it
	temp, _ := strconv.ParseInt(rnd, 10, 64)
	rnd = strconv.Itoa(int(temp + 1))
	runString = "-run" + run + "-rnd" + LeftPad(rnd, "0", 2)
	
	dir = filepath.Dir(inFileName)
	file = filepath.Base(inFileName)
	extension = filepath.Ext(inFileName)
	baseName = strings.TrimSuffix(file, extension)
	baseName = strings.TrimPrefix(baseName, "out-")
	// FIXME: use regexp to check the name
	baseName = baseName[:strings.LastIndex(baseName, "-rnd")] // to remove the last round number
// 	outFileName = filepath.Join(dir, "ics-" + baseName + runString + extension) //FIXME detectare nOfFiles
	outFileName = filepath.Join(dir, "ics-" + conf.BaseName() + runString + extension)
	return outFileName
}
