package main 

import (
	"time"

	"github.com/brunetto/goutils/debug"
	
	"bitbucket.org/brunetto/sltools/slt"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	slt.SimClean()
}
	
