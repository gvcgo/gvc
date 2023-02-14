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

func (that *version) Greater(v *version) {

}

type VComparator struct {
	Versions []string
	v        []*version
}

func NewVComparator(vs []string) *VComparator {
	vList := []*version{}
	var vresult []string
	for _, v := range vs {
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
	return &VComparator{Versions: []string{}, v: vList}
}
