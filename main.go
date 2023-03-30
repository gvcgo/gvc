package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/moqsien/gvc/pkgs/vctrl/vproxy"
)

func main() {
	c := cmd.New()
	ePath, _ := os.Executable()
	if !strings.HasSuffix(ePath, "gvc") && !strings.HasSuffix(ePath, "gvc.exe") && !strings.HasSuffix(ePath, "g") {
		c := confs.New()
		// c.SetupWebdav()
		c.Reset()
		// v := vctrl.NewGoVersion()
		// v.ShowRemoteVersions(vctrl.ShowStable)
		// v.UseVersion("1.19.6")
		// v.ShowInstalled()
		// v := vctrl.NewNVim()
		// v.Install()
		// fmt.Println(utils.JoinUnixFilePath("/abc", "d", "/a/", "abc"))
		// g := vctrl.NewGoVersion()
		// g.SearchLibs("json", 1)
		// fmt.Println(c.Proxy.GetSubUrls())
		// v := vctrl.NewProxy()
		// v.GetProxyList()
		// v := vproxy.NewProxyer()
		// v.RunXray()
		v := vproxy.NewXrayClient()
		fmt.Println(v)
		// v.UseVersion("v18.14.0")
		// v.UseVersion("v19.8.0")
		// v.ShowInstalled()
		// v.RemoveVersion("v18.14.0")
		// v.ShowInstalled()
		// v.UseVersion("java19")
	} else if len(os.Args) < 2 {
		self := vctrl.NewSelf()
		self.Install()
		self.ShowInstallPath()
	} else {
		c.Run(os.Args)
	}
}
