package main

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		inFileName string
	)
	
	
	if len(os.Args) < 2 {
		log.Fatal("Provide a STDOUT file to check")
	}
	
	inFileName = os.Args[1]

	slt.CheckSnapshot(inFileName)
	
	fmt.Println()
}
