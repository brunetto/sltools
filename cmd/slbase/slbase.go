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
	"github.com/spf13/cobra"
)

var (
	inFileName string
)

var slBaseCmd = &cobra.Command{
	Use:   "slbase",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type slbase help for help.")
	},
}

var loadNSmergersCmd = &cobra.Command{
	Use:   "",
	Short: "Load NS mergers",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a file containing the NS mergers in the format [TODO] to load")
		}
		loadNSmergers(inFileName)
	},
}

var modifydbCmd = &cobra.Command{
	Use:   "",
	Short: "Correlate DB collections to complete the binary data",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		modifydb()
	},
}

var loadPositionsCmd = &cobra.Command{
	Use:   "",
	Short: "Load positions",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a file containing the positions in the format [TODO] to load")
		}
		loadPositions(inFileName)
	},
}

var loadSimsCmd = &cobra.Command{
	Use:   "",
	Short: "Load simulations",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a file containing the simulated binaries in the format [TODO] to load")
		}
		loadSims(inFileName)
	},
}

func main() {
	defer debug.TimeMe(time.Now())

	initCommands()
	slBaseCmd.Execute()

}

func initCommands() {
	slBaseCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Input file")
	slBaseCmd.AddCommand(loadSimsCmd)
	slBaseCmd.AddCommand(loadPositionsCmd)
	slBaseCmd.AddCommand(modifydbCmd)
	slBaseCmd.AddCommand(loadNSmergersCmd)
}

type NSBinaryData struct {
	BinaryId string `bson:"binaryid"`
}

// Load NS mergers form file to DB
func loadNSmergers(nsMergerFile string) {
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
		binary       *NSBinaryData
	)

	// 	if len(os.Args) < 2 {
	// 		log.Fatal("Provide the file with the NS merger.")
	// 	}

	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(nsMergerFile); err != nil {
		log.Fatal(os.Stderr, "Can't open %s: error: %s\n", nsMergerFile, err)
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

		binary = &NSBinaryData{}
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

// TODO: merger structs using optional fields
type CompleteBinaryData struct {
	BinaryId   string   `bson:"binaryid"`
	Comb       int64    `bson:"comb"`
	N          int64    `bson:"n"`
	Metalicity float64  `bson:"z"` // TODO: fix database for this key
	Rv         float64  `bson:"rv"`
	Fpb        float64  `bson:"fpb"`
	W          float64  `bson:"w"`
	Tf         string   `bson:"tf"`
	SysTime    int64    `bson:"systime"`
	PhysTime   float64  `bson:"phystime"`
	Objs       string   `bson:"objs"`
	Hflag      string   `bson:"hflag"`
	Types      []string `bson:"types"`
	M0         float64  `bson:"m0"`
	M1         float64  `bson:"m1"`
	Sma        float64  `bson:"sma"`
	Period     float64  `bson:"period"`
	Tgw        float64  `bson:"tgw"`
	Mchirp     float64  `bson:"mchirp"`
	Ecc        float64  `bson:"ecc"`
	X          float64  `bson:"X"`
	Y          float64  `bson:"Y"`
	Z          float64  `bson:"Z"`
	R          float64  `bson:"R"`
	Merge      bool     `bson:"merge"`
}

func modifydb() {
	log.Fatal("To be finished...")
	// 	var (
	// 		err        error
	// 		session    *mgo.Session
	// 		database   *mgo.Database
	// 		cBinaries  *mgo.Collection
	// 		cPositions *mgo.Collection
	// 		cNSmergers *mgo.Collection
	// 		binary     CompleteBinaryData
	// 		binaries   []CompleteBinaryData
	// 		query      bson.M
	// 		q          *mgo.Query
	// 	)
	//
	// 	log.Println("Open the database")
	// 	// TODO: generalize
	// 	if session, err = mgo.Dial("mongodb://localhost:27017/"); err != nil {
	// 		log.Fatal("Error dialing db: ", err)
	// 	}
	// 	defer session.Close()
	//
	// 	// Optional. Switch the session to a monotonic behavior.
	// 	session.SetMode(mgo.Monotonic, true)
	//
	// 	log.Println("Open database")
	// 	database = session.DB("phd")
	//
	// 	log.Println("Open collection")
	// 	cBinaries = database.C("binaries")
	// 	cPositions = database.C("positions")
	// 	cNSmergers = database.C("ns_mergers")
	//
	// 	log.Println("Retrieve each binary, check for position and NS merger and update the document")

	// 1. Iterate over the binaries that are record of binaries collections
	// 2. For each binary search if it is a NS merger, in case, update
	// 3. Search if it has a position, in case update

	// 	query = bson.M{"$and": []bson.M{bson.M{"binaryid": "c76n4a103540b3540"}, bson.M{"systime": 13}}}
	// 	q = collection.Find(query) //.Sort("phystime")
	// 	q.All(&binaries)
	// 	var i int
	// 	var item CompleteBinaryData
	// 	for i, item = range binaries {
	// 		fmt.Printf("%v: %+v\n", i, item)
	// 	}

}

type PositionBinaryData struct {
	SysTime  int64   `bson:"systime"`
	PhysTime float64 `bson:"phystime"`
	BinaryId string  `bson:"binaryid"`
	X        float64 `bson:"x"`
	Y        float64 `bson:"y"`
	Z        float64 `bson:"z"`
	R        float64 `bson:"r"`
}

func loadPositions(positionFile string) {
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
		binary       *PositionBinaryData
	)

	if positionFile == "" {
		log.Fatal("Provide the file with the positions.")
	}

	// Open the file
	log.Println("Opening data file...")
	if fileObj, err = os.Open(positionFile); err != nil {
		log.Fatal(os.Stderr, "Can't open %s: error: %s\n", positionFile, err)
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
			log.Println("Done reading ", line, " lines from file ", positionFile, "  with err", err)
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

		binary = &PositionBinaryData{}
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

// Load simulation CVS file to database
// TODO: cod on required CSV file
func loadSims(simulationsCSVfile string) {
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
	if fileObj, err = os.Open(simulationsCSVfile); err != nil {
		log.Fatal(os.Stderr, "Can't open %s: error: %s\n", simulationsCSVfile, err)
	}
	defer fileObj.Close()

	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)

	log.Println("Open the database")
	// TODO: generalize with DB choice
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
			log.Println("Done reading ", line, " lines from file ", simulationsCSVfile, "  with err", err)
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

		if line%1000 == 0 {
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
