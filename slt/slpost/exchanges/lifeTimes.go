package exchanges

import (
	"fmt"
	"log"
	"time"
// 	"reflect"
)

func main() {
	
	tGlob0 := time.Now()
	
	var inPath string
	var inFile string
	var starMap StarMapType
	
	// Regexp string for all_the_fishes
// 	var regStringAllFishes = `^(\d{3})\s+` +// GROUP 1: Z eg 001
// 								`(\d{3})\s+` + // GROUP 2: n eg 001
// 								`(\S+)\s+`+ // GROUP 3: binary_ids Z001n001idsa12550b2550
// 								`(\d+)\s+` + // GROUP 4: sys_time eg 0
// 								`(\d+\.\d+)\s+`+ // GROUP 5: phys_time [Myr] 0.0slt
// 								`(\S+)\s+`+ // GROUP 6: objects_ids eg  2550|12550
// 								`(\S)\s+`+ // GROUP 7: hardflag eg H 
// 								`(\S+)\s+`+ // GROUP 8: types eg ns++|ns++
// 								`(\S+\.\S+)\s+`+ // GROUP 9: masse[0] eg 10.3837569427
// 								`(\S+\.\S+)\s+`+// GROUP 10: mass[1] eg 9.2141789593
// 								`(\S+\.\S+)\s+`+// GROUP 11: sma eg  3.6333e-05
// 								`(\S+\.\S+)\s+`+// GROUP 12: sma eg  4.6156152408e-06
// 								`(\S+\.*\S*)`// GROUP 13: sma eg   0.680846
// 												// NOTE: Maybe ecc is zero...
	
	var regStringDBHAll = `^(\d{3})\,\s+` +// GROUP 1: Z eg 001
								`(\d{1,3})\,\s+` + // GROUP 2: n eg 001
								`(\S+)\,\s+`+ // GROUP 3: binary_ids Z001n001idsa12550b2550
								`(\d+)\,\s+` + // GROUP 4: sys_time eg 0
								`(\d+\.\d+)\,\s+`+ // GROUP 5: phys_time [Myr] 0.0
								`(\S+)\,\s+`+ // GROUP 6: objects_ids eg  2550|12550
								`(\S)\,\s+`+ // GROUP 7: hardflag eg H 
								`(\S+)\,\s+`+ // GROUP 8: types eg ns++|ns++
								`(\S+\.\S+)\,\s+`+ // GROUP 9: masse[0] eg 10.3837569427
								`(\S+\.\S+)\,\s+`+// GROUP 10: mass[1] eg 9.2141789593
								`(\S+\.\S+)\,\s+`+// GROUP 11: sma eg  3.6333e-05
								`(\S+\.\S+)\,\s+`+// GROUP 12: sma eg  4.6156152408e-06
								`(\S+\.*\S*)`// GROUP 13: sma eg   0.680846
// // 												NOTE: Maybe ecc is zero...

	inPath = "../../../data/2013-10-10-analysis/03-final/"
// 	inFile = "all_the_fishes.txt"
	// because we want to know the lifetimes of BBHs 
	inFile = "bh-bh_all.txt"
	
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
