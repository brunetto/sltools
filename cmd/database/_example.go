// From https://docs.google.com/document/d/1GloUFQ1JK0SYRWIukBZyY3iSs6pbsUp4dHAs-Nj_xfU/pub
// and http://www.reddit.com/r/golang/comments/2mqhrz/mongodb_mgo_example/
package main
import (
        "fmt"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)
type Name struct {
        First string `bson:"first"`
        Last  string `bson:"last"`
}
type Worker struct {
        Id         bson.ObjectId `bson:"_id"`
        WorkerName Name          `bson:"worker_name"`
        Counts     []int         `bson:"counts"`
        Job        string        `bson:"job,omitempty"`
}
type M map[string]interface{}
func main() {
        dbSession, err := mgo.Dial("localhost")
        if err != nil {
                panic(err)
        }
        defer dbSession.Close()
        db := dbSession.DB(“test”)
        collectionName := “workers”
        dbWorkers := db.C(collectionName)
// --- delete all records in workers collection -------------
        changeInfo, _ := dbWorkers.RemoveAll(nil)
        fmt.Println("removed ", changeInfo.Removed)
// --- add 2 worker records ----------------------
        var newWorkers [2]Worker
        recid := bson.NewObjectId()
        newWorkers[0] = Worker{recid, Name{"Alf", "Smith"}, []int{55, 28, 33}, "welder"}
        recid = bson.NewObjectId()
        newWorkers[1] = Worker{recid, Name{"Tuf", "Taylor"}, []int{77, 58, 49, 60}, ""}
        err = dbWorkers.Insert(newWorkers[0], newWorkers[1])
        if err != nil {
                panic(err)
        }
        
// ---  query, where first name = Tuf, limit result to 1 record  ------------
// ---  update returned record, change value of job to “boss”  --------------
           result := new(Worker)
        var findParm, setParm M
        findParm = M{"worker_name.first": "Tuf"}
        qry := dbWorkers.Find(findParm).Limit(1)
        if cnt, _ := qry.Count(); cnt == 0 {
                fmt.Println("Tuf is missing")
        } else {
                qry.One(result)        
                findParm = M{"_id": result.Id}
                setParm = M{"$set": M{"job": "boss"}}  // don’t use $set if replacing doc
                err = dbWorkers.Update(findParm, setParm)
        }
// --- query, get all worker recs, sort by last name ------------------
        var total int
        iter := dbWorkers.Find(nil).Sort("worker_name.last").Iter()
        for iter.Next(result) {
                total = 0
                for _, v := range result.Counts {
                        total += v
                }
                fmt.Println(result.WorkerName, total)
        }
// ---  query, get worker recs where job field exists ------------------
        findParm = M{"job": M{"$exists": true}}
        iter = dbWorkers.Find(findParm).Iter()
        for iter.Next(result) {
                fmt.Println(result.WorkerName.Last, result.Job)
        }
