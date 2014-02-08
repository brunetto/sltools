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

var SlToolsCmd = &cobra.Command{
	Use:   "sltools",
	Short: "Tools for StarLab simulation management",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type sltools help for help.")
	},
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slt",
	Long:  `All software has versions. This is sltools'`,
	Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("StarLab Tools v0.3")
	},
}

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

var CreateICCmd = &cobra.Command{
	Use:   "createICs",
	Short: "Create ICs",
	Long:  `Create initial conditions from the JSON config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		CreateICs(conf)
	},
}

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
		Out2ICs(inFileName)
	},
}


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

var (
	inFileTmpl string
)

var StichOutputCmd = &cobra.Command{
	Use:   "stichOutput",
	Short: "Only download SL",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		StichOutput (inFileTmpl)
	},
}

func InitCommands() () {

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	SlToolsCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlToolsCmd.PersistentFlags().StringVarP(&ConfName, "confName", "c", "", "Name of the JSON config file")
	
	SlToolsCmd.AddCommand(ReadConfCmd)
	
	SlToolsCmd.AddCommand(CreateICCmd)
	
	SlToolsCmd.AddCommand(ContinueCmd)
	ContinueCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")
	
	SlToolsCmd.AddCommand(Out2ICsCmd)
	Out2ICsCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")	
	
	SlToolsCmd.AddCommand(CreateStartScriptsCmd)
	CreateStartScriptsCmd.Flags().StringVarP(&icsName, "icsName", "i", "", "ICs file name")
	CreateStartScriptsCmd.Flags().StringVarP(&simTime, "simTime", "t", "", "Remaining simulation time provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&randomNumber, "random", "r", "", "Init random seed provided by the out2ics command")
	
	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().StringVarP(&inFileTmpl, "inTmpl", "i", "", 
			"STDOUT template name (the STDOUT name without the extention and the )")
	
}

