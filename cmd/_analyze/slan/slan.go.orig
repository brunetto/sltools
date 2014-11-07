package main

import (
	"time"

	"bitbucket.org/brunetto/slan"

	"github.com/brunetto/goutils/debug"
)

func main() {
	defer debug.TimeMe(time.Now())

	var ()

	slan.Debug = false

	data := slan.New()
// 	data.Populate("../test", "all_the_fishes.txt")
	data.Populate("/home/ziosi/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/data/2013-10-10-analysis/03-final", "all_the_fishes.txt")
// 	data.Populate("/home/ziosi/Dropbox/Research/PhD_Mapelli/1-DCOB_binaries/Analysis/data/2013-10-10-analysis/03-final", "bh-bh_all.txt")
	data.Stars.Print()

}
