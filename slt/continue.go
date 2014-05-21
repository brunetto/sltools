package slt

import (
	"log"
	"time"

	"github.com/brunetto/goutils/debug"
)

// Continue provide a lazy function to prepare a simulation for he next run.
// It will convert the last snapshot from a StarLab stdout file in an ICs file
// by calling Out2ICs and then it will create the needed scripts for launching
// the simulation (kiraLaunch and PBSlaunch) calling CreateStartScripts.
// It needs a valid configuration file.
func Continue(inFileName string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		machine, remainingTime, randomSeed string
		nFileNameChan = make(chan string, 1)
		cssInfo = make(chan map[string]string, 1)
	)
	
	
	go slt.Out2ICs(inFileNameChan, cssInfo)
	go slt.CreateStartScripts(cssInfo, machine)
	
	inFileNameChan <- inFileName
	close(inFileNameChan)	

}

