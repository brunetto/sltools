package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	// 	"gopkg.in/mgo.v2/bson"
	
	"github.com/brunetto/goutils/debug"
)

type PartialBinaryData struct {
	BinaryId string
	Comb     int64
	N        int64
	Z        float64
	Rv       float64
	Fpb      float64
	W        float64
	Tf       string
	SysTime  int64
	PhysTime float64
	Objs     string
	Hflag    string
	Types    []string
	M0       float64
	M1       float64
	Sma      float64
	Period   float64
	Tgw      float64
	Mchirp   float64
	Ecc      float64
	Merge    bool
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
	SysTime  int
	PhysTime int
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

func main() {
defer debug.TimeMe(time.Now())
	var (
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
			`(\d+\.*\d*e*[-\+]*\d*),` + // "phystime": 10,
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
		regGroups = RegGroupsStruct{Id: 1,
			Comb:     2,
			N:        3,
			Z:        4,
			Rv:       5,
			Fpb:      6,
			W:        7,
			Tf:       8,
			SysTime:  9,
			PhysTime: 10,
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

		binaryRegexp = regexp.MustCompile(regStringCSVsims)
		regexResult  []string
		fileObj      *os.File
		err          error
		nReader      *bufio.Reader
		session      *mgo.Session
		database     *mgo.Database
		collection   *mgo.Collection
		readLine     string
		binary       *PartialBinaryData
	)

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
	collection = database.C("binaries")
	
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
			log.Println("Reg is:", regStringCSVsims)
			log.Fatal("Line is:", readLine)
		}

		binary = &PartialBinaryData{}
		binary.BinaryId = regexResult[regGroups.Id]
		binary.Comb, err = strconv.ParseInt(regexResult[regGroups.Comb], 10, 64)
		binary.N, err = strconv.ParseInt(regexResult[regGroups.N], 10, 64)
		binary.Z, err = strconv.ParseFloat(regexResult[regGroups.Z], 64)
		binary.Rv, err = strconv.ParseFloat(regexResult[regGroups.Rv], 64)
		binary.Fpb, err = strconv.ParseFloat(regexResult[regGroups.Fpb], 64)
		binary.W, err = strconv.ParseFloat(regexResult[regGroups.W], 64)
		binary.Tf = regexResult[regGroups.Tf]
		binary.SysTime, err = strconv.ParseInt(regexResult[regGroups.SysTime], 10, 64)
		binary.PhysTime, err = strconv.ParseFloat(regexResult[regGroups.PhysTime], 64)
		binary.Objs = regexResult[regGroups.Objs]
		binary.Hflag = regexResult[regGroups.Hflag]
		binary.Types = strings.Split(regexResult[regGroups.Types], "|")
		binary.M0, err = strconv.ParseFloat(regexResult[regGroups.M0], 64)
		binary.M1, err = strconv.ParseFloat(regexResult[regGroups.M1], 64)
		binary.Sma, err = strconv.ParseFloat(regexResult[regGroups.Sma], 64)
		binary.Period, err = strconv.ParseFloat(regexResult[regGroups.Period], 64)
		binary.Tgw, err = strconv.ParseFloat(regexResult[regGroups.Tgw], 64)
		binary.Mchirp, err = strconv.ParseFloat(regexResult[regGroups.Mchirp], 64)
		binary.Ecc, err = strconv.ParseFloat(regexResult[regGroups.Ecc], 64)
		binary.Merge = false

		if err != nil {
			log.Fatalf("Error parsing binary %+v: %v \n", binary, err)
		}
		
		if err = collection.Insert(binary); err != nil {
			log.Fatalf("Error inserting binary %+v into database: %v \n", binary, err)
		}
		
		if line % 1000 == 0 {
			fmt.Printf("\r Inserted line #: %v with phystime %v", line, binary.PhysTime)
		}
		line++

	}

	fmt.Println()
	// 	result := Person{}
	// 	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	fmt.Println("Phone:", result.Phone)
}
