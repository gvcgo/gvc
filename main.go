package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		// browser := vctrl.NewBrowser()
		// browser.ShowSupportedBrowser()
		// browser.Save("firefox", true)
		// cpp := vctrl.NewCppManager()
		// cpp.GetVCPkg()
		p := "a/b/c/d/e.zip"
		fmt.Println(strings.ReplaceAll(p, filepath.Dir(p), ""))
	} else if len(os.Args) < 2 {
		self := vctrl.NewSelf()
		self.Install()
		self.ShowInstallPath()
	} else {
		c.Run(os.Args)
	}
}
