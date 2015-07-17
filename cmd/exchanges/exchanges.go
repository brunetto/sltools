package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
)

const Debug bool = false

func main() {
	defer debug.TimeMe(time.Now())

	const hidePromiscuous bool = true

	var (
		err       error
		inPath    string
		inFile    string
		outFile   *os.File
		starMap   = StarMapType{}
		regGroups = RegGroupsStruct{Id: 1,
			Comb:     2,
			N:        3,
			Z:        4,
			Rv:       5,
			Fpb:      6,
			W:        7,
			Tf:       8,
			Systime:  9,
			Phystime: 10,
			Objs:     11,
			Hflag:    12,
			Types:    13,
			M0:       14,
			M1:       15,
			Sma:      16,
			Period:   17,
			Ecc:      18,
			Tgw:      19,
			Mchirp:   20,
		}

		// Regexp string for all_the_fishes
		regStringCSVsims string = `^(\S+?),` + // "id": 1,
			`(\d+?),` + // "comb": 2,
			`(\d+?),` + // "n": 3,
			`(\d+\.*\d*),` + // "z": 4,
			`(\d+\.*\d*),` + // "rv": 5,
			`(\d+\.*\d*),` + // "fpb": 6,
			`(\d+?),` + // "w": 7,
			`(\S+?),` + // "tf": 8,
			`(\d+?),` + // "systime": 9,
			`(\d+\.*\d*),` + // "phystime": 10,
			`(\S+\|\S+),` + // "objs": 11,
			`(\S),` + // "hflag": 12,
			`(\S+\|\S+?),` + // "types": 13,
			`(\d+\.*\d*e*[-\+]*\d*),` + // "m0": 14,
			`(\d+\.*\d*e*[-\+]*\d*?),` + // "m1": 15,
			`(\d+\.*\d*e*[-\+]*\d*?),` + // "sma": 16,
			`(\d+\.*\d*e*[-\+]*\d*?),` + // "period": 17,
			`(\d+\.*\d*e*[-\+]*\d*?),` + // "ecc": 18,
			`(\d+\.*\d*e*[-\+]*\d*?),` + // "tgw": 19,
			`(\d+\.*\d*e*[-\+]*\d*?)` // "mchirp": 20,

	// Single reg string for testing and checking purpose
	// ^(\S+?),(\d+?),(\d+?),(\d+\.*\d*),(\d+\.*\d*),(\d+\.*\d*),(\d+?),(\S+?),(\d+?),(\d+\.*\d*),(\S+\|\S+),(\S),(\S+\|\S+?),(\d+\.*\d*e*[-\+]*\d*),(\d+\.*\d*e*[-\+]*\d*?),(\d+\.*\d*e*[-\+]*\d*?),(\d+\.*\d*e*[-\+]*\d*?),(\d+\.*\d*e*[-\+]*\d*?),(\d+\.*\d*e*[-\+]*\d*?),(\d+\.*\d*e*[-\+]*\d*?)
	)
	
	if len(os.Args) < 2{
		log.Fatal("Please provide a data CSV file.")
	}
	
	inPath = ""
	inFile = os.Args[1]

	// Read data about binaries and populate data structure
	starMap.Populate(inPath, inFile, regStringCSVsims, regGroups, hidePromiscuous)

	// Compute exchanges
	starMap.CountExchanges()

	// Compute lifetimes
	starMap.ComputeLifeTimes()

	// 	starMap.Print(os.Stdout)

	if outFile, err = os.Create("exchanges-" + filepath.Base(inFile) + ".log"); err != nil {
		log.Fatal("Can't create file exchanges.log with err: ", err)
	}
	defer outFile.Close()

	long := true
	starMap.Print(outFile, long)
	starMap.CsvSummary("Exchanges-" + filepath.Base(os.Args[1]))

	// 	TODO
	// 	We keep track of every binary containing a star that will become a compact object (ns or bh, until now, wd required?).
	// 	From another point of view we keep track of every star that will be a CO when it is in a binary.
	//
	// 	To print this summary with all the properties and exchanges, in a file, and
	// 	a csv summary with all the properties but without every single exchange.
	//
	// 	Remember to print lifetimes.
	//
	// FIXME: if one binary exists only for a single time, it is not counted in the lifetimes

}

/*
 *******************************************************************
 *******************************************************************
 **               Type, variables and functions                   **
 *******************************************************************
 *******************************************************************
 */

