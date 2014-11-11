package slt

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brunetto/goutils"
	"github.com/spf13/cobra"
)

var (
	All               bool
	as                bool = false
	Debug             bool
	ConfName          string // ConfName is the name of the JSON configuration file.
	err               error
	endOfSimMyrString string = "110"
	force             bool   = false
	inFileName        string
	intTime           string
	machine           string
	noGPU             bool = false
	noBinaries        bool = false
	OnlyOut           bool
	OnlyErr           bool
	randomNumber      string
	RunICC            bool //
	selectedSnapshot  string
	simTime           string
	tf                bool = false
	Verb              bool
)

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
		fmt.Println("StarLab Tools v1.1")
	},
}

// ***
var CacCmd = &cobra.Command{
	Use:   "cac",
	Short: "Check and continue, will check the last simulations outputs, prepare the restat and restart.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		CAC()
	},
}

// ***
var CheckEndCmd = &cobra.Command{
	Use:   "checkEnd",
	Short: "Check the number of timesteps necessary to reach a given time in Myr. Need the files to have standard names.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var endOfSimMyr float64 = 100
		if inFileName == "" || endOfSimMyrString == "" {
			log.Fatal("Provide a STDOUT file and a time in Myr to try to find the final timestep")
		} else {
			if endOfSimMyr, err = strconv.ParseFloat(endOfSimMyrString, 64); err != nil {
				log.Fatal(err)
			}
		}
		CheckEnd(inFileName, endOfSimMyr)
	},
}

// ***
var CheckSnapshotCmd = &cobra.Command{
	Use:   "checkSnapshot",
	Short: "Check the snapshot for being OK.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a STDOUT from which to check")
		}
		CheckSnapshot(inFileName)
	},
}

// ***
var CheckStatusCmd = &cobra.Command{
	Use:   "checkStatus",
	Short: "Check the status of a folder of simulations.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		CheckStatus()
	},
}

// ***
var ComOrbitCmd = &cobra.Command{
	Use:   "comorbit",
	Short: "Extract the center-of-mass coordinates from a STDOUT file.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a STDOUT from which to extract the center of mass coordinates for the orbit")
		}
		ComOrbit(inFileName)
	},
}

// Out2ICsCmd + CreateStartScriptsCmd
var ContinueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Prepare the new ICs and start scripts from all the last STDOUTs",
	Long: `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file. It also write the necessary 
	start scripts.
	Use like:
	sltools continue -o out-cineca-comb19-NCM10000-fPB005-W9-Z010-run08-rnd01.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		if machine == "" {
			if ConfName != "" {
				conf := InitVars(ConfName)
				machine = conf.Machine
			} else {
				log.Fatal("I need to know the machine name by CLI flag or conf file.")
			}
		}
		if All {
			inFileName = "all"
		}
		Continue(inFileName, machine)
	},
}

// CreateICsCmd will launch the functions to create the ICs from JSON configuration file.
var CreateICsCmd = &cobra.Command{
	Use:   "createICs",
	Short: "Create ICs from the JSON config file.",
	Long: `Create initial conditions from the JSON config file.
	Use like:
	sltools createICs -c conf21.json -v -C
	sltools createICs -v -C -A # to create folders and ICs for all the config files.
	I suggest not to use the -C flag and to create manually the ICs from the autogenerated bash scripts with 
	for $script in $(ls create*); do bash $script; done`,
	Run: func(cmd *cobra.Command, args []string) {
		if All {
			log.Println("Create all ICs following all the .json config files in this folder")
			CreateICsWrap("all", RunICC)
		} else {
			CreateICsWrap(ConfName, RunICC)
		}
	},
}

// CreateStartScriptsCmd create start scripts: kiraLaunch and PBSlaunch
var CreateStartScriptsCmd = &cobra.Command{
	Use:   "css",
	Short: "Prepare the scripts to start a run on a cluster with PBS ",
	Long: `StarLab can restart a simulation from the last complete output.
	"css" write the necessary start scripts to start a 
	simulation from the ICs specifying the machine name, the simulation time in timesteps and a ICs, 
	or with -A will do this for all the ICs in the folder.
	Use like:
	sltools createStartScripts -i ics-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.txt -t 500 -m eurora [-r 36541656]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			done             = make(chan struct{})
			cssInfo          = make(chan map[string]string, 1)
			pbsLaunchChannel = make(chan string, 100)
		)
		if machine == "" {
			if ConfName == "" {
				log.Fatal("You must provide a machine name or a valid config file")
			} else {
				conf := InitVars(ConfName)
				machine = conf.Machine
			}
		}
		// Consumes pbs file names
		go func(pbsLaunchChannel chan string) {
			for _ = range pbsLaunchChannel {
			}
		}(pbsLaunchChannel)
		go CreateStartScripts(cssInfo, machine, pbsLaunchChannel, done)

		if All {
			runs, runMap, mapErr := FindLastRound("*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt")
			log.Println("Selected to create start scripts for all the runs in the folder")
			log.Println("Found: ")
			for _, run := range runs {
				if mapErr != nil && len(runMap[run]["ics"]) == 0 {
					continue
				}
				fmt.Printf("%v\n", runMap[run]["ics"][len(runMap[run]["ics"])-1])
			}
			fmt.Println()
			// Fill the channel with the last round of each run
			for _, run := range runs {
				if mapErr != nil && len(runMap[run]["ics"]) == 0 {
					continue
				}
				cssInfo <- map[string]string{
					"remainingTime":  simTime,
					"randomSeed":     "",
					"newICsFileName": runMap[run]["ics"][len(runMap[run]["ics"])-1],
				}
			}

		} else {
			cssInfo <- map[string]string{
				"remainingTime":  simTime,
				"randomSeed":     randomNumber,
				"newICsFileName": inFileName,
			}
		}
		close(cssInfo)
		<-done
		close(done)
	},
}

