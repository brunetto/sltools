I created two different binaries in case the names are not standard.    

### Restart from here

This is an attempt to overcome the problem of a stalled simulation 
because of the (binary integration stalling???).    

Running     

````bash
restartFromHere <STDOUT> <number of snapshot> 
````

you will obtain new ICs starting form <number of snapshot>.    

I recommend to cut the simulation at least one or two snapshot before it stalled.

### Cut STDERR

Running     

````bash 
cutStderr <STDERR> <number of snapshot>
````

you should obtain a new STDERR where the last snapshot is <number of snapshot>.