package main 

import (
// 	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	
	got "github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	var (
		inFileNames []string
		inFileName string
		outFileName string
		err error
	)
	
	if inFileNames, err = filepath.Glob("*.gz"); err != nil {
			log.Fatal("Error globbing for stiching all the run outputs in this folder: ", err)
	}
	
	for _, inFileName = range inFileNames {
// 		outFileName = Rename(inFileName)
		if strings.Contains(inFileName, "fPB00") {
			outFileName = strings.Replace(inFileName, "fPB00", "fPB000", 1)
		} else {
			outFileName = strings.Replace(inFileName, "fPB01", "fPB010", 1)
		}
		if got.Exists(outFileName) {
			log.Fatal(outFileName, " from ", inFileName, " already exists!")
		}
		log.Println("Renaming", inFileName, " to ", outFileName)
		if err = os.Rename(inFileName, outFileName); err != nil {
			log.Fatal(err)
		}
	}
}

var prefixes = map[string]string {
	"ew": "err",
	"new": "out",
	"TF": "ics",
}

var regMap = map[string]*regexp.Regexp{
	"Run": regexp.MustCompile(`n(\d{2})`),
	"TF": regexp.MustCompile(`TF(\w*?)_`),
	"Fpb" : regexp.MustCompile(`frac(\d{2})`),
	"Rnd" : regexp.MustCompile(`txt(\d).gz`),
}

var resMap = map[string]string {
	"Run": "",
	"TF":  "",
	"Fpb" :  "",
	"Rnd" : "",
}


func Rename(inFileName string) (outFileName string) {
	var newprefix string
	// File prefix -> filetype
	for oldp, newp := range prefixes {
		if strings.HasPrefix(inFileName, oldp) {
			newprefix = newp
		}
	}
	
	for name, reg := range regMap {
		if res := reg.FindStringSubmatch(inFileName); res != nil {
			resMap[name] = res[1]
		} else {
			if name == "Rnd"{
				resMap["Rnd"] = "01"
			}
		}
	}
	
	outFileName = newprefix +  
					"-cineca-comb00"  + "-TF" + resMap["TF"] + "-NCM5000-fPB" + resMap["Fpb"] + "-" + 
					"W5-Z010-run" + resMap["Run"] + "-rnd" + got.LeftPad(resMap["Rnd"], "0", 2) + ".txt.gz"
	return "../renamed/" + outFileName
}