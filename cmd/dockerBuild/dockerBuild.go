package main 

import (
	"log"
	"os"
	"strings"
	"time"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	if len(os.Args) < 2 {
		log.Fatal("Please give me a target host to build for.")
	}
	
	if strings.Contains(os.Args[1], "help") || strings.Contains(os.Args[1], "Help") {
		log.Fatal(`Use like: ./dockerBuild <hostname>
		where hostname can be: spritz, longisland, uno
		`)
	}
	
	slt.DockerBuild(os.Args[1])
	
}
