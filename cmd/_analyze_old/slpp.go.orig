package main

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/brunetto/slpp/sla"
)

func main() {
	
	tGlob0 := time.Now()
	fmt.Fprintf(os.Stderr, "===========================================================\n")
	fmt.Fprintf(os.Stderr, "========================== START ==========================\n")
	fmt.Fprintf(os.Stderr, "===========================================================\n")
	
	sla.InitCommands()
	sla.SlPpCmd.Execute()
	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
	fmt.Println()
	
	fmt.Print("\x07") // Beep when finish!!:D
	fmt.Fprintf(os.Stderr, "===========================================================\n")
	fmt.Fprintf(os.Stderr, "=========================== END ===========================\n")
	fmt.Fprintf(os.Stderr, "===========================================================\n")
} 