//Store the min and max time of the binaries in the file
type Time struct {
	Min float64
	Max float64
}

// Data of the star in binary and that changes with time (companion, orbital properties, ...)
type BinaryData struct {
	Time      float64
	BinaryId  string
	Companion string
	Hardness  string
	Types     []string
	M0        float64
	M1        float64
	Sma       float64
	Period    float64
	Tgw       float64
	Mchirp    float64
	Ecc       float64
}

// Print function for an exchange datum
func (binaryData *BinaryData) Print(outDest io.Writer) {
	fmt.Fprintln(outDest, "\t",
		binaryData.BinaryId,
		binaryData.Companion,
		binaryData.Hardness,
		binaryData.Types,
		binaryData.Ecc,
		binaryData.Types,
		binaryData.M0,
		binaryData.M1,
		binaryData.Sma,
		binaryData.Period,
		binaryData.Tgw,
		binaryData.Mchirp,
		binaryData.Ecc)
}

// NOTE: use a pointer otherwise structs will be unchangeble!!!!!!!
type ExchangesMap map[float64]*BinaryData

// Store the unchangeble star data we are interestbinaryData in
type StarData struct {
	StarId string
	Z      string
	NFile  string
	// DCOB flag
	DCOB string
	// Zero eccentricity flag for this star at some time
	ZeroEcc bool
	// Promiscuity flag for this star at some time
	Promiscuous bool
	// Start primordial flag
	Primordial bool
	// The key is the sys_time: I've tought it was unique -> now it is the phystime
	Exchanges ExchangesMap
	// After counting exchanges in Exchanges
	ExchangesStats *ExchangesStats
	// Lifetimes
	LifeTimes      LifeTimeStruct
	LifeTimesStats LifeTimeSummary
	Comb           string
	Rv             float64
	Fpb            float64
	W              int64
	Tf             string
}

func (starData *StarData) Keys() []float64 {
	// Retrieve and sort Exchanges map keys
	excTimes := make([]float64, len(starData.Exchanges))
	idx := 0
	for key, _ := range starData.Exchanges {
		excTimes[idx] = key
		idx++
	}
	sort.Float64s(excTimes)
	return excTimes
}

func (exchanges ExchangesMap) Keys() []float64 {
	// Retrieve and sort Exchanges map keys
	excTimes := make([]float64, len(exchanges))
	idx := 0
	for key, _ := range exchanges {
		excTimes[idx] = key
		idx++
	}
	sort.Float64s(excTimes)
	return excTimes
}

// Print function for StarData
func (starData *StarData) Print(outDest io.Writer, long bool) {
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "#       Star Data       #")
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "StarId = ", starData.StarId)
	fmt.Fprintln(outDest, "NFile = ", starData.NFile)
	fmt.Fprintln(outDest, "DCOB = ", starData.DCOB)
	fmt.Fprintln(outDest, "Primordial = ", starData.Primordial)
	fmt.Fprintln(outDest, "Z = ", starData.Z)
	fmt.Fprintln(outDest, "Comb", starData.Comb)
	fmt.Fprintln(outDest, "Rv", starData.Rv)
	fmt.Fprintln(outDest, "Fpb", starData.Fpb)
	fmt.Fprintln(outDest, "W", starData.W)
	fmt.Fprintln(outDest, "Tf", starData.Tf)
	starData.ExchangesStats.Print(outDest)
	starData.LifeTimesStats.Print(outDest)
	if long {
		starData.Exchanges.Print(outDest)
	}
	if long {
		starData.LifeTimes.Print(outDest)
	}
}

// Print function for StarData
func (starData *StarData) CsvSummary(format string) (outString string) {
	outString = fmt.Sprintf(format,
		starData.StarId,
		starData.Comb,
		starData.NFile,
		starData.DCOB,
		starData.Primordial,
		starData.Z,
		starData.Rv,
		starData.Fpb,
		starData.W,
		starData.Tf,
		starData.ExchangesStats.HardExchangesNumber,
		starData.ExchangesStats.SoftExchangesNumber,
		starData.ExchangesStats.TotalExchangesNumber,
		starData.LifeTimesStats.AllSum,
		starData.LifeTimesStats.AllMean,
		starData.LifeTimesStats.HardDBHSum,
		starData.LifeTimesStats.HardDBHMean,
		starData.LifeTimesStats.SoftDBHSum,
		starData.LifeTimesStats.SoftDBHMean,
		starData.LifeTimesStats.HardDNSSum,
		starData.LifeTimesStats.HardDNSMean,
		starData.LifeTimesStats.SoftDNSSum,
		starData.LifeTimesStats.SoftDNSMean,
		starData.LifeTimesStats.HardBHNSSum,
		starData.LifeTimesStats.HardBHNSMean,
		starData.LifeTimesStats.SoftBHNSSum,
		starData.LifeTimesStats.SoftBHNSMean,
	)
	return outString
}

