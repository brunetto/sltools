package slan

import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"time"
// 
// 	"github.com/brunetto/goutils/debug"
)

// func RunAll(inFileName string) {
// 	if Debug {
// 		defer debug.TimeMe(time.Now())
// 	}
// 	var (
// 		exchOutFile   *os.File
// 		lTOutFile     *os.File
// 		wasInDBHFile  *os.File
// 		lastInDBHFile *os.File
// 		err           error
// 	)
// 
// 	tGlob0 := time.Now()
// 
// 	inPath := filepath.Dir(inFileName)
// 	inFile := filepath.Base(inFileName)
// 
// 	starMap := make(StarMapType)
// 
// 	starMap.Populate(inPath, inFile)
// 
// 	// Open files
// 	if exchOutFile, err = os.Create("files/exchanges-from-" + inFile + ".dat"); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer exchOutFile.Close()
// 	if lTOutFile, err = os.Create("lifetimes.dat"); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer lTOutFile.Close()
// 
// 	if wasInDBHFile, err = os.Create("files/wasInDBH-from-" + inFile + ".dat"); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer wasInDBHFile.Close()
// 
// 	if lastInDBHFile, err = os.Create("files/lastInDBH-from-" + inFile + ".dat"); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer lastInDBHFile.Close()
// 
// 	starMap.ExecOnAll("CountExchanges")
// 	// 	starMap.SaveExch(exchOutFile)
// 	// 	starMap.PrintExchStats(os.Stderr)
// 
// 	starMap.ExecOnAll("ComputeLifeTimes")
// 
// 	// 	allLT := starMap.CollectLifeTimes()
// 	// 	allLT.SaveLifeTimes(lTOutFile)
// 	// 	allLT.PrintLTStats()
// 
// 	// Only DBH
// 	wasInDBH := starMap.WasInType("bh|bh")
// 	wasInDBH.SaveExch(wasInDBHFile)
// 
// 	// 	lastInDBH := starMap.ExtrcLastType("bh|bh")
// 	// 	lastInDBH.SaveExch(lastInDBHFile)
// 
// 	promiscouosMap := wasInDBH.WasPromiscuous()
// 
// 	log.Println("PStars ", len(promiscouosMap))
// 
// 	tGlob1 := time.Now()
// 	fmt.Println()
// 	log.Println("Wall time for reading the file ", tGlob1.Sub(tGlob0))
// 	fmt.Println()
// }
