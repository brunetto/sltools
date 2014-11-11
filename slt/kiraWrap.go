package slt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
		
	"github.com/capnm/sysinfo" // TODO: can be thrown using syscall!!
	
	
)

func KiraWrap(icsFileName, intTime, randomNumber string, noGPU bool) () {
	var (
		timeLimit string
		err error
		kiraString string
		kiraArgs []string
		kiraWrappedCmd *exec.Cmd
		pathName, icsName, outName, errName, ext string
		icsFile, outFile, errFile *os.File
		regRes map[string]string
		done = make(chan string, 1)
		randomSeed string = ""
		u  *user.User
		host, wd string
	)
	
	pathName = filepath.Dir(icsFileName)
	icsName = filepath.Base(icsFileName)
	timeLimit = intTime
	randomSeed = randomNumber
	
	fmt.Println("###################################################")
	
	// Extract fileNameBody, round and ext
	log.Println("Extract files names")
	ext = filepath.Ext(icsName)
	regRes, err = Reg(icsName)
	if err != nil {
		log.Println("Can't derive standard names from STDOUT => wrap it!!")
		errName = "err-" + icsName + ext
		outName = "out-" + icsName + ext
	} else {
		if regRes["prefix"] != "ics" {
			log.Fatalf("Please specify a STDIN file, found %v prefix", regRes["prefix"])
		}
		
		// Creating new filenames
		errName = "err-" + regRes["baseName"] + "-run" + regRes["run"] + "-rnd" + regRes["rnd"] + ext
		outName = "out-" + regRes["baseName"] + "-run" + regRes["run"] + "-rnd" + regRes["rnd"] + ext
	}
	
	if icsName, err = filepath.Abs(filepath.Join(pathName, icsName)); err != nil {
		log.Fatal("Error composing icsName: ", err)
	}
	if errName, err = filepath.Abs(filepath.Join(pathName, errName)); err != nil {
		log.Fatal("Error composing errName: ", err)
	}
	if outName, err = filepath.Abs(filepath.Join(pathName, outName)); err != nil {
		log.Fatal("Error composing outName: ", err)
	}
	
	if icsFile, err = os.Open(icsName); err != nil {log.Fatal(err)}
	defer icsFile.Close()
	if outFile, err = os.Create(outName); err != nil {log.Fatal(err)}
	defer outFile.Close()
	if errFile, err = os.Create(errName); err != nil {log.Fatal(err)}
	defer errFile.Close()
	
	errFile.WriteString("\n#==============================\n")	
	errFile.WriteString(fmt.Sprintf("\n#   %v Start with kiraWrap.\n", time.Now().Format(time.RFC850)))
	errFile.WriteString("\n#==============================\n")	
	
	outFile.WriteString("\n#==============================\n")	
	outFile.WriteString(fmt.Sprintf("\n#   %v Start with kiraWrap.\n", time.Now().Format(time.RFC850)))
	outFile.WriteString("\n#==============================\n")	
		
	
	kiraArgs =  []string{"-t", timeLimit,// +  // number of timesteps to compute
		"-d", "1",// +  // log output interval
		"-D", "1",// +  // snapshot interval
		// 				"-b", "1",// +  // frequency of full binary output
		"-f", "0",// +  // dynamical friction (0 = no friction, 1 = friction)
		"-n", "10",// +  // terminate if the cluster remains with only 10 particles
		"-e", "0.000",// + // softening 
		// 				"-B",// // switch on binary evolution
		//"-s 36543" // random seed 
	}
	
	if noGPU {
		log.Println("Selected the no GPU integration.")
		log.Println("Assuming kira is in $HOME/bin/kira-no-GPU, if not, please copy it there... for sake of simplicity!:P")
		kiraString = filepath.Join(os.Getenv("HOME"), "/bin/", "kira-no-GPU")
	} else if tf {
		log.Println("Selected TF versionn.")
		log.Println("Assuming kira is in $HOME/bin/kiraTF, if not, please copy it there... for sake of simplicity!:P")
		kiraString = filepath.Join(os.Getenv("HOME"), "/bin/", "kiraTF")
// 		kiraArgs = append(kiraArgs, "-F") // F now means dynamical friction
	} else if as {
		log.Println("Selected TF versionn.")
		log.Println("Assuming kira is in $HOME/bin/kira-AS, if not, please copy it there... for sake of simplicity!:P")
		kiraString = filepath.Join(os.Getenv("HOME"), "/bin/", "kira-AS")
	} else {
		log.Println("Assuming kira is in $HOME/bin/kira, if not, please copy it there... for sake of simplicity!:P")
		kiraString = filepath.Join(os.Getenv("HOME"), "/bin/", "kira")
	}
	
	if !noBinaries {
		log.Println("Binary evolution on, adding flags to kira")
		kiraArgs = append(kiraArgs, "-b", "1", "-B")
	} else {
		kiraArgs = append(kiraArgs, "-S")
		log.Println("Binary evolution off")
	}
	
	// Add the random seed if specified
	if randomSeed != "" {
		kiraArgs = append(kiraArgs, "-s", randomSeed)
	}
	
	kiraWrappedCmd = exec.Command(kiraString, kiraArgs...)
	if kiraWrappedCmd.Stdin = icsFile; err != nil {log.Fatal("Error connecting ICs to kira STDIN: ", err)}
	if kiraWrappedCmd.Stdout = outFile; err != nil {log.Fatal("Error connecting ICs to kira STDOUT: ", err)}
	if kiraWrappedCmd.Stderr = errFile; err != nil {log.Fatal("Error connecting ICs to kira STDERR: ", err)}
	
	if u, err = user.Current(); err != nil {
		log.Fatal("Can't retrieve username: ", err)
	}
	if host, err = os.Hostname(); err != nil {
		log.Fatal("Can't retrieve hostname: ", err)
	}
	if wd, err = os.Getwd(); err != nil {
		log.Fatal("Can't retrieve working dir: ", err)
	}
	
	log.Println("Run summary:")
	
	log.Printf("Username: %v (%v)\n", u.Username, u.Name)
	log.Println("Hostname: ", host)
	log.Println("Working dir: ", wd)
	log.Println("LD_LIBRARY_PATH: ", os.Getenv("LD_LIBRARY_PATH"))
	
	log.Println("Command: ", kiraString, kiraArgs)
	log.Println("STDIN = ", icsName)
	log.Println("STDOUT = ", outName)
	log.Println("STDERR = ", errName)
	
	log.Println("Ready... steady... Go!")
	
	if err = kiraWrappedCmd.Start(); err != nil {
		log.Fatal("Error starting kiraWrappedCmd: ", err)
	}
	
	log.Println("Waiting kira to finish while checking for problems...")
	// Wait for the process to end normally
	go waitProcess(kiraWrappedCmd, done)
	// Check for pp3-stalling situations
	go killTrigger(errName, kiraWrappedCmd, done)	
	
	if err = errors.New(<-done); err.Error() != "" {
		errFile.WriteString("\n"+err.Error()+"\n")
	}
	
	errFile.WriteString("\n#==============================\n")
	errFile.WriteString(fmt.Sprintf("\n#   %v Done with kiraWrap.\n", time.Now().Format(time.RFC850)))	
	errFile.WriteString(fmt.Sprintf("\n#   Username: %v (%V)\n", u.Username, u.Name))
	errFile.WriteString(fmt.Sprintf("\n#   Hostname: %v\n", host))
	errFile.WriteString("\n==============================\n")
	
	outFile.WriteString("\n#==============================\n")
	outFile.WriteString(fmt.Sprintf("\n#   %v Done with kiraWrap.\n", time.Now().Format(time.RFC850)))	
	outFile.WriteString("\n#==============================\n")
	
	fmt.Print("\x07") // Beep when finish!!:D
}

