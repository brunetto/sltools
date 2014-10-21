package sla

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

// ExecOnAll execute the starData function f on all the stars
func (starMap StarMapType) ExecOnAll(fName string/*func()*/) () {
	tGlob0 := time.Now()
	// Compute f
	log.Print("Running",fName , " for all the stars")
	keys := starMap.Keys()
	
	// Check if fName is valid on the first element only, the other are the same
	var v reflect.Value
	value := starMap[keys[0]]	
	if v = reflect.ValueOf(value).MethodByName(fName); v.IsValid() == false {
			log.Println("No valid function name in ", Whoami(false), ": ", fName)
			log.Fatal("Try with: CountExchanges, ComputeLifeTimes, Print, PrintWithExchs")
		}
	v.Call([]reflect.Value{})
	
	nStar := 2
	nStars := len(starMap)
	for _, key := range keys[1:] {
		value := starMap[key]
		fmt.Fprintf(os.Stderr, "\rDone: %v %%", (100 * nStar) / nStars)
		nStar++
		reflect.ValueOf(value).MethodByName(fName).Call([]reflect.Value{})
	}
	
	fmt.Fprint(os.Stderr, "\n")
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for ", Whoami(false), "(", fName, ") :", tGlob1.Sub(tGlob0))
}

// Old function, to backup!!
/*
// CountExchanges counts all the star's exchanges
func (starMap StarMapType) CountExchanges() {
	tGlob0 := time.Now()
	// Compute exchanges
	log.Print("Counting exchanges for all the stars")
	nStar := 1
	nStars := len(starMap)
	keys := starMap.Keys()
	for _, key := range keys {
		value := starMap[key]
		fmt.Fprintf(os.Stderr, "\rDone: %v %%", (100 * nStar) / nStars)
		nStar++
		value.CountExchanges()
	}
	fmt.Fprint(os.Stderr, "\n")
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for ", Whoami(false), " ", tGlob1.Sub(tGlob0))
}
*/
/*
// ComputeLifeTimes counts the lifetimes of all the stars between the exchanges
func (starMap StarMapType) ComputeLifeTimes() {
	tGlob0 := time.Now()
	// Compute lifetimes
	log.Print("Computing lifetimes for all the stars")
	nStar := 1
	nStars := len(starMap)
	keys := starMap.Keys()
	for _, key := range keys {
		value := starMap[key]
		fmt.Fprintf(os.Stderr, "\rDone: %v %%", (100 * nStar) / nStars)
		nStar++
		value.ComputeLifeTimes()
	}
	fmt.Fprint(os.Stderr, "\n")
	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for ", Whoami(false), " ", tGlob1.Sub(tGlob0))
}
*/

/*
  Print prints all the stars' data
func (starMap StarMapType) Print() {
	keys := starMap.Keys()
	for _, key := range keys {
// 		fmt.Print(key, " ")
		starMap[key].Print()
		fmt.Println("===========================================================")
	}
}
 */

