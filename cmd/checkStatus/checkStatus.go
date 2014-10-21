package main

import (
	"fmt"
	"time"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	slt.CheckStatus()
		
	fmt.Print("\x07") // Beep when finish!!:D
}
