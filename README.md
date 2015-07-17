# Sltools

Release 1.1 [![DOI](https://zenodo.org/badge/doi/10.5281/zenodo.12299.png)](http://dx.doi.org/10.5281/zenodo.12299)    
Release 1.2 [![DOI](https://zenodo.org/badge/doi/10.5281/zenodo.18980.png)](http://dx.doi.org/10.5281/zenodo.18980) 

SlTools incorporates several piece of software I developed to manage and analyse simulations for my [PhD project](http://brunettoziosi.eu/pages/research/phd/). Most of the tools have been developed to deal with  [StarLab](http://www.sns.ias.edu/~starlab/index.html) problems (struggling with the installation? try have a look at [this](http://http://brunettoziosi.eu/posts/starlab-gpu-installation/)) and Cineca clusters. Some other try to automatize some tasks of running hundreds of simulations. There are also some tools I used to analyze the data (in Python, I am slowly rewriting them in Go the embed them in sltools).

See also [this post](http://brunettoziosi.eu/posts/starlab-simulations-with-sltools) and [this guide on how to create ICs with Dockerized StarLab](http://brunettoziosi.eu/posts/dockerized-starlab).

The dev repo is [here](https://github.com/brunetto/sltools-dev).

---

SlTools woud provide an easy way to do the most common operation needed to run StarLab simulations.
All the commands are available as standalone commands or as subcommands of the `sltools` one. 
For now I'm still work in progress and some things can be a bit confused!:P

## Usage

````
SlTools would help in running simulations with StarLab.
It can create the inital conditions if StarLab is compiled and the 
necessary binaries are available.
SlTools can also prepare ICs from the last snapshot and stich the 
output.

Usage: 
  sltools [flags]
  sltools [command]

Available Commands: 
  version                   Print the version number of slt
  cac                       Check and continue, will check the last simulations outputs, prepare the restat and restart.
  checkEnd                  Check the number of timesteps necessary to reach a given time in Myr. Need the files to have standard names.
  checkSnapshot             Check the snapshot for being OK.
  checkStatus               Check the status of a folder of simulations.
  comorbit                  Extract the center-of-mass coordinates from a STDOUT file.
  continue                  Prepare the new ICs and start scripts from all the last STDOUTs
  createICs                 Create ICs from the JSON config file.
  css                       Prepare the scripts to start a run on a cluster with PBS 
  cutsim                    Shorten a give snapshot to a certain timestep
        Because I don't now how perverted names you gave to your files, 
        you need to fix the STDOUT and STDERR by your own.
        You can do this by running 

        cutsim out --inFile <STDOUT file> --cut <snapshot where to cut>
        cutsim err --inFile <STDERR file> --cut <snapshot where to cut>

        The old STDERR will be saved as STDERR.bck, check it and then delete it.
        It is YOUR responsible to provide the same snapshot name to the two subcommands
        AND I suggest you to cut the simulation few timestep before it stalled.
  kiraWrap                  Wrapper for the kira integrator
  out2ics                   Prepare the new ICs from the last STDOUT
  pbsLaunch                 submit all the PBS files in a folder
  readConf                  Read and print the configuration file
  relaunch                  relaunch all the simulations in a folder, that is, clean the folder, 
        check, continue and submit. It also check for finished runs and put them in the Rounds folder. 
        If all the runs are finished, it writes a "complete" file.
  restartFromHere           Prepare a pp3-stalled simulation to be restarted
  simClean                  Clean the folder
  slrecompile               Recompile starlab
  stichOutput               Stich output, only for one simulation or for all in the folder
  help [command]            Help about any command

 Available Flags:
  -A, --all=false: Run command on all the relevant files in the local folder
  -c, --confName="": Name of the JSON config file
  -d, --debug=false: Debug output
  -e, --endOfSimMyr="": Time in Myr to try to find the final timestep
  -h, --help=false: help for sltools
  -i, --inFile="": STDOUT from which to try to find the final timestep
  -v, --verb=false: Verbose and persistent output

Use "sltools help [command]" for more information about that command.

````

## Example of use
### Configuration file

Create one or more `json` configuration files for your runs, like:

```bash
$ cat conf86.json 
{
    "Runs": 10,
    "Comb": 86,
    "Ncm" : 10000,
    "Fpb" : 0.05,
    "W"   : 3,
    "Z"   : 0.10,
    "Rv": 5,
    "EndTime" : 500,
    "Machine" : "g2",
    "UserName" : "bziosi",
    "PName": <name of the project if you use PBS or similar>,
    "BinFolder": "/home/bziosi/bin/"
}
```

Run

```bash
sltools createICs -v -A
```

to generate, for each configuration file, a folder containing the script to create the initial conditions.

Run the script (unless you ran `createICs` with `-C` flags (not working very well):

```bash
for ICscript in $(ls create*); do bash $ICscript; done;
```

Now you can transfer you fresh ICs to the machine where you will run your jobs.

It is now time to create the scripts to submit the jobs:

```bash
sltools css -A -m <machine>
```

At the moment the available machines are `eurora`, `g2swin`, `plx`.

Submit all the PBS scripts with 

```bash
sltools pbsLaunch
```

and wait for them to complete.

Relaunch them with 

```
sltools relaunch
```

if necessary.

At the end run 

```bash
sltools stichOutput -A
```

to have your outputs files nicely checked and stiched.

## TODO

* Check errors in [stichOutput](https://github.com/brunetto/sltools/blob/master/slt/stichOutput.go), [readSnapshot](https://github.com/brunetto/sltools/blob/master/slt/readSnapshot.go) and all dependants, it seems that returning the error is not useful since I halt the program if the snapshot has problems. 