type LifeTimeStruct struct {
	All      f64Slice
	HardDBH  f64Slice
	SoftDBH  f64Slice
	HardDNS  f64Slice
	SoftDNS  f64Slice
	HardBHNS f64Slice
	SoftBHNS f64Slice
}

type f64Slice []float64

func (s *f64Slice) Sum() (sum float64) {
	sum = 0
	for _, value := range *s {
		sum += value
	}
	return sum
}

func (s *f64Slice) Mean() (mean float64) {
	mean = 0
	sum := s.Sum()
	if sum != 0 {
		mean = s.Sum() / float64(len(*s))
	}
	return mean
}

type LifeTimeSummary struct {
	AllSum       float64
	AllMean      float64
	HardDBHSum   float64
	HardDBHMean  float64
	SoftDBHSum   float64
	SoftDBHMean  float64
	HardDNSSum   float64
	HardDNSMean  float64
	SoftDNSSum   float64
	SoftDNSMean  float64
	HardBHNSSum  float64
	HardBHNSMean float64
	SoftBHNSSum  float64
	SoftBHNSMean float64
}

func (lt *LifeTimeSummary) Print(outDest io.Writer) {
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "#   Lifetime Summary    #")
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "AllSum          ", lt.AllSum)
	fmt.Fprintln(outDest, "AllMean         ", lt.AllMean)
	fmt.Fprintln(outDest, "HardDBHSum      ", lt.HardDBHSum)
	fmt.Fprintln(outDest, "HardDBHMean     ", lt.HardDBHMean)
	fmt.Fprintln(outDest, "SoftDBHSum      ", lt.SoftDBHSum)
	fmt.Fprintln(outDest, "SoftDBHMean     ", lt.SoftDBHMean)
	fmt.Fprintln(outDest, "HardDNSSum      ", lt.HardDNSSum)
	fmt.Fprintln(outDest, "HardDNSMean     ", lt.HardDNSMean)
	fmt.Fprintln(outDest, "SoftDNSSum      ", lt.SoftDNSSum)
	fmt.Fprintln(outDest, "SoftDNSMean     ", lt.SoftDNSMean)
	fmt.Fprintln(outDest, "HardBHNSSum     ", lt.HardBHNSSum)
	fmt.Fprintln(outDest, "HardBHNSMean    ", lt.HardBHNSMean)
	fmt.Fprintln(outDest, "SoftBHNSSum     ", lt.SoftBHNSSum)
	fmt.Fprintln(outDest, "SoftBHNSMean    ", lt.SoftBHNSMean)
}

func (lts *LifeTimeStruct) Print(outDest io.Writer) {
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "#      Lifetimes        #")
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "All          ", lts.All)
	fmt.Fprintln(outDest, "HardDBH      ", lts.HardDBH)
	fmt.Fprintln(outDest, "SoftDBH      ", lts.SoftDBH)
	fmt.Fprintln(outDest, "HardDNS      ", lts.HardDNS)
	fmt.Fprintln(outDest, "SoftDNS      ", lts.SoftDNS)
	fmt.Fprintln(outDest, "HardBHNS     ", lts.HardBHNS)
	fmt.Fprintln(outDest, "SoftBHNS     ", lts.SoftBHNS)
}

