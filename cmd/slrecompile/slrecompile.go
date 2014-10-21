package main

import (
	"log"
	"os"
	"strings"
	"time"
	
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/sltools/slt"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var confArgs []string

	if len(os.Args) < 2 {
		log.Println("Configure without flags!!!!!!!")
		log.Println("Assuming --with-f77=no")
		confArgs = []string{"--with-f77=no"}
	} else {
		if strings.Contains(os.Args[1], "help") {
			log.Fatal("Run as 'slrecompile <configure arguments>'")
		}
		confArgs = os.Args[1:]
	}
	slt.SLrecompile(confArgs)
}


