package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func PathIsDir(fPath string) (ok bool) {
	info, err := os.Stat(fPath)
	if err != nil {
		return
	}
	return info.IsDir()
}

func FileIsImage(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()
	_, _, err = image.Decode(file)
	return err == nil
}
