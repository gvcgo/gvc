package sorts

import (
	"strconv"
	"strings"
)

type jVersion struct {
	VName string
	V     int
}

func (that *jVersion) Greater(item Item) bool {
	v, ok := item.(*jVersion)
	if !ok {
		panic("unknown item")
	}
	if that.V > v.V {
		return true
	}
	return false
}

func (that *jVersion) String() string {
	return that.VName
}

func SortJDKVersion(vList []string) []string {
	jList := make([]Item, 0)
	for _, v := range vList {
		strList := strings.Split(v, "k")
		if len(strList) < 2 {
			continue
		}
		vStr := strList[1]
		if strings.Contains(vStr, "-") {
			vStr = strings.Split(vStr, "-")[0]
		}
		item := &jVersion{}
		item.VName = v
		item.V, _ = strconv.Atoi(vStr)
		jList = append(jList, item)
	}
	return QuickSort(jList)
}
