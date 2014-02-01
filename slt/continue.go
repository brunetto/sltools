package slt

func Continue (inFileName, machine, userName, pName string) {
	var (
		simTime string
		randomNumber string
	)
	
	simTime, randomNumber, icsName = Out2ICs (inFileName/*, fileN*/)
	CreateScripts (icsName, machine, userName, randomNumber, simTime, pName)
		
}