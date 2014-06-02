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
		inFileName    string
		nFileNameChan = make(chan string, 1)
		cssInfo       = make(chan map[string]string, 1)
	)

	if len(os.Args) < 2 {
		log.Fatal("You need to provide a STDOUT file to parse!!")
	}

	inFileName = os.Args[1]

	go slt.Out2ICs(inFileNameChan, cssInfo)

	inFileNameChan <- inFileName
	close(inFileNameChan)

	<-cssInfo
}
