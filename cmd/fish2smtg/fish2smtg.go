package main

import (
	"bufio"
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
	
// 	hdf5 "github.com/sbinet/go-hdf5"
	
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/readfile"
)

func main () () {
	defer debug.TimeMe(time.Now())
	const comment = "#"
	var (
		err error
		inPath, inFileName, outFileName string
		dataSlice *DataSlice
		dirs []string
		dir string
		nLines int
		outFile *os.File
	)
	
	if len(os.Args) < 3 || strings.Contains(os.Args[1], "help") {
		log.Fatal(`Please provide the parent folder for the simulations 
		and the name for the outfile.
		The code assume a folder tree like:
		
		parent folder
			|- comb16-TFno-Rv1-NCM10000-fPB005-W5-Z010
					|- Analysis/03-final/all_the_fishes.txt
			|- comb17-TFno-Rv1-NCM10000-fPB010-W5-Z010
			|- comb19-TFno-Rv1-NCM10000-fPB005-W9-Z010
			|- comb20-TFno-Rv1-NCM10000-fPB010-W9-Z010
			|- comb28-TFno-Rv1-NCM10000-fPB005-W5-Z100
			|- comb29-TFno-Rv1-NCM10000-fPB010-W5-Z100
			|- ...`)
	}
	inPath = os.Args[1]
	outFileName = os.Args[2]

	if outFile, err = os.Create(outFileName); err != nil {
		log.Fatalf("Can't open %v with err: %v\n", outFileName, err)
	}
	
	_, err = outFile.WriteString("# Binary_ids, comb, N, Z, Rv, Fpb, W0, Tf, " + 
									"SysTime, PhysTime, Objects_ids, Hardflag, " +
									"Types, Mass_0, Mass_1, Sma, Period, Ecc, Tgw, Mchirp\n")
	if err != nil {
		log.Fatalf("Can't write to %v with error %v\n", outFileName, err)
	}
	outFile.Close()
	
	log.Println("Scanning for folders")
	if dirs, err = filepath.Glob(filepath.Join(inPath, "comb*")); err != nil {
		log.Fatal("Can't glob with error: ", err)
	}
	
	sort.Strings(dirs)
	
	fmt.Println("Found:\n")
	for _, dir = range dirs {
		fmt.Println(dir)
	}
	
	log.Println("Loop on folders")
	for _, dir = range dirs {
		
		fmt.Println("Work on ", dir)
		
		inFileName = filepath.Join(dir, "Analysis", "03-final", "all_the_fishes.txt")
	
		dataSlice = &DataSlice{}
		
		nLines = dataSlice.CollectData(inFileName)
		
		fmt.Printf("Read %v lines\n", nLines)
		
		fmt.Println("Save data to ", outFileName)
		if err = dataSlice.ToCsv(outFileName); err != nil {
			log.Fatalf("Can't save data to %v with err %v:\n", outFileName, err)
		}
		
		log.Printf("Wrote %v lines on %v\n", nLines, outFileName)
	}
}

type DataSlice []*Data

// DCOB data from all_the_fishes files + 
// data calculated or extrapolated from the folders
type Data struct {
	// ID of the binary (Comb+N+ids of the objects)
	Binary_ids string
	// Parameter combination, identifies the simulation
	Comb string
	// Random realization (simulation run/file) number
	N int64
	// Metallicity
	Z float64
	// ICs virial radius of the King model
	Rv float64
	// IC primordial binary fraction
	Fpb float64
	// IC King central adimensional potential
	W0 int64
	// External tidal field string identifier
	Tf string
	// System (=N-body) time-step
	SysTime int64
	// Physical time in [Myr]
	PhysTime float64
	// IDs of the objects in the form ID1|ID2
	Objects_ids string
	// Distinguish between H(ard) and S(oft) binaries accordingly to StarLab
	Hardflag string
	// Objects type: 
	// -- = non compact
	// wd = white dwarf
	// ns++ = will be neutron star 
	// bh++ = will be black hole
	// ns = neutron star
	// bh = black hole
	Types string
	// Mass of the first object in Objects 
	// (not necessary the primary component)
	Mass_0 float64
	// Mass of the second object
	Mass_1 float64
	// Semi-major axis in [pc]
	Sma float64
	// Period in [Myr]
	Period float64
	// Ecc [0:1]
	Ecc float64
	// Peters' coalescence timescale 
	// see http://journals.aps.org/pr/pdf/10.1103/PhysRev.136.B1224, Eq. 5.9
	Tgw float64
	// Chirp mass
	Mchirp float64
}

