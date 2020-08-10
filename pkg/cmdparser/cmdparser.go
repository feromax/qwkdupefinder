package cmdparser

import (
	"errors"
	"fmt"
	"qwkdupefinder/pkg/debug"
	"strings"
)

//if the arg is a case-insensitive match against at least
//one of the switches, return true
func isMatch(arg string, switches []string) bool {
	arg = strings.ToLower(arg)
	for _, sw := range switches {
		if arg == strings.ToLower(sw) {
			return true
		}
	}
	return false
}

//ParseCLIArgs looks for the following --
//* '-n' or '--no-recurse'
//* '-s' or '--size-only'
//* '-h' or '--help'
//* '-v' or '--verbose'
//* one or more files/dirs
//-- and returns their values, returning an error
//if the args weren't pleasing enough.
//(In case of -h or --h, never return error, just show the help.)
//PRECONDITION: 2 or more values in os.[]Args (where 1st is name of executable)
func ParseCLIArgs(args []string) (items []string, sizeOnly bool, noRecurse bool,
	verbose bool, help bool, err error) {
	//defaults
	sizeOnly = false
	noRecurse = false
	verbose = false
	help = false

	//omit 0th arg, the name of executable
	args = args[1:]

	for len(args) > 0 {
		cur := args[0]
		if isMatch(cur, []string{"-s", "--size-only"}) {
			sizeOnly = true
		} else if isMatch(cur, []string{"-n", "--no-recurse"}) {
			noRecurse = true
		} else if isMatch(cur, []string{"-v", "--verbose"}) {
			verbose = true
		} else if isMatch(cur, []string{"-h", "--help"}) {
			help = true
			return
		} else if isMatch(cur, []string{"--debug-hidden"}) {
			debug.SetDebug()
		} else {
			//this should be the final arg(s) that are only
			//files or directories
			for len(args) > 0 {
				cur = args[0]
				//fmt.Printf("\t\t%s", cur)
				if strings.Index(cur, "-") == 0 {
					//a switch that didn't match the expected ones?
					//ONLY OK for asking for help
					if isMatch(cur, []string{"-h", "--help"}) {
						help = true
						return
					}
					err = fmt.Errorf("bad switch or location (%s)", cur)
					return
				}
				items = append(items, cur)
				args = args[1:]
			}
			return
		}

		//after processing a commandline switch, strip off
		//the switch arg and proceed to next arg
		args = args[1:]
	}

	//won't reach here in case help was requested
	if len(items) == 0 {
		err = errors.New("no file or directory specified")
	}

	return
}
