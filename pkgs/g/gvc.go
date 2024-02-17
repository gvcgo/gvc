package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/pkgs/clis"
	"github.com/gvcgo/gvc/pkgs/confs"
	"github.com/gvcgo/gvc/pkgs/vctrl"
)

var (
	GitTag  string
	GitHash string
	GitTime string
)

func main() {
	c := clis.New()
	c.SetVersionInfo(GitTag, GitHash, GitTime)
	ePath, _ := os.Executable()

	if !strings.HasSuffix(ePath, "g") && !strings.HasSuffix(ePath, "g.exe") {
		// for test
		cfg := confs.New()
		cfg.Reset()

		// self := vctrl.NewSelf()
		// self.CheckLatestVersion("v1.6.4")

		// zig := vctrl.NewZig()
		// zig.GetZigList()

		/*
			https://github.com/zigtools/zls/releases/latest/
			https://github.com/neovim/neovim/releases/latest/
			https://github.com/protocolbuffers/protobuf/releases/latest/
			https://github.com/typst/typst/releases/latest/
			https://github.com/vlang/v/releases/latest/
			https://github.com/v-analyzer/v-analyzer/releases/latest/
			https://github.com/zigtools/zls/releases/latest/
			https://github.com/gvcgo/gvc/releases/latest/
			https://github.com/git-for-windows/git/releases/latest/
			https://github.com/neovide/neovide/releases/latest/
		*/
		// gh := vctrl.NewGhDownloader()
		// fmt.Printf("%+v\n", gh.ParseReleasesForGithubProject("https://github.com/neovide/neovide/releases/latest/"))

		nv := vctrl.NewNeoVim()
		// nv.InstallNeovim()
		nv.ToggleProxy()
	} else if len(os.Args) < 2 {
		/*
			GVC is allowed to be installed in ~/.gvc/ or $GOPATH/bin/ .
		*/
		goPath := os.Getenv("GOPATH")
		toInstall := true

		// Installed in $GOBIN
		if goPath != "" && strings.Contains(ePath, filepath.Join(goPath, "bin")) {
			os.MkdirAll(confs.GVCDir, 0777)
			toInstall = false
		}

		// Installed in ~/.gvc/
		if strings.Contains(ePath, confs.GVCDir) {
			toInstall = false
		}

		// Not installed yet.
		if toInstall {
			self := vctrl.NewSelf()
			self.Install()
			self.ShowPath()
		}
	} else {
		// run Clis
		s := &gutils.CtrlCSignal{}
		s.ListenSignal()
		c.Run()
	}
}
