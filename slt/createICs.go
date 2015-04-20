package slt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"strconv"
	"time"

	"github.com/brunetto/goutils/debug"
)

// CreateAllICs will create the ICs for all the JSON config files found in this folder.
// Here I use sync.WaitGroup, another (older?) way is to use (from http://golang-examples.tumblr.com/tagged/goroutine)
// doSomething := make(chan int)
//     quit := make(chan int)
//
//     go func() {
//             select {
//             case <- doSomething:
//                     fmt.Println("done")
//             case <- quit:
//                     fmt.Println("quit")
//             }
//     }()
//
//     close(quit) // stop the goroutine
// //It’s better than sending a variable to ‘quit’ channel like,
//
// quit <- 1 // stop the goroutine
func CreateICsWrap(confName string, runICC bool) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		nProcs int = 1
		err       error
		confFiles []string
		combs     StringSet
		confChan = make(chan *ConfigStruct, 1)
		done = make(chan struct{})
	)
	
	for idx:=0; idx<nProcs; idx++ {
		go CreateICs(confChan, done)
	}
	
	// Check whether to create all the ICs or only one
	if confName == "all" || confName == "*" || 
		confName == "" || strings.Contains(confName, "*") {
		// Read all the JSON configuration files
		if confFiles, err = filepath.Glob("conf*.json"); err != nil {
			log.Fatal("Error globbing for stiching all the run outputs in this folder: ", err)
		}

		// Create a set (list of unique objs) from the conf names
		combs = NewStringSetFromSlice(confFiles)

		if Verb {
			log.Println("Found ", len(combs), " unique config files:")
			fmt.Println(combs.String())
		}

		// Read the conf files and launch the ICs creation
		for _, comb := range combs.Sorted() {
			if Verb {
				log.Println("Launching ICs creation based on ", comb)
			}
			conf := InitVars(comb)
			conf.RunICC = runICC
			confChan <- conf
		}
	} else {
		conf := InitVars(confName)
			conf.RunICC = runICC
			confChan <- conf
	}
	// Close channel otherwise goroutines will wait forever
	close(confChan)
	
	// Wait the goroutines to finish
	for idx:=0; idx<nProcs; idx++ {
		<-done // wait the goroutine to finish
	}
	
}

