package slt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/brunetto/goutils"
	"github.com/brunetto/goutils/debug"
)

// CreateStartScripts create the start scripts (kira launch and PBS launch for the ICs).
func CreateStartScripts(cssInfo chan map[string]string, machine string, pbsLaunchChannel chan string, done chan struct{}) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		infoMap        map[string]string
		err            error    // common error container
		currentDir     string   // current local directory
		stdOutFile     string   // STDOUT file for the next run
		stdErrFile     string   // STDOUT file for the next run
		shortName      string   // id for the job
		queue          string   // name of the queue on wich we will run
		comb, run, rnd string   //combination, run and round number
		baseName       string   // common part of the name without the extension
		walltime       string   // max time we can run on the queue
		kiraString     string   // string to launch kira
		pbsString      string   // string to submit the job to PBS
		kiraFile       *os.File // where to save kiraString
		pbsFile        *os.File // where to save PBS string
		kiraOutName    string   // kira file name
		pbsOutName     string   // PBS file name
		home           string   // path to home on the cluster
		kiraBinPath    string   // path to kira binaries
		modules        string   // modules we need to load
		regRes         map[string]string
		randomString   string
		timeTest       int
		project        string
		tidalString string = ""
	)

	if home = os.Getenv("HOME"); home == "" {
		log.Fatal("Can't get $HOME variable and locate your home")
	}

	if !goutils.Exists(filepath.Join(home, "bin", "kiraWrap")) &&
		!goutils.Exists(filepath.Join(home, "bin", "kira")) {
		log.Fatal("Can't find kiraWrap or kira in ", filepath.Join(home, "bin"))
	}

	for infoMap = range cssInfo {
		if Debug {
			log.Println("Retrieved ", infoMap)
		}

		// empty map if no need to create css scripts
		if len(infoMap) == 0 {
			pbsLaunchChannel <- ""
			continue
		}

		if timeTest, err = strconv.Atoi(infoMap["remainingTime"]); err != nil {
			log.Fatal("Can't retrieve remaining time in createStartSCript: ", err)
		}

		if timeTest < 1 {
			log.Println("No need to create a new ICs, simulation complete.")
			continue
		}

		if regRes, err = Reg(infoMap["newICsFileName"]); err != nil {
			log.Fatal(err)
		}
		if regRes["prefix"] != "ics" {
			log.Fatalf("Please specify an ICs file, found %v prefix", regRes["prefix"])
		}

		if infoMap["randomSeed"] == "0" {
			randomString = ""
		} else if infoMap["randomSeed"] == "" {
			randomString = ""
		} else {
			randomString = "-s " + infoMap["randomSeed"]
		}

		baseName = regRes["baseName"]
		comb = regRes["comb"]
		run = regRes["run"]
		rnd = regRes["rnd"]

		shortName = "r" + comb + "-" + run + "-" + rnd

		if currentDir, err = os.Getwd(); err != nil {
			log.Fatal("Can't find current working folder!!")
		}
		if currentDir, err = filepath.Abs(currentDir); err != nil {
			log.Fatal("Can't find absolute path to current working folder!!")
		}

		stdOutFile = "out-" + baseName + "-run" + run + "-rnd" + rnd + ".txt"
		stdErrFile = "err-" + baseName + "-run" + run + "-rnd" + rnd + ".txt"
		kiraOutName = "kiraLaunch-" + baseName + "-run" + run + "-rnd" + rnd + ".sh"
		pbsOutName = "PBS-" + baseName + "-run" + run + "-rnd" + rnd + ".sh"

		if as {
				tidalString = " -a "
		}
		
		if machine == "eurora" {
			modules = "module purge\n" +
				"module load profile/advanced\n" +
				"module load gnu/4.6.3\n" +
				"module load boost/1.53.0--gnu--4.6.3\n" +
				"module load cuda\n\n" +
				"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:" +
				"/cineca/prod/compilers/cuda/5.0.35/none/lib64:" +
				"/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
				"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/eurora/home/userexternal/mmapelli/\n\n" +
				"LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/5.0.35/none/lib64:/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
				"export LD_LIBRARY_PATH\n"
			queue = "parallel"
			walltime = "4:00:00"
			project = "IscrC_SCmerge"			
			kiraString = "#echo $PWD\n" +
				"#echo $LD_LIBRARY_PATH\n" +
				"#echo $HOSTNAME\n" +
				"#date\n" +
				filepath.Join(home, "bin", "kiraWrap") + tidalString + " -i " +
				filepath.Join(currentDir, infoMap["newICsFileName"]) + " -t " +
				infoMap["remainingTime"] + " " +
				randomString + "\n"
			pbsString = "#!/bin/bash\n" +
				"#PBS -N r" + shortName + "\n" +
				"#PBS -A " + project + "\n" +
				"#PBS -q " + queue + "\n" +
				"#PBS -l walltime=" + walltime + "\n" +
				"#PBS -l select=1:ncpus=1:ngpus=2\n\n" +
				modules +
				"sh " + filepath.Join(currentDir, kiraOutName) + "\n"
			// 			kiraBinPath + " -t " + infoMap["remainingTime"] + " -d 1 -D 1 -b 1 -f 0 \\\n" +
			// 			" -n 10 -e 0.000 -B " + randomString + " \\\n" +
			// 			"<  " + filepath.Join(currentDir, icsName) + " \\\n" +
			// 			">  " + filepath.Join(currentDir, stdOutFile) + " \\\n" +
			// 			"2> " + filepath.Join(currentDir, stdErrFile) + " \n"
		} else if machine == "g2swin" {
			modules = "module purge\n" +
				// 					"module load profile/advanced\n" +
				"module load gcc/4.6.4\n" +
				"module load boost/x86_64/gnu/1.51.0-gcc4.6\n" +
				"module load cuda/4.0\n\n" //+
				// 					"LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/5.0.35/none/lib64:/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
				// 					"export LD_LIBRARY_PATH\n"
			queue = "gstar"
			walltime = "07:00:00:00"
			project = "p003_swin"

			kiraString = "#echo $PWD\n" +
				"#echo $LD_LIBRARY_PATH\n" +
				"#echo $HOSTNAME\n" +
				"#date\n" +
				filepath.Join(home, "bin", "kiraWrap") + tidalString + " -i " +
				filepath.Join(currentDir, infoMap["newICsFileName"]) + " -t " +
				infoMap["remainingTime"] + " " +
				randomString + "\n"
			pbsString = "#!/bin/bash\n" +
				"#PBS -N r" + shortName + "\n" +
				"#PBS -A " + project + "\n" +
				"#PBS -q " + queue + "\n" +
				"#PBS -l walltime=" + walltime + "\n" +
				"#PBS -l nodes=1:ppn=1:gpus=2\n\n" +
				modules +
				"sh " + filepath.Join(currentDir, kiraOutName) + "\n"
		} else if machine == "plx" {
			modules = "module purge\n" +
				"module load gnu/4.1.2\n" +
				"module load profile/advanced\n" +
				"module load boost/1.41.0--intel--11.1--binary\n" +
				"module load cuda/4.0\n\n" +
				"LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/4.0/none/lib64:" +
				"/cineca/prod/compilers/cuda/4.0/none/lib:/cineca/prod/" +
				"libraries/boost/1.41.0/intel--11.1--binary/lib:/cineca/" +
				"prod/compilers/intel/11.1/binary/lib/intel64\n" +
				"export LD_LIBRARY_PATH\n\n"
			queue = "longpar"
			walltime = "24:00:00"
			project = "IscrC_SCmerge"

			home = "/plx/userexternal/bziosi00"
			kiraBinPath = filepath.Join(home, "slpack", "starlab", "usr", "bin", "kira")
			kiraString = "echo $PWD\n" +
				"echo $LD_LIBRARY_PATH\n" +
				"echo $HOSTNAME\n" +
				"date\n" +
				kiraBinPath + " -t " + infoMap["remainingTime"] + " -d 1 -D 1 -b 1 -f 0 \\\n" +
				" -n 10 -e 0.000 -B " + randomString + " \\\n" +
				"<  " + filepath.Join(currentDir, infoMap["newICsFileName"]) + " \\\n" +
				">  " + filepath.Join(currentDir, stdOutFile) + " \\\n" +
				"2> " + filepath.Join(currentDir, stdErrFile) + " \n"
			pbsString = "#!/bin/bash\n" +
				"#PBS -N r" + shortName + "\n" +
				"#PBS -A " + project + "\n" +
				"#PBS -q " + queue + "\n" +
				"#PBS -l walltime=" + walltime + "\n" +
				"#PBS -l select=1:ncpus=1:ngpus=2\n\n" +
				modules +
				"sh " + filepath.Join(currentDir, kiraOutName) + "\n"
		} else {
			log.Fatal("Uknown machine name ", machine)
		}

		if kiraFile, err = os.Create(kiraOutName); err != nil {
			log.Fatal(err)
		}
		defer kiraFile.Close()
		fmt.Fprint(kiraFile, kiraString)

		if pbsFile, err = os.Create(pbsOutName); err != nil {
			log.Fatal(err)
		}
		defer pbsFile.Close()
		fmt.Fprint(pbsFile, pbsString)
		pbsLaunchChannel <- pbsOutName
	}
	// 	close(pbsLaunchChannel)
	done <- struct{}{}
}
