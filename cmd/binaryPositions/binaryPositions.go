package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/brunetto/goutils/readfile"
	"github.com/brunetto/goutils/debug"
)

var (
	reg          *regexp.Regexp = regexp.MustCompile(`[rva]\s*=\s*(\S*)\s*(\S*)\s*(\S*)`)
	res          []string
	idxReg       *regexp.Regexp = regexp.MustCompile(`i\s*=\s*(\S*)`)
	idxRes       []string
	nameReg      *regexp.Regexp = regexp.MustCompile(`name\s*=\s*(\S*)`)
	nameRes      []string
	nReg         *regexp.Regexp = regexp.MustCompile(`N\s*=\s*(\S*)`)
	nRes         []string
	combReg      *regexp.Regexp = regexp.MustCompile(`out-comb(\d+)-\S+-run(\d+)-all.txt\S*`)
	combRes      []string
	particleChan chan *Particle
	linesChan    chan string
	dbug         bool = false
)

// TODO: intercetto la size_scale
// then size scale = ceil ( 6.955e5 / (size_scale * 3.08...e13) ) = ceil ( 1 Rsun in km / (size_scale * 1 pc in km) )

func main() {
	if true {defer debug.TimeMe(time.Now())}

	var (
		err                     error
		inFileName, outFileName string
		inFile, outFile         *os.File
		nReader                 *bufio.Reader
		nWriter                 *bufio.Writer
		done                    = make(chan string, 3)
		ext string
		fZip *gzip.Reader
		binaryPrefix string
	)

	if len(os.Args) < 3 {
		log.Fatal("Provide a STDOUT and a outfile")
	}
	inFileName = os.Args[1]
	ext = filepath.Ext(inFileName)
	
	outFileName = os.Args[2]

	if combRes = combReg.FindStringSubmatch(inFileName); combRes == nil {
		log.Fatal("Can't reg: ", inFileName)
	}
	
	binaryPrefix = "c" + combRes[1] + "n" + combRes[2]
	
	fmt.Println("Reading from ", inFileName)
	fmt.Println("Writing to ", outFileName)

	log.Println("Creating files")
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
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
	default:
		{
			log.Println("Unrecognized file type", inFileName)
			log.Fatal("with extension ", ext)
		}
	}

	if outFile, err = os.Create(outFileName); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()

	particleChan = make(chan *Particle, 1000)
	linesChan = make(chan string, 1000)

	log.Println("Launching goroutines")
	go Readline(nReader, linesChan, done)
	go ClusterParsing(linesChan, particleChan, done)
	go ClusterPrinting(nWriter, binaryPrefix, particleChan, done)

	// Wait goroutines to finish
	for i := 0; i < 3; i++ {
		fmt.Print(<-done)
	}

	fmt.Println()
}

type Particle struct {
	Time       int64
	PhysTime   float64
	Idx        string
	Name       string
	X, Y, Z    float64
	Dist       float64
/*	Vx, Vy, Vz string
	Ax, Ay, Az string
	Mass*/       string
// 	Multiple bool
// 	Sons []&Particle
// 	HasParent bool
// 	Parent &Particle{}
}

func Readline(nReader *bufio.Reader, linesChan chan string, done chan string) {
	if dbug {defer debug.TimeMe(time.Now())}
	var (
		err  error
		line string
		msg  string
	)

	for {
		if line, err = readfile.Readln(nReader); err != nil {
			if dbug {
				fmt.Println()
				log.Println(err)
			}
			break
		}
		linesChan <- line
	}
	close(linesChan)
	if dbug {
		msg = "Done from Readline\n"
	} else {
		msg = ""
	}
	done <- msg
}

