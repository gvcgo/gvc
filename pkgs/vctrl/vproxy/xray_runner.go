package vproxy

import (
	"fmt"
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

func (that *XrayRunner) Start() {
	that.Verifier.RunVmess(false)
	that.Cron.AddFunc(that.Conf.Proxy.GetCrontabStr(), func() {
		that.Verifier.RunVmess(false)
	})
	that.RestartClient()

	<-StopSignal
	fmt.Println("Exiting xray...")
	os.Exit(0)
}

func (that *XrayRunner) Stop() {
	StopSignal <- struct{}{}
}

func (that *XrayRunner) RestartClient() {
	if that.Client != nil {
		that.Client.CloseClient()
		time.Sleep(500 * time.Millisecond)
	}
	xo := &XrayVmessOutbound{}
	p := that.Verifier.VmessResult.ChooseFastest()
	if p != nil {
		xo.ParseVmessUri(p.GetUri())
		that.Client = NewXrayVmessClient(&XrayInbound{Port: that.Conf.Proxy.InboundPort})
		that.Client.StartVmessClient(that.Verifier.VmessResult.ChooseFastest())
		fmt.Printf("Xray started @socks5://127.0.0.1:%v", that.Conf.Proxy.InboundPort)
	}
}

func (that *XrayRunner) FetchVmess() {
	that.Verifier.RunVmess(true)
	for {
		if !that.Verifier.IsAllClientsRunning() {
			break
		} else {
			time.Sleep(time.Millisecond * 200)
		}
	}
}
