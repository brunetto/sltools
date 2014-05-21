package main

import (
	"log"
	"os"
	"time"

	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main() {
	defer debug.TimeMe(time.Now())

	var (
		inFileName string
		nFileNameChan chan string
		cssInfo chan map[string][string]
	)

	if len(os.Args) < 2 {
		log.Fatal("You need to provide a STDOUT file to parse!!")
	}

	inFileName = os.Args[1]
	
	go slt.Out2ICs(inFileNameChan chan string, cssInfo chan map[string][string])
	
	inFileNameChan <- inFileName
	close(inFileNameChan)
	
	<-cssInfo
}
