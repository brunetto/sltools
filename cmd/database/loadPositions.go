package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	// 	"gopkg.in/mgo.v2/bson"

	"github.com/brunetto/goutils/debug"
)

type BinaryData struct {
	SysTime  int64   `bson:"systime"`
	PhysTime float64 `bson:"phystime"`
	BinaryId string  `bson:"binaryid"`
	X        float64 `bson:"x"`
	Y        float64 `bson:"y"`
	Z        float64 `bson:"z"`
	R        float64 `bson:"r"`
}

func main() {
	defer debug.TimeMe(time.Now())
	var (
		// Regexp string for all_the_fishes
		regStringPos string = `^(\d+),` + // 1: SysTime
			`(-*\d+\.*\d*e*[-\+]*\d*),` + // 2: PhysTime
			`(\S+?),` + // 3: BinaryId
			`(-*\d+\.*\d*e*[-\+]*\d*),` + // 4: X
			`(-*\d+\.*\d*e*[-\+]*\d*),` + // 5: Y
			`(-*\d+\.*\d*e*[-\+]*\d*),` + // 6: Z
			`(-*\d+\.*\d*e*[-\+]*\d*)` // 7: R

		binaryRegexp = regexp.MustCompile(regStringPos)
		regexResult  []string
		fileObj      *os.File
		err          error
		nReader      *bufio.Reader
		session      *mgo.Session
		database     *mgo.Database
		collection   *mgo.Collection
		readLine     string
		binary       *BinaryData
	)

	if len(os.Args) < 2 {
		log.Fatal("Provide the file with the positions.")
	}

	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(os.Args[1]); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], os.Args[1], err)
	}
	defer fileObj.Close()

	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)

	log.Println("Open the database")
	if session, err = mgo.Dial("mongodb://localhost:27017/"); err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	database = session.DB("phd")
	collection = database.C("positions")

	// Clean
	// 	if err = collection.DropCollection(); err != nil {
	// 		log.Println("Error droppping collection: ", err)
	// 	}
	//
	// 	collection = database.C("binaries")

	// Read the file and fill the starMap slice
	log.Println("Start reading file and populating the DB...")
	line := 0
	for {
		if readLine, err = nReader.ReadString('\n'); err != nil {
			fmt.Println()
			log.Println("Done reading ", line, " lines from file ", os.Args[1], "  with err", err)
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
			log.Println("Reg is:", regStringPos)
			log.Fatal("Line is:", readLine)
		}

		binary = &BinaryData{}
		binary.BinaryId = regexResult[3]
		binary.SysTime, err = strconv.ParseInt(regexResult[1], 10, 64)
		binary.PhysTime, err = strconv.ParseFloat(regexResult[2], 64)
		binary.X, err = strconv.ParseFloat(regexResult[4], 64)
		binary.Y, err = strconv.ParseFloat(regexResult[5], 64)
		binary.Z, err = strconv.ParseFloat(regexResult[6], 64)
		binary.R, err = strconv.ParseFloat(regexResult[7], 64)

		if err != nil {
			log.Fatalf("Error parsing binary %+v: %v \n", binary, err)
		}

		if err = collection.Insert(binary); err != nil {
			log.Fatalf("Error inserting binary %+v into database: %v \n", binary, err)
		}

		if line%1000 == 0 {
			fmt.Printf("\r Inserted line #: %v with id %v at time %v ", line, binary.BinaryId, binary.PhysTime)
		}
		line++

	}

	fmt.Println()
}
