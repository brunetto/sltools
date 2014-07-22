package main

import (
	"fmt"
	"log"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
	"time"
	
	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/notify"
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		configureCmd, cleanCmd, cleanCmd2, makeCmd, installCmd *exec.Cmd
		confArgs []string
		err error
		files []string
		file string
	)
	
	if strings.Contains(os.Args[1], "help") {
		log.Fatal("Run as 'slrecompile <configure arguments>'")
	}
	
	if !goutils.Exists("configure") {
		log.Fatal("Can't find configure file, are you sure you are in the starlab folder?")
	}
	
	if !goutils.Exists(filepath.Join(".", "usr", "bin")) {
		log.Fatal("Can't find usr/bin folder, are you sure you are in the starlab folder?")
	}
	
	confArgs = os.Args[1:]
	
	configureCmd = exec.Command("./configure", confArgs...)
	cleanCmd     = exec.Command("make", "clean")
	cleanCmd2     = exec.Command("make", "clean")
// 	rmBinCmd     = exec.Command("rm", filepath.Join(".", "usr", "bin", "*"))
	makeCmd      = exec.Command("make")
	installCmd   = exec.Command("make", "install")
	
	if configureCmd.Stderr = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd2.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
// 	if rmBinCmd.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd.Stderr      = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if installCmd.Stderr   = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	
	if configureCmd.Stdout = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)} 
	if cleanCmd.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd2.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
// 	if rmBinCmd.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd.Stdout      = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if installCmd.Stdout   = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	
	log.Println("Make clean")
	
	if err = cleanCmd.Start(); err != nil {
		log.Fatal("Error starting cleanCmd: ", err)
	}
	if err = cleanCmd.Wait(); err != nil {
		log.Fatal("Error waiting cleanCmd: ", err)
	}
	
	log.Println("Removing binaries")
	
	if files, err = filepath.Glob(filepath.Join("usr", "bin", "*")); err != nil  {
		log.Fatal("Can't Glob ", err)
	}
	
	log.Println(files)
	
	for _, file = range files {
		log.Println("Removing ", file)
		os.Remove(file)
	}
		
	log.Println("Configure with args: ", confArgs)
	
	if err = configureCmd.Start(); err != nil {
		log.Fatal("Error starting configureCmd: ", err)
	}
	if err = configureCmd.Wait(); err != nil {
		log.Fatal("Error waiting configureCmd: ", err)
	}
	
	log.Println("Make clean again")
	
	if err = cleanCmd2.Start(); err != nil {
		log.Fatal("Error starting cleanCmd: ", err)
	}
	if err = cleanCmd2.Wait(); err != nil {
		log.Fatal("Error waiting cleanCmd: ", err)
	}
	
	log.Println("Make")
	
	if err = makeCmd.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = makeCmd.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	log.Println("Make install")
	
	if err = installCmd.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = installCmd.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	notify.Notify("StarLab recompile", "Done recompile!!!")
	fmt.Print("\x07") // Beep when finish!!:D
}


