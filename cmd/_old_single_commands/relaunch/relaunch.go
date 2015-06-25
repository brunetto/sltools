package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	// Clean folder
	slt.SimClean()
	
	if !goutils.Exists("complete") {
		// Check and continue
		slt.CAC("")
		
		// Submit // already in cac
// 		if err := slt.PbsLaunch(); err != nil {
// 			log.Fatal(err)
// 		}
	} else {
		log.Println("'complete' file found, assume simulations are complete.")
	}
	
	fmt.Print("\x07") // Beep when finish!!:D
}



