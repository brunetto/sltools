package slan

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/brunetto/goutils/readfile"

	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/sets"
)

// AllData contains the maps ofall the stars and the binaries.
// The keys are the name of the star/binary, the values the struct containing
// the data
type AllDataType struct {
	Stars    StarMapType
	Binaries BinaryMapType
}

func /*(allData *AllDataType)*/ New() *AllDataType {
	return &AllDataType{
		Stars:    StarMapType{},
		Binaries: BinaryMapType{},
	}
}

func (allData *AllDataType) Init() *AllDataType {
	return &AllDataType{
		Stars:    StarMapType{},
		Binaries: BinaryMapType{},
	}
}

// Populate parse the input file and pupolate the star data
func (allData *AllDataType) Populate(inPath string, inFile string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		fileObj      *os.File
		nReader      *bufio.Reader
		readLine     string
		err          error
		starId       = make([]string, 2)
		sysTime      = TimeDomain{Min: 18446744073709551615, Max: 0}
		currentTime  uint64
		currentIds   []string
		binaryRegexp *regexp.Regexp
		regexResult  []string
		binaryId string
		zE bool
		physTime, timeUnit float64
	)
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
		log.Printf("%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
		log.Fatal(err)
	}
	defer fileObj.Close()

	totCount, _ := readfile.LinesCount(fileObj)
	log.Println(totCount, " lines")

	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)

	// Read the file and fill the allData.Stars slice
	log.Println("Start reading file and populating allData.Stars...")
	line := 0
	for {
		line++
		if readLine, err = readfile.Readln(nReader); err != nil {
			if err.Error() != "EOF" {
				log.Println("Done reading ", line, " lines from file ", inFile, "  with err", err)
			}
			break
		}
		// Skip header and comments
		if readLine[0] == '#' {
			log.Println("Header/comment detected, skip...")
			continue
		}

		// Progress visuzlization
		fmt.Fprintf(os.Stderr, "\rParsed: %v %%", (100*line)/totCount)

		regexResult = binaryRegexp.FindStringSubmatch(readLine)

		if regexResult == nil {
			log.Printf("\nWith regexp %v \n and line %v ", binaryRegexp, readLine)
			log.Fatal("no match, nil regex result on line ", line)
		}

		// Update all binaries time domain if necessary
		// Not used until now but it can be useful
		if currentTime, err = strconv.ParseUint(regexResult[4], 10, 64); err != nil {
			log.Fatalf("Can't parse current time from %v with error %v\n", regexResult[4], err)
		}
		if sysTime.Min > currentTime {
			sysTime.Min = currentTime
		} else if sysTime.Max < currentTime {
			sysTime.Max = currentTime
		}

		currentIds = strings.Split(regexResult[6], "|")
		// If not yet present, create two entries in the map,
		// one for each of the components of the binary
		// else only update the new timestep data&companion
		// FIXME: here is probably the key to have this piece general 
		// (multiple systems instead of binaries only)
		for i := 0; i < 2; i++ {
			// Compose single object id
			starId[i] = "Z" + regexResult[1] +
				"n" + goutils.LeftPad(regexResult[2], "0", 3) +
				"id" + currentIds[i]

			// Check existence of the single object in the map
			// if it exists then update the timestep/companion table,
			// else if not then create the entry
			if _, exists := allData.Stars[starId[i]]; !exists {
				// Create the new StarData value for the key allData.Stars in the map
				allData.Stars[starId[i]] = &StarData{
					// Fill StarData attribute of the new map
					StarId:      starId[i],
					Z:           regexResult[1],
					NFile:       goutils.LeftPad(regexResult[2], "0", 3),
					LastDCOB:    "--",
					DCOB:        sets.NewStringSet(),
					ZeroEcc:     false,
					TimeDom: TimeDomain{Min: currentTime}, // Init first time in binary
				}
								
				// FIXME: Check for primordial (~, in reality it is check for t=0)
				// binary
				if currentTime == 0 {
					allData.Stars[starId[i]].Primordial = true
				} else {
					allData.Stars[starId[i]].Primordial = false
				}
				allData.Stars[starId[i]].Exchanges = make(ExchangesMap)
			}
			if currentTime > allData.Stars[starId[i]].TimeDom.Max {
				allData.Stars[starId[i]].TimeDom.Max = currentTime
			}
			
			if currentTime > 0 {
				if physTime, err = strconv.ParseFloat(regexResult[5], 64); err != nil {
					log.Fatalf("Can't parse phys time in %v with error %v\n", regexResult[5], err)
				}
				// Fix round errors
				timeUnit = 1e-6 * (math.Trunc(1e6 * (physTime / float64(currentTime))))
// 				if allData.Stars[starId[i]].TimeUnit != 0 {
// 					if allData.Stars[starId[i]].TimeUnit != timeUnit {
// 						log.Fatalf("Time unit different, from %v to %v\n", allData.Stars[starId[i]].TimeUnit, timeUnit)
// 					}
// 				} else {
					allData.Stars[starId[i]].TimeUnit = timeUnit
// 				}
			}
			// Fill timestep companion and data
			// If the entry (=the timestep) already exists is a problem
			// it means that I have found the same star in two binaries
			// and it is not a hyerarchical system nor a triple etc
			// because this have to be catch binaryData before (TO BE CHECKED)
			if i == 0 {
				binaryId, zE = allData.Binaries.AddBinary(regexResult)
			}
			
			
			if _, exists := allData.Stars[starId[i]].Exchanges[currentTime]; !exists {
				allData.Stars[starId[i]].Exchanges[currentTime] = &ExchData{
																BinaryId: binaryId, 
																Companion: currentIds[len(currentIds)-i-1],
																}
				// Update properties (only if I consider this binary)
				allData.Stars[starId[i]].LastDCOB = regexResult[8]
				allData.Stars[starId[i]].DCOB.Add(regexResult[8])
				allData.Stars[starId[i]].ZeroEcc = zE
			} else if showPromiscuous == false {
				allData.Stars[starId[i]].Promiscuous = true
				continue // doesn't save promiscuous binaries
			} else { //save adding 1000 to the time
				allData.Stars[starId[i]].Promiscuous = true
				allData.Stars[starId[i]].Exchanges[currentTime+1000] = &ExchData{
																BinaryId: binaryId, 
																Companion: currentIds[len(currentIds)-i-1],
																}
				// Update properties (only if I consider this binary)
				allData.Stars[starId[i]].LastDCOB = regexResult[8]
				allData.Stars[starId[i]].DCOB.Add(regexResult[8])
				allData.Stars[starId[i]].ZeroEcc = zE
				// Emit promiscuous warning
				if Verb {
					// Promiscuous warning during execution
					log.Println("################################################################")
					log.Println("## Reading line ", line)
					log.Println("## Found ", starId[i], " two times in the same timestep", currentTime)
					log.Println("## ", allData.Stars[starId[i]].Exchanges[currentTime])
					log.Println("## ", allData.Stars[starId[i]].Exchanges[currentTime+1000])
					log.Println("## At the moment I do nothing and add it as 1000+currentTime")
					log.Println("################################################################")
				}
			}

		}
	}
	fmt.Fprint(os.Stderr, "\n")

} // End Populate
