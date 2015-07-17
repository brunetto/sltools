package slt

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func CheckSnapshot(inFileName string) {
	var (
		err     error
		inFile  *os.File
		nReader *bufio.Reader
		fZip    *gzip.Reader
		ext     string
		regRes                         map[string]string
		lengthUnit, length, nStars, fPB, ncm       float64
		simulationStop                 int64 
	)
	
	// 	log.Println("Checking ", inFileName)
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	
	// Extract fileNameBody, round and ext
	if regRes, err = DeepReg(inFileName); err == nil {
		
		ext = regRes["ext"]
		
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
	
	} else if regRes, err = Reg(inFileName); err == nil {	
		
		log.Println("Can't derive deep info from name, only take the names and go standard for the rest")
		
		ext = regRes["ext"]
		
		// Default
		length = float64(1)
		nStars = float64(5500)
		
		fmt.Printf("Set default parameters for cluster: \n")
		fmt.Printf("Radius: %v\n", length)
		fmt.Printf("Number of stars: %v\n", nStars)
	} else {
		log.Println("Can't derive standard names or deep info from STDOUT => wrap it!!")
		ext = ".txt"
		
		// Default
		length = float64(1)
		nStars = float64(5500)
		
		fmt.Printf("Set default parameters for cluster: \n")
		fmt.Printf("Radius: %v\n", length)
		fmt.Printf("Number of stars: %v\n", nStars)
	}
	
	lengthUnit = math.Sqrt(0.25*0.25*math.Pow(length, 3)*(5500./nStars))
	simulationStop = 1 + int64(math.Floor(110./lengthUnit))
	fmt.Printf("\tApprox length unit: %2.2f || simulationStop: %v\n", lengthUnit, simulationStop)
	
	switch ext {
	case ".txt":
		{
			nReader = bufio.NewReader(inFile)
		}
	case ".gz":
		{
			fZip, err = gzip.NewReader(inFile)
			if err != nil {
				log.Fatal("Can't open %s: error: %s\n", inFile, err)
			}
			nReader = bufio.NewReader(fZip)
		}
	case ".txt.gz":
	{
		fZip, err = gzip.NewReader(inFile)
		if err != nil {
			log.Fatal("Can't open %s: error: %s\n", inFile, err)
		}
		nReader = bufio.NewReader(fZip)
	}
	default:
		{
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extension ", ext)
		}
	}

	for {
		if _, err = ReadOutSnapshot(nReader, true); err != nil {
			break
		}
	}
}
