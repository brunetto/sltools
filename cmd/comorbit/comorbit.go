package main

import (
	"log"
	"os"
	"time"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())

	if len(os.Args) < 2 {
		log.Fatal("Please provide a STDOUT to scan for COM orbit.")
	}
	
	var inFileName string = os.Args[1]
	slt.ComOrbit(inFileName)
}




