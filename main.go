package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	c := cmd.New()
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") && !strings.HasSuffix(ePath, "g") && !strings.HasSuffix(ePath, "g.exe") {
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
		// cpp.InstallVCPkg()
		content := []byte("testejalkjfldjfakljdflkdjfklajfl")
		bar := tui.NewProgressBar("test", len(content))
		bar.Start()
		io.Copy(bar, bytes.NewBuffer(content))
		// p := "a/b/c/d/e.zip"
		// fmt.Println(strings.ReplaceAll(p, filepath.Dir(p), ""))
	} else if len(os.Args) < 2 {
		self := vctrl.NewSelf()
		self.Install()
		self.ShowInstallPath()
	} else {
		c.Run(os.Args)
	}
}