func killTrigger(errName string, kiraWrappedCmd *exec.Cmd, done chan string) () {
	const toGB = float64(1. / (1024*1024*1024))
	const maxStderrGB = float64(2)
	const minDiskGB = float64(5)
	const maxMemPerCent = float64(98)
	
	var (
		fileInfo os.FileInfo
		sysInfo *sysinfo.SI
		memAvail float64
		diskAvailGB float64 
		wd string 
		reason string
		err error
	)
	
	if wd, err = os.Getwd(); err != nil {
		log.Fatal("Can't retrieve local dir: ", err)
	}
	
	fsInfo := syscall.Statfs_t{}
	
	for {
		// Check STDERR file size
		if fileInfo, err = os.Stat(errName); err != nil {
			log.Fatal("Error checking STDERR file size, err")
		}
		// STDERR exceesing aloowed dimension of 2GB
		// probably the simulation is stalling because 
		// of pp3 locked on a binary
		if float64(fileInfo.Size()) * toGB > maxStderrGB {
			reason = "probable pp3 stalling"
			break
		} 
		
		// Check memory availability
		sysInfo = sysinfo.Get()
		memAvail = float64(sysInfo.FreeRam) / float64(sysInfo.TotalRam) * 100
		if memAvail > maxMemPerCent {
			reason = fmt.Sprintf("memory less than 95%% on the system: %2.2f",  memAvail)
			break
		}
		
		// Check HDD space availability
		if err = syscall.Statfs(wd, &fsInfo); err!= nil{
			log.Fatal("Cant't retrieve file system information: ", err)
		}
		if diskAvailGB = float64(fsInfo.Bavail) * float64(fsInfo.Frsize) * toGB; diskAvailGB < minDiskGB {
			reason = " available disk space less than 5 GB on the system"
			break
		}
		
		// Wait some time
		time.Sleep(time.Duration(2) * time.Minute)
	}
	log.Println("Kill kira because ", reason)
	if err := kiraWrappedCmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill: ", err)
	}
	done <- "killed because " + reason
}

func waitProcess(kiraWrappedCmd *exec.Cmd, done chan string) () {
	err := kiraWrappedCmd.Wait()
	log.Println("Process exited with error ", err)
	if err == nil {
		log.Println("Maybe the simulation is COMPLETE or the STD* are mixed and wrong.")
	}
	done <- ""	
}





