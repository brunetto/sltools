package slt

import (
	"log"
	"regexp"
)

func Reg (inFileName string) (map[string]string) {
	var (
		regString string         = `(\w{3})-(\S*-comb(\d*)-\S*)-run(\d*)-rnd(\d*)(\.\S*)`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
	)
	
	if regRes = regExp.FindStringSubmatch(inFileName); regRes == nil {
		log.Fatal("Can't extract info in ", inFileName)
	}
	
	return map[string]string{
		"baseName": regRes[2], 
		"ext": regRes[6],
		"prefix": regRes[1],
		"comb": regRes[3],
		"run": regRes[4],
		"rnd": regRes[5],
	}
}