package main

import (
	"log"
	"os"
	"os/exec"
)

func main () {
	bashCmd := exec.Command("/bin/bash", "create_IC-cineca-comb20-NCM10000-fPB010-W9-Z010-run01-rnd00.sh")
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	if err := bashCmd.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Something else")
}

