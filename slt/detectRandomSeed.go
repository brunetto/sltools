package slt

import (
	"bitbucket.org/brunetto/goutils/readfile"
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	
	
)

func DetectRandomSeed(inFileName string) (randomSeed string) {
	var (
		line string
		regRandomSeed = regexp.MustCompile(`initial random seed\s*=\s*(\d+)`)
		resRandomSeed []string
		inFile *os.File
		err error
		nReader *bufio.Reader
		stdErrName string
	)
	
	stdErrName = "err" + strings.TrimPrefix(inFileName, "out")
	
	// Open file & create reader
	if inFile, err = os.Open(stdErrName); err != nil {log.Fatal(err)}
	defer inFile.Close()
	nReader = bufio.NewReader(inFile)
	
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
	return 	randomSeed
}