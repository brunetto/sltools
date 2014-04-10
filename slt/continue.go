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
func Continue(inFileName string, conf *ConfigStruct) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		icsName      string
		simTime      string
		randomNumber string
	)

	log.Println("Preparing to continue from ", inFileName)
	// Create the new ICs from the las snapshot
	simTime, randomNumber, icsName = Out2ICs(inFileName, conf)
	// Create start scripts (kira launch and PBS)
	CreateStartScripts(icsName, randomNumber, simTime, conf)

}
