package main

import (
	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	h := vctrl.NewGoVersion()
	h.Download("1.20")
	// c := config.New()
	// c.Reset()
}
