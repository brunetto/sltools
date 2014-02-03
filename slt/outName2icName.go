package slt

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)


// Create the output file name that will be the new IC for the restart
func OutName2ICName (inFileName string/*, fileN string*/) (outFileName string) {
	var (
		extension string
		baseName string
		file string
		dir string
		outRegString string = `out-cineca-comb(\d+)-NCM(\d+)-fPB(\d+)` + 
							`-W(\d)-Z(\d+)-run(\d+)-rnd(\d+).txt`
		outRegexp *regexp.Regexp = regexp.MustCompile(outRegString)
		outRegResult []string
	)
	
	outRegResult = outRegexp.FindStringSubmatch(inFileName); 
	if outRegResult == nil {
		log.Fatal("Can't find parameters in out name")
	}
	
	// Retrieve the round number and increment it
	temp, _ := strconv.ParseInt(outRegResult[7], 10, 64)
	rnd  = strconv.Itoa(int(temp + 1))
	
	dir = filepath.Dir(inFileName)
	file = filepath.Base(inFileName)
	extension = filepath.Ext(inFileName)
	baseName = strings.TrimSuffix(file, extension)
	baseName = strings.TrimPrefix(baseName, "out-")
	// FIXME: use regexp to check the name
	baseName = baseName[:strings.LastIndex(baseName, "-rnd")] // to remove the last round number
	outFileName = filepath.Join(dir, "ics-" + baseName) + "-rnd" + LeftPad(rnd, "0", 2) + extension //FIXME detectare nOfFiles
	return outFileName
}

func LeftPad(str, pad string, length int) (string) {
	var repeat int
	if (length - len(str)) % len(pad) != 0 {
		log.Fatal("Can't pad ", str, " with ", pad, " to length ", length)
	} else {
		repeat = (length - len(str)) / len(pad)
	}
	return strings.Repeat(pad, repeat) + str
}