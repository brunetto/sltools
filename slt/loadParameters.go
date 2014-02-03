package slt

import (
	"bitbucket.org/brunetto/goutils/readfile"
	"bytes"
	"bufio"
	"fmt"
	"log"
	"os"
// 	"strconv"
	"strings"
	"text/template"
)

type ParamStruct struct{
	/*Comb int64
	Run int64
	Z float64
	NCM int64
	W int64
	Fpb float64
	Tidal int64*/
	Comb string
	Run string
	Z string
	NCM string
	W string
	Fpb string
	Tidal string
}

type ParamSlice []*ParamStruct


func LoadParams (paramFile string) (parameters ParamSlice) {
	if Debug {Whoami(true)}
	var (
		fileObj *os.File
		nReader *bufio.Reader 
		readLine string
		err error
		paramLine []string
	)
	
	if paramFile == "" {
		log.Fatal(`You need to provide the file containing the parameters in the form:
# comb, run, Z, NCM, W, fPB, tidal
16, 10, 0.10, 10000, 5, 0.05, 0
17, 10, 0.10, 10000, 5, 0.10, 0
18, 10, 0.10, 10000, 5, 0.20, 0
19, 10, 0.10, 10000, 9, 0.05, 0
...
`)
	}
	
	log.Println("Start reading parameter file ", paramFile)
	
	if fileObj, err = os.Open(paramFile); err != nil {
		log.Fatal("Can't open parameter file ", paramFile, " with error ", err)
	}
	defer fileObj.Close()
	
	parameters = make(ParamSlice, 0)
		
	nReader = bufio.NewReader(fileObj)
	
	for {
		if readLine, err = readfile.Readln(nReader); err != nil {break}
		// Detect parameter names
		if strings.HasPrefix(readLine, "#") {
			continue
		}
		// Retrieve parameters
		paramLine = strings.Split(strings.Replace(readLine, " ", "", -1), ",")
		
		if len(paramLine) == 0 {
			continue
		}
		
		parameters = append(parameters, new(ParamStruct))
		
		this := len(parameters) - 1

		/*parameters[this].Comb, _ = strconv.ParseInt(paramLine[0], 10, 64)
		parameters[this].Run, _ = strconv.ParseInt(paramLine[1], 10, 64)
		parameters[this].Z, _ = strconv.ParseFloat(paramLine[2], 64)
		parameters[this].NCM, _ = strconv.ParseInt(paramLine[3], 10, 64)
		parameters[this].W, _ = strconv.ParseInt(paramLine[4], 10, 64)
		parameters[this].Fpb, _ = strconv.ParseFloat(paramLine[5], 64)
		parameters[this].Tidal, _ = strconv.ParseInt(paramLine[6], 10, 64)*/
		parameters[this].Comb  = paramLine[0]
		parameters[this].Run   = paramLine[1]
		parameters[this].Z     = /*strings.Replace(*/paramLine[2]/*, ".", "", -1)*/
		parameters[this].NCM   = paramLine[3]
		parameters[this].W     = paramLine[4]
		parameters[this].Fpb   = /*strings.Replace(*/paramLine[5]/*, ".", "", -1)*/
		parameters[this].Tidal = paramLine[6]
	}
	log.Println("Parameter slice populated with ", len(parameters), " elements")
	return parameters
}


func (parameters ParamSlice) PrintParams () {
	const paramString = "{{.Comb}}\t{{.Run}}\t{{.Z}}\t{{.NCM}}\t{{.W}}\t{{.Fpb}}\t{{.Tidal}}"
	var b bytes.Buffer
	
	paramTmpl := template.Must(template.New("printParams").Parse(paramString))
	
	fmt.Println("#\tcomb\trun\tZ\tNCM\tW\tFpb\ttidal")
	
	
	
	for idx, paramSet := range parameters {
		fmt.Print(idx, "\t")
		if err := paramTmpl.Execute(&b, paramSet); err != nil {
			log.Println("Error while executing template:", err)
		}
		fmt.Println(b.String())
		b.Reset()
// 		fmt.Println(idx, paramSet.comb, paramSet.run, paramSet.Z, paramSet.NCM, paramSet.W, paramSet.Fpb, paramSet.tidal)
	}
}

