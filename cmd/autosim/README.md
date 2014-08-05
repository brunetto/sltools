# Autosim

Autosim try to manage the simulations for you.

Just run `autosim` and it will search for a file with the subfolder to visit, then it will:

* check for the queue to be empty
* read the folders to visit and enter them one after the other
* clean the folder
* check the status of the simulations
* remove bad runs
* if the run is bad for the second time, move it to quaratine (TODO)
* prepare the good runs for a new round
* submit the jobs
* check if there's an error submitting or the queue is full
* wait the queue to be empty
* 



