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
	var regStringAllFishes = `^(\d{3})\s+` +// GROUP 1: Z eg 001
								`(\d{3})\s+` + // GROUP 2: n eg 001
								`(\S+)\s+`+ // GROUP 3: binary_ids Z001n001idsa12550b2550
								`(\d+)\s+` + // GROUP 4: sys_time eg 0
								`(\d+\.\d+)\s+`+ // GROUP 5: phys_time [Myr] 0.0
								`(\S+)\s+`+ // GROUP 6: objects_ids eg  2550|12550
								`(\S)\s+`+ // GROUP 7: hardflag eg H 
								`(\S+)\s+`+ // GROUP 8: types eg ns++|ns++
								`(\S+\.\S+)\s+`+ // GROUP 9: masse[0] eg 10.3837569427
								`(\S+\.\S+)\s+`+// GROUP 10: mass[1] eg 9.2141789593
								`(\S+\.\S+)\s+`+// GROUP 11: sma eg  3.6333e-05
								`(\S+\.\S+)\s+`+// GROUP 12: sma eg  4.6156152408e-06
								`(\S+\.*\S*)`// GROUP 13: sma eg   0.680846
// 												// NOTE: Maybe ecc is zero...
/*	
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
*/
	inPath = "../../../data/2013-10-10-analysis/03-final/"
	// We want to know how many exchange a BH undergo before entering a binary
	// but also during its BBH life
	inFile = "all_the_fishes.txt"
// 	inFile = "bh-bh_all.txt"
	
	//FIXME: istogrammare i tempi di vita ma prima
	//TODO: far funzionare i moduli esterni
	
	// Map with key = starId and value StarData struct
	starMap = make(StarMapType)

	hidePromiscuous := true
	
	starMap.Populate(inPath, inFile, regStringAllFishes, hidePromiscuous)
// 	starMap.Populate(inPath, inFile, regStringAllFishes)
	
	// Compute exchanges
	starMap.CountExchanges()
	
	// Extract only DCOB binaries
	// DBH
	starMapDBH := make(StarMapType)
	for key, value := range starMap {
		if value.DCOB == "DBH" {
			starMapDBH[key] = value
		}
	}
	
	// DNS
	starMapDNS := make(StarMapType)
	for key, value := range starMap {
		if value.DCOB == "DNS" {
			starMapDNS[key] = value
		}
	}
	
	// BHNS
	starMapBHNS := make(StarMapType)
	for key, value := range starMap {
		if value.DCOB == "BHNS" {
			starMapBHNS[key] = value
		}
	}
	
	
	/* PRIMO OUTPUT: EXCHANGES
	log.Println("Obtained data for ", len(starMap), " objects")
	
	var keys []string
	var value *StarData
	var key string
	
	fmt.Printf("###################################################################################\n")
	fmt.Println("# Extracted ", len(starMapDBH), " DBH objects")
	
	fmt.Printf("###################################################################################\n")
	fmt.Printf("# starId\t\t\tZ\tLastType\tPrimordial\tHard\tSoft\tTotal\n")
	fmt.Printf("###################################################################################\n")
	
	keys = starMapDBH.Keys()
	for idx:=0; idx<len(keys); idx++ {
		key = keys[idx]
		value = starMapDBH[key]
		fmt.Printf("%-25v\t%v\t%v\t%v\t\t%v\t%v\t%v\n", key, value.Z, value.DCOB, value.ExchangesNumbers.Primordial, value.ExchangesNumbers.HardExchangesNumber,  value.ExchangesNumbers.SoftExchangesNumber, value.ExchangesNumbers.TotalExchanges)
	}
	
	fmt.Printf("###################################################################################\n")
	fmt.Println("# Extracted ", len(starMapDNS), " DNS objects")
	
	fmt.Printf("###################################################################################\n")
	fmt.Printf("# starId\t\t\tZ\tLastType\tPrimordial\tHard\tSoft\tTotal\n")
	fmt.Printf("###################################################################################\n")

	keys = starMapDNS.Keys()
	for idx:=0; idx<len(keys); idx++ {
		key = keys[idx]
		value = starMapDNS[key]
		fmt.Printf("%-25v\t%v\t%v\t%v\t\t%v\t%v\t%v\n", key, value.Z, value.DCOB, value.ExchangesNumbers.Primordial, value.ExchangesNumbers.HardExchangesNumber,  value.ExchangesNumbers.SoftExchangesNumber, value.ExchangesNumbers.TotalExchanges)
	}
	
	fmt.Printf("###################################################################################\n")
	fmt.Println("# Extracted ", len(starMapBHNS), " BHNS objects")
	
	fmt.Printf("###################################################################################\n")
	fmt.Printf("# starId\t\t\tZ\tLastType\tPrimordial\tHard\tSoft\tTotal\n")
	fmt.Printf("###################################################################################\n")
	
	keys = starMapBHNS.Keys()
	for idx:=0; idx<len(keys); idx++ {
		key = keys[idx]
		value = starMapBHNS[key]
		fmt.Printf("%-25v\t%v\t%v\t%v\t\t%v\t%v\t%v\n", key, value.Z, value.DCOB, value.ExchangesNumbers.Primordial, value.ExchangesNumbers.HardExchangesNumber,  value.ExchangesNumbers.SoftExchangesNumber, value.ExchangesNumbers.TotalExchanges)
	}
	*/
