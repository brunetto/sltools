package slt

import (
	"fmt"
	"log"
	"time"
)

func StichOutput (inFileTmpl string) () {
	
	var (
	)
	
	tGlob0 := time.Now()
	
	if inFileTmpl == "" {
		log.Fatal("You need to specify a STDOUT input file template with the -i flag!!!")
	}
	
	//
	// STDOUT
	//
	
	StdOutStich (inFileTmpl)
	
	//
	// STDERR
	//
	
	StdErrStich (inFileTmpl)
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for stich output ", tGlob1.Sub(tGlob0))
}