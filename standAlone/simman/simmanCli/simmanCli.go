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
		
		userChan = make(chan bool, 1)
	)
	
	go simman.MainLoop(wakeUp, messageChan, jobInfoChan)
	go printMessages(messageChan)
// 	go timer
	go usrInput(userChan)
	
	for {
		wakeUp <- conn
		
		queueLines := <- jobInfoChan
		
		fmt.Println()
		
		if len(queueLines) == 0 {
			log.Println("All job finished")
			break
		}
		
		log.Printf("%v active jobs: \n", len(queueLines))
// 		queueLines.Print()
	
		log.Println("Next check in ", waitingTime)
		fmt.Println("If you want details about the last check, write 'details'.")
// 		time.Sleep(waitingTime)
		
		select {
			case <-userChan:
				queueLines.Print()
			case <-time.Tick(waitingTime):
				continue
		}
		
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

// func timer (timerChan chan bool) () {
// 	
// 	
// }


func usrInput (userChan chan bool) () {
	var userInput string
	for {
		_, _ = fmt.Scan(&userInput)
		if userInput == "details" {
			userChan <- true
		} else {
			fmt.Println("Unknown command: ", userInput)
		}
	}
	
}
