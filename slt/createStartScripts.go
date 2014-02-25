package slt

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
// 	"time"
)

// CreateStartScripts create the start scripts (kira launch and PBS launch for the ICs).
func CreateStartScripts (icsName, randomNumber, simTime string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}

	var (
		home string
		scratch string
		run string
		rnd string
		runString string
		thisFolder string
		absFolderName string
		kiraOutName string
		pbsOutName string
		icsRegString string = `ics-`+conf.BaseName()+`-run(\d+)-rnd(\d+).txt`
		icsRegexp *regexp.Regexp = regexp.MustCompile(icsRegString)
		icsRegResult []string
	)
	
// 	tGlob0 := time.Now()
	
	if conf.Machine == "eurora" {
		home = "/eurora/home/userexternal/" + conf.UserName
	} else if conf.Machine == "plx" {
		home = "/plx/userexternal/" + conf.UserName
	} else {
		log.Println("Uknown machine name ", conf.Machine)
		log.Fatal("I don't know how to create the home folder path")
	}
	
	scratch = "/gpfs/scratch/userexternal/" + conf.UserName
	
	log.Println("Extracting parameters from ICs name assuming regexp:")
	fmt.Println(icsRegString)
	icsRegResult = icsRegexp.FindStringSubmatch(icsName); 
	if icsRegResult == nil {
		log.Fatal("Can't find parameters in ICs name ", icsName)
	}
	
	run  = icsRegResult[1]
	rnd  = icsRegResult[2]
	
	thisFolder = "cineca-comb" + conf.CombStr() + "-run1_10-NCM" + conf.NcmStr() + "-fPB" + 
					conf.FpbCmpStr() + "-W" + conf.WStr() + "-Z" + conf.ZCmpStr()
	absFolderName = filepath.Join(scratch, conf.Machine + "-parameterSpace", thisFolder)
	
	runString = "-run" + run + "-rnd" + rnd
	kiraOutName = "kiraLaunch-" + conf.BaseName() + runString + ".sh"
	pbsOutName = "PBS-" + conf.BaseName() + runString + ".sh"
	
	log.Println("Creating kira and PBS scripts")
	
	CreateKira (kiraOutName, absFolderName, home, run, rnd, randomNumber, simTime, conf)
	CreatePBS (pbsOutName, kiraOutName, absFolderName, run, rnd, conf)
	
// 	tGlob1 := time.Now()
// 	fmt.Println()
// 	log.Println("Wall time for creating kira and PBS scripts ", tGlob1.Sub(tGlob0))
// 	fmt.Println()
}