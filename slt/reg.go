package slt

import (
	"fmt"
	"log"
	"regexp"
	
	"github.com/brunetto/goutils/debug"
)

func Reg (inFileName string) (map[string]string, error) {
	var (
		regString string         = `(\w{3})-(\S*-comb(\d*)-\S*)-run(\d*)-rnd(\d*)(\.\S*)`
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
