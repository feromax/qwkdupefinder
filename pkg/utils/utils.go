package utils

import (
	"os"
	"qwkdupefinder/pkg/debug"
)

//count of times a message was displayed
var oldMsgs map[string]int

func alreadyWarned(s string) bool {
	if oldMsgs == nil {
		oldMsgs = make(map[string]int)
	}

	count := oldMsgs[s]
	oldMsgs[s]++
	if count > 0 {
		return true
	}
	return false
}

//Warn prints a message to stderr while squelching duplicates
//(unless we're debugging)
func Warn(s string) {
	if !alreadyWarned(s) || debug.DEBUG {
		os.Stderr.WriteString("[WARN] " + s + "\n")
	}
}
