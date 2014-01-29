package slt


import (
// 	"bytes"
// 	"fmt"
	"log"
// 	"os"
// 	"path/filepath"
// 	"text/template"
)

func CreateICs (paramFile string) () {

	log.Println("Sorry, I'm not ready yet")

	/*
	var (
		parameters ParamSlice
		folderString string
		folderName string
		folderBuffer bytes.Buffer
		folderTmpl *template.Template
	)
	
	// Load parameters
	parameters = LoadParams(paramFile)
	
	// Review parameters
	fmt.Println()
	fmt.Println("Chosen parameters:\n")
	parameters.PrintParams()
	
	// Create folder name template to be filled with parameters
	folderString = "cineca-comb{{.Comb}}-run1_{{.Run}}-NCM{{.NCM}}-fPB{{.Fpb}}-W{{.W}}-Z{{.Z}}"
	folderTmpl = template.Must(template.New("folders").Parse(folderString))
	
	icScriptString = `./bin/makeking -n {{.NCM}} -w {{.W}} -i -u \
| ./bin/makemass -f 8  -l 0.1 -u 150 \
| ./bin/makesecondary -f {{.Fpb}} -q -l 0.1 \
| ./bin/add_star -R 1 -Z 0.10 \
| ./bin/scale -R 1 -M 1\
| ./bin/makebinary -f 2 -o 1 -l 1 -u 107836.09 \
> """ + os.path.join(folderName, name) + "-IC.txt`
	
	// Create folders
	fmt.Println()
	log.Println("Creating folders:\n")
	for idx := 0; idx < len(parameters); idx++ {
		if err := folderTmpl.Execute(&folderBuffer, parameters[idx]); err != nil {
			log.Fatal("Error while executing folder template:", err)
		}
		folderName = folderBuffer.String()
		folderBuffer.Reset()
		fmt.Println(folderName)
		if err := os.Mkdir(folderName, 0700); err != nil {
			log.Fatal("Can't create folder ", err)
		}
	}
	
	*/
	
	
	
	
}
