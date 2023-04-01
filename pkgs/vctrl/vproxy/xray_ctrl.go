package vproxy

import (
	"fmt"
	"strings"

	"github.com/moqsien/goktrl"
)

type XrayCtrl struct {
	Ktrl     *goktrl.Ktrl
	Runner   *XrayRunner
	sockName string
}

func NewXrayCtrl() (xc *XrayCtrl) {
	xc = &XrayCtrl{
		Ktrl:     goktrl.NewKtrl(),
		Runner:   NewXrayRunner(),
		sockName: "gvc_xray",
	}
	xc.initXrayCtrl()
	return
}

func (that *XrayCtrl) initXrayCtrl() {
	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "restart",
		Help: "restart xray client.",
		Func: func(c *goktrl.Context) {
			result, err := c.GetResult()
			if err != nil && strings.Contains(err.Error(), "connect:") {
				fmt.Println("Please Start An Xray Client.")
				return
			}
			fmt.Println(string(result))
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Runner.RestartClient()
			c.Send("Xray client restarted.", 200)
		},
		SocketName: that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "start",
		Help: "Start an Xray Client.",
		Func: func(c *goktrl.Context) {
			// TODO: command line
			fmt.Println("start xray client")
		},
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "stop",
		Help:        "Stop an Xray Client.",
		Func:        func(c *goktrl.Context) {},
		KtrlHandler: func(c *goktrl.Context) {},
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "vmess",
		Help:        "Fetch proxies from vmess sources. ",
		Func:        func(c *goktrl.Context) {},
		KtrlHandler: func(c *goktrl.Context) {},
	})
}
