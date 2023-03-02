package sorts

import (
	"strconv"
	"strings"
)

type gVersion struct {
	Major  int
	Minor  int
	Patch  int
	Beta   int
	RC     int
	Origin string
}

func (that *gVersion) Greater(item Item) bool {
	v, ok := item.(*gVersion)
	if !ok {
		panic("unknown item")
	}
	if that.Major > v.Major {
		return true
	}
	if that.Major < v.Major {
		return false
	}
	if that.Minor > v.Minor {
		return true
	}
	if that.Minor < v.Minor {
		return false
	}
	if that.Patch > v.Patch {
		return true
	}
	if that.Patch < v.Patch {
		return false
	}
	if that.RC != v.RC {
		if (that.RC > v.RC && v.RC != 0) || (that.RC == 0 && that.Beta == 0) {
			return true
		}
	}
	if that.Beta != v.Beta {
		if (that.Beta > v.Beta && v.Beta != 0) || that.Beta == 0 {
			return true
		}
	}
	return false
}

func (that *gVersion) String() string {
	return that.Origin
}

func SortGoVersion(vs []string) []string {
	vList := []Item{}
	var vresult []string
	m := make(map[string]struct{}, 50)
	for _, v := range vs {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		vs_ := gVersion{}
		if strings.Contains(v, "beta") {
			result := strings.Split(v, "beta")
			vresult = strings.Split(result[0], ".")
			vs_.Beta, _ = strconv.Atoi(result[1])

		} else if strings.Contains(v, "rc") {
			result := strings.Split(v, "rc")
			vresult = strings.Split(result[0], ".")
			vs_.RC, _ = strconv.Atoi(result[1])
		} else {
			vresult = strings.Split(v, ".")

		}
		vs_.Major, _ = strconv.Atoi(vresult[0])
		switch len(vresult) {
		case 2:
			vs_.Minor, _ = strconv.Atoi(vresult[1])
		case 3:
			vs_.Minor, _ = strconv.Atoi(vresult[1])
			vs_.Patch, _ = strconv.Atoi(vresult[2])
		}
		vs_.Origin = v
		vList = append(vList, &vs_)
	}
	return QuickSort(vList)
}
