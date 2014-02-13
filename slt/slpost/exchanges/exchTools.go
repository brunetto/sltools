package exchanges

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

)


//Store the min and max time of the binaries in the file
type SysTime struct {
	Min uint64
	Max uint64
}

// Data of the star in binary and that changes with time (companion, orbital properties, ...)
type BinaryData struct {
	BinaryId string
	Companion string
	Hardness string
	Types []string
	Ecc float64
}
// Print function for an exchange datum
func (binaryData *BinaryData) Print() {
	fmt.Println("\t", binaryData.BinaryId, binaryData.Companion, binaryData.Hardness, binaryData.Types, binaryData.Ecc)
}

// NOTE: use a pointer otherwise structs will be unchangeble!!!!!!!
type ExchangesMap map[uint64]*BinaryData

// Store the unchangeble star data we are interestbinaryData in
type StarData struct {
	StarId string
	Z string
	NFile string
	// DCOB flag
	DCOB string
	// Zero eccentricity flag for this star at some time
	ZeroEcc bool
	// Promiscuity flag for this star at some time
	Promiscuous bool
	// Start primordial flag
	Primordial bool
	// The key is the sys_time: I've tought it was unique
	Exchanges ExchangesMap
	// After counting exchanges in Exchanges
	ExchangesNumbers *ExchangesData
	// Lifetimes
	LifeTimes LifeTimeStruct
}

func (starData *StarData) Keys() (keys []uint64) {
	// Retrieve and sort Exchanges map keys
	excTimes := make([]uint64, len(starData.Exchanges))
	idx := 0 
	for key, _ := range starData.Exchanges {
        excTimes[idx] = key
        idx++
    }
    sort.Sort(uint64arr(excTimes))
	keys = uint64arr(excTimes)
	return keys
}

// Print function for StarData
func (starData *StarData) Print() {
	fmt.Println("#########################")
	fmt.Println("#       Star Data       #")
	fmt.Println("#########################")
	fmt.Println("StarId = ", starData.StarId)
	fmt.Println("Z = ", starData.Z)
	fmt.Println("NFile = ", starData.NFile)
	fmt.Println("DCOB = ", starData.DCOB)
	fmt.Println("Primordial = ", starData.Primordial)
	fmt.Println("#########################")
	fmt.Println("#       Exchanges       #")
	fmt.Println("#########################")
	// Retrieve and sort Exchanges map keys
	excTimes := starData.Keys()
	fmt.Println("SysTime, BinaryId, Companion, Hardness, Types, Ecc")
	for _, key := range excTimes {
		fmt.Print(key, " ")
		starData.Exchanges[uint64(key)].Print()
	}	
}

type LifeTimeStruct struct {
	All []uint64
	HardDBH []uint64
	SoftDBH []uint64
	HardDNS []uint64
	SoftDNS []uint64
	HardBHNS []uint64
	SoftBHNS []uint64
}

func (starData *StarData) ComputeLifeTimes() () {
// 	log.Println("Counting lifetimes ...")
	idx := 0
	time0 := uint64(0)
	time1 := uint64(0)
	binary := "0"
	delta := uint64(0)
	excTimes := starData.Keys()
	for _, key := range excTimes {
		idx++
		if (binary != starData.Exchanges[key].BinaryId) || (idx == len(excTimes)) {
			delta = time1-time0 + 1
// 			starData.LifeTimes = append(starData.LifeTimes, delta)
			time0 = key
			// this can happen considering promiscuous binaries with time=1000+time
			if delta > 400 {
				starData.Print()
				log.Fatal("Delta ", delta, starData.StarId,  starData.Exchanges[key].Types, starData.Exchanges[key].Hardness)
			}
			if starData.Exchanges[key].Hardness == "H" {
				if starData.Exchanges[key].Types[0] == "bh" && starData.Exchanges[key].Types[1] == "bh" {
					starData.LifeTimes.HardDBH = append(starData.LifeTimes.HardDBH, delta)
// 					log.Println("hardDBH")
				}
				if starData.Exchanges[key].Types[0] == "ns" && starData.Exchanges[key].Types[1] == "ns" {
					starData.LifeTimes.HardDNS = append(starData.LifeTimes.HardDNS, delta)
// 					log.Println("hardDNS")
				}
				if (starData.Exchanges[key].Types[0]+starData.Exchanges[key].Types[1] == "bhns" || starData.Exchanges[key].Types[0]+starData.Exchanges[key].Types[1] == "nsbh") {
					starData.LifeTimes.HardBHNS = append(starData.LifeTimes.HardBHNS, delta)
// 					log.Println("hardmix")
				}
			} else if starData.Exchanges[key].Hardness == "S" {
				if starData.Exchanges[key].Types[0] == "bh" && starData.Exchanges[key].Types[1] == "bh" {
					starData.LifeTimes.SoftDBH = append(starData.LifeTimes.SoftDBH, delta)
// 					log.Println("softDBH")
				}
				if starData.Exchanges[key].Types[0] == "ns" && starData.Exchanges[key].Types[1] == "ns" {
					starData.LifeTimes.SoftDNS = append(starData.LifeTimes.SoftDNS, delta)
// 					log.Println("softDNS")
				}
				if (starData.Exchanges[key].Types[0]+starData.Exchanges[key].Types[1] == "bhns" || starData.Exchanges[key].Types[0]+starData.Exchanges[key].Types[1] == "nsbh") {
					starData.LifeTimes.SoftBHNS = append(starData.LifeTimes.SoftBHNS, delta)
// 					log.Println("softmix")
				}
			} 
// 			log.Println(starData.Exchanges[key].Hardness, starData.Exchanges[key].Types)
			starData.LifeTimes.All = append(starData.LifeTimes.All, delta)
// 			log.Println("all")
			
		}
// 		log.Println(lifeTimes)
		time1 = key
		binary = starData.Exchanges[key].BinaryId
// 		fmt.Println(time0, time1, binary)
	}
}

