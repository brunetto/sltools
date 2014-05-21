package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)

func main() {
	if true {
		defer debug.TimeMe(time.Now())
	}
	var (
		icsName, machine, remainingTime, randomSeed string
		cssInfo chan map[string]string
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
	
	go slt.CreateStartScripts(cssInfo chan map[string]string, machine string)

	cssInfo <- map[string]string{
			"remainingTime": remainingTime,
			"randomSeed": randomSeed,
			"newICsFileName": icsName
		}
	
	close(cssInfo)
	
	// meglio avere un done chan struct{}{}?
}
