package util

import "os"

func PathExists(path string) (exists bool, isdir bool, err error) {
	f, err := os.Stat(path)
	if err == nil {
		return true, f.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, false, nil
	}
	return false, false, err
}
