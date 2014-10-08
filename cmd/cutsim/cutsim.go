package main

import (
	"fmt"
	"time"
	
	"github.com/spf13/cobra"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())
	initCommands()
	cutsimCmd.Execute()
	
}

var (
	inFileName string
	selectedSnapshot string
)

var cutsimCmd = &cobra.Command {
	Use:   "cutsim",
	Short: `Shorten a give snapshot to a certain timestep
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	cutsim out --inFile <STDOUT file> --cut <snapshot where to cut>
	cutsim err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var bothCutCmd = &cobra.Command {
	Use:   "cut",
	Short: "cut <STDOUT or STDERR>",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		slt.CutStdBoth(inFileName, selectedSnapshot)
	},	
}

var stdOutCutCmd = &cobra.Command {
	Use:   "out",
	Short: "cut STDOUT",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		slt.CutStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrCutCmd = &cobra.Command {
	Use:   "err",
	Short: "cut STDERR",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		slt.CutStdErr(inFileName, selectedSnapshot)
	},	
}	

func initCommands () {
	cutsimCmd.AddCommand(stdOutCutCmd)
	cutsimCmd.AddCommand(stdErrCutCmd)
	
	cutsimCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	cutsimCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
}