func timeAndScale (linesChan chan string) (timeStp int64, timeScale, sizeScale float64) {
	
	var (
		err error
		regSysTime = regexp.MustCompile(`system_time\s*=\s*(\d+)`)
		resSysTime []string
		regTimeScale = regexp.MustCompile(`time_scale\s+=\s+(\d\.\d*e*[\+-]*\d*)`)
		resTimeScale []string
		regSizeScale = regexp.MustCompile(`size_scale\s+=\s+(\d\.\d*e*[\+-]*\d*)`)
		resSizeScale []string
	)
	
	timeStp = 0
	timeScale = 0
	sizeScale = 0
	
	for line := range linesChan {
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			if timeStp, err = strconv.ParseInt(resSysTime[1], 10, 64); err != nil {
				log.Fatal("Error parsing timestep: ", err)
			}
		}
		if resTimeScale = regTimeScale.FindStringSubmatch(line); resTimeScale != nil {
			if timeScale, err = strconv.ParseFloat(resTimeScale[1], 64); err != nil {
				log.Fatal("Error parsing time scale: ", err)
			}
		}
		if resSizeScale = regSizeScale.FindStringSubmatch(line); resSizeScale != nil {
			if sizeScale, err = strconv.ParseFloat(resSizeScale[1], 64); err != nil {
				log.Fatal("Error parsing size scale: ", err)
			}
		}
		if float64(timeStp) * timeScale * sizeScale > 0 {
			break
		}
	}
	return timeStp, 1./timeScale, math.Ceil( 6.955e5 / (sizeScale * 3.08567758e13)) 
}

func ClusterParsing(linesChan chan string, particleChan chan *Particle, done chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	var (
		sentParticles int = 0
		filled        bool
		particle      *Particle
		err           error
		msg           string
		nestingLevel int = 0
		timeStp int64 = 0
		sizeScale, timeScale float64
	)

	log.Println("Start reading loop")
	for line := range linesChan {
		// Start reading a particle
		if strings.Contains(line, "(Particle") {
			nestingLevel++
			if nestingLevel == 1 {
				// Root particle
				timeStp, timeScale, sizeScale = timeAndScale(linesChan)
			}
			if nestingLevel == 2 {
				if particle, filled, err = ParseParticle(linesChan, timeStp, timeScale, sizeScale ); err != nil {
					log.Println("Error in filling the particles detected in ClusterParsing ", err)
					break
				}
				if filled == true {
					if dbug {log.Println("Send particle: ", particle.Idx)}
					sentParticles++
					particleChan <- particle
					if sentParticles % 100 == 0 {
						fmt.Printf("\rProcessed %v particles", sentParticles)
					}
				}
			}
		} else if strings.Contains(line, ")Particle") {
			nestingLevel--
		}
	}
	close(particleChan)
	if dbug {
		msg = "Done from ClusterParsing\n"
	} else {
		msg = ""
	}
	done <- msg
}

func ClusterPrinting(nWriter *bufio.Writer, binaryPrefix string, particleChan chan *Particle, done chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	fmt.Fprint(nWriter, "# time, physTime, name, x, y, z, r\n")//, vx, vy, vz, ax, ay, az, mass\n")
	for particle := range particleChan {

		orderedNames := strings.Split(strings.Trim(particle.Name, "()"), ",")
		sort.Strings(orderedNames)

		fmt.Fprintf(nWriter, "%v,%v,%v,%v,%v,%v,%v\n", //%v,%v,%v,%v,%v,%v,%v\n",
			particle.Time,
			particle.PhysTime, 
// 			particle.Name,
			binaryPrefix + "a" + orderedNames[0] + "b" + orderedNames[1],
			particle.X,
			particle.Y,
			particle.Z,
			particle.Dist,
// 			particle.Vx,
// 			particle.Vy,
// 			particle.Vz,
// 			particle.Ax,
// 			particle.Ay,
// 			particle.Az,
// 			particle.Mass,
		)
	}
	var msg string
	if dbug {
		msg = "Done from ClusterPrinting\n"
	} else {
		msg = ""
	}
	done <- msg
}

