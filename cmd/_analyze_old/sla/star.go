package sla

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
	"time"
	
	"bitbucket.org/brunetto/goutils/readfile"
)


// StarMapType is the struct containing all the star data
type StarMapType map[string]*StarData

// FIXME: check the order of function calling
// StarData store the star data we are interest in.
type StarData struct {
	// StarId uniquely identifies the star among all the simulations.
	StarId string
	// Z is the metallicity in units of Zsun = 0.019.
	Z string
	// NFile is the file (run) number.
	NFile string
	// DCOB last type flag (DBH, DNS, BHNS).
	LastDCOB string
	// 
	DCOB StringSet
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
	TimeDom  TimeDomain
}

// Populate parse the input file and pupolate the star data
func (starMap StarMapType) Populate(inPath string, inFile string) (){
	var (
		fileObj *os.File
		nReader *bufio.Reader
		readLine string
		err error
		starId = make([]string, 2)
		sysTime = TimeDomain{Min: 18446744073709551615, Max: 0}
		currentTime uint64
		currentIds []string
		binaryRegexp *regexp.Regexp
		regexResult []string
	)
	tGlob0 := time.Now()
	// showPromiscuous is a package global variable, just like Verb
	if showPromiscuous {
		log.Println("Showing binaries containing a star already in another binary (showPromiscuous flag: ", showPromiscuous, ")")
	} else {
		log.Println("Hiding binaries containing a star already in another binary (showPromiscuous flag: ", showPromiscuous, ")")
	}
	
	// FIXME: decide decent names for files
	// Check infile
	if inFile == "all_the_fishes.txt" {
		binaryRegexp = regexp.MustCompile(regStringAllFishes)
	} else if strings.HasSuffix(inFile, "_all.txt") {
		binaryRegexp = regexp.MustCompile(regStringDBHAll)
	} else {
		log.Fatal("Unrecognized file type ", inFile)
	}
	
	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(filepath.Join(inPath, inFile)); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
		log.Fatal(err)
	}
	defer fileObj.Close()
	
	totCount, _ := readfile.LinesCount(fileObj)
	
	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)
			
	// Read the file and fill the starMap slice
	log.Println("Start reading file and populating starMap...")
	line := 0
	for {
		line++
		if readLine, err =  readfile.Readln(nReader); err != nil {
			if err.Error() != "EOF" {
				log.Println("Done reading ", line, " lines from file " , inFile, "  with err", err)
			}
			break
		}
		// Skip header and comments
		if readLine[0] == '#' {
			log.Println("Header/comment detected, skip...")
			continue
		}
		
		// Progress visuzlization
		fmt.Fprintf(os.Stderr, "\rParsed: %v %%", (100 * line) / totCount)
		
		regexResult = binaryRegexp.FindStringSubmatch(readLine)
		
		if regexResult == nil {
			log.Println("With regexp ", binaryRegexp)
			log.Fatal("no match, nil regex result on line ", line)
		} 
		
		// Update all binaries time domain if necessary
		// Not used until now but it can be useful
		currentTime, _ = strconv.ParseUint(regexResult[4], 10, 64)
		if sysTime.Min > currentTime {
			sysTime.Min = currentTime
		} else if sysTime.Max < currentTime {
			sysTime.Max = currentTime
		}
		
		currentIds = strings.Split(regexResult[6], "|")
		// If not yet present, create two entries in the map, 
		// one for each of the components of the binary
		// else only update the new timestep data&companion
		// FIXME: here is probably the key to have this piece general (multiple systems instead of binaries only)
		for i := 0; i < 2; i++ { 
			// Retrieve single object id
			
			
			// Compose single object id
			starId[i] = "Z" + regexResult[1] +
						"n" + LeftPad(regexResult[2], "0", 3) +
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
				starMap[starId[i]].NFile = LeftPad(regexResult[2], "0", 3)
				starMap[starId[i]].LastDCOB = "--"
				starMap[starId[i]].DCOB = NewStringSet()
				starMap[starId[i]].ZeroEcc = false
				starMap[starId[i]].TimeDom.Min = currentTime // Init first time in binary
				// FIXME: Check for primordial (~, in reality it is check for t=0) 
				// binary
				if currentTime == 0 {
					starMap[starId[i]].Primordial = true
				} else {
					starMap[starId[i]].Primordial = false
				}
				starMap[starId[i]].Exchanges = make(ExchangesMap)
			}
			if currentTime > starMap[starId[i]].TimeDom.Max {
				starMap[starId[i]].TimeDom.Max = currentTime
			}
			// Fill timestep companion and data
			// If the entry (=the timestep) already exists is a problem
			// it means that I have found the same star in two binaries 
			// and it is not a hyerarchical system nor a triple etc
			// because this have to be catch binaryData before (TO BE CHECKED)
			binData, zE := AddBinaryToExch(regexResult, currentIds, i)
			
			if _, exists := starMap[starId[i]].Exchanges[currentTime]; !exists {
				starMap[starId[i]].Exchanges[currentTime] = binData
				// Update properties (only if I consider this binary)
				starMap[starId[i]].LastDCOB = regexResult[8]
				starMap[starId[i]].DCOB.Add(regexResult[8])
				starMap[starId[i]].ZeroEcc = zE
			} else if showPromiscuous == false {
				starMap[starId[i]].Promiscuous = true
				continue // doesn't save promiscuous binaries
			} else { //save adding 1000 to the time
				starMap[starId[i]].Promiscuous = true
				starMap[starId[i]].Exchanges[currentTime+1000] = binData
				// Update properties (only if I consider this binary)
				starMap[starId[i]].LastDCOB = regexResult[8]
				starMap[starId[i]].DCOB.Add(regexResult[8])
				starMap[starId[i]].ZeroEcc = zE
				// Emit promiscuous warning
				if Verb {
					// Promiscuous warning during execution
					log.Println("################################################################")
					log.Println("## Reading line ", line)
					log.Println("## Found ", starId[i], " two times in the same timestep", currentTime)
					log.Println("## ", starMap[starId[i]].Exchanges[currentTime])
					log.Println("## ", starMap[starId[i]].Exchanges[currentTime+1000])
					log.Println("## At the moment I do nothing and add it as 1000+currentTime")
					log.Println("################################################################")
				}
			}
			
		}
	}
	fmt.Fprint(os.Stderr, "\n")
	tGlob1 := time.Now()
	log.Println("starMap populated with ", len(starMap), " stars in ", tGlob1.Sub(tGlob0))
} // End Populate