func (starData *StarData) ComputeLifeTimes() {
	if Debug {
		log.Println("Counting lifetimes ...")
	}

	var (
		// Variables to store previous companion properties
		type0, type1 string
		hardness     string
		// 		binaryId string
		// Generic variables
		time0    float64 = 0
		delta    float64 = 0
		binary   string  = "0"
		excTimes []float64
	)

	excTimes = starData.Keys()

	for idx, key := range excTimes {

		if Debug {
			log.Println("Time, prev. binary, new binary: ", key, binary, starData.Exchanges[key].BinaryId)
		}

		if (binary == starData.Exchanges[key].BinaryId) && (idx != len(excTimes)-1) {
			if Debug {
				log.Println("Nothing changed, continue")
			}
			continue
		}

		// If star has a new companion, store data of the previous companion
		delta = key - time0 //+ 1
		if Debug {
			log.Println("New binary or end of star history, delta: ", delta)
		}

		// If starting a new star, no need to save the initial 0
		// (binary id is different from zero so it would store a delta=0 at the first loop
		if idx != 0 {
			// This can happen considering promiscuous binaries with time=1e10+time -> now it shouldn't
			if delta > 1e9 {
				log.Println("Huge delta: ", delta)
				long := true
				starData.Print(os.Stderr, long)
				log.Fatal("Maybe promiscuous binary or some problem.")
			}
			if Debug {
				log.Println("Store prev. binary: ", hardness, type0, type1, excTimes[idx-1], delta)
			}
			starData.storeLifetimes(hardness, type0, type1, excTimes[idx-1], delta)
		}
		if Debug {
			log.Printf("Update new binary data with %v at time %v\n", binary, key)
		}
		// Update variables to the new variable
		time0 = key
		hardness = starData.Exchanges[key].Hardness
		type0 = starData.Exchanges[key].Types[0]
		type1 = starData.Exchanges[key].Types[1]
		binary = starData.Exchanges[key].BinaryId
	}

	// Summary
	starData.LifeTimesStats.AllSum = starData.LifeTimes.All.Sum()
	starData.LifeTimesStats.AllMean = starData.LifeTimes.All.Mean()
	starData.LifeTimesStats.HardDBHSum = starData.LifeTimes.HardDBH.Sum()
	starData.LifeTimesStats.HardDBHMean = starData.LifeTimes.HardDBH.Mean()
	starData.LifeTimesStats.SoftDBHSum = starData.LifeTimes.SoftDBH.Sum()
	starData.LifeTimesStats.SoftDBHMean = starData.LifeTimes.SoftDBH.Mean()
	starData.LifeTimesStats.HardDNSSum = starData.LifeTimes.HardDNS.Sum()
	starData.LifeTimesStats.HardDNSMean = starData.LifeTimes.HardDNS.Mean()
	starData.LifeTimesStats.SoftDNSSum = starData.LifeTimes.SoftDNS.Sum()
	starData.LifeTimesStats.SoftDNSMean = starData.LifeTimes.SoftDNS.Mean()
	starData.LifeTimesStats.HardBHNSSum = starData.LifeTimes.HardBHNS.Sum()
	starData.LifeTimesStats.HardBHNSMean = starData.LifeTimes.HardBHNS.Mean()
	starData.LifeTimesStats.SoftBHNSSum = starData.LifeTimes.SoftBHNS.Sum()
	starData.LifeTimesStats.SoftBHNSMean = starData.LifeTimes.SoftBHNS.Mean()
}

func (starData *StarData) storeLifetimes(hardness, type0, type1 string, key, delta float64) {
	// if starData.Exchanges[key].Hardness == "H" {
	if hardness == "H" {
		if type0 == "bh" && type1 == "bh" {
			starData.LifeTimes.HardDBH = append(starData.LifeTimes.HardDBH, delta)
		}
		if type0 == "ns" && type1 == "ns" {
			starData.LifeTimes.HardDNS = append(starData.LifeTimes.HardDNS, delta)
		}
		if type0+type1 == "bhns" || type0+type1 == "nsbh" {
			starData.LifeTimes.HardBHNS = append(starData.LifeTimes.HardBHNS, delta)
		}
	} else if starData.Exchanges[key].Hardness == "S" {
		if type0 == "bh" && type1 == "bh" {
			starData.LifeTimes.SoftDBH = append(starData.LifeTimes.SoftDBH, delta)
		}
		if type0 == "ns" && type1 == "ns" {
			starData.LifeTimes.SoftDNS = append(starData.LifeTimes.SoftDNS, delta)
		}
		if type0+type1 == "bhns" || type0+type1 == "nsbh" {
			starData.LifeTimes.SoftBHNS = append(starData.LifeTimes.SoftBHNS, delta)
		}
	}
	starData.LifeTimes.All = append(starData.LifeTimes.All, delta)
}

