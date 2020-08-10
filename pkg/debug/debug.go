package debug

import "fmt"

//DEBUG global debug flag
var DEBUG bool = false

//SetDebug sets a flag for all packages to see
func SetDebug() {
	DEBUG = true
}

//ShowRetvals for debugging
func ShowRetvals(vals ...interface{}) {
	if DEBUG {
		printStr := "__DEBUG RETVALS: "
		for _, val := range vals {
			printStr = fmt.Sprintf("%s [ %v ] ", printStr, val)
		}
		fmt.Printf("%s\n", printStr)
	}
}

//ShowCallsAndArgs for debugging
func ShowCallsAndArgs(vals ...interface{}) {
	if DEBUG {
		printStr := ""
		for i, val := range vals {
			if i == 0 {
				printStr = fmt.Sprintf("\n%s__DEBUG\t%v()\n", printStr, val)
			} else {
				printStr = fmt.Sprintf("%s__DEBUG\t  %v\n", printStr, val)
			}
		}
		fmt.Printf(printStr)
	}

}
