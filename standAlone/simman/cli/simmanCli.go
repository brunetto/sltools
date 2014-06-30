package main

import (
	"fmt"
	"log"
	"time"
	
	"bitbucket.org/brunetto/sltools/standAlone/simman"
	
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/notify"
)



func main () () {
	defer debug.TimeMe(time.Now())
	
	var (
		waitingTime time.Duration = time.Duration(10) * time.Minute
		wakeUp = make(chan map[string]string, 1)
		messageChan = make(chan string, 1)
		jobInfoChan = make(chan simman.JobMap, 1)
		conn = map[string]string{
			"usr": "bziosi00",
			"server": "login.eurora.cineca.it:22",
			"pathToKey": "",
			}
	)
	
	go simman.MainLoop(wakeUp, messageChan, jobInfoChan)
	go printMessages(messageChan)
	
	for {
		wakeUp <- conn
		
		queueLines := <- jobInfoChan
		
		fmt.Println()
		
		if len(queueLines) == 0 {
			log.Println("All job finished")
			break
		}
		
		log.Printf("%v active jobs: \n", len(queueLines))
		queueLines.Print()
	
		log.Println("Next check in ", waitingTime)
		time.Sleep(waitingTime)
	}
	
	close(wakeUp)
	close(messageChan)
	close(jobInfoChan)
	
	notify.Notify("EURORA", "All job are done!!!")
	fmt.Print("\x07") // Beep when finish!!:D
}

func printMessages (messageChan chan string) () {
	for message := range messageChan {
		log.Printf(message)
	}
}

