package slt

import (
	"log"
)

func CreateKira () () {
	kiraBinPath := "/eurora/home/userexternal/bziosi00/starlabjune19_2013/usr/bin/kira"
	randomString := "-s 1361557926"
	folderName := "/gpfs/scratch/userexternal/bziosi00/parameterSpace/thisFolder/"
	initCondFile := "pippoIC.txt"
	stdOutFile := "pippoOut.txt"
	stdErrFile := "pippoErr.txt"
	
	// I know I can use `` but I don't like the string not to be align with the 
	// rest of the code
	kiraString := "echo $PWD\n" + 
				  "echo $LD_LIBRARY_PATH\n" + 
				  kiraBinPath + " -t 500 -d 1 -D 1 -b 1 -f 0.3 \\\n"+ " -n 10 -e 0.000 -B " + randomString + " \\\n" +
				  "< " + folderName + " " + initCondFile + " \\\n"
				  "> " + folderName + " " + stdOutFile + " \\\n"
				  "2> " + folderName + " " + stdErrFile + " \n"
						
	log.Println(kiraString)
	
}