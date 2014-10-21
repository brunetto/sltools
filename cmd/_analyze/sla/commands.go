package sla

import (
	"fmt"

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
// showPromiscuous flags not tohide binaries containing a star already in another binary
var showPromiscuous bool

// SlPpCmd is the main command.
var SlPpCmd = &cobra.Command{
	Use:   "slpp",
	Short: "Tools for analysing StarLab simulations",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Choose a sub-command or type sltools help for help.")
	RunAll(inFileName)
	},
}

// VersionCmd print the sltpp version.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slpp",
	Long:  `All software has versions. This is sltpp's one.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("StarLab Analysis Tools v0.0")
	},
}

var (
	inFileName string
	dcob bool
)
/*
// ExchangesCmd compute the exchanges
var ExchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		Exchanges(inFileName)
	},
}

// LifeTimeCmd compute lifetimes
var LifeTimeCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		LifeTime(inFileName)
	},
}
*/
// Init commands and attach flags
func InitCommands() () {

	SlPpCmd.AddCommand(VersionCmd)
	SlPpCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	SlPpCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlPpCmd.Flags().StringVarP(&inFileName, "infile", "i", "", "Infile to parse for the exchange")
	SlPpCmd.Flags().BoolVarP(&dcob, "dcob", "D", false, "Work on the files listing only the DCOB life of binaries")
	SlPpCmd.Flags().BoolVarP(&showPromiscuous, "sprom", "s", false, "Not hide binaries containing a star already in another binary")
	
// 	SlPpCmd.AddCommand(ExchangesCmd)
	
// 	SlPpCmd.AddCommand(LifeTime)

}

