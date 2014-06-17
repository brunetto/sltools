## kiraWrap

kiraWrap is a wrapper for the kira integrator.
It runs kira in the way we usually do (tidal fields on request!:P) and you only 
need to specify the ICs file, the integration time and the random seed if present.
Your CLI should resemble:

````bash
kiraWrap <icsFileName> <integration time> <random seed if present>
````

`kiraWrap` assumes the kira binary is in `$HOME/bin/` for simplicity.


