package slan

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/sets"
)

// StarMapType is the struct containing all the star data
// map[starID]*starData, es: starID = Z010n003id102460
type StarMapType map[string]*StarData

// FIXME: check the order of function calling
// StarData store the star data we are interest in.
type StarData struct {
	// Combination number to identify the parameters set
	Comb string	
	// StarId uniquely identifies the star among all the simulations.
	StarId string
	// Z is the metallicity in units of Zsun = 0.019.
	Z string
	// NFile is the file (run) number.
	NFile string
	// DCOB last type flag (DBH, DNS, BHNS).
	LastDCOB string
	//
	DCOB sets.StringSet
	// ZeroEcc flags if the star is found in a
	// zero-eccentricity binary at any time.
	ZeroEcc bool
	// Promiscuous flags if the star reside in two binaries at the same time
	// at any time.
	Promiscuous bool
	// Primordial flags if the star is in a binary at t=0
	Primordial bool
	// Exchanges tracks all the star's exchanges:
	// key is the uint64 sys_time: I've tought it was unique
	// value is of type *BinaryData
	Exchanges ExchangesMap
	// ExchangeSummary summarizes the exchanges data/stats after counting
	// them in CountExchanges()
	ExchangeSummary *ExchangeStats
	// LifeTimesStats summarizes the lifetimes data/stats after counting
	// them in ComputeLifeTimes
	LifeTimesStats LifeTimeMap
	// TimeDomain store the first and last time a star is found in binaryData
	TimeDom TimeDomain
	// Time unit of the simulation
	TimeUnit float64
}

// Keys returns the sorted keys of the maps with all the stars
func (starMap StarMapType) Keys() (keys []string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	keys = make([]string, len(starMap))
	idx := 0
	for key, _ := range starMap {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)
	return keys
}

// ExtrcType return a map with only stars of a given last type.
// [bh|bh, ns|ns, bh|ns]
func (starMap StarMapType) ExtrcLastType(mapType string) StarMapType {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	if mapType != "bh|bh" && mapType != "ns|ns" && mapType != "bh|ns" {
		log.Fatal("Wrong type in ", debug.FName(false), " function, ", mapType, " not in [bh|bh, ns|ns, bh|ns]")
	}
	starMapPerType := make(StarMapType)
	for key, value := range starMap {
		if value.LastDCOB == mapType {
			starMapPerType[key] = value
		} else if mapType == "bh|ns" {
			if value.LastDCOB == "ns|bh" {
				starMapPerType[key] = value
			}
		}
	}
	return starMapPerType
}

// WasInType return a map with only stars that were in a given binary type.
// [bh|bh, ns|ns, bh|ns]
func (starMap StarMapType) WasInType(mapType string) StarMapType {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	if mapType != "bh|bh" && mapType != "ns|ns" && mapType != "bh|ns" {
		log.Fatal("Wrong type in ", debug.FName(false), " function, ", mapType, " not in [bh|bh, ns|ns, bh|ns]")
	}
	starMapPerType := make(StarMapType)
	for key, value := range starMap {
		if value.DCOB.Exists(mapType) {
			starMapPerType[key] = value
		} else if mapType == "bh|ns" {
			if value.DCOB.Exists("ns|bh") {
				starMapPerType[key] = value
			}
		}
	}
	return starMapPerType
}

// WasInType return a map with only stars that were in a given binary type.
// [bh|bh, ns|ns, bh|ns]
func (starMap StarMapType) WasPromiscuous() StarMapType {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	starMapPerType := make(StarMapType)
	for key, value := range starMap {
		if value.Promiscuous {
			starMapPerType[key] = value
		}
	}
	return starMapPerType
}

// Print prints StarData data
func (starData *StarData) Print() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	fmt.Println("#########################")
	fmt.Println("#       Star Data       #")
	fmt.Println("#########################")
	fmt.Println("StarId = ", starData.StarId)
	fmt.Println("Z = ", starData.Z)
	fmt.Println("NFile = ", starData.NFile)
	fmt.Println("TimeUnit = ", starData.TimeUnit)
	fmt.Println("Last DCOB = ", starData.LastDCOB)
	fmt.Println("DCOB = ", starData.DCOB.String())
	fmt.Println("Primordial = ", starData.Primordial)
	fmt.Println("First time in binary = ", starData.TimeDom.Min)
	fmt.Println("Last time in binary = ", starData.TimeDom.Max)
}

func (starMap StarMapType) Summary() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	fmt.Println("Number of stars: ", len(starMap))
}

// Print prints StarData data
func (starMap StarMapType) Print() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	starMap.Summary()
	for _, star := range starMap.Keys() {
		starMap[star].Print()
	}
	fmt.Println("\n=== Summary ===\n")
	starMap.Summary()
}

// // PrintWithExchs prints StarData data with exchanges
// func (starData *StarData) PrintWithExchs() {
// 	if Debug {
// 		defer debug.TimeMe(time.Now())
// 	}
// 	starData.Print()
// 	fmt.Println("#########################")
// 	fmt.Println("#       Exchanges       #")
// 	fmt.Println("#########################")
// 	starData.Exchanges.Print()
// }
