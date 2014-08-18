package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"path/filepath"
	"os"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		inFileName string
		err error
		inFile *os.File
		nReader  *bufio.Reader
		fZip                           *gzip.Reader
		ext      string 
	)
	
	
	if len(os.Args) < 2 {
		log.Fatal("Provide a STDOUT file to check")
	}
	
	inFileName = os.Args[1]
	
	if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
	defer inFile.Close()
	
	ext = filepath.Ext(inFileName)
	
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
		case ".txt.gz":
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
				log.Fatal("with extention ", ext)
			}
	}
	
	for {
		if _, err = slt.ReadOutSnapshot(nReader); err != nil {break}
	}
	fmt.Println()
}