// ***
var CutSimCmd = &cobra.Command{
	Use: "cutsim",
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

var stdOutCutCmd = &cobra.Command{
	Use:   "out",
	Short: "cut STDOUT",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		CutStdOut(inFileName, selectedSnapshot)
	},
}

var stdErrCutCmd = &cobra.Command{
	Use:   "err",
	Short: "cut STDERR",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		CutStdErr(inFileName, selectedSnapshot)
	},
}

var KiraWrapCmd = &cobra.Command{
	Use:   "kiraWrap",
	Short: "Wrapper for the kira integrator",
	Long: `Wrap the kira integrator providing
	environment monitoring.
	The "no-GPU" flag allow you to run the non GPU version 
	of kira if you installed kira-no-GPU in $HOME/bin/.
	Run with:
	
	kiraWrap (--no-GPU)
	
	You can also specify you want our modify version with Allen-Santillan 
	tidal field provided that you have that version of kira, named kiraTF in your
	~/bin/ folder. Run with 
	
	kiraWrap -f.`,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" || intTime == "" {
			log.Fatal("Provide an ICs file and the integration time.")
		}
		KiraWrap(inFileName, intTime, randomNumber, noGPU)
	},
}

// Out2ICsCmd creates new ICs from STDOUT to restart the simulation
var Out2ICsCmd = &cobra.Command{
	Use:   "out2ics",
	Short: "Prepare the new ICs from the last STDOUT",
	Long: `StarLab can restart a simulation from the last complete output.
	The out2ics command prepare the new ICs parsing the last STDOUT and writing
	the last complete snapshot to the new input file.
	Use like:
	sltools out2ics -i out-cineca-comb16-NCM10000-fPB005-W5-Z010-run06-rnd00.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cssInfo        = make(chan map[string]string, 1)
			inFileNameChan = make(chan string, 1)
		)
		if force {
			log.Println("Force to run even if end-of-simulation detected")
		}
		go Out2ICs(inFileNameChan, cssInfo)
		inFileNameChan <- inFileName
		close(inFileNameChan)
		<-cssInfo
	},
}

// ***
var PbsLaunchCmd = &cobra.Command{
	Use:   "pbsLaunch",
	Short: "submit all the PBS files in a folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := PbsLaunch(); err != nil {
			log.Fatal(err)
		}
	},
}

// ReadConfCmd load the JSON configuration file.
var ReadConfCmd = &cobra.Command{
	Use:   "readConf",
	Short: "Read and print the configuration file",
	Long: `Read and print the configuration specify by the -c flag.
	It must be in the form of a JSON file like:

	{
		"Runs": 50,
		"Comb": 60, 
		"Ncm" : 10000,
		"Fpb" : 0.10,
		"W"   : 5,
		"Z"   : 1.00,
		"Rv"  : 5, 
		"EndTime" : 500,
		"Machine" : "eurora",
		"UserName" : "bziosi00",
		"PName": "IscrC_SCmerge", 
		"BinFolder": "$HOME/bin/"
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

// ***
var ReLaunchCmd = &cobra.Command{
	Use: "relaunch",
	Short: `relaunch all the simulations in a folder, that is, clean the folder, 
	check, continue and submit. It also check for finished runs and put them in the Rounds folder. 
	If all the runs are finished, it writes a "complete" file.`,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Clean folder
		SimClean()

		if !goutils.Exists("complete") {
			// Check and continue
			CAC()

			// Submit: already included in CAC
			// 		if err := slt.PbsLaunch(); err != nil {
			// 			log.Fatal(err)
			// 		}
		} else {
			log.Println("'complete' file found, assume simulations are complete.")
		}
	},
}

var RestartFromHereCmd = &cobra.Command{
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
	It is YOUR responsibility to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},
}

var stdOutRestartCmd = &cobra.Command{
	Use:   "out",
	Short: "Prepare a pp3-stalled stdout to restart the simulation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RestartStdOut(inFileName, selectedSnapshot)
	},
}

