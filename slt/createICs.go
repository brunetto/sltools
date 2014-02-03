package slt


import (
	
	"log"
	"strconv"
)

func CreateICs (confName string) () {
	
	var (
		conf *Config
		
		
	)
	log.Println("Read configuration form ", confName)
	conf.ReadConf(confName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}
	
	log.Println("Create folder ", )
	// Create folder
	runs = strconv.ItoA(conf.Runs)
	comb = strconv.ItoA(conf.Comb)
	ncm = strconv.ItoA(conf.Ncm)
	fpb = strings.Replace(strconv.FormatFloat(conf.Fpb, 'f', 2, 32), ".",  "", -1)
	w = strconv.ItoA(conf.W)
	z = strings.Replace(strconv.FormatFloat(conf.Z, 'f', 2, 32), ".",  "", -1)
	
	folderName = "cineca-comb" + comb + "-run1_" + runs + "-NCM" + ncm + "-fPB" + fpb + "-W" + w + "-Z" + z
	if err = os.Mkdir(folderName, 0700); err != nil {log.Fatal("Can't create folder ", err)}

	// Create bash script
	
	
	// Run bash script
		
	// Creating initial kira script
	
	// Creating initial PBS script
}
