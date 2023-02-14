package utils

import (
	"strconv"
	"strings"
)

type version struct {
	Major  int
	Minor  int
	Patch  int
	Beta   int
	RC     int
	Origin string
}

func (that *version) Greater(v *version) bool {
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

type VComparator struct {
	Versions []string
	v        []*version
}

func NewVComparator(vs []string) *VComparator {
	vList := []*version{}
	var vresult []string
	m := make(map[string]struct{}, 50)
	for _, v := range vs {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		vs_ := version{}
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

	return &VComparator{Versions: make([]string, 0), v: vList}
}

func (that *VComparator) sort(vList []*version) (r []*version) {
	if len(vList) < 1 {
		return vList
	}
	mid := vList[0]
	left := make([]*version, 0)
	right := make([]*version, 0)
	for i := 1; i < len(vList); i++ {
		if mid.Greater(vList[i]) {
			left = append(left, vList[i])
		} else {
			right = append(right, vList[i])
		}
	}
	left, right = that.sort(left), that.sort(right)
	r = append(r, left...)
	r = append(r, mid)
	r = append(r, right...)
	return r
}

func (that *VComparator) Order() []string {
	vs := that.sort(that.v)
	for _, v := range vs {
		that.Versions = append(that.Versions, v.Origin)
	}
	return that.Versions
}
