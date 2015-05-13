package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"
	// 	"gopkg.in/mgo.v2/bson"

	"github.com/brunetto/goutils/debug"
)

type BinaryData struct {
	BinaryId string  `bson:"binaryid"`
}

func main() {
	defer debug.TimeMe(time.Now())
	var (
		// Regexp string for all_the_fishes
		regStringNSMgr string = `^(\S+)` // 1: BinaryId

		binaryRegexp = regexp.MustCompile(regStringNSMgr)
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
		log.Fatal("Provide the file with the NS merger.")
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
	collection = database.C("ns_mergers")

	// Clean
	if err = collection.DropCollection(); err != nil {
		log.Println("Error droppping collection: ", err)
	}

	collection = database.C("ns_mergers")

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
			log.Println("Reg is:", regStringNSMgr)
			log.Fatal("Line is:", readLine)
		}

		binary = &BinaryData{}
		binary.BinaryId = regexResult[1]

		if err != nil {
			log.Fatalf("Error parsing binary %+v: %v \n", binary, err)
		}

		if err = collection.Insert(binary); err != nil {
			log.Fatalf("Error inserting binary %+v into database: %v \n", binary, err)
		}

		fmt.Printf("\r Inserted line #: %v with id %v", line, binary.BinaryId)
		
		line++

	}

	fmt.Println()
}
