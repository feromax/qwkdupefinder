package main

import (
	"fmt"
	"os"
	"qwkdupefinder/pkg/cmdparser"
	"qwkdupefinder/pkg/qwkdupefinder"
)

func main() {
	if len(os.Args) <= 1 {
		usage()
	} //exits

	items, sizeOnly, noRec, verbose, help, err := cmdparser.ParseCLIArgs(os.Args)
	if err != nil {
		fmt.Printf("Error:  %s\n\n", err.Error())
		fmt.Printf("For help:  %s -h\n\n", os.Args[0])
		os.Exit(1)
	}

	if help {
		usage()
	} //exits

	//in case of...bugs.
	//debug.SetDebug()

	err = qwkdupefinder.DoQwkDupeFind(items, sizeOnly, noRec, verbose)
	if err != nil {
		fmt.Printf("Error\n%s\n", err.Error())
		os.Exit(1)
	}

}

func usage() {
	cmd := os.Args[0]
	fmt.Printf(`
%s identifies likely duplicate files from the given files or directories by
comparing file sizes and then, for files with matching size, sampling a number of
1KB file segments and comparing their checksums.  This tool is well suited to finding
duplicate files in media collections, where 
(a) files tend to be large (and thus examining full file contents may take too long), and
(b) opening reported duplicate files and visually determining if the files are 
    actually duplicates is likely feasible.

%s is designed to operate quickly over large files where more exhaustive comparisons
done by reading entire files (e.g., when using hashing commands like md5sum and
sha[1,256,512]) may take too long.

This tool's REPORT MAY BE INACCURATE since only parts of files reporrted as duplicates
have been compared!  Use with care.

Last point:  Symbolic links are ignored.
___

USAGE:  %s [-h] [-v] [-n] [-s] { FILE | DIRECTORY } ...

SWITCHES
-h or --help
	shows this help

-v or --verbose
	provide detailed output

-n or --no-recursion
	do not compare files in subdirectories of directories provided

-s or --size-only
	for a faster report, only compare file sizes rather than also comparing portions
	of files that have the same size; THIS REDUCES THE ACCURACY OF THE RESULTS!

`, cmd, cmd, cmd)
	os.Exit(1)
}
