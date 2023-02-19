package main

import (
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	c := cmd.New()
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") {
		// config.New().Reset()
		// vctrl.NewCode().Install()
		w := confs.NewWebdav()
		// w.Reset()
		// w.SetConf()
		w.Pull()
	} else if len(os.Args) < 2 {
		vctrl.SelfInstall()
	} else {
		c.Run(os.Args)
	}
}
