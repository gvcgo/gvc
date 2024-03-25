package asciinema

import (
	"os"
	"strings"
)

var descardingList []string = []string{
	`?\u001b\\\u001b[6n`,
}

func verify(line string) bool {
	for _, s := range descardingList {
		if strings.Contains(line, s) {
			return false
		}
	}
	return true
}

func FixCast(fPath string) {
	content, _ := os.ReadFile(fPath)
	if len(content) > 0 {
		sList := strings.Split(string(content), "\n")
		data := []string{}
		for _, line := range sList {
			if verify(line) {
				data = append(data, line)
			}
		}
		if len(data) > 0 {
			s := strings.Join(data, "\n")
			os.WriteFile(fPath, []byte(s), os.ModePerm)
		}
	}
}
