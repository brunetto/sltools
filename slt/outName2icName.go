package slt

import (
	"log"
	"path/filepath"
	"strings"
)


// Create the output file name that will be the new IC for the restart
func OutName2ICName (inFileName, fileN string) (outFileName string) {
	var (
		extension string
		baseName string
		file string
		dir string
	)
	
	dir = filepath.Dir(inFileName)
	file = filepath.Base(inFileName)
	extension = filepath.Ext(inFileName)
	baseName = strings.TrimSuffix(file, extension)
	baseName = strings.TrimPrefix(baseName, "out-")
	log.Println(baseName)
	log.Println(strings.LastIndex(baseName, "-rnd"))
	// FIXME: use regexp to check the name
	baseName = baseName[:strings.LastIndex(baseName, "-rnd")] // to remove the last round number
	outFileName = filepath.Join(dir, "ics-" + baseName) + "-rnd" + LeftPad(fileN, "0", 2) + extension //FIXME detectare nOfFiles
	return outFileName
}

func LeftPad(str, pad string, length int) (string) {
	var repeat int
	if (length - len(str)) % len(pad) != 0 {
		log.Fatal("Can't pad ", str, " with ", pad, " to length ", length)
	} else {
		repeat = (length - len(str)) / len(pad)
	}
	return strings.Repeat(pad, repeat)	
}