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
	Fpb float32
	W int
	Z float32
	Machine string
	UserName string
	PName string
}

func (conf *Config) ReadConf (confName string) () {
	var (
		confFile []byte
		err error
	)

	if confFile, err = ioutil.ReadFile(confName); err != nil {log.Fatal(err)}
	conf = new(Config)
	if err = json.Unmarshal(confFile, conf); err != nil {log.Fatal("parse config: ", err)}
}

func (conf *Config) Print () () {
	fmt.Println("Numer of realizations:\t", conf.Runs)
	fmt.Println("Combination number:\t", conf.Comb)
	fmt.Println("Number of center of mass:\t", conf.Ncm)
	fmt.Println("Number of primordial binaries:\t", conf.Fpb)
	fmt.Println("Central adim. potential:\t", conf.W)
	fmt.Println("Metallicity:\t", conf.Z)
	fmt.Println("Simulation end:\t", conf.EndTime)
	fmt.Println("Machine name:\t", conf.Machine)
	fmt.Println("UserName:\t", conf.UserName)
}