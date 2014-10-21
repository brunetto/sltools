[![DOI](https://zenodo.org/badge/doi/10.5281/zenodo.12299.png)](http://dx.doi.org/10.5281/zenodo.12299) 

SlTools woud provide an easy way to do the most common operation needed to run StarLab simulations.
All the commands are available as standalone commands or as subcommands of the `sltools` one. 
For now I'm still work in progress and some things can be a bit confused!:P

## Usage

````
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
  comorbit                  Extract the center-of-mass coordinates from a STDOUT file.
  out2ics                   Prepare the new ICs from the last STDOUT
  pbsLaunch                 submit all the PBS files in a folder
  relaunch                  relaunch all the simulations in a folder, that is, clean the folder, 
        check, continue and submit. It also check for finished runs and put them in the Rounds folder. 
        If all the runs are finished, it writes a "complete" file.
  restartFromHere           Prepare a pp3-stalled simulation to be restarted
  simClean                  Clean the folder
  comorbit                  Extract the center-of-mass coordinates from a STDOUT file.
  slrecompile               Recompile starlab
  readConf                  Read and print the configuration file
  createICs                 Create ICs from the JSON config file.
  continue                  Prepare the new ICs and start scripts from all the last STDOUTs
  out2ics                   Prepare the new ICs from the last STDOUT
  css                       Prepare the scripts to start a run on a cluster with PBS 
  stichOutput               Stich output, only for one simulation or for all in the folder
  help [command]            Help about any command

 Available Flags:
  -A, --all=false: Run command on all the relevant files in the local folder
  -c, --confName="": Name of the JSON config file
  -d, --debug=false: Debug output
      --help=false: help for sltools
  -v, --verb=false: Verbose and persistent output

Use "sltools help [command]" for more information about that command.


````
