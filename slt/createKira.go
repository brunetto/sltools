package slt

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func CreateKira (randomNumber, simTime string) () {
	var (
		kiraFile *os.File
		kiraWriter *bufio.Writer
		kiraString string
		err error
		randomString string
		kiraBinPath string
		stdOutFile string
		stdErrFile string
	)
	
	if rnd == "00" {
		randomString = ""
	} else {
		randomString = "-s " + randomNumber
	}
	
	kiraBinPath = filepath.Join(home, "slpack", "starlab", "usr", "bin", "kira")
	
	stdOutFile = "out-" + baseName + ".txt"
	stdErrFile = "err-" + baseName + ".txt"
	
	// I know I can use `` but I don't like the string not to be align with the 
	// rest of the code
	kiraString = "echo $PWD\n" + 
				  "echo $LD_LIBRARY_PATH\n" + 
				  kiraBinPath + " -t " + simTime + " -d 1 -D 1 -b 1 -f 0.3 \\\n" + 
				  " -n 10 -e 0.000 -B " + randomString + " \\\n" +
				  "<  " + filepath.Join(absFolderName, icsName) + " \\\n" +
				  ">  " + filepath.Join(absFolderName, stdOutFile) + " \\\n" +
				  "2> " + filepath.Join(absFolderName, stdErrFile) + " \n"
	
	if kiraFile, err = os.Create(kiraOutName); err != nil {log.Fatal(err)}
	defer kiraFile.Close()
	
	kiraWriter = bufio.NewWriter(kiraFile)
	defer kiraWriter.Flush()
	
	kiraWriter.WriteString(kiraString)
}









