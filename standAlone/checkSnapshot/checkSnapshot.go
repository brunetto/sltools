package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	
	"bitbucket.org/brunetto/sltools/slt"
)

func main () () {
	var (
		inFileName string
		err error
		inFile *os.File
		nReader  *bufio.Reader
	)
	
	
	if len(os.Args) < 2 {
		log.Fatal("Provide a file to check")
	}
	
	inFileName = os.Args[1]
	
	if inFile, err = os.Open(inFileName); err != nil {log.Fatal(err)}
	defer inFile.Close()
	
	nReader = bufio.NewReader(inFile)
	
	for {
		if _, err = slt.ReadOutSnapshot(nReader); err != nil {break}
	}
	fmt.Println()
}