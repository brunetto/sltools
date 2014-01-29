package main

import (
	"./slt"
)

func main() {
	
	slt.InitCommands()
	slt.SltCmd.Execute()
	
} // END MAIN

