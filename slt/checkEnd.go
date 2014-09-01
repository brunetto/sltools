package slt

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	
	"github.com/brunetto/goutils/debug"
)

func CheckEnd (inFileName string, endOfSimMyr float64) () {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	var (
		err error
		regRes map[string]string
		length, fPB, ncm, nStars, timeUnit float64
		simulationStop int64
	)
	
	// Extract fileNameBody, round and ext
	if regRes, err = DeepReg(inFileName); err == nil {
		if regRes["prefix"] != "out" {
			log.Fatalf("Please specify a STDOUT file, found %v prefix", regRes["prefix"])
		}
		
		if length, err = strconv.ParseFloat(regRes["Rv"], 64); err != nil {
			log.Fatalf("Can't convert cluster length %v to float64: %v\n", regRes["Rv"], err)
		}
		if fPB, err = strconv.ParseFloat(regRes["fPB"][:1]+"."+regRes["fPB"][1:], 64); err != nil {
			log.Fatalf("Can't convert cluster fPB %v to float64: %v\n", regRes["fPB"], err)
		}
		if ncm, err = strconv.ParseFloat(regRes["NCM"], 64); err != nil {
			log.Fatalf("Can't convert cluster NCM %v to float64: %v\n", regRes["NCM"], err)
		}
		nStars = ncm * (1 + fPB)
	} else {
		log.Fatal("Can't derive standard names or deep info from STDOUT")
	}
	
	// Now simulation stop is calculated scaling that of the clusters of the first simlations
	// with the formula
	//
	// maxTimeStep2 = ( maxTime / sqrt(timeUnit1**2 * (length2/length1)**3 * (m1 / m2)))
	//
	// with
	//
	// maxTime ~ 100 Myr
	// timeUnit1 ~ 0.25 Myr
	// length1 = 1 pc
	// m1 / m2 approximated with the number of stars, so m2 = NCM * (1 + fPB) and m1 = 5500
	// FIXME maybe leng should be timeUnit or something similar????
	
	timeUnit = math.Sqrt(0.25*0.25*math.Pow(length, 3)*(5500./nStars))
	simulationStop = 1 + int64(math.Floor(endOfSimMyr/timeUnit))
	fmt.Printf("\tApprox time unit: %2.2f || simulationStop: %v\n", timeUnit, simulationStop)
}

