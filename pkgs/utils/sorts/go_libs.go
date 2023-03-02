package sorts

import (
	"time"
)

const (
	ByUpdate   int = 1
	ByImported int = 0
)

type GoLibrary struct {
	Name     string
	Version  string
	Imported int
	Update   string
	UpdateAt time.Time
	SortType int
}

func (that *GoLibrary) String() string {
	return that.Name
}

func (that *GoLibrary) Greater(item Item) bool {
	v, ok := item.(*GoLibrary)
	if !ok {
		panic("unknown item")
	}
	if that.SortType == ByUpdate {
		return that.UpdateAt.After(v.UpdateAt)
	} else {
		if that.Imported > v.Imported {
			return true
		}
	}
	return false
}

func SortGoLibs(iList []Item) (r []*GoLibrary) {
	result := QSort(iList)
	for _, v := range result {
		if value, ok := v.(*GoLibrary); ok {
			r = append(r, value)
		}
	}
	return
}
