package sla
/*
import (
	"fmt"
	"log"
	"time"
// 	"reflect"
)

func LifeTime() {
	
	tGlob0 := time.Now()
	
	var starMap StarMapType
	
	//FIXME: istogrammare i tempi di vita ma prima
	//TODO: far funzionare i moduli esterni
	
	// Map with key = starId and value StarData struct
	starMap = make(StarMapType)

	hidePromiscuous := true
	
	starMap.Populate(inPath, inFile, regStringDBHAll, hidePromiscuous)
// 	starMap.Populate(inPath, inFile, regStringAllFishes)
	
	// Compute exchanges
	starMap.CountExchanges()
		
	var (lifeTimes001A, lifeTimes001H, lifeTimes001S, 
		 lifeTimes010A, lifeTimes010H, lifeTimes010S, 
		 lifeTimes100A, lifeTimes100H, lifeTimes100S []uint64) 
	
	for _, value := range starMap {
		value.ComputeLifeTimes()
	}
	
	for _, value := range starMap {
		if value.Z == "001" {
			lifeTimes001A = append(lifeTimes001A, value.LifeTimes.All...)
			lifeTimes001H = append(lifeTimes001H, value.LifeTimes.HardDBH...)
			lifeTimes001S = append(lifeTimes001S, value.LifeTimes.SoftDBH...)
		} else if value.Z == "010" {
			lifeTimes010A = append(lifeTimes010A, value.LifeTimes.All...)
			lifeTimes010H = append(lifeTimes010H, value.LifeTimes.HardDBH...)
			lifeTimes010S = append(lifeTimes010S, value.LifeTimes.SoftDBH...)
		} else if value.Z == "100" {
			lifeTimes100A = append(lifeTimes100A, value.LifeTimes.All...)
			lifeTimes100H = append(lifeTimes100H, value.LifeTimes.HardDBH...)
			lifeTimes100S = append(lifeTimes100S, value.LifeTimes.SoftDBH...)
		} else {
			log.Fatal("Found strange Z")
		}
	}
	
// 	fmt.Println("####################################")
	fmt.Print("\nZ001 All ")
	for _, item := range lifeTimes001A {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ001 Hard ")
	for _, item := range lifeTimes001H {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ001 Soft ")
	for _, item := range lifeTimes001S {
		fmt.Print(item, " ")
	}
	
// 	fmt.Println("####################################")
	fmt.Print("\nZ010 All ")
	for _, item := range lifeTimes010A {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ010 Hard ")
	for _, item := range lifeTimes010H {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ010 Soft ")
	for _, item := range lifeTimes010S {
		fmt.Print(item, " ")
	}
	
// 	fmt.Println("####################################")
	fmt.Print("\nZ100 All ")
	for _, item := range lifeTimes100A {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ100 Hard ")
	for _, item := range lifeTimes100H {
		fmt.Print(item, " ")
	}
	fmt.Print("\nZ100 Soft ")
	for _, item := range lifeTimes100S {
		fmt.Print(item, " ")
	}

	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
} //END MAIN
*/