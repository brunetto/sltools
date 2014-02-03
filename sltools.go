package main

import (
	"bitbucket.org/brunetto/sltools/slt"
)

func main() {
	
	slt.InitCommands()
	slt.SlToolsCmd.Execute()
	
} 

