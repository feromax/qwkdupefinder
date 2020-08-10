package qwkdupefinder

import (
	"fmt"
	"qwkdupefinder/pkg/debug"
	"qwkdupefinder/pkg/filecomparator"
	"qwkdupefinder/pkg/fileutils"
	"qwkdupefinder/pkg/utils"
)

//track files w/same size in bytes
var sizeToFiles map[int64][]string = make(map[int64][]string)
var matches int

//doFsCrawl does the initial traversal of directories, tracking file sizes
//and paths for subsequent analysis
func doFsCrawl(items []string, noRec bool, verbose bool) error {
	for _, item := range items {
		isForD, err := fileutils.IsFileOrDir(item)
		if err != nil {
			utils.Warn(err.Error())
		} else if !isForD {
			if fileutils.IsLink(item) {
				utils.Warn(item + " is a symlink, SKIPPING.")
			} else {
				utils.Warn(item + " is is not a valid file or directory, SKIPPING.")
			}
		}

		if isF, _ := fileutils.IsFile(item); isF {
			storeSize(item, verbose)
		} else {
			//name is a dir
			//get list of files, process them...
			files, err := fileutils.GetListOfFiles(item)

			if err != nil {
				utils.Warn("Cannot read contents of directory " + item + ", SKIPPING.")
				continue
			}
			for _, file := range files {
				storeSize(file, verbose)
			}

			if !noRec { //i.e., recurse (default)
				dirNames, err := fileutils.GetListOfDirs(item)
				if err != nil {
					return err
				}
				err = doFsCrawl(dirNames, noRec, verbose)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil

}

func printMatches(verified bool, size int64, files []string) {
	matches++
	verStr := "✓"
	if !verified {
		verStr = " "
	}
	for _, f := range files {
		fmt.Printf("MATCH:%d %s\tSIZE:%d\t%q\n", matches, verStr, size, f)
	}
	fmt.Println()
}

//once sizeToFiles has been populated, produce a smaller map
//with sizes that are repeated, referring to arrays of files len > 1
func getSizeMatches(szOnly bool, verbose bool) (szMatches map[int64][]string) {
	if szOnly {
		fmt.Printf(`
=======================================================
DUPLICATES REPORT (ON BASIS OF MATCHING FILESIZE, ONLY)
=======================================================
`)
	} else if verbose {
		fmt.Printf("\n[Interim size-only match report ]\n")
	}

	szMatches = make(map[int64][]string)
	for size, files := range sizeToFiles {
		if len(files) > 1 {
			szMatches[size] = files
			if verbose || szOnly {
				printMatches(false, size, files)
			}
		}
	}
	return szMatches

}

//DoQwkDupeFind takes a list of file/dir names and whether to recursively
//operate on subdirs.  Any errors returned to this function are
//returned to the caller after having exhausted all possible files/dirs
//for inspection.
func DoQwkDupeFind(items []string, sizeOnly bool, noRec bool, verbose bool) error {
	debug.ShowCallsAndArgs("DoQwkDupeFind", items, sizeOnly, noRec, verbose)

	err := doFsCrawl(items, noRec, verbose)
	if err != nil {
		utils.Warn(err.Error())
	}

	szMatchMap := getSizeMatches(sizeOnly, verbose)

	if !sizeOnly {
		fmt.Printf(`
=================================================================================================================
DUPLICATES REPORT ON BASIS OF FILESIZE AND SAMPLING CONTENTS -- ✓ indicates verified duplicate (for files <100KB)
=================================================================================================================
`)
		matches = 0
		for size, files := range szMatchMap {
			verified, matchSets := filecomparator.GetSampleMatches(size, files, verbose)

			for _, files := range matchSets {
				if len(files) > 1 {
					printMatches(verified, size, files)
				}
			}
		}

	}
	return nil
}

func storeSize(fName string, verbose bool) {
	debug.ShowCallsAndArgs("storeSize", fName, verbose)
	size := fileutils.GetFileSize(fName)
	sizeToFiles[size] = append(sizeToFiles[size], fName)
	if verbose {
		fmt.Printf("%s [size %s]", fName, fileutils.PrettySize(size))
		matches := len(sizeToFiles[size])
		if matches > 1 {
			fmt.Printf("...%d size match(es) found so far\n", matches)
		} else {
			fmt.Println()
		}
	}
}
