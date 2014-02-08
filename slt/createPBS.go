package slt

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func CreatePBS (pbsOutName string, kiraOutName string, absFolderName string, run string, rnd string, conf *ConfigStruct) () {
	if Debug {Whoami(true)}
	
	var (
		pbsFile *os.File
		pbsWriter *bufio.Writer
		pbsString string
		err error
		modules string
	)
	
	if conf.Machine == "eurora" {
		modules = "module purge\n" +
					"module load profile/advanced\n" +
					"module load gnu/4.6.3\n" +
					"module load boost/1.53.0--gnu--4.6.3\n" +
					"module load cuda\n\n" +
					"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:" +
					"/cineca/prod/compilers/cuda/5.0.35/none/lib64:" +
					"/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
					"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:" +
					"/eurora/home/userexternal/mmapelli/\n\n"
	} else if conf.Machine == "plx" {
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
	} else {
		log.Fatal("Uknown machine name ", conf.Machine)
	}
	
	pbsString = "#!/bin/bash\n" +
				"#PBS -N r" + conf.CombStr() + "-" + run + "-" + rnd + "\n" +
				"#PBS -A " + conf.PName + "\n" +
				"#PBS -q longpar\n" +
				"#PBS -l walltime=24:00:00\n" +
				"#PBS -l select=1:ncpus=1:ngpus=2\n\n" +
				modules +
				"sh "+ filepath.Join(absFolderName, kiraOutName)
				
	log.Println("Write PBS launch script to ", pbsOutName)
	if pbsFile, err = os.Create(pbsOutName); err != nil {log.Fatal(err)}
	defer pbsFile.Close()
	
	pbsWriter = bufio.NewWriter(pbsFile)
	defer pbsWriter.Flush()
	
	pbsWriter.WriteString(pbsString)
}