var stdErrRestartCmd = &cobra.Command{
	Use:   "err",
	Short: "Prepare a pp3-stalled stderr so that it is synced with the stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RestartStdErr(inFileName, selectedSnapshot)
	},
}

// ***
var SimCleanCmd = &cobra.Command{
	Use:   "simClean",
	Short: "Clean the folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		SimClean()
	},
}

// ***
var SLRecompileCmd = &cobra.Command{
	Use:   "slrecompile",
	Short: "Recompile starlab",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var confArgs = []string{"--with-f77=no"}
		SLrecompile(confArgs)
	},
}

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
		if All {
			log.Println("Stich all!")
			StichThemAll(inFileName)
		} else {
			StichOutputSingle(inFileName)
		}
	},
}

// Init commands and attach flags
func InitCommands() {
	SlToolsCmd.PersistentFlags().BoolVarP(&All, "all", "A", false, "Run command on all the relevant files in the local folder")
	SlToolsCmd.PersistentFlags().StringVarP(&ConfName, "confName", "c", "", "Name of the JSON config file")
	SlToolsCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlToolsCmd.PersistentFlags().StringVarP(&endOfSimMyrString, "endOfSimMyr", "e", "", "Time in Myr to try to find the final timestep")
	SlToolsCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT from which to try to find the final timestep")
	SlToolsCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.AddCommand(CacCmd)

	SlToolsCmd.AddCommand(CheckEndCmd)

	SlToolsCmd.AddCommand(CheckSnapshotCmd)
	SlToolsCmd.AddCommand(CheckStatusCmd)

	SlToolsCmd.AddCommand(ComOrbitCmd)

	SlToolsCmd.AddCommand(ContinueCmd)
	ContinueCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")
	ContinueCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")

	SlToolsCmd.AddCommand(CreateICsCmd)
	CreateICsCmd.Flags().BoolVarP(&RunICC, "runIcc", "C", false, "Run the creation of the ICs instad of only create scripts")

	SlToolsCmd.AddCommand(CreateStartScriptsCmd)
	CreateStartScriptsCmd.PersistentFlags().BoolVarP(&as, "as", "a", false, "Run Allen-Santillan version of kira (debug strings).")
	CreateStartScriptsCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")
	CreateStartScriptsCmd.Flags().StringVarP(&randomNumber, "random", "r", "", "Init random seed provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&simTime, "simTime", "t", "500", "Remaining simulation time provided by the out2ics command")

	SlToolsCmd.AddCommand(CutSimCmd)
	CutSimCmd.AddCommand(stdOutCutCmd)
	CutSimCmd.AddCommand(stdErrCutCmd)
	CutSimCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cutTime", "t", "", "At which timestep stop")

	SlToolsCmd.AddCommand(KiraWrapCmd)
	KiraWrapCmd.PersistentFlags().BoolVarP(&as, "as", "a", false, "Run Allen-Santillan version of kira (debug strings).")
	KiraWrapCmd.PersistentFlags().BoolVarP(&noBinaries, "no-binaries", "b", false, "Switch off binary evolution.")
	KiraWrapCmd.PersistentFlags().BoolVarP(&tf, "tf", "f", false, "Run TF version of kira (debug strings).")
	KiraWrapCmd.PersistentFlags().BoolVarP(&noGPU, "no-GPU", "n", false, "Run without GPU support if kira-no-GPU installed in $HOME/bin/.")
	KiraWrapCmd.PersistentFlags().StringVarP(&randomNumber, "random", "s", "", "Random number.")
	KiraWrapCmd.PersistentFlags().StringVarP(&intTime, "time", "t", "", "Number of timestep to integrate before stop the simulation.")

	SlToolsCmd.AddCommand(Out2ICsCmd)
	Out2ICsCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "Last STDOUT to be used as input")
	Out2ICsCmd.Flags().BoolVarP(&force, "force", "f", false, "Disable end-of-simulaiton check")

	SlToolsCmd.AddCommand(PbsLaunchCmd)
	SlToolsCmd.AddCommand(ReadConfCmd)

	SlToolsCmd.AddCommand(ReLaunchCmd)
	ReLaunchCmd.PersistentFlags().BoolVarP(&as, "as", "a", false, "Run Allen-Santillan version of kira (debug strings).")

	SlToolsCmd.AddCommand(RestartFromHereCmd)
	RestartFromHereCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cutTime", "t", "", "At which timestep stop")
	RestartFromHereCmd.AddCommand(stdOutRestartCmd)
	RestartFromHereCmd.AddCommand(stdErrRestartCmd)

	SlToolsCmd.AddCommand(SimCleanCmd)
	SlToolsCmd.AddCommand(SLRecompileCmd)

	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().BoolVarP(&OnlyOut, "onlyOut", "O", false, "Only stich STDOUTs")
	StichOutputCmd.Flags().BoolVarP(&OnlyErr, "onlyErr", "E", false, "Only stich STDERRs")
}
