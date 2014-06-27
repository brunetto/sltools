package main

import (
	"bytes"
// 	"fmt"
	"log"
	"regexp"
	"strings"
	
	"code.google.com/p/go.crypto/ssh"
	"github.com/brunetto/goutils/connection"
)

// see http://kiyor.us/2013/12/29/golang-ssh-example/

func main () () {
	
	var (
		usr = "bziosi00"
		server = "login.eurora.cineca.it:22"
		config *ssh.ClientConfig
		queueRegString string         = `\S+\.\S+\s+(\S+)\s+(\S+)\s+(\S+)\s+\S+\s+\d+\s+\d+\s+(\S+)\s+(\d+):(\d+)\s+(\S)\s+\S+`
// 		group 1: job name
// 		group 2: user
// 		group 3  queue 
// 		group 4: job name 2
// 		group 5: requested ram
// 		group 6: remaining hours
// 		group 7: remaining minutes
// 		group 8: status
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
		queueLines map[string]map[string]string
	)
	
	config = connection.PubKeyClientConfig(usr, "")
	
	log.Println("Try to connect to ", server)
			
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Println("Failed to dial: " + err.Error())
	}
	
	log.Println("Start new session")
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: " + err.Error())
	}
	defer session.Close()

	jobQueue = retrieveQueue(usr, session)
	
	tmp := strings.Split(jobQueue, "\n")
	
	// Remove non job lines
	for idx, line := range queueLines {
		if queueRegRes = queueRegExp.FindStringSubmatch(line); queueRegRes != nil {
			// 		group 1: job name
// 		group 2: user
// 		group 3  queue 
// 		group 4: job name 2
// 		group 5: requested ram
// 		group 6: remaining hours
// 		group 7: remaining minutes
// 		group 8: status
			
			
			// 		jobRegString string         = `\S+(\d{2})-(\d{2})-(\d{2})`
// 		group 1: comb
// 		group 2: run
// 		group 3  round 
// 		jobRegExp    *regexp.Regexp = regexp.MustCompile(jobRegString)
// 		jobRegRes []string
			queueLines = append(queueLines[:idx], queueLines[idx+1:]...)
		}
	}
	
	log.Printf("Retrieved %v job lines ", len(queueLines))
	
	log.Println("Done")
	
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
