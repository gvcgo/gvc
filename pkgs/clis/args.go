package clis

import (
	"strings"

	"github.com/gvcgo/gvc/pkgs/clis/langs"
)

/*
See "github.com/gvcgo/gvc/pkgs/clis/langs"
*/
func HandleArgs(args ...string) (aList []string) {
	if len(args) < 4 {
		return args
	}
	// for "gvc go build"
	if (args[1] == "go" || args[1] == "g") && (args[2] == "build" || args[2] == "bui" || args[2] == "b") {
		for _, v := range args {
			if strings.HasPrefix(v, "-") && !strings.Contains(v, " ") {
				aList = append(aList, strings.Replace(v, "-", langs.TempChar, 1))
			} else {
				aList = append(aList, v)
			}
		}
		return aList
	}
	return args
}
