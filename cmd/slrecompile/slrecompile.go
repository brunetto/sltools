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
		syncCmd, configureCmd, configureCmd2, cleanCmd, cleanCmd2, cleanCmd3, makeCmd, makeCmd2, makeCmd3, rmBinCmd, installCmd *exec.Cmd
		confArgs []string
		err error
		files []string
		file string
		home, wd, wdBase, workDir string
	)
	
	if len(os.Args) < 2 {
		log.Println("Configure without flags!!!!!!!")
		log.Println("Assuming --with-f77=no")
		confArgs = []string{"--with-f77=no"}
	} else {
	
	if strings.Contains(os.Args[1], "help") {
		log.Fatal("Run as 'slrecompile <configure arguments>'")
	}
	confArgs = os.Args[1:]
	}
	
	if !goutils.Exists("configure") {
		log.Fatal("Can't find configure file, are you sure you are in the starlab folder?")
	}
	
	if !goutils.Exists(filepath.Join(".", "usr", "bin")) {
		log.Fatal("Can't find usr/bin folder, are you sure you are in the starlab folder?")
	}
	
	if home = os.Getenv("HOME"); home == "" {
		log.Fatal("Can't get home dir ")
	}
	if wd, err = os.Getwd(); err != nil {
		log.Fatal("Can't get local dir with error: ", err)
	}
	wdBase = filepath.Base(filepath.Dir(wd))
	workDir = filepath.Join(home, "Code", "Research", wdBase, "starlab")
	
	syncCmd = exec.Command("rsync", "-a", "-v", "-u", "-h", "-z",  "--progress", ".", workDir)
	configureCmd = exec.Command("./configure", confArgs...)
	configureCmd2 = exec.Command("./configure", confArgs...)
	cleanCmd     = exec.Command("make", "clean")
	cleanCmd2     = exec.Command("make", "clean")
	cleanCmd3     = exec.Command("make", "clean")
	rmBinCmd     = exec.Command("rm", filepath.Join(".", "usr", "bin", "*"))
	makeCmd      = exec.Command("make")
	makeCmd2      = exec.Command("make")
	makeCmd3      = exec.Command("make")
	installCmd   = exec.Command("make", "install")
	
	if syncCmd.Stderr = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if configureCmd.Stderr = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if configureCmd2.Stderr = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd2.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd3.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if rmBinCmd.Stderr     = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd.Stderr      = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd2.Stderr      = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd3.Stderr      = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if installCmd.Stderr   = os.Stderr; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	
	if syncCmd.Stdout = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)} 
	if configureCmd.Stdout = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)} 
	if configureCmd2.Stdout = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)} 
	if cleanCmd.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd2.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if cleanCmd3.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if rmBinCmd.Stdout     = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd.Stdout      = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd2.Stdout      = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if makeCmd3.Stdout      = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	if installCmd.Stdout   = os.Stdout; err != nil {log.Fatal("Can't pipe STD* to os.Std*: ", err)}
	
	
	log.Printf("Sync from %v to %v\n", wd, workDir)
	
	if err = syncCmd.Start(); err != nil {
		log.Fatal("Error starting syncCmd: ", err)
	}
	
	if err = syncCmd.Wait(); err != nil {
		log.Fatal("Error waiting syncCmd: ", err)
	}
	
	if err = os.Chdir(workDir); err != nil {
		log.Fatalf("Can't cd to %v with error %v \n", filepath.Join(home, "Code", "Mapelli", wdBase), err)
	}
	
	log.Println("Configure with args: ", confArgs)
	
	if err = configureCmd.Start(); err != nil {
		log.Fatal("Error starting configureCmd: ", err)
	}
	if err = configureCmd.Wait(); err != nil {
		log.Fatal("Error waiting configureCmd: ", err)
	}
	
	log.Println("Make clean")
	
	if err = cleanCmd.Start(); err != nil {
		log.Fatal("Error starting cleanCmd: ", err)
	}
	if err = cleanCmd.Wait(); err != nil {
		log.Fatal("Error waiting cleanCmd: ", err)
	}
	
	if err = cleanCmd2.Start(); err != nil {
		log.Fatal("Error starting cleanCmd: ", err)
	}
	if err = cleanCmd2.Wait(); err != nil {
		log.Fatal("Error waiting cleanCmd: ", err)
	}
	
	if err = cleanCmd3.Start(); err != nil {
		log.Fatal("Error starting cleanCmd: ", err)
	}
	if err = cleanCmd3.Wait(); err != nil {
		log.Fatal("Error waiting cleanCmd: ", err)
	}
	
	log.Println("Removing binaries")
	
	if files, err = filepath.Glob(filepath.Join("usr", "bin", "*")); err != nil  {
		log.Println("Can't Glob ", err)
	}
	
	log.Println(files)
	
	if len(files) > 0 {
		for _, file = range files {
			log.Println("Removing ", file)
			os.Remove(file)
		}
	}
		
	log.Println("Configure with args: ", confArgs)
	
	if err = configureCmd2.Start(); err != nil {
		log.Fatal("Error starting configureCmd: ", err)
	}
	if err = configureCmd2.Wait(); err != nil {
		log.Fatal("Error waiting configureCmd: ", err)
	}
	
	log.Println("Make")
	
	if err = makeCmd.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = makeCmd.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	if err = makeCmd2.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = makeCmd2.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	if err = makeCmd3.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = makeCmd3.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	log.Println("Make install")
	
	if err = installCmd.Start(); err != nil {
		log.Fatal("Error starting makeCmd: ", err)
	}
	if err = installCmd.Wait(); err != nil {
		log.Fatal("Error waiting makeCmd: ", err)
	}
	
	if err = os.Chdir(wd); err != nil {
		log.Fatalf("Can't cd to %v with error %v ", wd, err)
	}
	
	notify.Notify("StarLab recompile", "Done recompile!!!")
	fmt.Print("\x07") // Beep when finish!!:D
}


