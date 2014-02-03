package main

import (
	"log"
	"os"
	"os/exec"
)

func main () {
	bashCmd := exec.Command("bash", "cineca-comb18-run1_10-NCM10000-fPB020-W5-Z010/create_IC-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.sh")
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		log.Fatal(err)
	}
}