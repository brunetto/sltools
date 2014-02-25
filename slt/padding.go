package slt

import(
	"log"
	"strings"
)

// LeftPad returns the string padded filling remaining left spaces to `length` with `pad`.
func LeftPad(str, pad string, length int) (string) {
	if Debug {Whoami(true)}
	var repeat int
	if (length - len(str)) % len(pad) != 0 {
		log.Fatal("Can't pad ", str, " with ", pad, " to length ", length)
	} else {
		repeat = (length - len(str)) / len(pad)
	}
	return strings.Repeat(pad, repeat) + str
}