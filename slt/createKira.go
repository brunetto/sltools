package slt

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brunetto/goutils/debug"
)

// CreateKira create the kira script to launch the simulation.
func CreateKira(kiraOutName, absFolderName, home, run, rnd, randomNumber, simTime string, conf *ConfigStruct) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		kiraFile     *os.File
		kiraWriter   *bufio.Writer
		kiraString   string
		err          error
		randomString string
		kiraBinPath  string
		stdOutFile   string
		stdErrFile   string
		icsName      string
	)

	if rnd == "00" {
		randomString = ""
		simTime = "500"
	} else {
		randomString = "-s " + randomNumber
	}

	kiraBinPath = "/eurora/home/userexternal/bziosi00/starlabjune19_2013/usr/bin/kira"
	//filepath.Join(home, "slpack", "starlab", "usr", "bin", "kira")

	runString := "-run" + run + "-rnd" + rnd

	stdOutFile = "out-" + conf.BaseName() + runString + ".txt"
	stdErrFile = "err-" + conf.BaseName() + runString + ".txt"
	icsName = "ics-" + conf.BaseName() + runString + ".txt"

	// I know I can use `` but I don't like the string not to be align with the
	// rest of the code
	kiraString = "echo $PWD\n" +
		"echo $LD_LIBRARY_PATH\n" +
		kiraBinPath + " -t " + simTime + " -d 1 -D 1 -b 1 -f 0 \\\n" +
		" -n 10 -e 0.000 -B " + randomString + " \\\n" +
		"<  " + filepath.Join(absFolderName, icsName) + " \\\n" +
		">  " + filepath.Join(absFolderName, stdOutFile) + " \\\n" +
		"2> " + filepath.Join(absFolderName, stdErrFile) + " \n"

	log.Println("Write kira launch script to ", kiraOutName)
	if kiraFile, err = os.Create(kiraOutName); err != nil {
		log.Fatal(err)
	}
	defer kiraFile.Close()

	kiraWriter = bufio.NewWriter(kiraFile)
	defer kiraWriter.Flush()

	kiraWriter.WriteString(kiraString)
}
