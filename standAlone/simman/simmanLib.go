package simman

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	
	"code.google.com/p/go.crypto/ssh"
	
	"github.com/brunetto/goutils/connection"
)

// see http://kiyor.us/2013/12/29/golang-ssh-example/

func MainLoop (wakeUp chan map[string]string, messageChan chan string, jobInfoChan chan JobMap) () {
	
	var (
		usr, server, pathToKey string
		exists bool
		session *ssh.Session
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
		
		err error
		jobQueue string
		queueLines = JobMap{}
	)
	
	for conn := range wakeUp {
		
		messageChan <- "Check connection data...\n"
		
		if usr, exists = conn["usr"]; !exists {log.Fatal("usr not found on wakeUp chanel")}
		if server, exists = conn["server"]; !exists {log.Fatal("server not found on wakeUp chanel")}
		if pathToKey, exists = conn["pathToKey"]; !exists {log.Fatal("pathToKey not found on wakeUp chanel")}
		
		messageChan <- "Start connection...\n"
		if session, err = connection.SshSessionWithKey(server, usr, pathToKey); err != nil {
			messageChan <- "Failed to create session: " + err.Error()
		}
	// 	defer session.Close()
		messageChan <- "Connection opened...\n"
		messageChan <- "Retrieving data...\n"

		jobQueue = retrieveQueue(usr, session)
		
		messageChan <- "Data retrieved...\n"
		
		session.Close()
		
		messageChan <- "Connection closed...\n"
		
		
		// Remove non job lines and create map
		for _, line := range strings.Split(jobQueue, "\n") {
			if queueRegRes = queueRegExp.FindStringSubmatch(line); queueRegRes != nil {
				// Check for duplicated job ID (impossible)
				if _, exists = queueLines[queueRegRes[1]]; exists {
					messageChan <- "Two jobs with the same id: " + queueRegRes[1]
				}
				
				if jobRegRes = jobRegExp.FindStringSubmatch(queueRegRes[4]); queueRegRes == nil {
					messageChan <- "Can't retrieve job info from job name: " + queueRegRes[4]
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
		
		messageChan <- "Data parsed and sent.\n"
		jobInfoChan <- queueLines
		
		// Clean job map
		queueLines = JobMap{}
		
	}
// 	log.Printf("Retrieved %v job lines ", len(queueLines))
// 	queueLines.Print()
// 	log.Println("Done")
	
}

func retrieveQueue (usr string, session *ssh.Session) (string) {
	var ( 
		b bytes.Buffer
		cmd string
		stdout string
		err error
	)
	
	session.Stdout = &b
	cmd = "qstat -u " + usr
	log.Println("Run: ", cmd)
	if err = session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	
	stdout = b.String()
	b.Reset()
	
	return stdout
}

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

type JobMap map[string]map[string]string


