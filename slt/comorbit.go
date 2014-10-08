package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	
	"github.com/brunetto/goutils/debug"
)

func ComOrbit (inFileName string) () {
	defer debug.TimeMe(time.Now())
	
	var (
		err error
		outFileName string
		inFile, outFile *os.File
		ext string
		nReader *bufio.Reader
		fZip *gzip.Reader
		snap *DumbSnapshot = new(DumbSnapshot)
		line string
		idx int
		coordReg = regexp.MustCompile(`r  =  (\S+\s+\S+\s+\S+)`)
		coordRes []string
		coords []string = []string{}
	)
	
	outFileName = "coords-"+inFileName
	
	ext = filepath.Ext(inFileName)
	
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
	
	for {
		if snap, err = ReadOutSnapshot(nReader); err != nil {
			break
		}
		// Loop over the snap lines
		for idx, line = range snap.Lines {
			// Assuming root node is the first of the snap
			if !strings.Contains(line, "(Dynamics") {continue}
			// Assuming coords are 3 lines after Dynamics section beginning
			if coordRes = coordReg.FindStringSubmatch(snap.Lines[idx+3]); coordRes == nil {
				log.Fatalf("Can't find coordinates in %v: %v\n", snap.Lines[idx+3], coordRes)
			}
// 			fmt.Println(coordRes[1])
			break
		}
		// For each snap we have one coord set for the COM
		coords = append(coords, coordRes[1])	
		coordRes = []string{}
	}
	
	if outFile, err = os.Create(outFileName); err != nil {
		log.Fatal("Can't create outfile with error: ", err)
	}
	for _, line = range coords {
		outFile.WriteString(line+"\n")
	}
	fmt.Println()
}




