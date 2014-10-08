package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())

	var (
		err error
		inFileName string
		endOfSimMyr float64 = 100
	)
	
	if len(os.Args) < 3 {
		log.Fatal("Provide a STDOUT file and a time in Myr to try to find the final timestep")
	} else {
		inFileName = os.Args[1]
		if endOfSimMyr, err = strconv.ParseFloat(os.Args[2], 64); err != nil {
			log.Fatal(err)
		}
	}
		
	slt.CheckEnd (inFileName, endOfSimMyr)
	
	fmt.Print("\x07") // Beep when finish!!:D
}