/*
	// TEST PRINT
	log.Println("######################################################")
	log.Println("Zero eccentricity binaries")
	log.Println("######################################################")
	
	// Retrieve and sort Exchanges map keys
	keys := starMap.Keys()
	
	for _, key := range keys {
		if starMap[key].ZeroEcc == true {
// 			fmt.Print(key, " ")
			starMap[key].Print()
		}
	}
	
	log.Println("######################################################")
	log.Println("Promiscuous binaries")
	log.Println("######################################################")
	
	// Promiscuous binaries
	for _, key := range keys {
		if starMap[key].Promiscuous == true {
// 			fmt.Print(key, " ")
			starMap[key].Print()
		}
	}
	

	starMap["Z001n001id854"].Print()
	starMap["Z001n001id854"].CountExchanges().Print()
	starMap["Z001n001id854"].ExchangesNumbers.Print()
	*/


	/*
	 SECONDO OUTPUT: NUMBER OF BINARIES TOUCHED BY DCOB+ OBJECTS 
	

	touchedBinaries := make(map[string]string)
	for _, value := range starMap {
		if ((value.DCOB == "DBH") || (value.DCOB == "DNS") || (value.DCOB == "BHNS")) {
			fmt.Println(value.DCOB)
			for _, bData := range value.Exchanges {
				if _, exists := touchedBinaries[bData.BinaryId]; !exists {
					touchedBinaries[bData.BinaryId] = ""
				}
			}
		}
	}
	
	log.Println("Number of DCOB+ touched binaries ", len(touchedBinaries))
	 */
	
	//TERZO OUTPUT, LIFE TIMES
	/*	 */
	log.Println("========================================================")
	log.Println("========================================================")
	log.Println("LIFE TIMES")
	log.Println("========================================================")
	log.Println("========================================================")

	var sum uint64
	
	var lifeTimesAll []uint64
	for _, value := range starMap {
		lifeTimesAll = append(lifeTimesAll, value.LifeTimes().all...)
	}
// 	log.Println(lifeTimesAll)
	log.Println(len(lifeTimesAll))
	sum = uint64(0)
	for _, value := range(lifeTimesAll) {
		sum = sum + value
	}
	log.Println("lifeTimesAll", 0.254 * float64(sum) / float64(len(lifeTimesAll)))
	log.Println("========================================================")
	
	var lifeTimesHardDBH []uint64
	for _, value := range starMap {
		lifeTimesHardDBH = append(lifeTimesHardDBH, value.LifeTimes().hardDBH...)
	}
