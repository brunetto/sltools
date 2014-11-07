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
		regString string         = `(\w{3})-cineca-comb(\d*)-NCM(\d+)-fPB(\d+)-W(\d+)-Z(\d+)-run(\d*)-rnd(\d*)(\.\S*\.*\S*)`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
		dry bool = false
	)
		
	if len(os.Args) > 1 {
		if os.Args[1] == "--dry" {
			dry = true
		}
	}
	
	if files, err = filepath.Glob("*cineca-comb*.*"); err != nil {
		log.Fatal("Can't glob files")
	}
		
	for _, file = range files {
		if regRes = regExp.FindStringSubmatch(file); regRes == nil {
			fmt.Printf("Can't reg %v\n", file)
			continue
		}
		prefix := regRes[1]
		comb := regRes[2]
		ncm := regRes[3]
		fPB := regRes[4]
		w := regRes[5]
		z := regRes[6]
		run := regRes[7]
		rnd := regRes[8]
		ext := regRes[9]
		rv := "1"
		
		newfile = prefix + "-" + "comb" + comb + 
					"-TFno-Rv" + rv + 
					"-NCM" + ncm + 
					"-fPB" + fPB + 
					"-W" + w + 
					"-Z" + z + 
					"-run" + run + 
					"-rnd" + rnd + ext
		
		fmt.Printf("Renaming %v in %v \n", file, newfile)
		if !dry {
			if err = os.Rename(file, newfile); err != nil {
				log.Fatalf("Can't rename %v with error %v\n", file, err)
			}
		}
	}
}