func (d *Data) Fill (lineRes []string, addData map[string]string) (err error) {
	const PERIOD_UNIT float64 = 1e6
	const DISTANCE_UNIT float64 = 206264.806 // 1 Parsec = 206 264.806 Astronomical Units
	const PC2M = 3.0857e16
	const PC2AU float64 = 206264.806
	const SECONDS_IN_A_YEAR float64 = 60*60*24*365
	const LIGHT_SPEED float64 = 299792458 // m/s
	const G float64 = 6.67398e-11 // m^3 kg^-1 s^-2
	const M_SUN float64 = 1.98855e30
	const YR2GYR float64 = 1e9
	
	var CONSTANT float64 = (5. * (math.Pow(LIGHT_SPEED, 5)) * (math.Pow(PC2M, 4))) / (256 * (math.Pow(G, 3)) * (SECONDS_IN_A_YEAR * YR2GYR) * (math.Pow(M_SUN, 3))) 
	
	if d.Z, err = parseZ(lineRes[1]); err != nil {
		return fmt.Errorf("Can't parse into float Z %v\n", lineRes[1])
	}
    
    if d.N, err = strconv.ParseInt(lineRes[2], 10, 0); err != nil {
		return fmt.Errorf("Can't parse %v into Data.N\n", lineRes[2])
	}
    
    if d.SysTime, err = strconv.ParseInt(lineRes[4], 10, 0); err != nil {
		return fmt.Errorf("Can't parse %v into Data.SysTime\n", lineRes[4])
	}
    
    if d.PhysTime, err = strconv.ParseFloat(lineRes[5], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.PhysTime\n", lineRes[5])
	}

	d.Objects_ids = lineRes[6] 
    
    d.Hardflag = lineRes[7] 
    
    d.Types = lineRes[8] 
    
    if d.Mass_0, err = strconv.ParseFloat(lineRes[9], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Mass_0\n", lineRes[9])
	}
    
    if d.Mass_1, err = strconv.ParseFloat(lineRes[10], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Mass_1\n", lineRes[10])
	}
    
    if d.Sma, err = strconv.ParseFloat(lineRes[11], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Sma\n", lineRes[11])
	}
    
    if d.Period, err = strconv.ParseFloat(lineRes[12], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Period\n", lineRes[12])
	}
    
    if d.Ecc, err = strconv.ParseFloat(lineRes[13], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Ecc\n", lineRes[13])
	}
		
	d.Comb = addData["comb"]
	
	if d.Rv, err = strconv.ParseFloat(addData["rv"], 64); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Rv\n", addData["rv"])
	}

	if d.Fpb, err = parseFpb(addData["fpb"]); err != nil {
		return fmt.Errorf("Can't parse %v into Data.Fpb\n", addData["fpb"])
	}
	
	if d.W0, err = strconv.ParseInt(addData["W0"], 10, 0); err != nil {
		return fmt.Errorf("Can't parse %v into Data.W0\n", lineRes[4])
	}
	
	d.Tf = addData["tf"]
	
    d.Binary_ids = "c" + d.Comb + "n" + strconv.FormatInt(d.N, 10) + strings.Split(lineRes[3], "ids")[1]
	
	d.Tgw = CONSTANT * math.Pow(d.Sma, 4) * math.Pow((1 - math.Pow(d.Ecc,2)), (7./2)) / (d.Mass_0 * d.Mass_1 * (d.Mass_0 + d.Mass_1))
	
	d.Mchirp = ChirpMass(d.Mass_0, d.Mass_0)
		
	return nil
}

