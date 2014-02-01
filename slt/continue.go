package slt

func Continue (inFileName, machine, userName, pName string) {
	var (
		simTime string
		randomNumber string
	)
	
	simTime, randomNumber = Out2ICs (inFileName, fileN)
	// creare il nome delle ics da quello dello stdout
	icsName = ...
	CreateScripts (icsName, machine, userName, randomNumber, simTime, pName string)
	
	
	
}