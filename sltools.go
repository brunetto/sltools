package main

import (
	"bitbucket.org/brunetto/sltools/slt"
)

func main() {
	
	slt.InitCommands()
	slt.SltCmd.Execute()
	
} // END MAIN