// Save data to csv file
func (ds *DataSlice) ToCsv (outFileName string) (err error) {
	var (
		outFile *os.File
		nWriter *bufio.Writer
		outputLine string
	)
	
	// Write data to file
	if !goutils.Exists(outFileName) {
		log.Fatal("outfile already exists")
		if outFile, err = os.Create(outFileName); err != nil {
			log.Fatal(err)
		}
	} else {
		// os.O_RDWR needed to not have an error using if in this way (don't know why)
		if outFile, err = os.OpenFile(outFileName, os.O_APPEND|os.O_RDWR, 0666); err != nil { 
			log.Fatal(err)
		}
	}
	defer outFile.Close()
	 
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()
	
	fmt.Println("Start writing lines: ", len(*ds))
	for _, record := range *ds {
		outputLine = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
										record.Binary_ids, 
										record.Comb, 
										record.N, 
										record.Z, 
										record.Rv, 
										record.Fpb, 
										record.W0, 
										record.Tf, 
										record.SysTime, 
										record.PhysTime, 
										record.Objects_ids, 
										record.Hardflag, 
										record.Types, 
										record.Mass_0, 
										record.Mass_1, 
										record.Sma, 
										record.Period, 
										record.Ecc, 
										record.Tgw,
										record.Mchirp,
								)
		_, err = nWriter.WriteString(outputLine)
		if err != nil {
			log.Fatalf("Can't write %v to %v with error: %v\n", outputLine, outFileName, err)
		}
	}
	return nil
}

// Recreate the float version of Z
func parseZ (z string) (z1 float64, err error) {
	if z1, err = strconv.ParseFloat(z[:1] + "." + z[1:], 64); err != nil {
		return z1, err
	}
	return z1, err
}

// Recreate the float version of Z
func parseFpb (fpb string) (fpb1 float64, err error) {
	if fpb1, err = strconv.ParseFloat(fpb[:1] + "." + fpb[1:], 64); err != nil {
		return fpb1, err
	}
	return fpb1, err
}

// Collect data from the file
func  (ds *DataSlice) CollectData (inFileName string) (nLines int) {
	var (
		inFile *os.File
		nReader *bufio.Reader
		err error
		line string
		dirRegString string = `\S*comb(\d+)` +					// Group(1): comb
								`-TF(\S*)` +					// Group(2): tf
								`-Rv(\d+)` +					// Group(3): Rv
								`-NCM(\d+)` +					// Group(4): NCM, not used
								`-fPB(\d+)` +					// Group(5): fpb
								`-W(\d+)` +						// Group(6): W0
								`-Z(\d+)`						// Group(7): Z
		dirReg *regexp.Regexp = regexp.MustCompile(dirRegString)
		dirRes []string
		
		lineRegString string = `^(\d+)\s+` +					// Group(1): Z
								`(\d+)\s+` +					// Group(2): N
								`(\S+)\s+` +					// Group(3): Binary_ids
								`(\d+)\s+` +					// Group(4): SysTime
								`(\d+\.*\d*e*-*\d*)\s+` +		// Group(5): PhysTime
								`(\S+\|\S+)\s+` +				// Group(6): Objects_ids
								`([H-S])\s+` +					// Group(7): Hardflag
								`(\S+)\s+` +					// Group(8): Types
								`(\d+\.*\d*e*-*\d*)\s+` +		// Group(9): Mass_0
								`(\d+\.*\d*e*-*\d*)\s+` +		// Group(10): Mass_1
								`(\d+\.*\d*e*-*\d*)\s+` +		// Group(11): Sma
								`(\d+\.*\d*e*-*\d*)\s+` +		// Group(12): Period
								`(\d+\.*\d*e*-*\d*)`			// Group(13): Ecc
		lineReg *regexp.Regexp = regexp.MustCompile(lineRegString)
		lineRes []string
		newData *Data
	)
	
	if dirRes = dirReg.FindStringSubmatch(inFileName); dirRes == nil {
		log.Fatalf("Can't reg comb, rv, ... in %v\n: ", inFileName)
	}
	
	if inFile, err = os.Open(inFileName); err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	nReader = bufio.NewReader(inFile)
	
	// Read the file
	log.Println("Start reading ", inFileName)
	for nLines = 1; ; nLines++ { 
		if line, err = readfile.Readln(nReader); err != nil {
			if err.Error() != "EOF" {
				log.Fatal("Failed reading with err ", err)
			} 
			break
		}
		// Skips comments
		if strings.HasPrefix(line, "#") { 
			nLines--
			continue 
		}

		// Regexp quantities from the line
		if lineRes = lineReg.FindStringSubmatch(line); lineRes == nil {
			log.Fatalf("Can't reg %v on line %v: ", line, nLines)
		}
		newData = &Data{}
		if err = newData.Fill(lineRes, map[string]string{
			"comb": dirRes[1],
			"rv": dirRes[3],
			"fpb": dirRes[5],
			"tf": dirRes[2],
			"W0":dirRes[6],
		}); err != nil {
			log.Fatal("Can't fill with err: ", err)
		}
		*ds = append(*ds, newData)
	}
	fmt.Println("Done reading")
	return nLines
}

