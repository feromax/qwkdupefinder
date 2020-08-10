package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"qwkdupefinder/pkg/debug"
	"qwkdupefinder/pkg/utils"
	"strings"
)

//GetFileSize returns the size of a file whose path is provided as a string
func GetFileSize(f string) int64 {
	debug.ShowCallsAndArgs("GetFileSize", f)
	info, _ := os.Stat(f)
	return info.Size()
}

//IsFile returns true if the provided string is a path to a file
func IsFile(f string) (bool, error) {
	debug.ShowCallsAndArgs("IsFile", f)
	info, err := os.Stat(f)
	if err != nil {
		//may be a broken symlink; warn and continue
		utils.Warn(err.Error())
		return false, nil
	}
	retVal := info.Mode().IsRegular()
	debug.ShowRetvals(retVal)
	return retVal, nil
}

//IsDir returns true if the provided string is a path to a directory
func IsDir(f string) (bool, error) {
	debug.ShowCallsAndArgs("IsDir", f)
	info, err := os.Stat(f)
	if err != nil {
		return false, err
	}
	retVal := info.Mode().IsDir() && !IsLink(f)
	debug.ShowRetvals(retVal, nil)

	return retVal, nil
}

//IsLink returns true if the provided string is symlink
func IsLink(f string) bool {
	infoL, _ := os.Lstat(f)
	return infoL.Mode()&os.ModeSymlink == os.ModeSymlink
}

//IsFileOrDir returns true if the string is a valid dir or file path
func IsFileOrDir(x string) (bool, error) {
	debug.ShowCallsAndArgs("IsFileOrDir", x)
	isDir, err := IsDir(x)
	if err != nil {
		return false, err
	}
	isFile, err := IsFile(x)
	if err != nil {
		return false, err
	}

	return (isDir || isFile) && !IsLink(x), nil

}

//PrettySize converts a number of bytes into a human-readable
//size, appending the most appropriate unit (KB, MB, GB, or TB)
func PrettySize(nBytes int64) string {
	if nBytes < 1024 {
		return fmt.Sprintf("%dB", nBytes)
	}

	var unit string
	var sizeDbl float32 = float32(nBytes)

	if sizeDbl/1024 > 1 {
		unit = "KB"
		sizeDbl /= 1024
	}
	if sizeDbl/1024 > 1 {
		unit = "MB"
		sizeDbl /= 1024
	}
	if sizeDbl/1024 > 1 {
		unit = "GB"
		sizeDbl /= 1024
	}
	if sizeDbl/1024 > 1 {
		unit = "TB"
		sizeDbl /= 1024
	}
	return fmt.Sprintf("%.2f%s", sizeDbl, unit)

}

func stripOffMBSuffix(str string) string {
	if len(str) >= 2 {
		lastTwo := strings.ToUpper(str[len(str)-2:])
		if lastTwo == "MB" {
			str = str[:len(str)-2]
		}
	}
	return str
}

//GetListOfFiles returns a slice of filenames with
//given dirName is prefix
//PRECONDITION:  dirName is a valid directory
func GetListOfFiles(dirName string) ([]string, error) {
	debug.ShowCallsAndArgs("GetListOfFiles", dirName)

	var files []string
	items, err := ioutil.ReadDir(dirName)
	if err != nil {
		return files, err
	}

	// dirName might be "files/" or "." or "/"; in cases
	// where there's no trailing slash, add one
	if dirName[len(dirName)-1:] != "/" {
		dirName = fmt.Sprintf("%s/", dirName)
	}

	for _, item := range items {
		var itemName string
		itemName = fmt.Sprintf("%s%s", dirName, item.Name())

		isF, err := IsFile(itemName)
		if err != nil {
			return files, err
		} else if isF {
			files = append(files, itemName)
		}
	}
	debug.ShowRetvals(files, nil)

	return files, nil

}

//GetListOfDirs returns a slice of dirnames witih
//dirName as prefix
//PRECONDITION:  dirName is a valid directory
func GetListOfDirs(dirName string) ([]string, error) {
	debug.ShowCallsAndArgs("GetListOfDirs", dirName)
	var dirs []string

	// dirName might be "files/" or "." or "/"; in cases
	// where there's no trailing slash, add one
	if dirName[len(dirName)-1:] != "/" {
		dirName = fmt.Sprintf("%s/", dirName)
	}

	items, err := ioutil.ReadDir(dirName)
	if err != nil {
		return dirs, err
	}
	for _, item := range items {
		itemName := fmt.Sprintf("%s%s", dirName, item.Name())
		isD, err := IsDir(itemName)
		if err != nil {
			utils.Warn(err.Error())
		} else if isD {
			dirs = append(dirs, itemName)
		}
	}
	debug.ShowRetvals(dirs, nil)

	return dirs, nil

}
