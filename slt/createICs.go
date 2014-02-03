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
	
	var (
		err error
		runs string
		comb string
		ncm string
		fpb string
		w string
		z string	
		folderName string
		icsBaseCmd string
		baseName string
		icsCmd string
		outIcsName string
		outIcsScriptName string
		icsScriptFile *os.File
		icsScriptWriter *bufio.Writer
	)
	
	if ConfName == "" {
		log.Fatal("Provide a JSON configuration file")
	}
	
	log.Println("Read configuration form ", ConfName)
	conf = new(Config)
	conf.ReadConf(ConfName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}

	if binFolder == "" {
		if conf.BinFolder == "" {
			log.Fatal("I need to know where binaries for ICs are.")
		} else {
			binFolder = conf.BinFolder
		}
	}
	
	binFolder = "/home/ziosi/Dropbox/Research/PhD_Mapelli/4-ParameterSpace/Simulations/CINECA_SCRATCH/bin/"

	runs = strconv.Itoa(conf.Runs)
	comb = strconv.Itoa(conf.Comb)
	ncm = strconv.Itoa(conf.Ncm)
	fpb = strings.Replace(strconv.FormatFloat(conf.Fpb, 'f', 2, 64), ".",  "", -1)
	w = strconv.Itoa(conf.W)
	z = strings.Replace(strconv.FormatFloat(conf.Z, 'f', 2, 64), ".",  "", -1)
	
	makeking        := filepath.Join(binFolder, "makeking")
	makemass        := filepath.Join(binFolder, "makemass")
	makesecondary   := filepath.Join(binFolder, "makesecondary")
	add_star        := filepath.Join(binFolder, "add_star")
	scale           := filepath.Join(binFolder, "scale")
	makebinary      := filepath.Join(binFolder, "makebinary")
	
	icsBaseCmd = "#!/bin/bash\n" +
			makeking + " -n " + ncm + " -w " + w + " -i -u \\\n" +
			"| " + makemass + " -f 8  -l 0.1 -u 150 \\\n" +
			"| " + makesecondary + " -f " + strconv.FormatFloat(conf.Fpb, 'f', 2, 64) + 
																	" -q -l 0.1 \\\n" +
			"| " + add_star + " -R 1 -Z " + strconv.FormatFloat(conf.Z, 'f', 2, 64) + 
																			" \\\n" +
			"| " + scale + " -R 1 -M 1\\\n" +
			"| " + makebinary + " -f 2 -o 1 -l 1 -u 107836.09 \\\n" +
			"> "  
	
	folderName = "cineca-comb" + comb + "-run1_" + runs + "-NCM" + ncm + "-fPB" + fpb + "-W" + w + "-Z" + z
	log.Println("Create folder and change to it:", folderName)
	if err = os.Mkdir(folderName, 0700); err != nil {log.Fatal("Can't create folder ", err)}
	if err = os.Chdir(folderName); err!= nil {
		log.Println("Error while entering in folder ", folderName)
		log.Fatal(err)
	}
	
	for runIdx :=1; runIdx<2/*conf.Runs+1*/; runIdx++ {
		baseName = "cineca-comb" + comb + "-NCM" + ncm + "-fPB" +
					fpb + "-W" + w + "-Z" + z + "-run" +  
					LeftPad(strconv.Itoa(runIdx), "0", 2) + "-rnd00"
		outIcsName = filepath.Join(/*folderName, */"ics-" + baseName + ".txt")
		outIcsScriptName = filepath.Join(/*folderName,*/ "create_IC-" + baseName + ".sh")
		icsCmd = icsBaseCmd + outIcsName
		
		if icsScriptFile, err = os.Create(outIcsScriptName); err != nil {log.Fatal(err)}
		defer icsScriptFile.Close()
		
		icsScriptWriter = bufio.NewWriter(icsScriptFile)
		defer icsScriptWriter.Flush()
		
		if _, err = icsScriptWriter.WriteString(icsCmd); err != nil {
			log.Fatal("Error while writing ", outIcsScriptName, err)
		}
		
		log.Println("Creating ICs files with: bash", outIcsScriptName)
		bashCmd := exec.Command("bash", outIcsScriptName)
		bashCmd.Stdout = os.Stdout
		bashCmd.Stderr = os.Stderr
		if err := bashCmd.Run(); err != nil {
			log.Fatal(err)
		}
		
		CreateStartScripts ("ics-" + baseName + ".txt", conf.Machine, conf.UserName, "", strconv.Itoa(conf.EndTime), conf.PName)
		
	}

	
	
	
}
