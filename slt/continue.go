package slt

import (
	"time"

	"github.com/brunetto/goutils/debug"
)

// Continue provide a lazy function to prepare a simulation for he next run.
// It will convert the last snapshot from a StarLab stdout file in an ICs file
// by calling Out2ICs and then it will create the needed scripts for launching
// the simulation (kiraLaunch and PBSlaunch) calling CreateStartScripts.
// It needs a valid configuration file.
func Continue(inFileName, machine string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		inFileNameChan = make(chan string, 1)
		cssInfo = make(chan map[string]string, 1)
		done = make(chan struct{})
	)
	
	
	go Out2ICs(inFileNameChan, cssInfo)
	go CreateStartScripts(cssInfo, machine, done)
	
	inFileNameChan <- inFileName
	close(inFileNameChan)
	<-done // wait the goroutine to finish

}

