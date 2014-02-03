package slt


import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CreateICs (confName string, binFolder string) () {
	
	var (
		err error
		conf *Config
		runs string
		comb string
		ncm string
		fpb string
		w string
		z string	
		folderName string
		icsBaseCmd string
		baseName string
		outIcsName string
		outIcsScriptName string
		icsScriptFile *os.File
		icsScriptWriter *bufio.Writer
		err error
	)
	
	log.Println("Read configuration form ", confName)
	conf = new(Config)
	conf.ReadConf(confName)
	if Verb {
		log.Println("Loaded:")
		conf.Print()
	}

	if binFolder == "" {log.Fatal("I need to know where binaries for ICs are.")}
	
	log.Println("Create folder ", )
	// Create folder
	runs = strconv.Itoa(conf.Runs)
	comb = strconv.Itoa(conf.Comb)
	ncm = strconv.Itoa(conf.Ncm)
	fpb = strings.Replace(strconv.FormatFloat(conf.Fpb, 'f', 2, 64), ".",  "", -1)
	w = strconv.Itoa(conf.W)
	z = strings.Replace(strconv.FormatFloat(conf.Z, 'f', 2, 64), ".",  "", -1)
	
	makeking        := filepath.Join(binFolder, "makeking")
	makemass        := filepath.Join(binFolder, "makemass")
	makesecondary   := filepath.Join(binFolder, "makesecondary")
	add_star        := filepath.Join(binFolder, "add_star")
	scale           := filepath.Join(binFolder, "scale")
	makebinary      := filepath.Join(binFolder, "makebinary")
	
	icsBaseCmd = makeking + " -n " + ncm + " -w " + w + " -i -u \\\n" +
			"| " + makemass + " -f 8  -l 0.1 -u 150 \\\n" +
			"| " + makesecondary + " -f " + strconv.FormatFloat(conf.Fpb, 'f', 2, 64) + 
																	" -q -l 0.1 \\\n" +
			"| " + add_star + " -R 1 -Z " + strconv.FormatFloat(conf.Z, 'f', 2, 64) + 
																			" \\\n" +
			"| " + scale " -R 1 -M 1\\\n" +
			"| " + makebinary " -f 2 -o 1 -l 1 -u 107836.09 \\\n" +
			"> "  
	
	folderName = "cineca-comb" + comb + "-run1_" + runs + "-NCM" + ncm + "-fPB" + fpb + "-W" + w + "-Z" + z
	if err = os.Mkdir(folderName, 0700); err != nil {log.Fatal("Can't create folder ", err)}
	
	for runIdx :=1; runIdx<conf.Runs+1; runIdx++ {
		baseName = "cineca-comb" + comb + "-NCM" + ncm + "-fPB" +
					fpb + "-W" + w + "-Z" + z + "-run" +  
					LeftPad(strconv.Itoa(runIdx), 0, 2) + "-rnd00"
		outIcsName = filepath.Join(folderName, "ics-" + baseName)+".txt"
		outIcsScriptName = filepath.Join(folderName, "create_IC-" + baseName)+".sh"
		icsCmd = icsBaseCmd + outIcsName
		
		if icsScriptFile, err = os.Create(outIcsScriptName); err != nil {log.Fatal(err)}
		defer icsScriptFile.Close()
		
		icsScriptWriter = bufio.NewWriter(icsScriptFile)
		defer icsScriptWriter.Flush()
		
		icsScriptWriter.WriteString(icsCmd)
				 
				 
	}

	// Create bash script
	
	
	// Run bash script
		
	// Creating initial kira script
	
	// Creating initial PBS script
	
}
