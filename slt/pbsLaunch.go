package slt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"time"
	
	"github.com/brunetto/goutils/debug"
)

func PbsLaunch () (error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		globName string = "PBS*.sh"
		regString string         = `PBS-\S*-run(\d*)-rnd(\d*)\.sh`
		regExp    *regexp.Regexp = regexp.MustCompile(regString)
		regRes []string
		inFiles  []string
		pbsFiles = map[string]string{}
		pbsFile string
		err error
		keys []string
		key string
		pbsCmd *exec.Cmd
		exists bool
		stdo, stde bytes.Buffer 
	)
	
	log.Println("Searching for files in the form: ", globName)
	if inFiles, err = filepath.Glob(globName); err != nil {
		log.Fatal("Error globbing files in this folder: ", err)
	}
	
	for _, pbsFile = range inFiles {
		if regRes = regExp.FindStringSubmatch(pbsFile); regRes == nil {
			log.Fatal("%v can't extract name info from %v")
		}
		if _, exists = pbsFiles[regRes[1]]; exists {
			log.Println("Found two PBS files ()round of the same run: ", regRes[1])
			log.Fatal("Be sure to delete PBS files form previous rounds")
		}
		pbsFiles[regRes[1]] = pbsFile
	}
	
	// Sort runs
	keys = make([]string, len(pbsFiles))
	idx := 0 
	for key, _ := range pbsFiles {
        keys[idx] = key
        idx++
    }
    sort.Strings(keys)
	
	for _, key = range keys {
		pbsCmd = exec.Command("qsub", pbsFiles[key])
		if pbsCmd.Stdout = &stdo; err != nil {log.Fatal("Error connecting STDOUT: ", err)}
		if pbsCmd.Stderr = &stde; err != nil {log.Fatal("Error connecting STDERR: ", err)}
		if err = pbsCmd.Start(); err != nil {
			log.Fatal("Error starting pbsCmd: ", err)
		}
		log.Println("Execute ", "qsub ", pbsFiles[key])
		if err = pbsCmd.Wait(); err != nil {
			log.Fatal("Error while waiting for pbsCmd: ", err)
		}
		fmt.Println(stdo.String())
		fmt.Println(stde.String())
		if stde.Len() != 0 {
			return errors.New(stde.String())
		}
		stdo.Reset()
		stde.Reset()
	}
	return nil
}

func PbsLaunchOnTheFly (pbsLaunchChannel chan string) (error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	
	var (
		pbsFile string
		err error
		pbsCmd *exec.Cmd
		stdo, stde bytes.Buffer 
	)
	
	for pbsFile = range pbsLaunchChannel {
		pbsCmd = exec.Command("qsub", pbsFile)
		if pbsCmd.Stdout = &stdo; err != nil {log.Fatal("Error connecting STDOUT: ", err)}
		if pbsCmd.Stderr = &stde; err != nil {log.Fatal("Error connecting STDERR: ", err)}
		if err = pbsCmd.Start(); err != nil {
			log.Fatal("Error starting pbsCmd: ", err)
		}
		log.Println("Execute ", "qsub ", pbsFile)
		if err = pbsCmd.Wait(); err != nil {
			log.Fatal("Error while waiting for pbsCmd: ", err)
		}
		fmt.Println(stdo.String())
		fmt.Println(stde.String())
		if stde.Len() != 0 {
			return errors.New(stde.String())
		}
		stdo.Reset()
		stde.Reset()
	}
	return nil
}


