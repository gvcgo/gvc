package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/vctrl"
	"github.com/shirou/gopsutil/v3/process"
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
		_ = vctrl.NewTypstVersion()

		// cc := vchat.NewChatGpt()
		// cc.SetAppKey("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")
		// cc.SetProxyPort(2019)

		// cc.Chat("用go写一个桶排序")
		// v, _ := mem.VirtualMemory()
		pp, _ := process.Pids()
		for _, p := range pp {
			// fmt.Println(p)
			proc := process.Process{Pid: p}
			name, _ := proc.Cmdline()
			// name := strings.Join(names, " ")
			if strings.Contains(name, "xray") {
				fmt.Println(name)
			}
		}
		// almost every return value is a struct
		// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

		// convert to JSON. String() is also implemented
		// fmt.Println(v)
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
