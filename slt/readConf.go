package slt

import (
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Runs int
	Comb int
	Ncm int
	Fpb float64
	W int
	Z float64
	Machine string
	UserName string
	PName string
	EndTime int
	BinFolder string
}

func (conf *Config) ReadConf (confName string) () {
	if Debug {Whoami(true)}
	var (
		confFile []byte
		err error
	)
	if confName == "" {log.Fatal("You must specify a JSON config file")}
	if confFile, err = ioutil.ReadFile(confName); err != nil {log.Fatal(err)}
	if err = json.Unmarshal(confFile, conf); err != nil {log.Fatal("Parse config: ", err)}
}

func (conf *Config) Print () () {
	if Debug {Whoami(true)}
	fmt.Println("Numer of realizations:\t\t", conf.Runs)
	fmt.Println("Combination number:\t\t", conf.Comb)
	fmt.Println("Number of center of mass:\t", conf.Ncm)
	fmt.Println("Number of primordial binaries:\t", conf.Fpb)
	fmt.Println("Central adim. potential:\t", conf.W)
	fmt.Println("Metallicity:\t\t\t", conf.Z)
	fmt.Println("Timesteps:\t\t\t", conf.EndTime)
	fmt.Println("Machine name:\t\t\t", conf.Machine)
	fmt.Println("UserName:\t\t\t", conf.UserName)
}