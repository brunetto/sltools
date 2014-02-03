package slt

import (
)

func Continue (inFileName, machine, userName, pName string) {
	if Debug {Whoami(true)}
	var (
		simTime string
		randomNumber string
	)
	
	simTime, randomNumber, icsName = Out2ICs (inFileName/*, fileN*/)
	CreateStartScripts (icsName, machine, userName, randomNumber, simTime, pName)
		
}