package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/config"
	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	c := cmd.New()
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") {
		config.New().Reset()
		fmt.Println("test")
		vctrl.NewCode().Run()
	} else if len(os.Args) < 2 {
		vctrl.SelfInstall()
	} else {
		c.Run(os.Args)
	}
}
