package main

import (
	"bitbucket.org/brunetto/sltools/slt"
	"fmt"
	"log"
	"time"
)

func main() {
	
	tGlob0 := time.Now()
	fmt.Println("===========================================================")
	fmt.Println("========================== START ==========================")
	fmt.Println("===========================================================")
	
	slt.InitCommands()
	slt.SlToolsCmd.Execute()
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
	fmt.Println()
	
	fmt.Print("\x07") // Try to beep!!:)
	fmt.Println("===========================================================")
	fmt.Println("=========================== END ===========================")
	fmt.Println("===========================================================")
} 

