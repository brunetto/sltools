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
	particleChan chan *Particle
	linesChan    chan string
	dbug         bool = false
)

func main() {
	if true {
		defer debug.TimeMe(time.Now())
	}

	var (
		err                     error
		inFileName, outFileName string
		inFile, outFile         *os.File
		nReader                 *bufio.Reader
		nWriter                 *bufio.Writer
		done                    = make(chan string, 3)
		ext                     string
		fZip                    *gzip.Reader
	)

	if len(os.Args) < 3 {
		log.Fatal("Provide infile and outfile")
	}
	inFileName = os.Args[1]
	ext = filepath.Ext(inFileName)

	outFileName = os.Args[2]

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
	go ClusterPrinting(nWriter, particleChan, done)

	// Wait goroutines to finish
	for i := 0; i < 3; i++ {
		fmt.Print(<-done)
	}

	fmt.Println()
}

type Cluster struct {
	NStars int
	OriginalLog string
	
}

type Particle struct {
	Time       string
	Idx        string
	Name       string
	X, Y, Z    string
	Dist       float64
	Vx, Vy, Vz string
	Ax, Ay, Az string
	Mass       string
	Type string
	Multiple   bool
	Sons       []*Particle
	// 	HasParent bool
	// 	Parent &Particle{}
}

func (p *Particle) PrintSL () () {
	
	
	
}


func (p *Particle) PrintAscii () () {
	
	
	
}


func (p *Particle) SaveHDF5 () () {
	
	
	
}


func Readline(nReader *bufio.Reader, linesChan chan string, done chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
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
		nestingLevel  int    = 0
		timeStp       string = ""
	)

	log.Println("Start reading loop")
	for line := range linesChan {
		// Start reading a particle
		if strings.Contains(line, "(Particle") {
			nestingLevel++
			if nestingLevel == 1 {
				// Root particle
				var (
					m, x, y, z, vx, vy, vz, ax, ay, az string
					err error
				)
				
				if timeStp, m, x, y, z, vx, vy, vz, ax, ay, az, err = ReadPDynamics(linesChan); err != nil {
					return &Particle{}, false, err
				}
			}
			if nestingLevel == 2 {
				if particle, filled, err = ParseParticle(linesChan, timeStp); err != nil {
					log.Println("Error in filling the particles detected in ClusterParsing ", err)
					break
				}
				if filled == true {
					if dbug {
						log.Println("Send particle: ", particle.Idx)
					}
					sentParticles++
					particleChan <- particle
					if sentParticles%100 == 0 {
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

func ClusterPrinting(nWriter *bufio.Writer, particleChan chan *Particle, done chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	fmt.Fprint(nWriter, "# time, idx, name, x, y, z, r, vx, vy, vz, ax, ay, az, mass\n")
	for particle := range particleChan {
		fmt.Fprintf(nWriter, "%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n",
			particle.Time,
			particle.Idx,
			particle.Name,
			particle.X,
			particle.Y,
			particle.Z,
			particle.Dist,
			particle.Vx,
			particle.Vy,
			particle.Vz,
			particle.Ax,
			particle.Ay,
			particle.Az,
			particle.Mass,
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

func ParseParticle(linesChan chan string, timeStp string) (*Particle, bool, error) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	var (
		n, name, idx, m, x, y, z, vx, vy, vz, ax, ay, az string
		dist                                             float64
		// 		line                                             string
		err error
	)
	for {
		if n, name, idx = ReadPBasics(linesChan); n != "1" {
			// 			return &Particle{}, false, nil
		}

		ReadPLog(linesChan)
		if _, m, x, y, z, vx, vy, vz, ax, ay, az, err = ReadPDynamics(linesChan); err != nil {
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
			X:    x,
			Y:    y,
			Z:    z,
			Dist: dist,
			Vx:   vx,
			Vy:   vy,
			Vz:   vz,
			Ax:   ax,
			Ay:   ay,
			Az:   az,
			Mass: m,
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
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	for line := range linesChan {
		if strings.Contains(line, ")Log") {
			break
		}
	}
}

func ReadPDynamics(linesChan chan string) (time, m, x, y, z, vx, vy, vz, ax, ay, az string, err error) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	time = "--" // Only present in root node
	
	var (
		regSysTime = regexp.MustCompile(`system_time\s*=\s*(\d+)`)
		resSysTime []string
	)

		
	return timeStp
	
	for line := range linesChan {
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			time = resSysTime[1]
		}
		if strings.Contains(line, "m  =") {
			m = strings.Split(line, "  = ")[1]
		}
		if strings.Contains(line, "r  =") {
			if res = reg.FindStringSubmatch(line); res != nil {
				x = res[1]
				y = res[2]
				z = res[3]
			}
		}
		if strings.Contains(line, "v  =") {
			if res = reg.FindStringSubmatch(line); res != nil {
				vx = res[1]
				vy = res[2]
				vz = res[3]
			}
		}
		if strings.Contains(line, "a  =") {
			if res = reg.FindStringSubmatch(line); res != nil {
				ax = res[1]
				ay = res[2]
				az = res[3]
			}
		}
		if strings.Contains(line, "(Hydro") {
			return time, m, x, y, z, vx, vy, vz, ax, ay, az, nil
		}
	}
	err = errors.New("Stop reading")
	return "", "", "", "", "", "", "", "", "", "", err
}

func ReadPHydro(linesChan chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	for line := range linesChan {
		if strings.Contains(line, ")Hydro") {
			break
		}
	}
}

func ReadPStar(linesChan chan string) {
	if dbug {
		defer debug.TimeMe(time.Now())
	}
	for line := range linesChan {
		if strings.Contains(line, ")Star") {
			break
		}
	}
}
