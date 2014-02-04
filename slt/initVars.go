package slt

import (
	"log"
)

func InitVars() {
// Check for configuration file passed by -c flag
	if ConfName == "" {
		log.Fatal("Provide a JSON configuration file")
	}
	
	// Read conf file and create conf struct
	log.Println("Read configuration form ", ConfName)
	conf = new(Config)
	conf.ReadConf(ConfName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}
}
