package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/brunetto/sltools/slt"
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
	
	fmt.Print("\x07") // Beep when finish!!:D
	fmt.Println("===========================================================")
	fmt.Println("=========================== END ===========================")
	fmt.Println("===========================================================")
} 

