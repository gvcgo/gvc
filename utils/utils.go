package utils

import "os"

func PathIsDir(fPath string) (ok bool) {
	info, err := os.Stat(fPath)
	if err != nil {
		return
	}
	return info.IsDir()
}
