package slt

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
)

var (
	home string
	scratch string
	comb string
	ncm string
	w string
	z string
	run string
	rnd string
	fpb string
	thisFolder string
	absFolderName string
	baseName string
	kiraOutName string
	icsRegString string = `ics-cineca-comb(\d+)-NCM(\d+)-fPB(\d+)` + 
							`-W(\d)-Z(\d+)-run(\d+)-rnd(\d+).txt`
	icsRegexp *regexp.Regexp = regexp.MustCompile(icsRegString)
	icsRegResult []string
	)

func CreateScripts (icsName, machine, userName, randomNumber, simTime, pName string) () {
	
	if Verb {
		log.Println("You inserted: " )
		fmt.Println("icsName = ", icsName)
		fmt.Println("machine = ", machine)
		fmt.Println("random = ", randomNumber)
		fmt.Println("time = ", simTime)
		fmt.Println("userName = ", userName)
	}
	
	if icsName == "" ||
		machine == "" ||
		userName == "" ||
		pName == "" {
			log.Println("You must specify icsName, machine, user!!!")
			log.Fatal("Type 'sltools help createScripts' for help.")
		}
	
	if machine == "eurora" {
		home = "/eurora/home/userexternal/" + userName
	} else if machine == "plx" {
		home = "/plx/userexternal/" + userName
	} else {
		log.Fatal("Uknown machine name ", machine)
	}
	
	scratch = "/gpfs/scratch/userexternal/" + userName
	
	log.Println("Extracting parameters from ICs name")
	icsRegResult = icsRegexp.FindStringSubmatch(icsName); 
	if icsRegResult == nil {
		log.Fatal("Can't find parameters in ICs name ", icsName)
	}
	
	comb = icsRegResult[1]
	ncm  = icsRegResult[2]
	fpb  = icsRegResult[3]
	w    = icsRegResult[4]
	z    = icsRegResult[5]
	run  = icsRegResult[6]
	rnd  = icsRegResult[7]
	
	thisFolder = "cineca-comb" + comb + "-run1_10-NCM" + ncm + "-fPB" + 
					fpb + "-W" + w + "-Z" + z
	absFolderName = filepath.Join(scratch, machine+"-parameterSpace", thisFolder)
	baseName = "cineca-comb" + comb + "-NCM" + ncm + "-fPB" +
				 fpb + "-W" + w + "-Z" + z + "-run" + run + "-rnd" + rnd
				 
	kiraOutName = "kiraLaunch-" + baseName + ".sh"
	
	log.Println("Creating kira and PBS scripts")
	
	CreateKira (randomNumber, simTime)
	CreatePBS (pName)
	
}