// 	log.Println(lifeTimesHardDBH)
	log.Println(len(lifeTimesHardDBH))
	sum = uint64(0)
	for _, value := range(lifeTimesHardDBH) {
		sum = sum + value
	}
	log.Println("lifeTimesHardDBH", 0.254 * float64(sum) / float64(len(lifeTimesHardDBH)))
	log.Println("========================================================")
	
	var lifeTimesHardDNS []uint64
	for _, value := range starMap {
		lifeTimesHardDNS = append(lifeTimesHardDNS, value.LifeTimes().hardDNS...)
	}
// 	log.Println(lifeTimesHardDNS)
	log.Println(len(lifeTimesHardDNS))
	sum = uint64(0)
	for _, value := range(lifeTimesHardDNS) {
		sum = sum + value
	}
	log.Println("lifeTimesHardDNS", 0.254 * float64(sum) / float64(len(lifeTimesHardDNS)))
	log.Println("========================================================")
	
	var lifeTimesHardBHNS []uint64
	for _, value := range starMap {
		lifeTimesHardBHNS = append(lifeTimesHardBHNS, value.LifeTimes().hardBHNS...)
	}
// 	log.Println(lifeTimesHardBHNS)
	log.Println(len(lifeTimesHardBHNS))
	sum = uint64(0)
	for _, value := range(lifeTimesHardBHNS) {
		sum = sum + value
	}
	log.Println("lifeTimesHardBHNS", 0.254 * float64(sum) / float64(len(lifeTimesHardBHNS)))
	log.Println("========================================================")
	
	var lifeTimesSoftDBH []uint64
	for _, value := range starMap {
		lifeTimesSoftDBH = append(lifeTimesSoftDBH, value.LifeTimes().softDBH...)
	}
// 	log.Println(lifeTimesSoftDBH)
	log.Println(len(lifeTimesSoftDBH))
	sum = uint64(0)
	for _, value := range(lifeTimesSoftDBH) {
		sum = sum + value
	}
	log.Println("lifeTimesSoftDBH", 0.254 * float64(sum) / float64(len(lifeTimesSoftDBH)))
	log.Println("========================================================")
	
	var lifeTimesSoftDNS []uint64
	for _, value := range starMap {
		lifeTimesSoftDNS = append(lifeTimesSoftDNS, value.LifeTimes().softDNS...)
	}
// 	log.Println(lifeTimesSoftDNS)
	log.Println(len(lifeTimesSoftDNS))
	sum = uint64(0)
	for _, value := range(lifeTimesSoftDNS) {
		sum = sum + value
	}
	log.Println("lifeTimesSoftDNS", 0.254 * float64(sum) / float64(len(lifeTimesSoftDNS)))
	log.Println("========================================================")
	
	var lifeTimesSoftBHNS []uint64
	for _, value := range starMap {
		lifeTimesSoftBHNS = append(lifeTimesSoftBHNS, value.LifeTimes().softBHNS...)
	}
