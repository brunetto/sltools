package slt

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func CreateICs (conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	// This variables are private to this function
	var (
		err error
		folderName string				// will contain the realizations of this combination
		icsBaseCmd string				// common base command to create the ics
		icsCmd string					// complete ICs creation command (contains the output file name)
		outIcsName string				// final ICs name
		outIcsScriptName string			// name for the ICs creation script
		icsScriptFile *os.File			// file obj for the ICs script
		icsScriptWriter *bufio.Writer	// writer for the ICs script
	)
	
	// Check we know where the binaries for the ICs are... not checking its existance
	if conf.BinFolder == "" {
		log.Fatal("I need to know where binaries for ICs are, no folder found in conf struct")
	}
	
	// ICs binaries with path
	makeking        := filepath.Join(conf.BinFolder, "makeking")
	makemass        := filepath.Join(conf.BinFolder, "makemass")
	makesecondary   := filepath.Join(conf.BinFolder, "makesecondary")
	add_star        := filepath.Join(conf.BinFolder, "add_star")
	scale           := filepath.Join(conf.BinFolder, "scale")
	makebinary      := filepath.Join(conf.BinFolder, "makebinary")
	
	// Base ICs script commands in a string, it misses the ICs filename
	icsBaseCmd = "#!/bin/bash\n" +
			"set -xeu\n" +
			makeking + " -n " + conf.NcmStr() + 
					   " -w " + conf.WStr() + 
					   " -i -u \\\n" +
	 "| " + makemass + " -f 8  -l 0.1 -u 150 \\\n" +
	 "| " + makesecondary + " -f " + conf.FpbStr() + 
							" -q -l 0.1 \\\n" +
	 "| " + add_star + " -R 1 -Z " + conf.ZStr() + " \\\n" +
	 "| " + scale + " -R 1 -M 1\\\n" +
	 "| " + makebinary + " -f 2 -o 1 -l 1 -u 107836.09 \\\n" +
	 "> "  // Redirect output to the proper ICs file

	// Assemble folder name, create it and go into
	folderName = "cineca-comb" + conf.CombStr() + 
				"-run1_" + conf.RunsStr() + 
				"-NCM" + conf.NcmStr() + 
				"-fPB" + conf.FpbCmpStr() + 
				"-W" + conf.WStr() + 
				"-Z" + conf.ZCmpStr()
					
	log.Println("Create folder and change to it:", folderName)
	if err = os.Mkdir(folderName, 0700); err != nil {log.Fatal("Can't create folder ", err)}
	
	// Copy config file inside folder to be read and for backup
	_, err = CopyFile(ConfName, filepath.Join(folderName, ConfName))
	// Go into the new folder
	if err = os.Chdir(folderName); err!= nil {
		log.Println("Error while entering in folder ", folderName)
		log.Fatal(err)
	}
	
	// Create the scripts and run them
	for runIdx :=1; runIdx<conf.Runs+1; runIdx++ {
		// Basename suffix
		runString := "-run" +  	LeftPad(strconv.Itoa(runIdx), "0", 2) + "-rnd00"
		// ICs final name
		outIcsName = "ics-" + conf.BaseName() + runString + ".txt"
		// Add ICs final file name to ICs creation command
		icsCmd = icsBaseCmd + outIcsName
		// ICs creation script name
		outIcsScriptName = "create_IC-" + conf.BaseName() + runString + ".sh"
		
		// Write the script file
		if icsScriptFile, err = os.Create(outIcsScriptName); err != nil {log.Fatal(err)}
		defer icsScriptFile.Close()
		icsScriptWriter = bufio.NewWriter(icsScriptFile)
		defer icsScriptWriter.Flush()
		if _, err = icsScriptWriter.WriteString(icsCmd); err != nil {
			log.Fatal("Error while writing ", outIcsScriptName, err)
		}
		
		// Run it
		if _, err := os.Stat(outIcsScriptName); err == nil {
			log.Println(outIcsScriptName, " exists, try to run")
		}
		// BUG: this does not work, damn it!!!:(
		bashCmd := exec.Command("bash", "-x", outIcsScriptName)
// 		bashCmd := exec.Command("/bin/bash", "-c", "sleep 3000") // this works, bastard!!!
		bashCmd.Stdout = os.Stdout
		bashCmd.Stderr = os.Stderr
		if err := bashCmd.Run(); err != nil {
			log.Fatal(err)
		}
	
		// Create kiraLaunch and PBSlaunch scripts with the same functions used in Continue
		icsRandomSeed := "" // let SL decide it
		CreateStartScripts("ics-" + conf.BaseName() + runString + ".txt", icsRandomSeed, conf.EndTimeStr(), conf)
		
	}
}
