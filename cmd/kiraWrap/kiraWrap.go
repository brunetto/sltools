package main

import (
	"log"
	"time"
	
	"github.com/spf13/cobra"
	
	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/sltools/slt"
)


var (
	noGPU, tf, as, noBinaries bool
	icsFileName string
	intTime string
	randomNumber string
)
	
var kiraWrapCmd = &cobra.Command{
	Use:   "kiraWrap",
	Short: "Wrapper for the kira integrator",
	Long: `Wrap the kira integrator providing
	environment monitoring.
	The "no-GPU" flag allow you to run the non GPU version 
	of kira if you installed kira-no-GPU in $HOME/bin/.
	Run with:
	
	kiraWrap (--no-GPU)`,
	Run: func(cmd *cobra.Command, args []string) {
		if icsFileName == "" || intTime == "" {
			log.Fatal("Provide an ICs file and the integration time.")
		}
		slt.KiraWrap(icsFileName, intTime, randomNumber, noGPU)
	},
}

func InitCommands() {
	kiraWrapCmd.PersistentFlags().BoolVarP(&noGPU, "no-GPU", "n", false, "Run without GPU support if kira-no-GPU installed in $HOME/bin/.")
	kiraWrapCmd.PersistentFlags().BoolVarP(&tf, "tf", "f", false, "Run TF version of kira (debug strings).")
	KiraWrapCmd.PersistentFlags().BoolVarP(&as, "as", "a", false, "Run Allen-Santillan version of kira (debug strings).")
	kiraWrapCmd.PersistentFlags().BoolVarP(&noBinaries, "no-binaries", "b", false, "Switch off binary evolution.")
	kiraWrapCmd.PersistentFlags().StringVarP(&icsFileName, "ics", "i", "", "ICs file to start with.")
	kiraWrapCmd.PersistentFlags().StringVarP(&intTime, "time", "t", "", "Number of timestep to integrate before stop the simulation.")
	kiraWrapCmd.PersistentFlags().StringVarP(&randomNumber, "random", "s", "", "Random number.")
}

func main () () {
	defer debug.TimeMe(time.Now())
	
	InitCommands()
	kiraWrapCmd.Execute()
}