func (starData *StarData) CountExchanges() *ExchangesStats {
	keys := starData.Keys()
	excStats := &ExchangesStats{}
	excStats.StarId = starData.StarId
	excStats.Primordial = starData.Primordial
	hardCompanions := []string{}
	softCompanions := []string{}
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
	excStats.HardCompanions = hardCompanions
	excStats.SoftCompanions = softCompanions
	excStats.HardExchangesNumber = len(excStats.HardCompanions)
	excStats.SoftExchangesNumber = len(excStats.SoftCompanions)
	// Correct for the first entry in binary
	if starData.Primordial {
		if starData.Exchanges[keys[0]].Hardness == "H" {
			excStats.HardExchangesNumber--
		} else if starData.Exchanges[keys[0]].Hardness == "S" {
			excStats.SoftExchangesNumber--
		}
	}
	excStats.TotalExchangesNumber = excStats.HardExchangesNumber + excStats.SoftExchangesNumber
	starData.ExchangesStats = excStats
	return excStats
}

type ExchangesStats struct {
	StarId               string
	Primordial           bool
	HardCompanions       []string
	HardExchangesNumber  int
	SoftCompanions       []string
	SoftExchangesNumber  int
	TotalExchangesNumber int
}

func (exchangesStats *ExchangesStats) Print(outDest io.Writer) {
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "#   Exchanges Summary   #")
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "Data for ", exchangesStats.StarId)
	fmt.Fprintln(outDest, "Primordial binary ", exchangesStats.Primordial)
	fmt.Fprintln(outDest, "HardCompanions = ", exchangesStats.HardCompanions)
	fmt.Fprintln(outDest, "SoftCompanions = ", exchangesStats.SoftCompanions)
	fmt.Fprintln(outDest, "HardExchangesNumber = ", exchangesStats.HardExchangesNumber)
	fmt.Fprintln(outDest, "SoftExchangesNumber = ", exchangesStats.SoftExchangesNumber)
	fmt.Fprintln(outDest, "TotalExchangesNumber = ", exchangesStats.TotalExchangesNumber)
}

func (exchanges ExchangesMap) Print(outDest io.Writer) {
	fmt.Fprintln(outDest, "#########################")
	fmt.Fprintln(outDest, "#       Exchanges       #")
	fmt.Fprintln(outDest, "#########################")
	// Retrieve and sort Exchanges map keys
	excTimes := exchanges.Keys()
	fmt.Fprintln(outDest, "PhysTime, BinaryId, Companion, Hardness, Types, Ecc, Types, M0, M1, Sma, Period, Tgw, Mchirp, Ecc")
	for _, key := range excTimes {
		fmt.Fprint(outDest, key, " ")
		exchanges[key].Print(outDest)
	}
}

type StarMapType map[string]*StarData

