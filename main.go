package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	c := cmd.New()
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") {
		// c := confs.New()
		// c.SetupWebdav()
		// c.Reset()
		// v := vctrl.NewGoVersion()
		// v.ShowRemoteVersions(vctrl.ShowStable)
		// v.UseVersion("1.19.6")
		// v.ShowInstalled()
		// v := vctrl.NewNVim()
		// v.Install()
		fmt.Println(utils.JoinUnixFilePath("/abc", "d", "/a/", "abc"))
	} else if len(os.Args) < 2 {
		vctrl.SelfInstall()
	} else {
		c.Run(os.Args)
	}
}
