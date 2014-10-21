package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/readfile"
	"github.com/brunetto/goutils/debug"
)

var (
	queueRegString string         = `(\S+\.\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+\S+\s+\d+\s+\d+\s+(\S+)\s+(\d+):(\d+)\s+(\S)\s+(\S+):*(\S+)`
	// 		group 1: jobID
	// 		group 2: user
	// 		group 3  queue 
	// 		group 4: jobName
	// 		group 5: requestedRam
	// 		group 6: requiredHours
	// 		group 7: requiredMinutes
	// 		group 8: status
	//		group 9: elapsedHours
	//		group 10: elapsedMinutes
	queueRegExp    *regexp.Regexp = regexp.MustCompile(queueRegString)
	queueRegRes []string
	jobRegString string         = `\S+(\d{2})-(\d{2})-(\d{2})`
	// 		group 1: comb
	// 		group 2: run
	// 		group 3  round 
	jobRegExp    *regexp.Regexp = regexp.MustCompile(jobRegString)
	jobRegRes []string
)

func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		user string = os.Getenv("USER")
		waitingTime time.Duration = time.Duration(1) * time.Hour
		baseFolder string
		folders = []string{}
		inFile *os.File
		nReader *bufio.Reader 
		err error
		line string
	)
	
	defer debug.TimeMe(time.Now())
	
	baseFolder, err = os.Getwd()
	
	////////////////////////
	//  ETERNAL LOOP
	////////////////////////
	for {
		////////////////////////
		// Wait the queue to finish
		////////////////////////
		Wait(user, waitingTime)
		
		// Read again folders file, maybe something changed
		if inFile, err = os.Open("folders.txt"); err != nil {
			log.Fatal("Can't open folders file: ", err)
		}
		defer inFile.Close()
		nReader = bufio.NewReader(inFile)
		
		if line, err = readfile.Readln(nReader); err != nil {
			if err.Error() == "EOF" {
				log.Fatal("Non EOF error reading folder list: ", err)
			}
		}
		folders = append(folders, line)

		////////////////////////
		// Folders loop
		////////////////////////
		for _, folder := range folders {
			
			////////////////////////
			// Enter in folder
			////////////////////////
			log.Println("Entering ", folder)
			if err = os.Chdir(folder); err != nil {
				log.Fatalf("Can't enter in %v, error: %v\n", folder, err)
			}
			
			////////////////////////
			// Clean folder
			////////////////////////
			log.Println("Clean")
			slt.SimClean()
			
			////////////////////////
			// Continue good runs
			////////////////////////
			log.Println("Check and continue good runs")
			slt.CAC()
			
			////////////////////////
			// Launch new runs
			////////////////////////
			log.Println("Submit runs")
			if err = slt.PbsLaunch(); err != nil {
				if err.Error() == "qsub: Job exceeds queue resource limits" {
					Wait(user, waitingTime)
				} else {
					log.Fatal(err)
				}
			}
			////////////////////////
			// Back to the base folder
			////////////////////////
			if err = os.Chdir(baseFolder); err != nil {
				log.Fatalf("Can't enter in %v, error: %v\n", folder, err)
			}
		}
	}	
	
	fmt.Print("\x07") // Beep when finish!!:D
}


func QueueCheck (user string) (JobMap) {
	var (
		stdo, stde bytes.Buffer 
		queueCmd *exec.Cmd
		err error
		queue string
		exists bool
		queueLines = JobMap{}
	)
	
	queueCmd = exec.Command("qstat", "-u", user)
	if queueCmd.Stdout = &stdo; err != nil {log.Fatal("Error connecting STDOUT: ", err)}
	if queueCmd.Stderr = &stde; err != nil {log.Fatal("Error connecting STDERR: ", err)}
	log.Println("Execute ", "qstat ", " -u ", user)
	if err = queueCmd.Start(); err != nil {
		log.Fatal("Error starting queueCmd: ", err)
	}
	
	if err = queueCmd.Wait(); err != nil {
		log.Fatal("Error while waiting for queueCmd: ", err)
	}
	fmt.Println(stdo.String())
	fmt.Println(stde.String())
	
	queue = stdo.String()
	stdo.Reset()
	stde.Reset()
	
	// Remove non job lines and create map
	for _, line := range strings.Split(queue, "\n") {
		if queueRegRes = queueRegExp.FindStringSubmatch(line); queueRegRes != nil {
			// Check for duplicated job ID (impossible)
			if _, exists = queueLines[queueRegRes[1]]; exists {
				log.Println("Two jobs with the same id: " + queueRegRes[1])
			}
			
			if jobRegRes = jobRegExp.FindStringSubmatch(queueRegRes[4]); queueRegRes == nil {
				log.Println("Can't retrieve job info from job name: " + queueRegRes[4])
			}
			
			queueLines[queueRegRes[1]] = map[string]string{
											"user": queueRegRes[2], 
											"queue": queueRegRes[3], 
											"jobName": queueRegRes[4], 
											"requestedRam": queueRegRes[5], 
											"requestedHours": queueRegRes[6], 
											"requestedMinutes": queueRegRes[7], 
											"status": queueRegRes[8],
											"elapsedHours": queueRegRes[9], 
											"elapsedMinutes": queueRegRes[10], 
											"comb": jobRegRes[1], 
											"run": jobRegRes[2], 
											"round": jobRegRes[3], 
			}
		}
	}
		
	return queueLines	
}

type JobMap map[string]map[string]string

func (m JobMap) Print () () {
	var PBSStatus = map[string]string{
			"E": "Job is exiting after having run.",
			"H": "Job is held.",
			"Q": "job is queued, eligable to run or routed.",
			"R": "job is running.",
			"T": "job is being moved to new location.",
			"W": "job is waiting for its	execution time (-a option) to	be reached.",
			"S": "(Unicos only) job is suspend.",
		}
		
	for key, value := range m {
		fmt.Printf("jobID: %v\nuser: %v\nqueue: %v\njobName: %v\nrequestedRam: %v\nrequestedTime: %v:%v\nelapsedTime: %v:%v\nstatus: %v\ncomb: %v\nrun: %v\nround: %v\n\n",
			key,
			value["user"],
			value["queue"],
			value["jobName"], 
			value["requestedRam"], 
			value["requestedHours"], 
			value["requestedMinutes"], 
			value["elapsedHours"], 
			value["elapsedMinutes"], 
			value["status"] + " = " + PBSStatus[value["status"]], 
			value["comb"], 
			value["run"], 
			value["round"])
	}
}

// Wait for the queue to be empty
func Wait (user string, waitingTime time.Duration) ()	{	
	for {
		jobMap := QueueCheck(user)
		// Stop waiting if all jobs are finished
		// and restart working in folders
		if len(jobMap) <= 95 {
			log.Println("Queue empty, start working on runs")
			break
		}
		log.Printf("Still %v runs to go, waiting %v...\n", len(jobMap), waitingTime)
		time.Sleep(waitingTime)
	}
}
	