// 	log.Println(lifeTimesSoftBHNS)
	log.Println(len(lifeTimesSoftBHNS))
	sum = uint64(0)
	for _, value := range(lifeTimesSoftBHNS) {
		sum = sum + value
	}
	log.Println("lifeTimesSoftBHNS", 0.254 * float64(sum) / float64(len(lifeTimesSoftBHNS)))
	log.Println("========================================================")
	/*in ordina solo hard, solo soft, tutte (40.4 Myr, 4.2 Myr, 14.0 Myr)... da controllare...
	 
	 ziosi@spritz:~/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/scripts/exchanges$ go run exchanges.go > lifetimes.log
	2013/10/17 15:35:39 Header detected, skip...
	2013/10/17 15:35:39 Done reading  451  lines from file with err EOF
	2013/10/17 15:35:39 ========================================================
	2013/10/17 15:35:39 ========================================================
	2013/10/17 15:35:39 LIFE TIMES
	2013/10/17 15:35:39 ========================================================
	2013/10/17 15:35:39 ========================================================
	2013/10/17 15:35:39 [1 20 110 1 19 1 131 1 1 19 1 118 1 131 1 20 110]
	2013/10/17 15:35:39 17
	2013/10/17 15:35:39 40.35294117647059
	2013/10/17 15:35:39 Wall time for all  79.004078ms
	ziosi@spritz:~/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/scripts/exchanges$ go run exchanges.go > lifetimes.log
	2013/10/17 15:35:54 Header detected, skip...
	2013/10/17 15:35:54 Done reading  451  lines from file with err EOF
	2013/10/17 15:35:54 ========================================================
	2013/10/17 15:35:54 ========================================================
	2013/10/17 15:35:54 LIFE TIMES
	2013/10/17 15:35:54 ========================================================
	2013/10/17 15:35:54 ========================================================
	2013/10/17 15:35:54 [1 1 1 1 1 1 1 1 1 33 1 1 1 1 2 1 1 1 1 1 1 1 7 1 79 1 1 1 44 1 6 1 5 1 1 1 2 1 3 3 1 1 1 8 6 2 1 1 1 1 1 1 1 1 1 1 1]
	2013/10/17 15:35:54 57
	2013/10/17 15:35:54 4.280701754385965
	2013/10/17 15:35:54 Wall time for all  72.435945ms
	ziosi@spritz:~/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/scripts/exchanges$ go run exchanges.go > lifetimes.log
	2013/10/17 15:39:25 Header detected, skip...
	2013/10/17 15:39:25 Done reading  451  lines from file with err EOF
	2013/10/17 15:39:25 ========================================================
	2013/10/17 15:39:25 ========================================================
	2013/10/17 15:39:25 LIFE TIMES
	2013/10/17 15:39:25 ========================================================
	2013/10/17 15:39:25 ========================================================
	2013/10/17 15:39:25 [1 20 110 1 1 1 1 1 1 1 1 6 1 1 1 1 1 1 1 2 1 3 3 1 1 1 8 6 2 1 1 1 1 1 1 118 1 19 1 1 5 1 1 1 33 1 1 20 110 1 1 1 1 19 1 1 1 1 131 1 131 1 1 1 1 2 1 7 1 79 1 6 118 1 44]
	2013/10/17 15:39:25 75
	2013/10/17 15:39:25 14.04
	2013/10/17 15:39:25 Wall time for all  79.689705ms
	ziosi@spritz:~/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/scripts/exchanges$ 

	 
	 */
	
	
	/*
	 2013/10/17 16:52:18 LIFE TIMES
2013/10/17 16:52:18 ========================================================
2013/10/17 16:52:18 ========================================================
2013/10/17 16:52:23 65572
2013/10/17 16:52:23 lifeTimesAll 3.6411999328981883
2013/10/17 16:52:23 ========================================================
2013/10/17 16:52:28 821
2013/10/17 16:52:28 lifeTimesHardDBH 26.56573934226553
2013/10/17 16:52:28 ========================================================
2013/10/17 16:52:33 73
2013/10/17 16:52:33 lifeTimesHardDNS 79.3141095890411
2013/10/17 16:52:33 ========================================================
2013/10/17 16:52:37 44
2013/10/17 16:52:37 lifeTimesHardBHNS 76.67336363636365
2013/10/17 16:52:37 ========================================================
2013/10/17 16:52:42 7745
2013/10/17 16:52:42 lifeTimesSoftDBH 3.0651519690122657
2013/10/17 16:52:42 ========================================================
2013/10/17 16:52:47 0
2013/10/17 16:52:47 lifeTimesSoftDNS NaN
2013/10/17 16:52:47 ========================================================
2013/10/17 16:52:53 49
2013/10/17 16:52:53 lifeTimesSoftBHNS 3.4212244897959185
2013/10/17 16:52:53 ========================================================
2013/10/17 16:52:53 Wall time for all  1m48.620163634s

	 
	 */
	
	
	
	

	
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
} //END MAIN
