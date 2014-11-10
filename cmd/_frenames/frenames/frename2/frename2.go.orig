package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main () () {
	var (
		err error
		files []string
		file, newfile string
		regString string         = `\S*cineca(\d+)_bin_N(\d+)_frac(\d+)_W(\d+)_Z(\d+)\.*\S*(.\S+)`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
		dry bool = false
		comb, z string
	)
		
	if len(os.Args) > 1 {
		if os.Args[1] == "--dry" {
			dry = true
		}
	}
	
	if files, err = filepath.Glob("*cineca*"); err != nil {
		log.Fatal("Can't glob files")
	}
		
	for _, file = range files {
		if regRes = regExp.FindStringSubmatch(file); regRes == nil {
			fmt.Printf("Can't reg %v\n", file)
			continue
		}
			
		prefix := "err"
		ncm := "5000"
		fPB := "01"
		w := "5"
		zOrig := regRes[5]
// 		run := regRes[1]
		ext := ".txt.gz"
		rv := "1"
		
		if zOrig == "001" {
			comb = "1"
			z = "001"
		} else if zOrig == "01" {
			comb = "2"
			z = "010"			
		} else if zOrig == "1" {
			comb = "3"
			z = "100"
		}  
		
		newfile = prefix + "-" + "comb" + comb + 
					"-TFno-Rv" + rv + 
					"-NCM" + ncm + 
					"-fPB" + fPB + 
					"-W" + w + 
					"-Z" + z + 
					"-run" + regRes[1] + "-all" +
					/*"-rnd" + rnd +*/  ext
		
		fmt.Printf("Renaming %v in %v \n", file, newfile)
		if !dry {
			if err = os.Rename(file, newfile); err != nil {
				log.Fatalf("Can't rename %v with error %v\n", file, err)
			}
		}
	}
}


