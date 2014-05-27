package slt

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// Verb control the package-wise verbosity.
// Use with:
// if Verb { ...
var Verb bool

// Debug activate the package-wise debug verbosity.
// Use with:
// if Verb { ...
var Debug bool

// ConfName is the name of the JSON configuration file.
var ConfName string

// SlToolsCmd is the main command.
var SlToolsCmd = &cobra.Command{
	Use:   "sltools",
	Short: "Tools for StarLab simulation management",
	Long: `SlTools would help in running simulations with StarLab.
It can create the inital conditions if StarLab is compiled and the 
necessary binaries are available.
SlTools can also prepare ICs from the last snapshot and stich the 
output.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type sltools help for help.")
	},
}

// VersionCmd print the sltools version.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slt",
	Long:  `All software has versions. This is sltools' one.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("StarLab Tools v0.8")
	},
}

// ReadConfCmd load the JSON configuration file.
var ReadConfCmd = &cobra.Command{
	Use:   "readConf",
	Short: "Read and print the configuration file",
	Long: `Read and print the configuration specify by the -c flag.
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

var (
	RunICC bool
	ICsAll bool
)

// CreateICsCmd will launch the functions to create the ICs from JSON configuration file.
var CreateICsCmd = &cobra.Command{
	Use:   "createICs",
	Short: "Create ICs from the JSON config file.",
	Long: `Create initial conditions from the JSON config file.
	Use like:
	sltools createICs -c conf21.json -v -C
	sltools createICs -v -C -A # to create folders and ICs for all the config files`,
	Run: func(cmd *cobra.Command, args []string) {
		if ICsAll {
			log.Println("Create all ICs following all the .json config files in this folder")
			CreateAllICs()
		} else {
			conf := InitVars(ConfName)
			CreateICsSingleWrap(conf)
		}
	},
}

var (
	inFileName string
)

// Out2ICsCmd creates new ICs from STDOUT to restart the simulation
var Out2ICsCmd = &cobra.Command{
	Use:   "out2ics",
	Short: "Prepare the new ICs from the last STDOUT",
	Long: `StarLab can restart a simulation from the last complete output.
	The out2ics command prepare the new ICs parsing the last STDOUT and writing
	the last complete snapshot to the new input file.
	Use like:
	sltools out2ics -i out-cineca-comb16-NCM10000-fPB005-W5-Z010-run06-rnd00.txt -n 1`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cssInfo = make (chan map[string]string, 1)
			inFileNameChan = make (chan string, 1)
		)
		go Out2ICs(inFileNameChan, cssInfo)
		inFileNameChan <- inFileName
		close(inFileNameChan)
		<-cssInfo
	},
}

var (
	icsName      string
	randomNumber string
	simTime      string
	machine string
)

// CreateStartScriptsCmd create start scripts: kiraLaunch and PBSlaunch
var CreateStartScriptsCmd = &cobra.Command{
	Use:   "createStartScripts",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long: `StarLab can restart a simulation from the last complete output.
	The createStartScripts write the necessary start scripts to start a 
	simulation from the ICs.
	Use like:
	sltools createStartScripts -i ics-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.txt -c conf.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			done = make (chan struct{})
			cssInfo = make (chan map[string]string, 1)
		)
		go CreateStartScripts(cssInfo, machine, done)
		cssInfo <- map[string]string{
				"remainingTime": simTime,
				"randomSeed": randomNumber,
				"newICsFileName": icsName,
		}
		close(cssInfo)
		<- done
		close(done)
	},
}

// Out2ICsCmd + CreateStartScriptsCmd
var ContinueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long: `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file. It also write the necessary 
	start scripts.
	Use like:
	sltools continue -o out-cineca-comb19-NCM10000-fPB005-W9-Z010-run08-rnd01.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		Continue(inFileName, machine)
	},
}

var (
	OnlyOut  bool
	OnlyErr  bool
	StichAll bool
)

// StichOutputCmd stiches STDOUT and STDERR from different round of the same simulation
// (if you restarded your simulation). Can be run serially or in parallel on all the
// file in the folder
var StichOutputCmd = &cobra.Command{
	Use:   "stichOutput",
	Short: "Stich output, only for one simulation or for all in the folder",
	Long: `Stich STDOUT and STDERR from different round of the same simulation 
	(if you restarded your simulation). Can be run serially or in parallel on all the
	file in the folder.
	You just need to select one of the files to stich or the --all flag to stich 
	all the files in the folder accordingly to their names.
	Use like:
	sltools stichOutput -c conf19.json -i out-cineca-comb19-NCM10000-fPB005-W9-Z010-run09-rnd00.txt
	sltools stichOutput -c conf19.json -A # to stich all the outputs in the folder`,
	Run: func(cmd *cobra.Command, args []string) {
		if StichAll {
			log.Println("Stich all!")
			StichThemAll(inFileName)
		} else {
			StichOutputSingle(inFileName)
		}
	},
}

// Init commands and attach flags
func InitCommands() {

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	SlToolsCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlToolsCmd.PersistentFlags().StringVarP(&ConfName, "confName", "c", "", "Name of the JSON config file")

	SlToolsCmd.AddCommand(ReadConfCmd)

	SlToolsCmd.AddCommand(CreateICsCmd)
	CreateICsCmd.Flags().BoolVarP(&RunICC, "runIcc", "C", false, "Run the creation of the ICs instad of only create scripts")
	CreateICsCmd.Flags().BoolVarP(&ICsAll, "all", "A", false, "Create all the ICs according to the conf.json files in the local folder")

	SlToolsCmd.AddCommand(ContinueCmd)
	ContinueCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")
	ContinueCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")

	SlToolsCmd.AddCommand(Out2ICsCmd)
	Out2ICsCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")

	SlToolsCmd.AddCommand(CreateStartScriptsCmd)
	CreateStartScriptsCmd.Flags().StringVarP(&icsName, "icsName", "i", "", "ICs file name")
	CreateStartScriptsCmd.Flags().StringVarP(&simTime, "simTime", "t", "", "Remaining simulation time provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&randomNumber, "random", "r", "", "Init random seed provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")

	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT or STDERR name to find what to stich")
	StichOutputCmd.Flags().BoolVarP(&OnlyOut, "onlyOut", "O", false, "Only stich STDOUTs")
	StichOutputCmd.Flags().BoolVarP(&OnlyErr, "onlyErr", "E", false, "Only stich STDERRs")
	StichOutputCmd.Flags().BoolVarP(&StichAll, "all", "A", false, "Stich all the run outputs in the folder")
}
