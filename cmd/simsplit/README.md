# simsplit

Split a simulation output in its snapshots.

## Usage

Use with

```
simsplit -t <STD type> -i <file name> [-e <single timestep to be extracted>]
```

like:

```
simsplit -t out -i out-comb20-TFno-Rv1-NCM10000-fPB010-W9-Z010-run00-all.txt.gz
simsplit -t err -i err-comb20-TFno-Rv1-NCM10000-fPB010-W9-Z010-run00-all.txt.gz
simsplit -t err -e 18 -i err-comb20-TFno-Rv1-NCM10000-fPB010-W9-Z010-run00-all.txt.gz
```

