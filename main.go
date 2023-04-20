package main

import (
	"os"
	"strings"

	"github.com/moqsien/gvc/pkgs/cmd"
	"github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/tchat/tui"
	"github.com/moqsien/gvc/pkgs/vctrl"
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
		// co := chatgpt.NewChatGptConf()
		// co.Reset()
		// co.Restore()
		// co.GetOptions()
		// co.SetConfField("api_key", "xxx")
		// fmt.Println(*co)

		// r := vchatgpt.NewRunner()
		// r.Run()

		// ui.Window()
		m := tui.NewTui()
		m.Run()

		// type Test struct {
		// 	TestKey key.Binding
		// }
		// t := Test{
		// 	TestKey: key.NewBinding(key.WithKeys("ctrl+i"), key.WithHelp("ctrl+i", "test")),
		// }
		// val := reflect.ValueOf(t)
		// for i := 0; i < val.NumField(); i++ {
		// 	fmt.Println(val.Field(i).Type().Name())
		// 	k := val.FieldByName(val.Type().Field(i).Name).Interface()
		// 	n, ok := k.(key.Binding)
		// 	fmt.Println(ok, n)
		// }
		// cc := vchat.NewChatGptConf()
		// cc.Reset()
		// cc.Restore()
		// cc.ShowConfigOpts()
		// cc.SetConfField("api_key", "")
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