func (starMap StarMapType) Populate(inPath string, inFile string, regString string, regGroups RegGroupsStruct, hidePromiscuous bool) {
	defer debug.TimeMe(time.Now())

	var (
		fileObj      *os.File
		nReader      *bufio.Reader
		readLine     string
		err          error
		starId       = make([]string, 2)
		time         = Time{Min: 0, Max: 0}
		currentTime  float64
		currentIds   []string
		binaryRegexp = regexp.MustCompile(regString)
		regexResult  []string
		types        []string
		ecc          float64
	)

	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(filepath.Join(inPath, inFile)); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
	}
	defer fileObj.Close()

	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)

	// Read the file and fill the starMap slice
	log.Println("Start reading file and populating starMap...")
	line := 0
	for {
		if readLine, err = nReader.ReadString('\n'); err != nil {
			log.Println("Done reading ", line, " lines from file ", inFile, "  with err", err)
			break
		}
		if readLine[0] == '#' {
			log.Println("Header detected, skip...")
			continue
		}
		regexResult = binaryRegexp.FindStringSubmatch(readLine)
		if regexResult == nil {
			log.Println("With regexp ", binaryRegexp)
			log.Println("no match, nil regex result on line ", line)
			log.Println("Reg is:", regString)
			log.Println("Line is:", readLine)
		}

		// Update time domain if necesary
		if currentTime, err = strconv.ParseFloat(regexResult[regGroups.Phystime], 64); err != nil {
			log.Fatal("Error parsing current physical time: ", err)
		}

		// We don't want to follow single simulations reaching Gyrs when
		// the others stop at ~100
		if currentTime > 100 {
			continue
		}

		if time.Min > currentTime {
			time.Min = currentTime
		} else if time.Max < currentTime {
			time.Max = currentTime
		}

		// If not yet present, create two entries in the map,
		// one for each of the components of the binary
		// else only update the new timestep data&companion

		for i := 0; i < 2; i++ {
			// Retrieve single object id
			currentIds = strings.Split(regexResult[regGroups.Objs], "|")

			// Compose single object id
			starId[i] = "c" + regexResult[regGroups.Comb] +
				"n" + regexResult[regGroups.N] +
				"id" + currentIds[i]

			// Check existence of the single object in the map
			// if it exists then update the timestep/companion table,
			// else if not then create the entry
			if _, exists := starMap[starId[i]]; !exists {
				// Create the new StarData value for the key starMap in the map
				starMap[starId[i]] = NewStar(starId[i], regexResult, regGroups)

				// FIXME: Check for primordial (~, in reality it is check for t=0)
				// binary
				if currentTime == 0 {
					starMap[starId[i]].Primordial = true
				} else {
					starMap[starId[i]].Primordial = false
				}
				starMap[starId[i]].Exchanges = make(ExchangesMap)
			}
			// Fill timestep companion and data
			// If the entry (=the timestep) already exists is a problem
			// it means that I have found the same star in two binaries
			// and it is not a hyerarchical system nor a triple etc
			// because this have to be catchbinaryData before (TO BE CHECKED)
			if _, exists := starMap[starId[i]].Exchanges[currentTime]; !exists {
				starMap[starId[i]].Exchanges[currentTime] = NewBinary(regexResult, regGroups)
				starMap[starId[i]].Exchanges[currentTime].Companion = currentIds[len(currentIds)-i-1]
				starMap[starId[i]].Exchanges[currentTime].Time = currentTime

				types = starMap[starId[i]].Exchanges[currentTime].Types
				ecc = starMap[starId[i]].Exchanges[currentTime].Ecc

			} else if hidePromiscuous == true {
				continue // doesn't save promiscuous binaries
			} else { //save adding 1e10 to the time
				starMap[starId[i]].Promiscuous = true
				starMap[starId[i]].Exchanges[currentTime+1e10] = NewBinary(regexResult, regGroups)
				starMap[starId[i]].Exchanges[currentTime+1e10].Companion = currentIds[len(currentIds)-i-1]
				starMap[starId[i]].Exchanges[currentTime+1e10].Time = currentTime

				types = starMap[starId[i]].Exchanges[currentTime+1e10].Types
				ecc = starMap[starId[i]].Exchanges[currentTime+1e10].Ecc
			}
			// Set DCOB flag if appropriate
			/*
			 * This is not useful in classifying binaries because it shows only the
			 * type of the last binary in which the objects was found,
			 * BUT the presence of this flag make possible to know that the
			 * object was in a DCOB binary at some time at least one time
			 */
			if types[0] == "bh" && types[1] == "bh" {
				starMap[starId[i]].DCOB = "DBH"
			}
			if types[0] == "ns" && types[1] == "ns" {
				starMap[starId[i]].DCOB = "DNS"
			}
			if types[0]+types[1] == "bhns" || types[0]+types[1] == "nsbh" {
				starMap[starId[i]].DCOB = "BHNS"
			}

			// Zero eccentricity flag
			if ecc == 0 {
				starMap[starId[i]].ZeroEcc = true
			}
		}
		line += 1
	}
}

func (starMap StarMapType) Keys() (keys []string) {
	keys = make([]string, len(starMap))
	idx := 0
	for key, _ := range starMap {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)
	return keys
}

func (starMap StarMapType) Print(outDest io.Writer, long bool) {
	keys := starMap.Keys()
	for _, key := range keys {
		// fmt.Print(key, " ")
		starMap[key].Print(outDest, long)
		fmt.Fprintln(outDest, "===========================================================")
	}
}

// Count exchanges for all the stars
func (starMap StarMapType) CountExchanges() {
	defer debug.TimeMe(time.Now())
	// Compute exchanges
	for _, value := range starMap {
		value.CountExchanges()
	}
}

// Compute lifetimes for all the stars
func (starMap StarMapType) ComputeLifeTimes() {
	defer debug.TimeMe(time.Now())
	for _, value := range starMap {
		value.ComputeLifeTimes()
	}
}

type RegGroupsStruct struct {
	Id       int
	Comb     int
	N        int
	Z        int
	Rv       int
	Fpb      int
	W        int
	Tf       int
	Systime  int
	Phystime int
	Objs     int
	Hflag    int
	Types    int
	M0       int
	M1       int
	Sma      int
	Period   int
	Ecc      int
	Tgw      int
	Mchirp   int
}

