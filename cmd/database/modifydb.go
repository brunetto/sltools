package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BinaryData struct {
	BinaryId string   `bson:"binaryid"`
	Comb     int64    `bson:"comb"`
	N        int64    `bson:"n"`
	Z        float64  `bson:"z"`
	Rv       float64  `bson:"rv"`
	Fpb      float64  `bson:"fpb"`
	W        float64  `bson:"w"`
	Tf       string   `bson:"tf"`
	SysTime  int64    `bson:"systime"`
	PhysTime float64  `bson:"phystime"`
	Objs     string   `bson:"objs"`
	Hflag    string   `bson:"hflag"`
	Types    []string `bson:"types"`
	M0       float64  `bson:"m0"`
	M1       float64  `bson:"m1"`
	Sma      float64  `bson:"sma"`
	Period   float64  `bson:"period"`
	Tgw      float64  `bson:"tgw"`
	Mchirp   float64  `bson:"mchirp"`
	Ecc      float64  `bson:"ecc"`
	X        float64  `bson:"X"`
	Y        float64  `bson:"Y"`
	Z        float64  `bson:"Z"`
	R        float64  `bson:"R"`
	Merge    bool     `bson:"merge"`
}

func main() {

	var (
		err        error
		session    *mgo.Session
		database   *mgo.Database
		cBinaries  *mgo.Collection
		cPositions *mgo.Collection
		cNSmergers *mgo.Collection
		binary     BinaryData
		binaries   []BinaryData
		query      bson.M
		q          *mgo.Query
	)

	log.Println("Open the database")
	if session, err = mgo.Dial("mongodb://localhost:27017/"); err != nil {
		log.Fatal("Error dialing db: ", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	log.Println("Open database")
	database = session.DB("phd")

	log.Println("Open collection")
	cBinaries  = database.C("binaries")
	cPositions = database.C("positions")
	cNSmergers = database.C("ns_mergers")

	log.Println("Retrieve each binary, check for position and NS merger and update the document")

	// 1. Iterate over the binaries that are record of binaries collections
	// 2. For each binary search if it is a NS merger, in case, update
	// 3. Search if it has a position, in case update
	
// 	query = bson.M{"$and": []bson.M{bson.M{"binaryid": "c76n4a103540b3540"}, bson.M{"systime": 13}}}
// 	q = collection.Find(query) //.Sort("phystime")
// 	q.All(&binaries)
// 	var i int
// 	var item BinaryData
// 	for i, item = range binaries {
// 		fmt.Printf("%v: %+v\n", i, item)
// 	}

}
