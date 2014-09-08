package main

import (
	"log"
	"os"
	"time"

	"github.com/brunetto/goutils/debug"
	"bitbucket.org/brunetto/sltools/slt"
)

func main() {
	if true {
		defer debug.TimeMe(time.Now())
	}
	var (
		icsName, machine, remainingTime, randomSeed string
		cssInfo = make(chan map[string]string, 1)
		pbsLaunchChannel = make(chan string, 100)
		done = make(chan struct{})
	)

	if len(os.Args) < 5 {
		log.Fatal(`You MUST specify the ICs file name, the machine, the remainingTime 
		and the random seed for which to create start scripts!!!`)
	}
	icsName = os.Args[1]
	machine = os.Args[2]
	remainingTime = os.Args[3]
	
	if os.Args[4] == "0" {
		randomSeed = ""
	} else {
		randomSeed = os.Args[4]
	}
	
	
	go slt.CreateStartScripts(cssInfo, machine, pbsLaunchChannel, done)
	// Condumes pbs file names
	go func (pbsLaunchChannel chan string) {
		for _ = range pbsLaunchChannel {
		}
	} (pbsLaunchChannel)

	cssInfo <- map[string]string{
			"remainingTime": remainingTime,
			"randomSeed": randomSeed,
			"newICsFileName": icsName,
		}
	
	close(cssInfo)
	<-done // wait the goroutine to finish
}
