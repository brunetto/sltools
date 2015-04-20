package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/readfile"
)

const logBin = true

func main() {
	if true {
		defer debug.TimeMe(time.Now())
	}

	var (
		regString string = `\d+,\s+\d+,\s+\S+,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*(-*\d*\.*\d*e*-*\d*),\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*-*\d*\.*\d*e*-*\d*,\s*(-*\d*\.*\d*e*-*\d*)`
		// 		0, 1, --, 0.514488965035739598, 0.194670802686332967, -0.222068317660608999, 0.5932199880910123, 0.204271715528226705, 0.105948237834844639, -0.437641911596323208, , , ,  5.23946321410425484e-05
		regExp  *regexp.Regexp = regexp.MustCompile(regString)
		regRes  []string
		line    string
		inData  = []InData{}
		outData []OutData
		mass, radius float64
		err  error
		nBin int
		bin int = 0
		inFileName, outFileName string
		inFile, outFile *os.File
		nReader *bufio.Reader
		delta float64 = 0.5
		min   float64 = 0
		max   float64 = 0
	)

	if len(os.Args) < 2 {
		log.Fatal("Provide infile and outfile")
	}
	inFileName = os.Args[1]

	outFileName = "profiles-" + inFileName

	fmt.Println("Reading from ", inFileName)
	fmt.Println("Writing to ", outFileName)

	log.Println("Creating files")
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	nReader = bufio.NewReader(inFile)

	if outFile, err = os.Create(outFileName); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for {
		if line, err = readfile.Readln(nReader); err != nil {
			fmt.Println()
			log.Println(err)
			break
		}
		
		if strings.HasPrefix(line, "#") {continue}

		if regRes = regExp.FindStringSubmatch(line); regRes == nil {
			log.Fatalf("Can't reg %v\n", line)
		}

		if radius, err = strconv.ParseFloat(regRes[1], 64); err != nil {
			log.Fatalf("Can't parse radius from %v with err: %v\n", regRes[1], err)
		}
		
		if mass, err = strconv.ParseFloat(regRes[2], 64); err != nil {
			log.Fatalf("Can't parse mass from %v with err: %v\n", regRes[2], err)
		}
		
		if radius > max {
			max = radius
		}
	
		inData = append(inData, InData{Radius: radius, Mass: mass})
	}
	
// 	max = max+max*0.2
	max = 50

	if logBin {
		min = 0.001
		nBin = int(math.Ceil( 2. * math.Pow(float64(len(inData)), 1./3.)))
		delta = (math.Log10(1+max) - math.Log10(1+min)) / float64(nBin)
		fmt.Println(math.Log10(1+max), math.Log10(1+min), nBin)
// 		nBin = int(math.Ceil((math.Log10(max)-math.Log10(min)) / math.Log10(delta)))+1
		
	} else {
		nBin = int(math.Ceil((max-min) / delta))+1
		// Auto bins sucks at small radii
// 		nBin = int(math.Ceil( /*2. **/ math.Pow(float64(len(inData)), 1./3.)))
// 		delta = (max - min) / float64(nBin)
	}
	outData = make([]OutData, nBin+1)
	
	fmt.Println("Min ", min)
	fmt.Println("Max ", max)
	fmt.Println("delta ", delta)
	fmt.Println("nBin ", nBin)
	fmt.Println("Len outData: ", len(outData))
// 	fmt.Println("maxBin ", float64(nBin)*delta)
	
	if len(inData) == 0 {
		log.Fatal("No data retrieved")
	}

	for idx, _ := range inData {
		if logBin {
			bin = int(math.Log10(1+inData[idx].Radius) / delta)
		} else {
			bin = int(inData[idx].Radius / delta)
		}
// 		fmt.Println("idx, bin ", idx, bin)
		outData[bin].MassDiff = outData[bin].MassDiff + inData[idx].Mass
	}

	outData[0].LowerRadius = min
	outData[0].MassCum = outData[0].MassDiff
	// First shell is spherical
	outData[0].Density = (3. * outData[0].MassCum) / (4 * math.Pi * delta * delta * delta)
	outFile.WriteString("# BinN, LowerRadius, MassDiff, MassCum, Density\n")
	str := fmt.Sprintf("%v,%.2f,%v,%v,%v\n", 0, outData[0].LowerRadius,
			outData[0].MassDiff,
			outData[0].MassCum,
			outData[0].Density)
	outFile.WriteString(str)

	

	fmt.Println("Start print")
	
	for idx:=1; idx<len(outData); idx++ {
		if logBin {
			outData[idx].LowerRadius = math.Pow(10, float64(idx) * delta)-1
		} else {
			outData[idx].LowerRadius = float64(idx) * delta
		}
		outData[idx].MassCum = outData[idx-1].MassCum + outData[idx].MassDiff
		// Old way
// 		tmp := math.Pow(outData[idx].LowerRadius+delta, 3) - delta*delta*delta
		// New way, tmp = r2^3 - r1^3
		tmp := math.Pow(outData[idx].LowerRadius, 3) - math.Pow(outData[idx-1].LowerRadius, 3)
		// Average shell density
		outData[idx].Density = (3. * outData[idx].MassDiff) / (4 * math.Pi * tmp)  
		
		str := fmt.Sprintf("%v,%.2f,%v,%v,%v\n", idx, outData[idx].LowerRadius,
			outData[idx].MassDiff,
			outData[idx].MassCum,
			outData[idx].Density)
		outFile.WriteString(str)
	}

}

type InData struct {
	Radius float64
	Mass   float64
	// 	Bin int
}

type OutData struct {
	LowerRadius float64
	MassDiff    float64
	MassCum     float64
	Density     float64
}

