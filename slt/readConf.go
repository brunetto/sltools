package slt

import (
	"bytes"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// Config is the configuration struc for sltools.
// It contains all the variables common to all the simulations in a given folder.
// BUG: I haven't decided yet if it is global to the package or it has to 
// be passed.
type ConfigStruct struct {
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

// ReadConf load configuration parameters for this set of runs forom a json file.
func (conf *ConfigStruct) ReadConf (confName string) () {
	if Debug {Whoami(true)}
	
	var (
		confFile []byte
		err error
	)
	if confName == "" {log.Fatal("You must specify a JSON config file")}
	if confFile, err = ioutil.ReadFile(confName); err != nil {log.Fatal(err)}
	if err = json.Unmarshal(confFile, conf); err != nil {log.Fatal("Parse config: ", err)}
	
	log.Print("Checking loaded values: ")
	if conf.Runs <= 0 {log.Fatal("Runs field in configuation file is empty, zero or negative")}
	if conf.Comb < 0 {log.Fatal("Comb field in configuation file is negative")}
	if conf.Ncm  <= 0 {log.Fatal("Ncm field in configuation file is empty, zero or negative")}
	if conf.Fpb  < 0 {log.Fatal("Fpb field in configuation file is negative")}
	if conf.W    < 0 {log.Fatal("W field in configuation file is  negative")}
	if conf.Z    < 0 {log.Fatal("Z field in configuation file is empty or negative")}
	if conf.Machine == "" {log.Fatal("Machine field in configuation file is empty")}
	if conf.UserName == "" {log.Fatal("UserName field in configuation file is empty")}
	if conf.PName == "" {log.Fatal("PName field in configuation file is empty")}
	if conf.EndTime <= 0 {log.Fatal("EndTime field in configuation file is empty, zero or negative")}
	fmt.Print("OK!")
}

// Print prints the configuration parameters.
func (conf *ConfigStruct) Print () () {
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

// RunsStr return the number of runs to run in string form
func (conf *ConfigStruct) RunsStr () (string) {
	return strconv.Itoa(conf.Runs)
}

// CombStr return the combination number in string form
func (conf *ConfigStruct) CombStr () (string) {
	return strconv.Itoa(conf.Comb)
}

// NcmStr return the number of centres of mass in the simulation in string form
func (conf *ConfigStruct) NcmStr () (string) {
	return strconv.Itoa(conf.Ncm)
}

// FpbStr return the primordial binaries fraction in string form
func (conf *ConfigStruct) FpbStr () (string) {
	return strconv.FormatFloat(conf.Fpb, 'f', 2, 64)
}

// FpbCmpStr return the primordial binaries fraction in compact string form
// i.e. without the dot and filled with zeroes
func (conf *ConfigStruct) FpbCmpStr () (string) {
	return strings.Replace(conf.FpbStr(), ".",  "", -1)
}

// ZStr return the metallicity in string form
func (conf *ConfigStruct) ZStr () (string) {
	return strconv.FormatFloat(conf.Z, 'f', 2, 64)
}

// ZCmpStr return the metallicity in compact string form
// i.e. without the dot and filled with zeroes
func (conf *ConfigStruct) ZCmpStr () (string) {
	return strings.Replace(conf.ZStr(), ".",  "", -1)
}

// WStr return the central adimensional potential in string form
func (conf *ConfigStruct) WStr () (string) {
	return strconv.Itoa(conf.W)
}

// WStr return the EndTime of the simulation in string form
func (conf *ConfigStruct) EndTimeStr () (string) {
	return strconv.Itoa(conf.EndTime)
}
	
// BaseNameprovides the basename for all the files here
// by filling the basename template with the configuation parameter.
// Not the best way but just to try the template package!:P
func (conf *ConfigStruct) BaseName () (string) {
	if Debug {Whoami(true)}
	
	// baseName string
	const baseName = "cineca-comb{{.CombStr}}" + 
			   "-NCM{{.NcmStr}}" + 
			   "-fPB{{.FpbCmpStr}}" + 
			   "-W{{.WStr}}" + 
			   "-Z{{.ZCmpStr}}"
	// template to be filled
	var baseTmpl *template.Template = template.Must(template.New("baseNameTmpl").Parse(baseName))
	// buffer to write into
	var buf bytes.Buffer
	// execute the template and check for errors
	if err := baseTmpl.Execute(&buf, conf); err != nil {
		log.Println("Error while creating basename in conf.BaseName:", err)
	}
	// copy the buf content in a string
	ret := buf.String()
	buf.Reset() // reset the buffer just to be sure
	return ret
}
	
// InitVars provide the configuration struct to the package.
// It will check if a json configuration file is provided by the user via the -c flag
// or if there's a default config.json file in the present folder
func InitVars(ConfName string) (*ConfigStruct) {
	if Debug {Whoami(true)}
	
	// Check for configuration file passed by -c flag
	// or for the standard default config.json
	if ConfName == "" {
		log.Println("No JSON configuration file provided by the user")
	} else if _, err := os.Stat("conf.json"); err == nil {
		log.Println("Search for default json configuration file conf.json")
		log.Printf("conf.json exists, I will read it")
		ConfName = "conf.json"
	} else {
		log.Println("Search for default json configuration file conf.json")
		log.Fatal("No json configuration file proided via the -c flag nor default config.json found in this folder.")
	}
	
	// Read conf file and create conf struct
	log.Println("Read configuration form ", ConfName)
	var conf *ConfigStruct = new(ConfigStruct)
	conf.ReadConf(ConfName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}
	// Return a pointer to the new configuration structure
	return conf	
}

	
	
	
	
	
	
	