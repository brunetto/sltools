package slt

import (
	"log"
)

// Continue provide a lazy function to prepare a simulation for he next run.
// It will convert the last snapshot from a StarLab stdout file in an ICs file
// by calling Out2ICs and then it will create the needed scripts for launching
// the simulation (kiraLaunch and PBSlaunch) calling CreateStartScripts.
// It needs a valid configuration file.
func Continue (inFileName string, conf *ConfigStruct) {
	if Debug {Whoami(true)}

	var (
		icsName string
		simTime string
		randomNumber string
	)
	
	log.Println("Preparing to continue from ", inFileName)
	simTime, randomNumber, icsName = Out2ICs (inFileName)
	CreateStartScripts (icsName, randomNumber, simTime, conf)
		
}