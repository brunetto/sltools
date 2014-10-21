package slan

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/brunetto/goutils/debug"
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
type ExchangesMap map[uint64]*ExchData // FIXME: update all other references

type ExchData struct {
	BinaryId string
	Companion string
}

// ExchangeSummary summarizes the exchanges data/stats after counting
// them in CountExchanges()
type ExchangeStats struct {
	StarId              string
	Primordial          bool
	HardExchanges       []string
	HardExchangesNumber int
	SoftExchanges       []string
	SoftExchangesNumber int
	TotalExchanges      int
}

// uint64arr is a uint64 array type useful to sort it
type uint64arr []uint64

func (a uint64arr) Len() int           { return len(a) }
func (a uint64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a uint64arr) Less(i, j int) bool { return a[i] < a[j] }

// countExchanges counts a star's exchanges
func (data *AllDataType) CountExchanges(starId string) *ExchangeStats {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	if Verb {
		log.Println("Count exchanges for ", starId)
	}
	timesteps := data.Stars[starId].Exchanges.Keys()
	excData := new(ExchangeStats)
	excData.StarId = data.Stars[starId].StarId
	excData.Primordial = data.Stars[starId].Primordial
	hardCompanions := make([]string, 0)
	softCompanions := make([]string, 0)
	lastCompanion := "dummy"
	for _, timeStep := range timesteps {
		if data.Stars[starId].Exchanges[timeStep].Companion != lastCompanion {
			if data.Binaries[data.Stars[starId].Exchanges[timeStep].BinaryId].TimeProperties[timeStep].Hardness == "H" {
				lastCompanion = data.Stars[starId].Exchanges[timeStep].Companion
				hardCompanions = append(hardCompanions, lastCompanion)
			} else if data.Binaries[data.Stars[starId].Exchanges[timeStep].BinaryId].TimeProperties[timeStep].Hardness == "S" {
				lastCompanion = data.Stars[starId].Exchanges[timeStep].Companion
				softCompanions = append(softCompanions, lastCompanion)
			} else {
				log.Fatalf("Wrong hardness detected %v in %v ", data.Binaries[data.Stars[starId].Exchanges[timeStep].BinaryId].TimeProperties[timeStep], data.Stars[starId].StarId)
			}
		}
	}
	excData.HardExchanges = hardCompanions
	excData.SoftExchanges = softCompanions
	excData.HardExchangesNumber = len(excData.HardExchanges)
	excData.SoftExchangesNumber = len(excData.SoftExchanges)
	// Correct for the first entry inbinary
	if data.Stars[starId].Primordial {
		if data.Binaries[data.Stars[starId].Exchanges[timesteps[0]].BinaryId].TimeProperties[timesteps[0]].Hardness == "H" {
			excData.HardExchangesNumber--
		} else if data.Binaries[data.Stars[starId].Exchanges[timesteps[0]].BinaryId].TimeProperties[timesteps[0]].Hardness == "S" {
			excData.SoftExchangesNumber--
		}
	}
	excData.TotalExchanges = excData.HardExchangesNumber + excData.SoftExchangesNumber
	data.Stars[starId].ExchangeSummary = excData
	return excData
}

// SaveExch save all the stars' exchanges to file
func (starMap StarMapType) SaveExch(writer io.Writer) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var keys []string
	var value *StarData
	var key string

	fmt.Fprint(writer, "####################################################################################################################################\n")
	fmt.Fprintf(writer, "# Extracted %v objects\n", len(starMap))

	fmt.Fprint(writer, "####################################################################################################################################\n")
	// 	fmt.Print("# starId\t\t\tZ\tLastType\tPrimordial\tHard\tSoft\tTotal\tFirstInBin\tLastInBin\n")
	fmt.Fprintf(writer, "%-30v%+8v%+20v%+15v%+7v%+7v%+7v%+12v%+12v\n", "# starId", "Z", "LastType", "Primordial", "Hard", "Soft", "Total", "FirstInBin", "LastInBin")
	fmt.Fprint(writer, "####################################################################################################################################\n")

	keys = starMap.Keys()
	for idx := 0; idx < len(keys); idx++ {
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

func (starMap StarMapType) PrintExchStats(writer io.Writer) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		keys       []string
		value      *StarData
		key        string
		primordial int = 0
		hard       int = 0
		soft       int = 0
		total      int = 0
	)

	keys = starMap.Keys()
	for idx := 0; idx < len(keys); idx++ {
		key = keys[idx]
		value = starMap[key]
		if value.ExchangeSummary.Primordial {
			primordial++
		}
		hard = hard + value.ExchangeSummary.HardExchangesNumber
		soft = soft + value.ExchangeSummary.SoftExchangesNumber
		total = total + value.ExchangeSummary.TotalExchanges
	}

	nStars := float32(len(starMap))

	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "####################################################################################################################################\n")
	fmt.Fprintf(writer, "# EXCHANGES STATISTICS\n")
	fmt.Fprintf(writer, "####################################################################################################################################\n")

	fmt.Fprintf(writer, "# nStars: %v\n", nStars)
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "", "Primordial", "Hard", "Soft", "Total")
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "", "----------", "----", "----", "-----")
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "Average", "NA", float32(hard)/nStars, float32(soft)/nStars, float32(total)/nStars)
	fmt.Fprintf(writer, "# %+10v\t%+10v\t%+10v\t%+10v\t%+10v\n", "Total", primordial, hard, soft, total)
}

// func (exmap ExchangesMap) Print() {
// 	if Debug {
// 		defer debug.TimeMe(time.Now())
// 	}
// 	// Retrieve and sort Exchanges map keys
// 	excTimes := exmap.Keys()
// 	fmt.Println("SysTime, BinaryId, Companion, Hardness, Types, Ecc")
// 	for _, key := range excTimes {
// 		fmt.Print(key, " ")
// 		exmap[uint64(key)].Print()
// 	}
// }

// Print prints a summary of the exchanges stats for the star
func (exchangeSummary *ExchangeStats) Print() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
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
	if Debug {
		defer debug.TimeMe(time.Now())
	}

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
