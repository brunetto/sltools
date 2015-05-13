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
	Merge    bool     `bson:"merge"`
}

func main() {

	var (
		err        error
		session    *mgo.Session
		database   *mgo.Database
		collection *mgo.Collection
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
	collection = database.C("binaries")

	log.Println("Search for one binary")
	binary = BinaryData{}
	if err = collection.Find(bson.M{"binaryid": "c76n4a103540b3540"}).One(&binary); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("One binary: %+v\n", binary)

	log.Println("Search for merger candidates")
	query = bson.M{"$and": []bson.M{bson.M{"binaryid": "c76n4a103540b3540"}, bson.M{"systime": 13}}}

	q = collection.Find(query) //.Sort("phystime")

	q.All(&binaries)

	var i int
	var item BinaryData

	for i, item = range binaries {
		fmt.Printf("%v: %+v\n", i, item)
	}

}