func (starData *StarData) CountExchanges() (*ExchangesData) {
	keys := starData.Keys()
	excData := new(ExchangesData)
	excData.StarId = starData.StarId
	excData.Primordial = starData.Primordial
	hardCompanions := make([]string, 1)
	softCompanions := make([]string, 1)
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
	excData.HardExchanges = hardCompanions[1:]
	excData.SoftExchanges = softCompanions[1:]
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
	starData.ExchangesNumbers = excData
	return excData
}

type ExchangesData struct {
	StarId string
	Primordial bool
	HardExchanges []string
	HardExchangesNumber int
	SoftExchanges []string
	SoftExchangesNumber int
	TotalExchanges int
}

func (exchangesData *ExchangesData) Print() {
	fmt.Println("#########################")
	fmt.Println("#   Exchanges Summary   #")
	fmt.Println("#########################")
	fmt.Println("Data for ", exchangesData.StarId)
	fmt.Println("Primordial binary ", exchangesData.Primordial)
	fmt.Println("HardExchanges = ", exchangesData.HardExchanges)
	fmt.Println("SoftExchanges = ", exchangesData.SoftExchanges)
	fmt.Println("HardExchangesNumber = ", exchangesData.HardExchangesNumber)
	fmt.Println("SoftExchangesNumber = ", exchangesData.SoftExchangesNumber)
	fmt.Println("TotalExchanges = ", exchangesData.TotalExchanges)
}

type StarMapType map[string]*StarData

