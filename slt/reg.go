package slt

import (
	"fmt"
	"log"
	"regexp"
	
	"github.com/brunetto/goutils/debug"
)

func Reg (inFileName string) (map[string]string, error) {
	var (
		regString string         = `(\w{3})-(\S*comb(\S*?)-\S*)-run(\d*)-[a-z]*(\d*)(\.\S+\.*\S*)`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
		err error
	)
	
	if regRes = regExp.FindStringSubmatch(inFileName); regRes == nil {
		err = fmt.Errorf("%v can't extract name info from %v", debug.FName(false), inFileName)
		log.Println(err)
		return map[string]string{
		"baseName": "", 
		"ext": "",
		"prefix": "",
		"comb": "",
		"run": "",
		"rnd": "",
		}, err
	}
	
	return map[string]string{
		"baseName": regRes[2], 
		"ext": regRes[6],
		"prefix": regRes[1],
		"comb": regRes[3],
		"run": regRes[4],
		"rnd": regRes[5],
	}, nil
}


func DeepReg (inFileName string) (map[string]string, error) {
	var (
		regString string         = `(\w{3})-(\S*comb(\S*?)-TF(\S+)-Rv(\d+)-NCM(\d+)-fPB(\d+)-W(\d+)-Z(\d+))-run(\d*)-[a-z]*(\d*)(\.\S+\.*\S*)`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
		err error
	)
	
	if regRes = regExp.FindStringSubmatch(inFileName); regRes == nil {
		err = fmt.Errorf("%v can't extract name info from %v", debug.FName(false), inFileName)
		log.Println(err)
		return map[string]string{
		"baseName": "", 
		"ext": "",
		"prefix": "",
		"comb": "",
		"TF": "", 
		"Rv": "",
		"NCM": "",
		"fPB": "",
		"W": "",
		"Z": "",
		"run": "",
		"rnd": "",
		}, err
	}
	
	return map[string]string{
		"baseName": regRes[2], 
		"ext": regRes[12],
		"prefix": regRes[1],
		"comb": regRes[3],
		"TF": regRes[4], 
		"Rv": regRes[5],
		"NCM": regRes[6],
		"fPB": regRes[7],
		"W": regRes[8],
		"Z": regRes[9],
		"run": regRes[10],
		"rnd": regRes[11],
	}, nil
}

