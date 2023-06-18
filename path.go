package easyuse

import (
	"os"
)

func IsDir(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil {
		return fInfo.IsDir()
	} else {
		// os.IsExist(nil) 永远等于false， os.IsNotExist(nil)也是
		return os.IsExist(err) && fInfo.IsDir()
	}

}

func IsRegular(path string) bool {
	fInfo, err := os.Stat(path)
	if err == nil {
		return fInfo.Mode().IsRegular()
	} else {
		// os.IsExist(nil) 永远等于false, os.IsNotExist(nil)也是
		return os.IsExist(err) && fInfo.Mode().IsRegular()
	}
}