// CreateICs creates the ICs for a single run.
func CreateICs(confChan chan *ConfigStruct, doneParent chan struct{}) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
		
	// This variables are private to this function
	var (
		err              error
		folderName       string        // will contain the realizations of this combination
		icsBaseCmd       string        // common base command to create the ics
		icsCmd           string        // complete ICs creation command (contains the output file name)
		outIcsName       string        // final ICs name
		outIcsScriptName string        // name for the ICs creation script
		icsScriptFile    *os.File      // file obj for the ICs script
		icsScriptWriter  *bufio.Writer // writer for the ICs script
		outIcsFile       *os.File      // new ICs file
		nIcsWriter       *bufio.Writer // writer for the ICs file
		outIcsFileLog    *os.File      // new ICs file creation log
		nIcsWriterLog    *bufio.Writer // new ICs file creation log writer
		written          int           // written bytes
	)
		
	log.Println("Assuming kira is in $HOME/bin/kira, if not, please copy it there... for sake of simplicity!:P")
	
	for conf := range confChan {
		if Debug{log.Println("Retrieved: ", conf.FileName)}
		// Check we know where the binaries for the ICs are... not checking its existance
		if conf.BinFolder == "" {
			log.Fatal("I need to know where binaries for ICs are, no folder found in conf struct")
		}
		
		// FIXME: potebbe essere interessante usare http://godoc.org/labix.org/v2/pipe
		// oppure https://github.com/natefinch/sh
		// vedere anche http://play.golang.org/p/C608onvLWR

		// ICs binaries with path
		makeking := filepath.Join(/*os.Getenv("HOME"), "/bin/", */"makeking") //filepath.Join(conf.BinFolder, "makeking") 
		makemass := filepath.Join(/*os.Getenv("HOME"), "/bin/",*/ "makemass")
		makesecondary := filepath.Join(/*os.Getenv("HOME"), "/bin/",*/ "makesecondary")
		add_star := filepath.Join(/*os.Getenv("HOME"), "/bin/",*/ "add_star")
		scale := filepath.Join(/*os.Getenv("HOME"), "/bin/",*/ "scale")
		makebinary := filepath.Join(/*os.Getenv("HOME"), "/bin/",*/ "makebinary")

		// Base ICs script commands in a string, it misses the ICs filename
		icsBaseCmd = "#!/bin/bash\n" +
			"set -xeu\n" +
			makeking + " -n " + conf.NcmStr() +
			" -w " + conf.WStr() +
			" -i -u \\\n" +
			"| " + makemass + " -f 8  -l 0.1 -u 150 \\\n" +
			"| " + makesecondary + " -f " + conf.FpbStr() +
			" -q -l 0.1 \\\n" +
			"| " + add_star + " -R 1 -Z " + conf.ZStr() + " \\\n" +
			"| " + scale + " -R " + conf.RvStr() + " -M 1\\\n" +
			"| " + makebinary + " -f 2 -o 1 -l 1 -u 107836.09 \\\n" +
			"> " // Redirect output to the proper ICs file

		// Assemble folder name, create it and go into
// 		folderName = "cineca-comb" + conf.CombStr() +
// 			"-run1_" + conf.RunsStr() +
// 			"-NCM" + conf.NcmStr() +
// 			"-fPB" + conf.FpbCmpStr() +
// 			"-W" + conf.WStr() +
// 			"-Z" + conf.ZCmpStr()
		folderName = conf.BaseName()

		log.Println("Create folder and change to it:", folderName)
		if err = os.Mkdir(folderName, 0700); err != nil {
			log.Fatal("Can't create folder ", err)
		}

		// Copy config file inside folder to be read and for backup
		log.Println(fmt.Sprintf("Copy %v to %v", conf.FileName, filepath.Join(folderName, conf.FileName)))
		_, err = CopyFile(conf.FileName, filepath.Join(folderName, conf.FileName))
		// Go into the new folder
// 		if err = os.Chdir(folderName); err != nil {
// 			log.Println("Error while entering in folder ", folderName)
// 			log.Fatal("Can't cd into folder ", err)
// 		}

		// Start create start scripts
// 		go CreateStartScripts(cssInfo, conf.Machine, done)
		
		// Create the scripts
		for runIdx := 0; runIdx < conf.Runs; runIdx++ {
			/*
			* BASH SCRIPTS
			*/
			// Complete bash script with output file
			// Basename suffix
			runString := "-run" + LeftPad(strconv.Itoa(runIdx), "0", 2) + "-rnd00"
			// ICs final name
			outIcsName = "ics-" + conf.BaseName() + runString + ".txt\n"
			// Add ICs final file name to ICs creation command
			icsCmd = icsBaseCmd + outIcsName
			// ICs creation script name
			outIcsScriptName = "create_IC-" + conf.BaseName() + runString + ".sh"

			log.Println("Write ", outIcsScriptName)
			// Write the script file
			if icsScriptFile, err = os.Create(filepath.Join(folderName, outIcsScriptName)); err != nil {
				log.Fatal("Can't write to script file: ", err)
			}
			defer icsScriptFile.Close()
			icsScriptWriter = bufio.NewWriter(icsScriptFile)
			defer icsScriptWriter.Flush()
			if written, err = icsScriptWriter.WriteString(icsCmd); err != nil {
				log.Fatal("Error while writing ", outIcsScriptName, err)
			}
			icsScriptWriter.Flush()
			if Debug{log.Println("Written ", written, " bytes on ", outIcsScriptName)}

			// Create kiraLaunch and PBSlaunch scripts with the same functions used in Continue	
// 			cssInfo <- map[string]string{
// 					"remainingTime": "500",
// 					"randomSeed": "",
// 					"newICsFileName": fmt.Sprintf("%v", filepath.Join(folderName, "ics-"+conf.BaseName()+runString+".txt")),
// 			}
		}
		
// 		close(cssInfo)
// 		<-done // wait the goroutine to finish
		
		if conf.RunICC {
			log.Println("Also create ICs files running makeking etc")
			// Sometimes it crashes, untill I find why, I create the scripts
			// and the run the binaries only if -C flag is activated
			for runIdx := 0; runIdx < conf.Runs; runIdx++ {
				/*
				* ICs
				*/

				// Basename suffix
				runString := "-run" + LeftPad(strconv.Itoa(runIdx), "0", 2) + "-rnd00"
				// ICs final name
				outIcsName = "ics-" + conf.BaseName() + runString + ".txt"
				// Add ICs final file name to ICs creation command
				icsCmd = icsBaseCmd + outIcsName
				// ICs creation script name
				outIcsScriptName = "create_IC-" + conf.BaseName() + runString + ".sh"

				// REINIT PROCESSES BECAUSE EACH COMMAND IS A ONE-TIME CALL
				// Creating commands and pipes
				makekingCmd := exec.Command(makeking, "-n", conf.NcmStr(), "-w", conf.WStr(), "-i", "-u")
				makemassCmd := exec.Command(makemass, "-f", "8", "-l", "0.1", "-u", "150")
				makesecondaryCmd := exec.Command(makesecondary, "-f", conf.FpbStr(), "-q", "-l", "0.1")
				add_starCmd := exec.Command(add_star, "-R", conf.RvStr(), "-Z", conf.ZStr())
				scaleCmd := exec.Command(scale, "-R", "1", "-M", "1")
				makebinaryCmd := exec.Command(makebinary, "-f", "2", "-o", "1", "-l", "1", "-u", "107836.09")

				// makeking -> makemass
				if makemassCmd.Stdin, err = makekingCmd.StdoutPipe(); err != nil {
					log.Fatal("Create pipe to makemass: ", err)
				}
				// makemass -> makesecondary
				if makesecondaryCmd.Stdin, err = makemassCmd.StdoutPipe(); err != nil {
					log.Fatal("Create pipe to makesecondary: ", err)
				}
				// makesecondary -> add_star
				if add_starCmd.Stdin, err = makesecondaryCmd.StdoutPipe(); err != nil {
					log.Fatal("Create pipe to add_star: ", err)
				}
				// add_star -> scaleCmd
				if scaleCmd.Stdin, err = add_starCmd.StdoutPipe(); err != nil {
					log.Fatal("Create pipe to scale: ", err)
				}
				// scaleCmd -> makebinaryCmd
				if makebinaryCmd.Stdin, err = scaleCmd.StdoutPipe(); err != nil {
					log.Fatal("Create pipe to makebinary: ", err)
				}

				// Create ICs file and writer
				if outIcsFile, err = os.Create(filepath.Join(folderName, outIcsName)); err != nil {
					log.Fatal("Can't create ICs file: ", err)
				}
				defer outIcsFile.Close()
				nIcsWriter = bufio.NewWriter(outIcsFile)
				defer nIcsWriter.Flush()

				// Create ICs log file and writer
				if outIcsFileLog, err = os.Create(filepath.Join(folderName, "Create-" + outIcsName + ".log")); err != nil {
					log.Fatal("Can't create log file: ", err)
				}
				defer outIcsFileLog.Close()
				nIcsWriterLog = bufio.NewWriter(outIcsFileLog)
				defer nIcsWriterLog.Flush()

				makemassCmd.Stderr = nIcsWriterLog
				makesecondaryCmd.Stderr = nIcsWriterLog
				add_starCmd.Stderr = nIcsWriterLog
				scaleCmd.Stderr = nIcsWriterLog
				makebinaryCmd.Stderr = nIcsWriterLog

				// Attach the file writer to the cmd stdout
				makebinaryCmd.Stdout = nIcsWriter

				log.Println("Starting the creation of ", outIcsName)
				if err = makekingCmd.Start(); err != nil {
					log.Fatal("Start makeking: ", err)
				}
				if err = makekingCmd.Wait(); err != nil {
					log.Fatal("Wait makeking: ", err)
				}
				if err = makemassCmd.Start(); err != nil {
					log.Fatal("Start makemass: ", err)
				}
				if err = makemassCmd.Wait(); err != nil {
					log.Fatal("Wait makemass: ", err)
				}

				if err = makesecondaryCmd.Start(); err != nil {
					log.Fatal("Start makesecondary: ", err)
				}
				if err = makesecondaryCmd.Wait(); err != nil {
					log.Fatal("Wait makesecondary: ", err)
				}

				if err = add_starCmd.Start(); err != nil {
					log.Fatal("Start add_star: ", err)
				}
				if err = add_starCmd.Wait(); err != nil {
					log.Fatal("Wait add_Star: ", err)
				}

				if err = scaleCmd.Start(); err != nil {
					log.Fatal("Start scale: ", err)
				}
				if err = scaleCmd.Wait(); err != nil {
					log.Fatal("Wait scale: ", err)
				}

				if err = makebinaryCmd.Start(); err != nil {
					log.Fatal("Start makebinary: ", err)
				}
				if err = makebinaryCmd.Wait(); err != nil {
					log.Fatal("Wait makebinary: ", err)
				}

				nIcsWriter.Flush()
				nIcsWriterLog.Flush()

				/*
					In case of problems, Dave Cheney suggest
					(https://groups.google.com/d/msg/golang-nuts/pBa-6ywQE8c/V9JOsXMENrAJ)
					to lock the log while writing

					type W struct {
						w io.Writer
						sync.Mutex
					}

					func (w *W) Write(buf []byte) (int, error) {
						w.Lock()
						defer w.Unlock()
						return w.w.Write(buf)
					}
				*/

				log.Println("Wrote ", outIcsName)
			}
		} else {
			fmt.Println()
			log.Println("Created only ICs scripts")
			fmt.Println()
		}
		
		// Go back to the parent folder
// 		localFolder, err := os.Getwd()
// 		if err != nil {log.Fatal("Error while getting local folder: ", err)}
// 		
// 		parentFolder := filepath.Dir(localFolder)
// 		if err = os.Chdir(parentFolder); err != nil {
// 			log.Println("Error while coming back to the parent folder ", parentFolder)
// 			log.Fatal("Can't cd into folder ", err)
// 		}
	}
	doneParent <- struct{}{}
}
