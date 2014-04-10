package slt

import (
	"log"
)

// AbsInt is a self-made absolute value function for int64.
func AbsInt(n int64) int64 {
	if n >= 0 {
		return n
	} else if n < 0 {
		return -n
	} else {
		log.Fatal("Problem determining the sign of ", n)
		return 0
	}
}
