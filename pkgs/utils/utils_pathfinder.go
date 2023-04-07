package utils

import (
	"os"
	"path/filepath"
)

type DirFinder struct {
	StartDir       string
	ParentDirName  string
	UniqueFileName string
	path           string
}

func NewBinaryFinder(dir ...string) (bf *DirFinder) {
	bf = &DirFinder{}
	if len(dir) > 0 {
		bf.StartDir = dir[0]
	}
	if len(dir) > 1 {
		bf.ParentDirName = dir[1]
	}
	if len(dir) > 2 {
		bf.UniqueFileName = dir[2]
	}
	return
}

func (that *DirFinder) SetStartDir(dir string) {
	that.StartDir = dir
}

func (that *DirFinder) SetParentDirName(name string) {
	that.ParentDirName = name
}

func (that *DirFinder) SetUniqueFileName(name string) {
	that.UniqueFileName = name
}

func (that *DirFinder) String() string {
	that.getPath(that.StartDir)
	return that.path
}

func (that *DirFinder) getPath(dir string) {
	if rd, err := os.ReadDir(dir); err == nil {
		for _, d := range rd {
			if that.UniqueFileName != "" {
				if !d.IsDir() && d.Name() == that.UniqueFileName {
					that.path = dir
				} else if d.IsDir() {
					that.getPath(filepath.Join(dir, d.Name()))
				}
			} else {
				if d.IsDir() && d.Name() == that.ParentDirName {
					that.path = dir
				} else if d.IsDir() {
					that.getPath(filepath.Join(dir, d.Name()))
				}
			}
		}
	}
}
