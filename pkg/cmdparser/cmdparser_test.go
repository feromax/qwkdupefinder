package cmdparser

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
)

type cliTest struct {
	argStr    string
	exItems   []string
	exSzOnly  bool
	exNoRec   bool
	exVerbose bool
	exHelp    bool
	exErr     error
}

func arrContentsMatch(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)

	for i, aVal := range a {
		bVal := b[i]
		if aVal != bVal {
			return false
		}
	}
	return true
}

func TestParseCLIArgs(t *testing.T) {
	filesArrs := [][]string{
		nil,                        //0
		{"/"},                      //1
		{"apple"},                  //2
		{"/", "/tmp"},              //3
		{"apple", "./apple/banan"}, //4
	}

	cliTests := []cliTest{
		//                              []files, szOnly, noRec, verb, help, err

		//no files
		//{"./cmd_0", filesArrs[0], ...},  -- violates precondition
		{"./cmd_1 -h", filesArrs[0], false, false, false, true, nil}, //get help

		//one file/dir
		{"./cmd_2 FILES -h", filesArrs[1], false, false, false, true, nil},                 //get help, never error
		{"./cmd_3 -h FILES", filesArrs[2], false, false, false, true, nil},                 //get help, never error
		{"./cmd_4 --help FILES", filesArrs[2], false, false, false, true, nil},             //get help, never error
		{"./cmd_5 FILES", filesArrs[2], false, false, false, false, nil},                   //single file arg
		{"./cmd_6 FILES", filesArrs[3], false, false, false, false, nil},                   //multiple files
		{"./cmd_7 -z FILES", filesArrs[3], false, false, false, false, errors.New("x")},    //multiple files, invalid switch
		{"./cmd_8 FILES --z", filesArrs[4], false, false, false, false, errors.New("x")},   //multiple files, invalid switch in bad position
		{"./cmd_9 -n FILES -n", filesArrs[2], false, false, false, false, errors.New("x")}, //good switch, bad position
		{"./cmd_10 -n FILES", filesArrs[2], false, true, false, false, nil},                //good switch, good position
		{"./cmd_11 --no-recurse FILES", filesArrs[3], false, true, false, false, nil},      //good switch, good position
		{"./cmd_12 --size-only -n FILES", filesArrs[3], true, true, false, false, nil},     //good multi switch, good positions
		{"./cmd_13 -v -n FILES", filesArrs[4], false, true, true, false, nil},              //good multi switch, good positions
		{"./cmd_14 -s --verbose FILES", filesArrs[3], true, false, true, false, nil},       //good multi switch, good positions
	}

	for _, cliTest := range cliTests {
		filesAsStr := strings.Join(cliTest.exItems, " ")
		fmt.Println(filesAsStr)

		//template FILES using filesAsStr
		fullArgs := cliTest.argStr
		fullArgs = strings.ReplaceAll(fullArgs, "FILES", filesAsStr)

		actItems, actSizeOnly, actNoRec, actVerbose, actHelp, actErr := ParseCLIArgs(strings.Split(fullArgs, " "))
		actErrStr := "<nil>"
		if actErr != nil {
			actErrStr = actErr.Error()
		}
		t.Logf("args=\"%s\"\tActuals: items=%v, sz-only=%t, noRec=%t, v=%t, help=%t, err=%s",
			fullArgs, actItems, actSizeOnly, actNoRec, actVerbose, actHelp, actErrStr)

		if cliTest.exHelp && actHelp {
			//ok here
		} else if cliTest.exHelp != actHelp {
			t.Logf("Help called for when not expected, or vice versa")
			t.Fail()
		} else if cliTest.exErr != nil && actErr == nil || cliTest.exErr == nil && actErr != nil {
			t.Logf("Error expected and not returned, or vice-versa")
			t.Fail()
		} else if actErr != nil && cliTest.exErr != nil {
			//error as expected, ok
		} else if !arrContentsMatch(cliTest.exItems, actItems) {
			t.Logf("Filename not parsed correctly")
			t.Fail()
		} else if cliTest.exNoRec != actNoRec {
			t.Logf("recursion not parsed correctly")
			t.Fail()
		} else if cliTest.exVerbose != actVerbose {
			t.Logf("verbose not parsed correctly")
			t.Fail()
		} else if cliTest.exSzOnly != actSizeOnly {
			t.Logf("size-only not parsed correctly")
			t.Fail()
		}

	}

}
