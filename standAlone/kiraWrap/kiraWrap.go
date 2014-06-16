package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	
	"github.com/brunetto/goutils/debug"
	"bitbucket.org/brunetto/sltools/slt"
)


func main () () {
	defer debug.TimeMe(time.Now())
	
	if len(os.Args) < 3 {
		log.Fatal("Provide the ICs and how many timesteps you want to integrate")
	}
	
	var (
		timeLimit string
		err error
		kiraString string
		kiraArgs []string
		kiraWrappedCmd *exec.Cmd
		icsName, outName, errName, ext string
		icsFile, outFile, errFile *os.File
		regRes map[string]string
		done = make(chan string, 1)
	)
	
	icsName = os.Args[1]
	timeLimit = os.Args[2]
	
	// Extract fileNameBody, round and ext
	log.Println("Extract files names")
	ext = filepath.Ext(icsName)
	regRes, err = slt.Reg(icsName)
	if err != nil {
		log.Println("Can't derive standard names from STDOUT => wrap it!!")
		errName = "err-" + icsName + ext
		outName = "out-" + icsName + ext
	} else {
		if regRes["prefix"] != "ics" {
			log.Fatalf("Please specify a STDIN file, found %v prefix", regRes["prefix"])
		}
		
		// Creating new filenames
		errName = "err-" + regRes["baseName"] + "-run" + regRes["run"] + "-rnd" + regRes["rnd"] + ext
		outName = "out-" + regRes["baseName"] + "-run" + regRes["run"] + "-rnd" + regRes["rnd"] + ext
	}
	
	if icsFile, err = os.Open(icsName); err != nil {log.Fatal(err)}
	defer icsFile.Close()
	if outFile, err = os.Create(outName); err != nil {log.Fatal(err)}
	defer outFile.Close()
	if errFile, err = os.Create(errName); err != nil {log.Fatal(err)}
	defer errFile.Close()
		
	kiraString = "/home/ziosi/Code/Mapelli/slpack/starlab/usr/bin/kira"
	kiraArgs =  []string{"-t", timeLimit,// +  // number of timesteps to compute
				"-d", "1",// +  // log output interval
				"-D", "1",// +  // snapshot interval
				"-b", "1",// +  // frequency of full binary output
				"-f", "0",// +  // dynamical friction (0 = no friction, 1 = friction)
				"-n", "10",// +  // terminate if teh cluster remains with only 10 particles
				"-e", "0.000",// + // softening 
				"-B",// // switch on binary evolution
				//"-s 36543" // random seed 
	}
	kiraWrappedCmd = exec.Command(kiraString, kiraArgs...)
	if kiraWrappedCmd.Stdin = icsFile; err != nil {log.Fatal("Error connecting ICs to kira STDIN: ", err)}
	if kiraWrappedCmd.Stdout = outFile; err != nil {log.Fatal("Error connecting ICs to kira STDOUT: ", err)}
	if kiraWrappedCmd.Stderr = errFile; err != nil {log.Fatal("Error connecting ICs to kira STDERR: ", err)}
	
	log.Println("Going to run:")
	fmt.Println(kiraString, kiraArgs)
	fmt.Println("STDIN = ", icsName)
	fmt.Println("STDOUT = ", outName)
	fmt.Println("STDERR = ", errName)
	
	log.Println("Ready... steady... Go!")
	
	if err = kiraWrappedCmd.Start(); err != nil {
		log.Fatal("Error starting kiraWrappedCmd: ", err)
	}
	
	log.Println("Waiting kira to finish while checking for pp3-stalling...")
	// Wait for the process to end normally
	go waitProcess(kiraWrappedCmd, done)
	// Check for pp3-stalling situations
	go checkFileSize(errName, kiraWrappedCmd, done)	
	
	if <-done == "killed" {
		errFile.WriteString("\nKilled because pp3 stalling\n")
	}
}

func checkFileSize(errName string, kiraWrappedCmd *exec.Cmd, done chan string) () {
	var (
		fileInfo os.FileInfo
		err error
	)
	for {
		if fileInfo, err = os.Stat(errName); err != nil {
			log.Fatal("Error checking STDERR file size, err")
		}
		// STDERR exceesing aloowed dimension of 2GB
		// probably the simulation is stalling because 
		// of pp3 locked on a binary
		if fileInfo.Size() / (1024*1024*1024) > 2 {break} 
	}
	log.Println("Detected STDERR file with dimension ", fileInfo.Size())
	log.Println("Kill kira")
	if err := kiraWrappedCmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill: ", err)
	}
	done <- "killed"
}

func waitProcess(kiraWrappedCmd *exec.Cmd, done chan string) () {
	err := kiraWrappedCmd.Wait()
	log.Println("Process exited with error ", err)
	done <- "ended"	
}