func ChirpMass(m0, m1 float64) (chirpmass float64) {
	mTot := m0 + m1
	mu := (m0 * m1) / mTot

	chirpmass = math.Pow(mu, 3./5) * math.Pow(mTot, 2./5)
	return chirpmass
}

// func (ds DataSlice) ToH5 (outFileName string) (err error) {
// 	var (
// 		tname string = "data"
// 		compress int = 0
// 	)
// 	f, err := hdf5.CreateFile(outFileName, hdf5.F_ACC_TRUNC)
// 	if err != nil {
// 		panic(fmt.Errorf("CreateFile failed: %s", err))
// 	}
// 	defer f.Close()
// 	fmt.Printf(":: file [%s] created (id=%d)\n", outFileName, f.Id())
// 	table, err := f.CreateTableFrom(tname, Data{}, len(ds), compress)
// 	if err != nil {
// 		panic(fmt.Errorf("CreateTableFrom failed: %s", err))
// 	}
// 	defer table.Close()
// 	fmt.Printf(":: table [%s] created (id=%d)\n", tname, table.Id())
// 
// // 	if !table.IsValid() {
// // 		panic("table is invalid")
// // 	}
// // 
// // 	// write one packet to the packet table
// // 	if err = table.Append(ds[0]); err != nil {
// // 		panic(fmt.Errorf("Append failed with single packet: %s", err))
// // 	}
// // 
// // 	// write several packets
// // 	parts := ds[1:]
// // 	if err = table.Append(parts); err != nil {
// // 		panic(fmt.Errorf("Append failed with multiple packets: %s", err))
// // 	}
// // 
// // 	// get the number of packets
// // 	n, err := table.NumPackets()
// // 	if err != nil {
// // 		panic(fmt.Errorf("NumPackets failed: %s", err))
// // 	}
// // 	fmt.Printf(":: nbr entries: %d\n", n)
// // 	if n != len(ds) {
// // 		panic(fmt.Errorf(
// // 			"Wrong number of packets reported, expected %d but got %d",
// // 			len(ds), n,
// // 		))
// // 	}
// // 
// // 	// iterate through packets
// // 	for i := 0; i != n; i++ {
// // 		p := make([]Data, 1)
// // 		if err := table.Next(&p); err != nil {
// // 			panic(fmt.Errorf("Next failed: %s", err))
// // 		}
// // 		fmt.Printf(":: data[%d]: %v\n", i, p)
// // 	}
// 	return nil
// }
