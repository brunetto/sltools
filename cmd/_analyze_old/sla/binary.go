package sla

import (
	"fmt"
	"strconv"
)

// BinaryData store star's data in binary and that changes with time 
// (companion, orbital properties, ...).
type BinaryData struct {
	BinaryId string
	Companion string
	Hardness string
	Types string
	Ecc float64
}

// TimeDomain store the first and last time a star is found in binary
type TimeDomain struct {
	Min uint64
	Max uint64
}

func AddBinaryToExch(regexResult []string, currentIds []string, i int) (*BinaryData, bool) {
	binData := new(BinaryData)
	binData.BinaryId = regexResult[3]
	binData.Companion = currentIds[len(currentIds)-i-1]
	binData.Hardness = regexResult[7]
	binData.Types = regexResult[8]
	// Set DCOB flag if appropriate	
	binData.Ecc, _ = strconv.ParseFloat(regexResult[13], 64)
	var zeroEcc bool
	if binData.Ecc == 0 {
		zeroEcc = true
	} else {
		zeroEcc = false
	}
	return binData, zeroEcc
}

// Print printsfunction for an exchange datum.
func (binaryData *BinaryData) Print() {
	fmt.Println("\t", binaryData.BinaryId, binaryData.Companion, binaryData.Hardness, binaryData.Types, binaryData.Ecc)
}

