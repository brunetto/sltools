package slt 

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
)

func SimClean () () {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	var (
		files []string
		file string
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

	log.Println("Trash (PBS) files, found:")
	if files, err = filepath.Glob("r*"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Rename(file, filepath.Join(trashDir, file)); err != nil {
			log.Fatal("Error while moving ", file, ": ", err)
		}
	}
	
	log.Println("Scripts files, found:")
	if files, err = filepath.Glob("*.sh"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Rename(file, filepath.Join(scriptDir, file)); err != nil {
			log.Fatal("Error while moving ", file, ": ", err)
		}
	}
	
	log.Println("Log files, found:")
	if files, err = filepath.Glob("*.log"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Rename(file, filepath.Join(logDir, file)); err != nil {
			log.Fatal("Error while moving ", file, ": ", err)
		}
	}
	
	log.Println("Tmp files, found:")
	if files, err = filepath.Glob("*~"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Remove(file); err != nil {
			log.Fatal("Error while removing ", file, ": ", err)
		}
	}
	
	log.Println("Hidden files, found:")
	if files, err = filepath.Glob(".err*"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Remove(file); err != nil {
			log.Fatal("Error while removing ", file, ": ", err)
		}
	}
	
	if files, err = filepath.Glob(".out*"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Remove(file); err != nil {
			log.Fatal("Error while removing ", file, ": ", err)
		}
	}
	
	if files, err = filepath.Glob(".ics*"); err != nil {
		log.Fatal("Error globbing files to trash: ", err)
	}
	fmt.Println(files)
	for _, file = range files {
		if err = os.Remove(file); err != nil {
			log.Fatal("Error while removing ", file, ": ", err)
		}
	}
}
	
