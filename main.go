/*
a dev tool for multi-platforms
*/
package main

import (
	"os"
	"strings"

	utils "github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
)

var (
	GitHash string
	GitTime string
	GitTag  string
)

func main() {
	c := cmd.New()
	c.SetVersionInfo(GitTag, GitHash, GitTime)
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") && !strings.HasSuffix(ePath, "g") && !strings.HasSuffix(ePath, "g.exe") {
		/*
			this is for go run main.go.
		*/
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
		// fmt.Println(os.Environ()[0])
		// fmt.Println(os.Getwd())
		// gh := vctrl.NewGhDownloader()
		// gh.Download("https://github.com/moqsien/gvc/", false)

		// vl := vctrl.NewVlang()
		// vl.InstallVAnalyzerForVscode()

		// vp := vctrl.NewGSudo()
		// vp.Install(true)
		// vp.Install(true)
		// p := "a/b/c/d/e.zip"
		// fmt.Println(strings.ReplaceAll(p, filepath.Dir(p), ""))
		// gh := vctrl.NewGhDownloader()
		// gh.InstallGitForWindows()

		// br := vctrl.NewBrowser()
		// br.Save("edge", false)

		// w := vctrl.NewGVCWebdav()
		// w.DeploySSHFiles()
		self := vctrl.NewSelf()
		self.CheckLatestVersion("v1.6.4")
	} else if len(os.Args) < 2 {
		self := vctrl.NewSelf()
		self.Install()
		self.ShowPath()
	} else {
		s := &utils.CtrlCSignal{}
		s.ListenSignal()
		c.RunApp()
	}
}
