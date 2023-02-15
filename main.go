package main

import (
	"fmt"

	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	h := vctrl.NewGoVersion()
	h.UseVersion("1.20")
	h.UseVersion("1.20.1")
	h.ShowInstalled()
	fmt.Println("======")
	h.RemoveVersion("1.20")
	// h.ShowInstalled()
	// c := config.New()
	// c.Reset()
}
