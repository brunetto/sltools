package slt

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func CreateICs (binFolder string) () {
	if Debug {Whoami(true)}
	
	// This variables are private to this function
	var (
		err error
		runs string						// number of cluster of this type = number of run = number of stat realization
		comb string						// combination number, see table
		ncm string						// number of centres of masses, both single or binary
		fpb string						// fraction of primordial binaries
		w string						// central adimensional potential
		z string						// metallicity
		folderName string				// will contain the realizations of this combination
		icsBaseCmd string				// common base command to create the ics
		baseName string					// base name for the ICs and the scripts
		icsCmd string					// complete ICs creation command (contains the output file name)
		outIcsName string				// final ICs name
		outIcsScriptName string			// name for the ICs creation script
		icsScriptFile *os.File			// file obj for the ICs script
		icsScriptWriter *bufio.Writer	// writer for the ICs script
	)
	
	// Check for configuration file passed by -c flag
	if ConfName == "" {
		log.Fatal("Provide a JSON configuration file")
	}
	
	// Read conf file and create conf struct
	log.Println("Read configuration form ", ConfName)
	conf = new(Config)
	conf.ReadConf(ConfName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}

	// Check we know where the binaries for the ICs are... not checking its existance
	if binFolder == "" {
		if conf.BinFolder == "" {
			log.Fatal("I need to know where binaries for ICs are.")
		} else {
			binFolder = conf.BinFolder
		}
	}
	
	// Just here to be faster
	/*binFolder = "/home/brunetto/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/"*/

	// Convert to string the simulation parameter to be used to create the scripts
	runs = strconv.Itoa(conf.Runs)
	comb = strconv.Itoa(conf.Comb)
	ncm = strconv.Itoa(conf.Ncm)
	fpb = strings.Replace(strconv.FormatFloat(conf.Fpb, 'f', 2, 64), ".",  "", -1)
	w = strconv.Itoa(conf.W)
	z = strings.Replace(strconv.FormatFloat(conf.Z, 'f', 2, 64), ".",  "", -1)
	
	// ICs binaries with path
	makeking        := filepath.Join(binFolder, "makeking")
	makemass        := filepath.Join(binFolder, "makemass")
	makesecondary   := filepath.Join(binFolder, "makesecondary")
	add_star        := filepath.Join(binFolder, "add_star")
	scale           := filepath.Join(binFolder, "scale")
	makebinary      := filepath.Join(binFolder, "makebinary")
	
	// Base ICs script commands in a string
	icsBaseCmd = "#!/bin/bash\n" +
			"set -e -u\n" +
			makeking + " -n " + ncm + " -w " + w + " -i -u \\\n" +
			"| " + makemass + " -f 8  -l 0.1 -u 150 \\\n" +
			"| " + makesecondary + " -f " + strconv.FormatFloat(conf.Fpb, 'f', 2, 64) + 
																	" -q -l 0.1 \\\n" +
			"| " + add_star + " -R 1 -Z " + strconv.FormatFloat(conf.Z, 'f', 2, 64) + 
																			" \\\n" +
			"| " + scale + " -R 1 -M 1\\\n" +
			"| " + makebinary + " -f 2 -o 1 -l 1 -u 107836.09 \\\n" +
			"> "  

	// Assemble folder name, create it and go into
	folderName = "cineca-comb" + comb + "-run1_" + runs + "-NCM" + ncm + "-fPB" + fpb + "-W" + w + "-Z" + z
	log.Println("Create folder and change to it:", folderName)
	if err = os.Mkdir(folderName, 0700); err != nil {log.Fatal("Can't create folder ", err)}
	// Copy config file inside folder to be read and for backup
	_, err = CopyFile(ConfName, filepath.Join(folderName, ConfName))
	if err = os.Chdir(folderName); err!= nil {
		log.Println("Error while entering in folder ", folderName)
		log.Fatal(err)
	}
	
	// Create the scripts and run them
	for runIdx :=1; runIdx<2/*conf.Runs+1*/; runIdx++ {
		baseName = "cineca-comb" + comb + "-NCM" + ncm + "-fPB" +
					fpb + "-W" + w + "-Z" + z + "-run" +  
					LeftPad(strconv.Itoa(runIdx), "0", 2) + "-rnd00"
		outIcsName = filepath.Join(/*folderName, */"ics-" + baseName + ".txt")
		outIcsScriptName = filepath.Join(/*folderName,*/ "create_IC-" + baseName + ".sh")
		icsCmd = icsBaseCmd + outIcsName
		log.Println(icsCmd)
		
		// Write the script file
		if icsScriptFile, err = os.Create(outIcsScriptName); err != nil {log.Fatal(err)}
		defer icsScriptFile.Close()
		icsScriptWriter = bufio.NewWriter(icsScriptFile)
		defer icsScriptWriter.Flush()
		if _, err = icsScriptWriter.WriteString(icsCmd); err != nil {
			log.Fatal("Error while writing ", outIcsScriptName, err)
		}
		
		// Run it
// 		log.Println("Creating ICs files with: bash", outIcsScriptName)
		log.Println("Creating ICs files with: bash", "create_IC-cineca-comb20-NCM10000-fPB010-W9-Z010-run01-rnd00.sh")
		if _, err := os.Stat("create_IC-cineca-comb20-NCM10000-fPB010-W9-Z010-run01-rnd00.sh"); err == nil {
			log.Printf("file exists; processing...")
		}
		bashCmd := exec.Command("bash", "-x", outIcsScriptName)
// 		bashCmd := exec.Command("/bin/bash", "-c", "sleep 3000")
		bashCmd.Stdout = os.Stdout
		bashCmd.Stderr = os.Stderr
		if err := bashCmd.Run(); err != nil {
			log.Fatal(err)
		}
	
		// Create kiraLaunch and PBSlaunch scripts with the same functions used in Continue
		CreateStartScripts ("ics-" + baseName + ".txt", conf.Machine, conf.UserName, "", strconv.Itoa(conf.EndTime), conf.PName)
		
	}

	
	
	
}