func (starMap StarMapType) Populate(inPath string, inFile string, regString string, hidePromiscuous bool) (){
	var fileObj *os.File
	var nReader *bufio.Reader
	var readLine string
	var err error
	var starId = make([]string, 2)
	var sysTime = SysTime{Min: 0, Max: 0}
	var currentTime uint64
	var currentIds []string
	var binaryRegexp = regexp.MustCompile(regString)
	var regexResult []string
	
		
	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(filepath.Join(inPath, inFile)); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
		panic(err)
	}
	defer fileObj.Close()
	
	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)
		
	// Read the file and fill the starMap slice
	log.Println("Start reading file and populating starMap...")
	line := 0
	for {
		if readLine, err = nReader.ReadString('\n'); err != nil {
			log.Println("Done reading ", line, " lines from file " , inFile, "  with err", err)
			break
		}
		if readLine[0] == '#' {
			log.Println("Header detected, skip...")
			continue
		}
		regexResult = binaryRegexp.FindStringSubmatch(readLine)
		if regexResult == nil {
			log.Println("With regexp ", binaryRegexp)
			log.Fatal("no match, nil regex result on line ", line)
		}
		
		// Update time domain if necesary
		currentTime, _ = strconv.ParseUint(regexResult[4], 10, 64)
		if sysTime.Min > currentTime {
			sysTime.Min = currentTime
		} else if sysTime.Max < currentTime {
			sysTime.Max = currentTime
		}
		
		// If not yet present, create two entries in the map, 
		// one for each of the components of the binary
		// else only update the new timestep data&companion
		
		for i := 0; i < 2; i++ {
			// Retrieve single object id
			currentIds = strings.Split(regexResult[6], "|")
			
			// Compose single object id
			starId[i] = "Z" + regexResult[1] +
						"n" + regexResult[2] +
						"id" + currentIds[i]

			// Check existence of the single object in the map
			// if it exists then update the timestep/companion table,
			// else if not then create the entry
			if _, exists := starMap[starId[i]]; !exists {
				// Create the new StarData value for the key starMap in the map
				starMap[starId[i]] = new(StarData) // or var p = &Pippo
				// Fill StarData attribute of the new map
				starMap[starId[i]].StarId = starId[i]
				starMap[starId[i]].Z = regexResult[1]
				starMap[starId[i]].NFile = regexResult[2]
				starMap[starId[i]].DCOB = "--"
				starMap[starId[i]].ZeroEcc = false
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
				starMap[starId[i]].Exchanges[currentTime] = new(BinaryData)
				starMap[starId[i]].Exchanges[currentTime].BinaryId = regexResult[3]
				starMap[starId[i]].Exchanges[currentTime].Companion = currentIds[len(currentIds)-i-1]
				starMap[starId[i]].Exchanges[currentTime].Hardness = regexResult[7]
				types := strings.Split(regexResult[8], "|")
				starMap[starId[i]].Exchanges[currentTime].Types = types
				// Set DCOB flag if appropriate
				/*
				 This is not useful in classifying binaries because it shows only the 
				 type of the last binary in which the objects was found, 
				 BUT the presence of this flag make possible to know that the 
				 object was in a DCOB binary at some time at least one time				 
				 */
				if types[0] == "bh" && types[1] == "bh" {
					starMap[starId[i]].DCOB = "DBH"
				}
				if types[0] == "ns" && types[1] == "ns" {
					starMap[starId[i]].DCOB = "DNS"
				}
				if (types[0]+types[1] == "bhns" || types[0]+types[1] == "nsbh") {
					starMap[starId[i]].DCOB = "BHNS"
				}
				starMap[starId[i]].Exchanges[currentTime].Ecc, _ = strconv.ParseFloat(regexResult[13], 64)
				if starMap[starId[i]].Exchanges[currentTime].Ecc == 0 {
					starMap[starId[i]].ZeroEcc = true
				}
			} else if hidePromiscuous == true {
				continue // doesn't save promiscuous binaries
			} else { //save adding 1000 to the time
				starMap[starId[i]].Promiscuous = true
				starMap[starId[i]].Exchanges[currentTime+1000] = new(BinaryData)
				starMap[starId[i]].Exchanges[currentTime+1000].BinaryId = regexResult[3]
				starMap[starId[i]].Exchanges[currentTime+1000].Companion = currentIds[len(currentIds)-i-1]
				starMap[starId[i]].Exchanges[currentTime+1000].Hardness = regexResult[7]
				types := strings.Split(regexResult[8], "|")
				starMap[starId[i]].Exchanges[currentTime+1000].Types = types
				// Set DCOB flag if appropriate
				if types[0] == "bh" && types[1] == "bh" {
					starMap[starId[i]].DCOB = "DBH"
				}
				if types[0] == "ns" && types[1] == "ns" {
					starMap[starId[i]].DCOB = "DNS"
				}
				if (types[0]+types[1] == "bhns" || types[0]+types[1] == "nsbh") {
					starMap[starId[i]].DCOB = "BHNS"
				}
				/*
				// Promiscuous warning during execution
				log.Println("################################################################")
				log.Println("## Reading line ", line)
				log.Println("## Found ", starId[i], " two times in the same timestep", currentTime)
				log.Println("## ", starMap[starId[i]].Exchanges[currentTime])
				log.Println("## ", starMap[starId[i]].Exchanges[currentTime+1000])
// 				log.Println("Adieu...")
				log.Println("## At the moment I do nothing and add it as 1000+currentTime")
				log.Println("################################################################")
				*/
				if starMap[starId[i]].Exchanges[currentTime+1000].Ecc == 0 {
					starMap[starId[i]].ZeroEcc = true
				}
			}
		}
		line += 1
	}
	log.Println("starMap populated with ", len(starMap), " stars.")
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

func (starMap StarMapType) Print() {
	keys := starMap.Keys()
	for _, key := range keys {
// 		fmt.Print(key, " ")
		starMap[key].Print()
		fmt.Println("===========================================================")
	}
}

func (starMap StarMapType) CountExchanges() {
	// Compute exchanges
	for _, value := range starMap {
		value.CountExchanges()
	}
}

type uint64arr []uint64
func (a uint64arr) Len() int { return len(a) }
func (a uint64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a uint64arr) Less(i, j int) bool { return a[i] < a[j] }