func ChirpMass(m0, m1 float64) (chirpmass float64) {
	mTot := m0 + m1
	mu := (m0 * m1) / mTot

	chirpmass = math.Pow(mu, 3./5) * math.Pow(mTot, 2./5)
	return chirpmass
}

func Tgw(sma, ecc, m0, m1 float64) (tgw float64) {
	const PERIOD_UNIT float64 = 1000000
	const DISTANCE_UNIT float64 = 206264.806 // 1 Parsec = 206 264.806 Astronomical Units
	const PC2M = 3.0857e16
	const PC2AU float64 = 206264.806
	const SECONDS_IN_A_YEAR float64 = 60 * 60 * 24 * 365
	const LIGHT_SPEED float64 = 299792458 // m/s
	const G float64 = 6.67398e11          // m^3 kg^-1 s^-2
	const M_SUN float64 = 1.98855e30
	const YR2GYR float64 = 1000000000

	var CONSTANT float64 = (5. * (math.Pow(LIGHT_SPEED, 5)) * (math.Pow(PC2M, 4))) / (256 * (math.Pow(G, 3)) * (SECONDS_IN_A_YEAR * YR2GYR) * (math.Pow(M_SUN, 3)))

	tgw = CONSTANT * math.Pow(sma, 4) * math.Pow((1-math.Pow(ecc, 2)), (7./2)) / (m0 * m1 * (m0 + m1))

	return tgw
}

func NewStar(starId string, regexResult []string, regGroups RegGroupsStruct) (s *StarData) {
	s = &StarData{}
	// Fill StarData attribute of the new map
	s.StarId = starId
	s.Z = regexResult[regGroups.Z]
	s.NFile = regexResult[regGroups.N]
	s.DCOB = "--"
	s.ZeroEcc = false
	s.Comb = regexResult[regGroups.Comb]
	s.Rv, _ = strconv.ParseFloat(regexResult[regGroups.Rv], 64)
	s.Fpb, _ = strconv.ParseFloat(regexResult[regGroups.Fpb], 64)
	s.W, _ = strconv.ParseInt(regexResult[regGroups.W], 10, 64)
	s.Tf = regexResult[regGroups.Tf]
	return s
}

func NewBinary(regexResult []string, regGroups RegGroupsStruct) (b *BinaryData) {
	var (
		m0, m1, sma, ecc float64
	)

	b = &BinaryData{}
	b.BinaryId = regexResult[regGroups.Id]
	b.Hardness = regexResult[regGroups.Hflag]
	types := strings.Split(regexResult[regGroups.Types], "|")
	b.Types = types

	m0, _ = strconv.ParseFloat(regexResult[regGroups.M0], 64)
	m1, _ = strconv.ParseFloat(regexResult[regGroups.M1], 64)
	sma, _ = strconv.ParseFloat(regexResult[regGroups.Sma], 64)
	ecc, _ = strconv.ParseFloat(regexResult[regGroups.Ecc], 64)

	b.M0 = m0
	b.M1 = m1
	b.Sma, _ = strconv.ParseFloat(regexResult[regGroups.Sma], 64)
	b.Period, _ = strconv.ParseFloat(regexResult[regGroups.Period], 64)

	b.Mchirp = ChirpMass(m0, m1)
	b.Tgw = Tgw(sma, ecc, m0, m1)

	return b
}

func (starMap StarMapType) CsvSummary(outFileName string) {
	var (
		outFile *os.File
		nWriter *bufio.Writer
		err     error
	)

	// Open file for writing
	outFile, err = os.Create(outFileName)
	defer outFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()

	header := "#StarId,Comb,NFile,DCOB,Primordial,Z,Rv,Fpb,W0,Tf," +
		"HardExchangesNumber,SoftExchangesNumber,TotalExchangesNumber," +
		"AllSum,AllMean,HardDBHSum,HardDBHMean,SoftDBHSum,SoftDBHMean,HardDNSSum," +
		"HardDNSMean,SoftDNSSum,SoftDNSMean,HardBHNSSum,HardBHNSMean,SoftBHNSSum,SoftBHNSMean\n"

	format := "%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n"

	nWriter.WriteString(header)

	keys := starMap.Keys()
	for _, key := range keys {
		nWriter.WriteString(starMap[key].CsvSummary(format))
	}
}
