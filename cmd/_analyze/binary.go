package slan

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	
	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
)

type BinaryMapType map[string]*BinaryData

// BinaryData store star's data in binary and that changes with time
// (companion, orbital properties, ...).
type BinaryData struct {
	BinaryId       string
	Z              string
	NFile          string
	Comb           string
	TimeUnit string
	TimeProperties map[uint64]*BinaryChangingProperties
}

// As a function of time
type BinaryChangingProperties struct {
	Hardness string
	Types    string
	Ecc      float64
	Sma      float64
	Period   float64
	Masses   [2]float64
	ChirpMass float64
	TGW float64
}

// TimeDomain store the first and last time a star is found in binary
type TimeDomain struct {
	Min uint64
	Max uint64
}

func (binData BinaryMapType) AddBinary(regexResult []string) (string, bool) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	const PERIOD_UNIT float64 = 1000000
	const DISTANCE_UNIT float64 = 206264.806 // 1 Parsec = 206 264.806 Astronomical Units
	const PC2M = 3.0857e16
	const PC2AU float64 = 206264.806
	const SECONDS_IN_A_YEAR float64 = 60*60*24*365
	const LIGHT_SPEED float64 = 299792458 // m/s
	const G float64 = 6.67398e11 // m^3 kg^-1 s^-2
	const M_SUN float64 = 1.98855e30
	const YR2GYR float64 = 1000000000
	
	var CONSTANT float64 = (5. * (math.Pow(LIGHT_SPEED, 5)) * (math.Pow(PC2M, 4))) / (256 * (math.Pow(G, 3)) * (SECONDS_IN_A_YEAR * YR2GYR) * (math.Pow(M_SUN, 3))) 
	
	
	var (
		binaryId string
		currentTime uint64
		exists bool
		err error
		mTot, mu float64
		ecc, sma, period, mass0, mass1 float64
	)

	binaryId = regexResult[3]
	currentTime, _ = strconv.ParseUint(regexResult[4], 10, 64)

	if _, exists := binData[binaryId]; !exists {
		binData[binaryId] = &BinaryData{
			BinaryId: binaryId,
			Z:        regexResult[1],
			NFile:    goutils.LeftPad(regexResult[2], "0", 3),
			// 			Comb:
		}
		binData[binaryId].TimeProperties = make(map[uint64]*BinaryChangingProperties)
	}
	// Does this binary in this timestep already exist?
	if _, exists = binData[binaryId].TimeProperties[currentTime]; exists {
		fmt.Println()
		log.Printf("WARNING: Binary %v already exists at %v\n overwriting it!", binaryId, currentTime)
	}

	binData[binaryId].TimeProperties[currentTime] = &BinaryChangingProperties{
		Hardness: regexResult[7],
		Types:    regexResult[8],
	}
	
	
	if ecc, err = strconv.ParseFloat(regexResult[13], 64); err != nil {
		log.Fatal("Error parsing float: ", err)
	}
	if sma, err = strconv.ParseFloat(regexResult[11], 64); err != nil {
		log.Fatal("Error parsing float: ", err)
	}
	if period, err = strconv.ParseFloat(regexResult[12], 64); err != nil {
		log.Fatal("Error parsing float: ", err)
	}
	if mass0, err = strconv.ParseFloat(regexResult[9], 64); err != nil {
		log.Fatal("Error parsing float: ", err)
	}
	if mass1, err = strconv.ParseFloat(regexResult[10], 64); err != nil {
		log.Fatal("Error parsing float: ", err)
	}
	
	
	binData[binaryId].TimeProperties[currentTime].Ecc = ecc
	binData[binaryId].TimeProperties[currentTime].Sma = sma
	binData[binaryId].TimeProperties[currentTime].Period = period
	binData[binaryId].TimeProperties[currentTime].Masses[0] = mass0
	binData[binaryId].TimeProperties[currentTime].Masses[1] = mass1
	
	mTot = mass0 + mass1
	mu = (mass0 * mass1) / mTot
	
	binData[binaryId].TimeProperties[currentTime].ChirpMass = math.Pow(mu, 3./5) * math.Pow(mTot, 2./5)
	binData[binaryId].TimeProperties[currentTime].TGW = CONSTANT * math.Pow(sma, 4) * math.Pow((1 - math.Pow(ecc,2)), (7./2)) / (mass0 * mass1 * (mass0 + mass1))

	var zeroEcc bool
	if binData[binaryId].TimeProperties[currentTime].Ecc == 0 {
		zeroEcc = true
	} else {
		zeroEcc = false
	}
	return binaryId, zeroEcc
}
/*
// Print function for an exchange datum.
func (binaryData *BinaryData) Print() {
	fmt.Println("\t", binaryData.BinaryId, binaryData.Companion,
		binaryData.Hardness, binaryData.Types, binaryData.Ecc)
}*/
