package main

import (
	"log"
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
		snap *slt.DumbSnapshot
	)
	
	
	if len(os.Args) < 2 {
		log.Fatal("Provide a file to check")
	}
	
	inFileName = os.Args[1]
	
	if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
	defer inFile.Close()
	
	nReader = bufio.NewReader(inFile)
	
	for {
		if snap, err = slt.ReadOutSnapshot(nReader); err != nil {break}
	}
	
	

}

