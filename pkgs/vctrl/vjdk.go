package vctrl

import "fmt"

var AllowedSuffixes = map[string]struct{}{
	".zip":     {},
	".tar.gz":  {},
	".tar.bz2": {},
	".tar.xz":  {},
}

type JDKPackage struct {
	Url      string
	FileName string
	OS       string
	Arch     string
	Size     string
	Checksum string
}

type JDKVersion struct {
	IsOfficial bool
}

func (that *JDKVersion) SetIsOfficial() {
	fmt.Println("Choose JDK download URL: ")
	fmt.Println("1) injdk.cn (Faster in china and by default.)")
	fmt.Println("2) oracle.com (Only latest versions are available.)")
	choice := "1"
	fmt.Scan(&choice)
	switch choice {
	case "2":
		that.IsOfficial = true
	default:
		that.IsOfficial = false
	}
}
