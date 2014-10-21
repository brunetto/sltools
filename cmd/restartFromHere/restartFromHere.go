package main

import (
	"fmt"
	"time"
	
	"github.com/spf13/cobra"
	
	"github.com/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	initCommands()
	restartFromHereCmd.Execute()
	
}

var (
	inFileName string
	selectedSnapshot string
)

var restartFromHereCmd = &cobra.Command {
	Use:   "restartFromHere",
	Short: "Prepare a pp3-stalled simulation to be restarted",
	Long: `Too often StarLab stalled while integrating a binary,
	this tool let you easily restart a stalled simulation.
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	restartFromHere out --inFile <STDOUT file> --cut <snapshot where to cut>
	restartFromHere err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var stdOutCutCmd = &cobra.Command {
	Use:   "out",
	Short: "Prepare a pp3-stalled stdout to restart the simulation",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		slt.RestartStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrCutCmd = &cobra.Command {
	Use:   "err",
	Short: "Prepare a pp3-stalled stderr so that it is synced with the stdout",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		slt.RestartStdErr(inFileName, selectedSnapshot)
	},	
}	

func initCommands () {
	restartFromHereCmd.AddCommand(stdOutCutCmd)
	restartFromHereCmd.AddCommand(stdErrCutCmd)
	
	restartFromHereCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	restartFromHereCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
}


