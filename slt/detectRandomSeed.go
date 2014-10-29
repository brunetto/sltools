package slt

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"

	"github.com/brunetto/goutils/readfile"
)

// DetectRandomSeed read the initial random seed form the STDERR.
func DetectRandomSeed(inFileName string) (randomSeed string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		line          string
		regRandomSeed = regexp.MustCompile(`initial random seed\s*=\s*(\d+)`)
		resRandomSeed []string
		inFile        *os.File
		err           error
		nReader       *bufio.Reader
		stdErrName    string
		fZip                           *gzip.Reader
		ext      string = filepath.Ext(inFileName)
	)

	stdErrName = "err" + strings.TrimPrefix(inFileName, "out")

	// Open file & create reader
	if inFile, err = os.Open(stdErrName); err != nil {
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
				log.Fatalf("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
		}
	default:
		{
			log.Println("Unrecognized file type", inFileName+".bck")
			log.Fatal("with extension ", ext)
		}
	}
	
	for {
		if line, err = readfile.Readln(nReader); err != nil {
			log.Fatal("STDERR interrupted before the random seed was found!!!")
		}
		// Search for timestep number
		if resRandomSeed = regRandomSeed.FindStringSubmatch(line); resRandomSeed != nil {
			randomSeed = resRandomSeed[1]
			break
		}
	}
	return randomSeed
}
