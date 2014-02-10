package main

import (
	"bitbucket.org/brunetto/sltools/slt"
	"fmt"
	"log"
	"time"
)

func main() {
	
	tGlob0 := time.Now()
	
	slt.InitCommands()
	slt.SlToolsCmd.Execute()
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
} 

