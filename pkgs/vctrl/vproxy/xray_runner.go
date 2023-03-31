package vproxy

import (
	"os"
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/robfig/cron/v3"
)

var StopSignal chan struct{} = make(chan struct{})

type XrayRunner struct {
	Client   *XrayClient
	Verifier *XrayVerifier
	Conf     *config.GVConfig
	Cron     *cron.Cron
}

func NewXrayRunner() (xr *XrayRunner) {
	xr = &XrayRunner{
		Verifier: NewVerifier(),
		Conf:     config.New(),
		Cron:     cron.New(),
	}
	return
}

func (that *XrayRunner) Run(force bool) {
	that.Cron.AddFunc(that.Conf.Proxy.GetCrontabStr(), func() {
		that.Verifier.RunVmess(force)
	})
	p := that.Verifier.VmessResult.ChooseFastest()
	if that.Client != nil {
		that.Client.CloseClient()
		time.Sleep(time.Millisecond * 500)
	}
	if p != nil {
		xo := &XrayVmessOutbound{}
		xo.ParseVmessUri(p.GetUri())
		that.Client = NewXrayVmessClient(&XrayInbound{Port: that.Conf.Proxy.InboundPort})
		that.Client.StartVmessClient(that.Verifier.VmessResult.ChooseFastest())
	}

	<-StopSignal
	os.Exit(0)
}
