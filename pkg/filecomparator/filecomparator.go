package filecomparator

import (
	"fmt"
	"math/rand"
	"os"
	"qwkdupefinder/pkg/checksum"
	"qwkdupefinder/pkg/debug"
	"qwkdupefinder/pkg/fileutils"
	"qwkdupefinder/pkg/utils"
	"strings"
)

var maxSamples int = 10

//tune up or down on basis of accuracy/speed
func getSampleSize(size int64) int {
	if size < 1048576 { //<1MB
		return 5120 //5KB
	} else if size < 20971520 { //<20MB
		return 10240 //10KB
	} else if size < 209715200 { //<200MB
		return 102400 //100KB
	} else { //200MB+
		return 512000 //500KB
	}
}

func getSampleMatchesHelper(results *[][]string, size int64, stillMatching []string, verbose bool, sampleSize int, samplesLeft int) {
	debug.ShowCallsAndArgs("getSampleMatchesHelper", *results, size, stillMatching, verbose, sampleSize, samplesLeft)

	//base cases for recursion
	if len(stillMatching) == 0 {
		return
	}

	if samplesLeft == 0 {
		if verbose {
			fmt.Printf("After sampling %s spread over %d locations, declaring a match:  %s\n",
				fileutils.PrettySize(int64(maxSamples*sampleSize)), maxSamples, strings.Join(stillMatching, ", "))
		}
		*results = append(*results, stillMatching)
		return
	}
	//end base cases for recursion

	//take a sample of all stillMatching files; any subsets that match sent recursively to this fcn
	chksumToFiles := make(map[string][]string)

	//where shall we sample?  anywhere from byte 0 to (size - sampleSize)
	start := int64(rand.Float64() * float64(size-int64(sampleSize)))
	if debug.DEBUG {
		fmt.Printf("File size=%d, seeking to %d with sample size %d (%d+%d=%d)\n",
			size, start, sampleSize, start, sampleSize, start+int64(sampleSize))
	}

	for _, f := range stillMatching {
		fHandle, err := os.Open(f)
		if err != nil {
			utils.Warn(err.Error())
		}
		defer fHandle.Close()

		fileBytes := make([]byte, sampleSize)
		_, err = fHandle.ReadAt(fileBytes, start)
		if err != nil {
			utils.Warn(err.Error())
		} else {
			chksum := checksum.ComputeSha1(fileBytes)
			chksumToFiles[chksum] = append(chksumToFiles[chksum], f)
			if debug.DEBUG {
				fmt.Printf("%s :: %s\n", chksum, f)
			}
		}
	}

	//files w/matching checksums move on together
	for _, files := range chksumToFiles {
		if len(files) > 1 {
			if debug.DEBUG {
				fmt.Printf("Matches on sample proceed further:  %s\n\n",
					strings.Join(files, ", "))
				exclusions := ""
				for _, f1 := range stillMatching {
					matched := false
					for _, f2 := range files {
						if f1 == f2 {
							matched = true
							break
						}
					}
					if !matched {
						exclusions += f1 + "; "
					}
				}
				fmt.Printf("EXCLUSIONS: %v\n", exclusions)
			}

			//dive, dive, dive!
			getSampleMatchesHelper(results, size, files, verbose, sampleSize, samplesLeft-1)
		}
	}

}

//GetSampleMatches is given a list of files w/the same size, and the
//size of the files, and returns a set of arrays with matching files.
//For example, if the set of files with the same size arrives as
// [ a, b, c, d, e ]
//and it's found that a=c and b=e, return
// [ [a,c], [b,e] ]
//
//If a file set has file size < 100KB, hash the entirety of the files
//rather than sampling.
func GetSampleMatches(size int64, files []string, verbose bool) (verified bool, matchSubsets [][]string) {
	if size < 102400 {
		return getSmallFileMatches(size, files, verbose)
	}

	sampleSize := getSampleSize(size)

	//in order to not get all samples up to the limit and then weeding out non matches to see which
	//subsets still match -- which is very bad for performance if we can often filter out nonmatches
	//on the basis of just one or two samples -- we'll employ some recursion to successively go further
	//into the allowed number of samples.  that'll allow us to compare after each sample only the files
	//that still match.
	//
	//the helper, during its recursion, will modify our return array directly since we're handing off
	//a pointer to matchSubsets.
	getSampleMatchesHelper(&matchSubsets, size, files, verbose, sampleSize, maxSamples)

	//false return value below indicates less than 100% certainty that
	//the files reported as duplicates really are
	return false, matchSubsets
}

func getSmallFileMatches(size int64, files []string, verbose bool) (verified bool, matchSubsets [][]string) {
	verified = true //small files are verified identical by exhaustive checksum comparison
	if verbose {
		fmt.Printf("Comparing contents of <100KB size files with same size:  %v\n",
			strings.Join(files, ", "))
	}

	hashToFiles := make(map[string][]string)
	for _, f := range files {
		fHandle, err := os.Open(f)
		if err != nil {
			utils.Warn(err.Error())
		}
		defer fHandle.Close()

		fileBytes := make([]byte, size)
		numRead, err := fHandle.Read(fileBytes)
		if err != nil {
			utils.Warn(err.Error())
		} else if int64(numRead) != size {
			utils.Warn(fmt.Sprintf("Number bytes read (%d bytes) from %s does not match the file size, %d.", numRead, f, size))
		} else {
			chksum := checksum.ComputeSha1(fileBytes)
			hashToFiles[chksum] = append(hashToFiles[chksum], f)
		}
	}

	//collate the matches
	for _, files := range hashToFiles {
		matchSubsets = append(matchSubsets, files)
	}
	return
}
