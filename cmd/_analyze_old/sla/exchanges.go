package sla

import (
	"fmt"
	"io"
	"log"
// 	"os"
	"sort"
)

// ExchangesMap is a map containing all the exchanges of a star.
// It is a list of all the "binary data" for the binaries a star
// was a member of.
// key is the uint64 sys_time: I've tought it was unique
// value is of type *BinaryData
// NOTE: use a pointer otherwise structs will be unchangeble:
// a map returns a copy of the element, in this case a copy of 
// the pointer to access the data.
type ExchangesMap map[uint64]*BinaryData

// ExchangeSummary summarizes the exchanges data/stats after counting 
// them in CountExchanges()
type ExchangeStats struct {
	StarId string
	Primordial bool
	HardExchanges []string
	HardExchangesNumber int
	SoftExchanges []string
	SoftExchangesNumber int
	TotalExchanges int
}

// uint64arr is a uint64 array type useful to sort it
type uint64arr []uint64
func (a uint64arr) Len() int { return len(a) }
func (a uint64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a uint64arr) Less(i, j int) bool { return a[i] < a[j] }

// countExchanges counts a star's exchanges
func (starData *StarData) CountExchanges() (*ExchangeStats) {
	if Verb {
		log.Println("Count exchanges for ", starData.StarId)
	}
	keys := starData.Exchanges.Keys()
	excData := new(ExchangeStats)
	excData.StarId = starData.StarId
	excData.Primordial = starData.Primordial
	hardCompanions := make([]string, 0)
	softCompanions := make([]string, 0)
	lastCompanion := "dummy"
	for _, key := range keys {
		if starData.Exchanges[key].Companion != lastCompanion {
			if starData.Exchanges[key].Hardness == "H" {
				lastCompanion = starData.Exchanges[key].Companion
				hardCompanions = append(hardCompanions, lastCompanion)
			} else if starData.Exchanges[key].Hardness == "S" {
				lastCompanion = starData.Exchanges[key].Companion
				softCompanions = append(softCompanions, lastCompanion)
			} else {
					log.Fatalf("Wrong hardness detected %v in %v ", starData.Exchanges[key].Hardness, starData.StarId)
			}
		}
	}
	excData.HardExchanges = hardCompanions
	excData.SoftExchanges = softCompanions
	excData.HardExchangesNumber = len(excData.HardExchanges)
	excData.SoftExchangesNumber = len(excData.SoftExchanges)
	// Correct for the first entry inbinary
	if starData.Primordial {
		if starData.Exchanges[keys[0]].Hardness == "H" {
			excData.HardExchangesNumber-- 
		} else if starData.Exchanges[keys[0]].Hardness == "S" {
			excData.SoftExchangesNumber-- 
		}
	}
	excData.TotalExchanges = excData.HardExchangesNumber + excData.SoftExchangesNumber
	starData.ExchangeSummary = excData
	return excData
}

// SaveExch save all the stars' exchanges to file
func (starMap StarMapType) SaveExch(writer io.Writer) {
	var keys []string
	var value *StarData
	var key string
	
	fmt.Fprint(writer, "####################################################################################################################################\n")
	fmt.Fprintf(writer, "# Extracted %v objects\n", len(starMap))
	
	fmt.Fprint(writer, "####################################################################################################################################\n")
// 	fmt.Print("# starId\t\t\tZ\tLastType\tPrimordial\tHard\tSoft\tTotal\tFirstInBin\tLastInBin\n")
	fmt.Fprintf(writer, "%-30v%+8v%+20v%+15v%+7v%+7v%+7v%+12v%+12v\n","# starId", "Z", "LastType", "Primordial", "Hard", "Soft", "Total", "FirstInBin", "LastInBin")
	fmt.Fprint(writer, "####################################################################################################################################\n")
	
	keys = starMap.Keys()
	for idx:=0; idx<len(keys); idx++ {
		key = keys[idx]
		value = starMap[key]
		fmt.Fprintf(writer, "%-30v%+8v%+20v%+15v%+7v%+7v%+7v%+12v%+12v\n", key, value.Z, value.LastDCOB, 
					value.ExchangeSummary.Primordial, 
					value.ExchangeSummary.HardExchangesNumber, 
					value.ExchangeSummary.SoftExchangesNumber, 
					value.ExchangeSummary.TotalExchanges, 
					value.TimeDom.Min,
					value.TimeDom.Max)
	}
}


func (starMap StarMapType) PrintExchStats(writer io.Writer) () {
	var (
		keys []string
		value *StarData
		key string
		primordial int = 0
		hard int = 0
		soft int = 0
		total int = 0 
	)
	
	keys = starMap.Keys()
	for idx:=0; idx<len(keys); idx++ {
		key = keys[idx]
		value = starMap[key]
		if value.ExchangeSummary.Primordial {primordial++}
		hard  = hard  + value.ExchangeSummary.HardExchangesNumber
		soft  = soft  + value.ExchangeSummary.SoftExchangesNumber
		total = total + value.ExchangeSummary.TotalExchanges
	}
	
	nStars := float32(len(starMap))
	
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "####################################################################################################################################\n")
	fmt.Fprintf(writer, "# EXCHANGES STATISTICS\n")
	fmt.Fprintf(writer, "####################################################################################################################################\n")
	
	fmt.Fprintf(writer, "# nStars: %v\n", nStars)
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "",  "Primordial", "Hard", "Soft", "Total")
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "",  "----------", "----", "----", "-----")
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "Average",  "NA", float32(hard)/nStars, float32(soft)/nStars, float32(total)/nStars)
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "Total", primordial, hard, soft, total)
}

func (exmap ExchangesMap) Print () () {
	// Retrieve and sort Exchanges map keys
	excTimes := exmap.Keys()
	fmt.Println("SysTime, BinaryId, Companion, Hardness, Types, Ecc")
	for _, key := range excTimes {
		fmt.Print(key, " ")
		exmap[uint64(key)].Print()
	}
}


// Print prints a summary of the exchanges stats for the star
func (exchangeSummary *ExchangeStats) Print() {
	fmt.Println("#########################")
	fmt.Println("#   Exchanges Summary   #")
	fmt.Println("#########################")
	fmt.Println("Data for ", exchangeSummary.StarId)
	fmt.Println("Primordial binary ", exchangeSummary.Primordial)
	fmt.Println("HardExchanges = ", exchangeSummary.HardExchanges)
	fmt.Println("SoftExchanges = ", exchangeSummary.SoftExchanges)
	fmt.Println("HardExchangesNumber = ", exchangeSummary.HardExchangesNumber)
	fmt.Println("SoftExchangesNumber = ", exchangeSummary.SoftExchangesNumber)
	fmt.Println("TotalExchanges = ", exchangeSummary.TotalExchanges)
}

// Retrieve and sort Exchanges map keys
func (exc ExchangesMap) Keys() (keys []uint64) {
	
	excTimes := make([]uint64, len(exc))
	idx := 0 
	for key, _ := range exc {
        excTimes[idx] = key
        idx++
    }
    sort.Sort(uint64arr(excTimes))
	keys = uint64arr(excTimes)
	return keys
}




