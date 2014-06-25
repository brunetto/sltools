package main

import (
	"fmt"
	"log"
	"path/filepath"
	"text/template"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
// 	"github.com/brunetto/goutils/sets"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		err          error
		inFiles      []string
		idx, file string
		globName string = "*-cineca-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt"
		
	)
	
	log.Println("Searching for files in the form: ", globName)
		
	if inFiles, err = filepath.Glob(globName); err != nil {
		log.Fatal("Error globbing files in this folder: ", err)
	}
	
	for idx, file = range inFiles {
		fmt.Println(idx, file)
	}
	
	
	
}


// printf "\n"; pwd; printf "\n"; for (( c=0; c<=9; c++ )); do printf "$c "; ls -lah out-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; prStintf "  "; ls -lah err-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; printf "  "; cat $(ls err-*-run0$c-rnd0* | tail -n 1) | grep "Time = " | tail -n 1; done
