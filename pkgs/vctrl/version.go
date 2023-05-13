package vctrl

import (
	"fmt"

	"github.com/TwiN/go-color"
)

const (
	VERSION = "1.2.0"
)

type Version struct{}

func (that *Version) Show() {
	fmt.Println(color.InGreen("***========================================================***"))
	fmt.Println(color.InPurple("   GVC Version: ") + color.InYellow("v"+VERSION))
	fmt.Println(color.InPurple("   Github:      ") + color.InYellow("https://github.com/moqsien/gvc"))
	fmt.Println(color.InPurple("   Gitee:       ") + color.InYellow("https://gitee.com/moqsien/gvc_tools"))
	fmt.Println(color.InPurple("   Email:       ") + color.InYellow("moqsien@foxmail.com"))
	fmt.Println(color.InGreen("***========================================================***"))
}

func NewVersion() (v *Version) {
	return &Version{}
}
