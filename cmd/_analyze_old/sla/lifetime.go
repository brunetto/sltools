package sla

import(
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

// LifeTimeMap contains stats about star's lifetimes
// in binary
// type LifeTimeMap struct {
// 	All []uint64
// 	HardDBH []uint64
// 	SoftDBH []uint64
// 	HardDNS []uint64
// 	SoftDNS []uint64
// 	HardBHNS []uint64http://www.wolfram.com/mathematica/how-to-buy/?a=1
// 	SoftBHNS []uint64
// }

type LifeTimeMap map[string][]uint64

var lifetimesFields = []string{"All", "HardDBH", "SoftDBH", "HardDNS", "SoftDNS", "HardBHNS", "SoftBHNS"}

func (lifetimes LifeTimeMap) Init() () {
	for _, field := range lifetimesFields {
		lifetimes[field] = make([]uint64, 0)
	}
}

// ComputeLifeTimes counts the lifetimes of the star between the exchanges
func (starData *StarData) ComputeLifeTimes() () {
	idx := 0
	time0 := uint64(0)
	time1 := uint64(0)
	binary := "0"
	delta := uint64(0)
	excTimes := starData.Exchanges.Keys()
	for _, key := range excTimes {
		idx++
		if (binary != starData.Exchanges[key].BinaryId) || (idx == len(excTimes)) {
			delta = time1-time0 + 1
// 			starData.LifeTimesStats = append(starData.LifeTimesStats, delta)
			time0 = key
			// this can happen considering promiscuous binaries with time=1000+time
			if delta > 400 {
				if Verb{
					log.Printf("Delta %v grater than 400 timesteps for star %v at time %v with type %v and hardness %v\n", 
						   delta, starData.StarId, time0, starData.Exchanges[key].Types, starData.Exchanges[key].Hardness)
					log.Println("Probably we reached the promiscuou binaries, going to the next star")
				}
				return
			}
			starData.LifeTimesStats = make(LifeTimeMap)
			starData.LifeTimesStats.Init()
			// Time spent in a hard binary
			// Hard and soft refers to the "stable" criterion in starlab
			// see SPZ paper for details
			if starData.Exchanges[key].Hardness == "H" {
				// Time spent in a hard bh|bh binary
				if starData.Exchanges[key].Types == "bh|bh" {
					starData.LifeTimesStats["HardDBH"] = append(starData.LifeTimesStats["HardDBH"], delta)
				}
				// Time spent in a hard ns|ns binary
				if starData.Exchanges[key].Types == "ns|ns" {
					starData.LifeTimesStats["HardDNS"] = append(starData.LifeTimesStats["HardDNS"], delta)
				}
				// Time spent in a hard bh|ns binary
				if (starData.Exchanges[key].Types == "bh|ns" || starData.Exchanges[key].Types == "ns|bh") {
					starData.LifeTimesStats["HardBHNS"] = append(starData.LifeTimesStats["HardBHNS"], delta)
				}
				// Time spent in a soft binary
			} else if starData.Exchanges[key].Hardness == "S" {
				// Time spent in a soft bh|bh binary
				if starData.Exchanges[key].Types == "bh|bh" {
					starData.LifeTimesStats["SoftDBH"] = append(starData.LifeTimesStats["SoftDBH"], delta)
				}
				// Time spent in a soft bh|bh binary
				if starData.Exchanges[key].Types == "ns|ns" {
					starData.LifeTimesStats["SoftDNS"] = append(starData.LifeTimesStats["SoftDNS"], delta)
				}
				// Time spent in a soft bh|ns binary
				if (starData.Exchanges[key].Types == "bh|ns" || starData.Exchanges[key].Types == "ns|bh") {
					starData.LifeTimesStats["SoftBHNS"] = append(starData.LifeTimesStats["SoftBHNS"], delta)
				}
			} 
			starData.LifeTimesStats["All"] = append(starData.LifeTimesStats["All"], delta)
		}
		time1 = key
		binary = starData.Exchanges[key].BinaryId
	}
}

type AllLTMap map[string]LifeTimeMap

func (starMap StarMapType) CollectLifeTimes() (AllLTMap) {
	allLT := make(AllLTMap)
	// Open each star
	keys := starMap.Keys()
	for _, key := range keys {
		starData := starMap[key]
		// Init containers
		if _, exists := allLT[starData.Z]; !exists {
			if Verb {log.Println("Init lifetimes container for Z = ", starData.Z)}
			allLT[starData.Z] = make(LifeTimeMap)
			allLT[starData.Z].Init()
		}
		for _, field := range lifetimesFields {
			allLT[starData.Z][field] = append(allLT[starData.Z][field], starData.LifeTimesStats[field]...)
		}
	}
	return allLT
}

// SaveLifeTimes save lifetimes to file, separated per hardness, type of the binary and metallicity
func (allLT AllLTMap) SaveLifeTimes(writer io.Writer) () {
	metallicities := allLT.Keys()
	for _, z := range metallicities {
		for _, field := range lifetimesFields {
			fmt.Fprintf(writer, "Z%v %v ", z, field)
			// Range over the items in "All", "HardDBH", "SoftDBH", "HardDNS", ...
			for _, item := range allLT[z][field] {
				fmt.Fprintf(writer, "%v ", item)
			}
			fmt.Fprintf(writer, "\n")
		}
	}
}

// SaveLifeTimes save lifetimes to file, separated per hardness, type of the binary and metallicity
func (allLT AllLTMap) PrintLTStats() (map[string]map[string]uint64) {
	const st2myr = 0.25
	var average, total uint64
	var writer = os.Stderr
	
	metallicities := allLT.Keys()
	averagesMap := make(map[string]map[string]uint64)
	
	fmt.Fprint(writer, "####################################################################################################################################\n")
	fmt.Fprint(writer, "# LIFETIMES STATISTICS\n")
	fmt.Fprint(writer, "####################################################################################################################################\n")
	
	fmt.Fprintln(writer, "\nAverage lifetime for:")
	
	
	for _, z := range metallicities {
		for _, field := range lifetimesFields {
			// Range over the items in "All", "HardDBH", "SoftDBH", "HardDNS", ...
			fmt.Fprintf(writer, "%-5v%-10v:%+6v timesteps ~ %+7v Myr\n", z, field, average, st2myr*float64(average))
			sum := uint64(0)
			for _, item := range allLT[z][field] {
				sum += item
			}
			total = uint64(len(allLT[z][field]))
			if total == 0 {
				average = 0
			} else {
				average = sum / uint64(len(allLT[z][field]))
			}
			
			if _, exists := averagesMap[z]; !exists {
				averagesMap[z] = make(map[string]uint64)
			}
			if _, exists := averagesMap[z][field]; !exists {
				averagesMap[z][field] = average
			}
			fmt.Fprintf(writer, "%-5v%-10v:%+6v timesteps ~ %+7v Myr\n", z, field, average, st2myr*float64(average))
		}
	}
	return averagesMap
}

// Keys returns the sorted keys
func (allLT AllLTMap) Keys() (keys []string) {
	keys = make([]string, len(allLT))
	idx := 0 
	for key, _ := range allLT {
        keys[idx] = key
        idx++
    }
    sort.Strings(keys)
	return keys
}



