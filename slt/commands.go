package slt

import (
	"github.com/spf13/cobra"
	"fmt"
)

// Package-wise verbosity
// use with:
// if Verb { ...
var Verb bool 
var Debug bool
var ConfName string

// Root command
var SlToolsCmd = &cobra.Command{
	Use:   "sltools",
	Short: "Tools for StarLab simulation management",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type sltools help for help.")
	},
}

// Print version
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slt",
	Long:  `All software has versions. This is sltools'`,
	Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("StarLab Tools v0.6")
	},
}

// Load JSON configuration file
var ReadConfCmd = &cobra.Command{
	Use:   "readConf",
	Short: "Read and print the configuration file",
	Long:  `Read and print the configuration specify by the -c flag.
It must be in the form of a JSON file like:

{
	"Runs": 10,
	"Comb": 18, 
	"Ncm" : 10000,
	"Fpb" : 0.10,
	"W"   : 5,
	"Z"   : 0.20,
	"EndTime" : 500,
	"Machine" : "plx",
	"UserName" : "bziosi00",
	"PName": "IscrC_VMStars" 
}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		if Verb {
			fmt.Println("Config:")
			conf.Print()
		}
	},
}

// Create ICs from JSON configuration file
var RunICC bool
var CreateICsCmd = &cobra.Command{
	Use:   "createICs",
	Short: "Create ICs",
	Long:  `Create initial conditions from the JSON config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		CreateICs(conf)
// 		CreateICsOld(conf)
	},
}

// Create new ICs from STDOUT to restart the simulation
var (
	inFileName string
	fileN string
)
var Out2ICsCmd = &cobra.Command{
	Use:   "out2ics",
	Short: "Prepare the new ICs from the last STDOUT",
	Long:  `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing the last STDOUT and writing
	the last complete snapshot to the new input file.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		Out2ICs(inFileName, conf)
	},
}

// Create start scripts: kiraLaunch and PBSlaunch
var (
	icsName string
	randomNumber string
	simTime string
	)
var CreateStartScriptsCmd = &cobra.Command{
	Use:   "createStartScripts",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long:  `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		CreateStartScripts(icsName, randomNumber, simTime, conf)
	},
}

// Out2ICsCmd + CreateStartScriptsCmd
var ContinueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long:  `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		Continue(inFileName, conf)
	},
}

// Stich STDOUT and STDERR from different round of the same simulation 
// (if you restarded your simulation)
var (
	OnlyOut bool
	OnlyErr bool
	StichAll bool
)
var StichOutputCmd = &cobra.Command{
	Use:   "stichOutput",
	Short: "stich output, only for one simulation or for all in the folder",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		if StichAll {
			log.Println("Stich all!")
			StichThemAll (conf)
		} else {
			StichOutputSingle (inFileName, conf)
		}
	},
}


// Init commands and attach flags
func InitCommands() () {

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	SlToolsCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlToolsCmd.PersistentFlags().StringVarP(&ConfName, "confName", "c", "", "Name of the JSON config file")
	
	SlToolsCmd.AddCommand(ReadConfCmd)
	
	SlToolsCmd.AddCommand(CreateICsCmd)
	CreateICsCmd.Flags().BoolVarP(&RunICC, "runIcc", "C", false, "Run the creation of the ICs instad of only create scripts")
	
	SlToolsCmd.AddCommand(ContinueCmd)
	ContinueCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")
	
	SlToolsCmd.AddCommand(Out2ICsCmd)
	Out2ICsCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")	
	
	SlToolsCmd.AddCommand(CreateStartScriptsCmd)
	CreateStartScriptsCmd.Flags().StringVarP(&icsName, "icsName", "i", "", "ICs file name")
	CreateStartScriptsCmd.Flags().StringVarP(&simTime, "simTime", "t", "", "Remaining simulation time provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&randomNumber, "random", "r", "", "Init random seed provided by the out2ics command")
	
	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT or STDERR name to find what to stich")
	StichOutputCmd.Flags().BoolVarP(&OnlyOut, "onlyOut", "O", false, "Only stich STDOUTs")
	StichOutputCmd.Flags().BoolVarP(&OnlyErr, "onlyErr", "E", false, "Only stich STDERRs")
	StichOutputCmd.Flags().BoolVarP(&StichAll, "all", "A", false, "Stich all the run outputs in the folder")
}

