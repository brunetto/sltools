package main 

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		trashFiles, scriptFiles, logFiles []string
		trashFile, scriptFile, logFile string
		trashDir = "Trash"
		scriptDir = "Scripts"	
		logDir = "Scripts"
		err error
	)
	
	log.Println("Check dirs existance and in case create them")
	if !goutils.Exists(scriptDir) {
		if err = os.Mkdir(scriptDir, 0700); err != nil {
			log.Fatal("Can't create folder ", err)
		}
	}
	
	if !goutils.Exists(trashDir) {
		if err = os.Mkdir(trashDir, 0700); err != nil {
			log.Fatal("Can't create folder ", err)
		}
	}

	log.Println("Globbin files, found:")
	if trashFiles, err = filepath.Glob("r*"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	
	fmt.Println(trashFiles)
	
	if scriptFiles, err = filepath.Glob("*.sh"); err != nil {
		log.Fatal("Error globbing script files: ", err)
	}

	fmt.Println(scriptFiles)
	
	if logFiles, err = filepath.Glob("*.log"); err != nil {
		log.Fatal("Error globbing log files: ", err)
	}
	
	fmt.Println(logFiles)
	
	log.Println("Moving files")
	
	for _, trashFile = range trashFiles {
		if err = os.Rename(trashFile, filepath.Join(trashDir, trashFile)); err != nil {
			log.Fatal("Error while moving ", trashFile, ": ", err)
		}
	}

	for _, scriptFile = range scriptFiles {
		if err = os.Rename(scriptFile, filepath.Join(scriptDir, scriptFile)); err != nil {
			log.Fatal("Error while moving ", scriptFile, ": ", err)
		}
	}
	for _, logFile = range logFiles {
		if err = os.Rename(logFile, filepath.Join(logDir, logFile)); err != nil {
			log.Fatal("Error while moving ", logFile, ": ", err)
		}
	}
	
	log.Println("Deleting tmp files")
	
	if err = os.RemoveAll("*~"); err != nil {
		log.Fatal("Error while removing *~: ", err)
	}
	
}