// Keys returns the sorted keys of the maps with all the stars
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

// ExtrcType return a map with only stars of a given last type.
// [bh|bh, ns|ns, bh|ns]
func (starMap StarMapType) ExtrcLastType(mapType string) (StarMapType) {
	if mapType != "bh|bh" && mapType != "ns|ns" && mapType != "bh|ns" {
		log.Fatal("Wrong type in ",  Whoami(false), " function, ", mapType, " not in [bh|bh, ns|ns, bh|ns]")
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
func (starMap StarMapType) WasInType(mapType string) (StarMapType) {
	if mapType != "bh|bh" && mapType != "ns|ns" && mapType != "bh|ns" {
		log.Fatal("Wrong type in ",  Whoami(false), " function, ", mapType, " not in [bh|bh, ns|ns, bh|ns]")
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
func (starMap StarMapType) WasPromiscuous() (StarMapType) {
	starMapPerType := make(StarMapType)
	for key, value := range starMap {
		if value.Promiscuous {
			starMapPerType[key] = value
		} 
	}
	return starMapPerType
}

// Print prints StarData data
func (starData *StarData) Print() () {
	fmt.Println("#########################")
	fmt.Println("#       Star Data       #")
	fmt.Println("#########################")
	fmt.Println("StarId = ", starData.StarId)
	fmt.Println("Z = ", starData.Z)
	fmt.Println("NFile = ", starData.NFile)
	fmt.Println("Last DCOB = ", starData.LastDCOB)
	fmt.Println("DCOB = ", starData.DCOB.String())
	fmt.Println("Primordial = ", starData.Primordial)
	fmt.Println("First time in binary = ", starData.TimeDom.Min)
	fmt.Println("Last time in binary = ", starData.TimeDom.Max)
}

// PrintWithExchs prints StarData data with exchanges
func (starData *StarData) PrintWithExchs() () {
	starData.Print()
	fmt.Println("#########################")
	fmt.Println("#       Exchanges       #")
	fmt.Println("#########################")
	starData.Exchanges.Print()
}