func ParseParticle(linesChan chan string, timeStp int64, timeScale, sizeScale float64) (*Particle, bool, error) {
	if dbug {defer debug.TimeMe(time.Now())}
	var (
		n, name, idx, /*m,*/ x, y, z/*, vx, vy, vz, ax, ay, az*/ string
		dist                                             float64
// 		line                                             string
		err                                              error
	)
	for {
		// TODO: Fare un check su n=2
		if n, name, idx = ReadPBasics(linesChan); n != "2" {    //> indagare, sto cercando di pescare solo le binarie
			return &Particle{}, false, nil                      //> indagare
		}
		
		ReadPLog(linesChan)
		if /*m,*/ x, y, z, /*vx, vy, vz, ax, ay, az,*/ err = ReadPDynamics(linesChan); err != nil {
			return &Particle{}, false, err
		}
		ReadPHydro(linesChan)
		ReadPStar(linesChan)

		xf, _ := strconv.ParseFloat(x, 64)
		yf, _ := strconv.ParseFloat(y, 64)
		zf, _ := strconv.ParseFloat(z, 64)
		dist = math.Sqrt(xf*xf + yf*yf + zf*zf)

		// Insert a loop here to consume lines until )Particle is reached???
// 		line = <-linesChan
// 		if strings.Contains(line, ")Particle") {
		return &Particle{
			Idx:  idx,
			Name: name,
			Time: timeStp,
			PhysTime: float64(timeStp) * timeScale,
			X:    xf * sizeScale,
			Y:    yf * sizeScale,
			Z:    zf * sizeScale,
			Dist: dist * sizeScale,
// 			Vx:   vx,
// 			Vy:   vy,
// 			Vz:   vz,
// 			Ax:   ax,
// 			Ay:   ay,
// 			Az:   az,
// 			Mass: m,
		}, true, nil
// 		}
	}
}

func ReadPBasics(linesChan chan string) (n, name, idx string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	name = "--" // may be the name is empty
	for line := range linesChan {
		if idxRes = idxReg.FindStringSubmatch(line); idxRes != nil {
			idx = idxRes[1]
		}
		if nameRes = nameReg.FindStringSubmatch(line); nameRes != nil {
			name = nameRes[1]
		}
		if nRes = nReg.FindStringSubmatch(line); nRes != nil {
			n = nRes[1]
		}
		if strings.Contains(line, "(Log") {
			return n, name, idx
		}
	}
	return "-1", "", "" // error while reading
}

func ReadPLog(linesChan chan string) {
	if dbug {defer debug.TimeMe(time.Now())}
	for line := range linesChan {
		if strings.Contains(line, ")Log") {
			break
		}
	}
}

func ReadPDynamics(linesChan chan string) (/*m, */x, y, z/*, vx, vy, vz, ax, ay, az*/ string, err error) {
	if dbug {defer debug.TimeMe(time.Now())}
	for line := range linesChan {
// 		if strings.Contains(line, "m  =") {
// 			m = strings.Split(line, "  = ")[1]
// 		}
		if strings.Contains(line, "r  =") {
			if res = reg.FindStringSubmatch(line); res != nil {
				x = res[1]
				y = res[2]
				z = res[3]
			}
		}
// 		if strings.Contains(line, "v  =") {
// 			if res = reg.FindStringSubmatch(line); res != nil {
// 				vx = res[1]
// 				vy = res[2]
// 				vz = res[3]
// 			}
// 		}
// 		if strings.Contains(line, "a  =") {
// 			if res = reg.FindStringSubmatch(line); res != nil {
// 				ax = res[1]
// 				ay = res[2]
// 				az = res[3]
// 			}
// 		}
		if strings.Contains(line, "(Hydro") {
			return /*m,*/ x, y, z, /*vx, vy, vz, ax, ay, az,*/ nil
		}
	}
	err = errors.New("Stop reading")
	return "", "", "", /*"", "", "", "", "", "", "",*/ err
}

func ReadPHydro(linesChan chan string) {
	if dbug {defer debug.TimeMe(time.Now())}
	for line := range linesChan {
		if strings.Contains(line, ")Hydro") {
			break
		}
	}
}

func ReadPStar(linesChan chan string) {
	if dbug {defer debug.TimeMe(time.Now())}
	for line := range linesChan {
		if strings.Contains(line, ")Star") {
			break
		}
	}
}
