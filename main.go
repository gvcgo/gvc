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
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") && !strings.HasSuffix(ePath, "g") {
		c := confs.New()
		// c.SetupWebdav()
		c.Reset()

		// r := vchatgpt.NewRunner()
		// r.Run()

		// ui.Window()
		// m := tui.NewTui()
		// m.Run()
		// chatgpt.Run()
		browser := vctrl.NewBrowser()
		browser.ShowSupportedBrowser()
		browser.Save("firefox", true)

	} else if len(os.Args) < 2 {
		self := vctrl.NewSelf()
		self.Install()
		self.ShowInstallPath()
	} else {
		c.Run(os.Args)
	}
}
