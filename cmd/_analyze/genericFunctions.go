package slan

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/brunetto/goutils/debug"
)

// ExecOnAll execute the starData function f on all the stars
func (starMap StarMapType) ExecOnAll(fName string /*func()*/) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	// Compute f
	log.Print("Running", fName, " for all the stars")
	keys := starMap.Keys()

	// Check if fName is valid on the first element only, the other are the same
	var v reflect.Value
	value := starMap[keys[0]]
	if v = reflect.ValueOf(value).MethodByName(fName); v.IsValid() == false {
		log.Println("No valid function name in ", debug.FName(false), ": ", fName)
		log.Fatal("Try with: CountExchanges, ComputeLifeTimes, Print, PrintWithExchs")
	}
	v.Call([]reflect.Value{})

	nStar := 2
	nStars := len(starMap)
	for _, key := range keys[1:] {
		value := starMap[key]
		fmt.Fprintf(os.Stderr, "\rDone: %v %%", (100*nStar)/nStars)
		nStar++
		reflect.ValueOf(value).MethodByName(fName).Call([]reflect.Value{})
	}

	fmt.Fprint(os.Stderr, "\n")
	fmt.Println()
}
