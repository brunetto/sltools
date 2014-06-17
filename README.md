SlTools woud provide an easy way to do the most common operation needed to run StarLab simulations.

## Usage

````bash
Usage: 
  sltools [flags]
  sltools [command]

Available Commands: 
  version         :: Print the version number of slt
  readConf        :: Read and print the configuration file
  createICs       :: Create ICs
  continue        :: Prepare the new ICs from all the last STDOUTs
  out2ics         :: Prepare the new ICs from the last STDOUT
  createStartScripts :: Prepare the new ICs from all the last STDOUTs
  stichOutput     :: Only download SL
  help [command]  :: Help about any command

 Available Flags:
  -c, --confName="": Name of the JSON config file
  -d, --debug=false: Debug output
  -v, --verb=false: Verbose and persistent output

Use "sltools help [command]" for more information about that command.

````