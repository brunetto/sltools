


printf "\n"; pwd; printf "\n"; for (( c=0; c<=9; c++ )); do printf "$c "; ls -lah out-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; prStintf "  "; ls -lah err-*-run0$c-rnd0* | awk '{print $5"\t"$9}' | tail -n 1; printf "  "; cat $(ls err-*-run0$c-rnd0* | tail -n 1) | grep "Time = " | tail -n 1; done
