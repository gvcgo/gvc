package vproxy

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gookit/color"
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

func (that *XrayRunner) Start(idx ...string) {
	that.Verifier.RunVmess(true)
	that.Cron.AddFunc(that.Conf.Proxy.GetCrontabStr(), func() {
		that.Verifier.RunVmess(false)
	})
	that.Cron.Start()

	that.RestartClient(idx...)
	<-StopSignal
	fmt.Println("Exiting xray...")
	os.Exit(0)
}

func (that *XrayRunner) Stop() {
	StopSignal <- struct{}{}
}

func (that *XrayRunner) RestartClient(idx ...string) (pStr string) {
	if that.Client != nil {
		that.Client.CloseClient()
		time.Sleep(500 * time.Millisecond)
	}
	index := -1
	if len(idx) > 0 {
		index, _ = strconv.Atoi(idx[0])
	}

	that.Verifier.VmessResult.CheckFilePath()

	xo := &XrayVmessOutbound{}
	p := that.Verifier.VmessResult.ChooseByIndex(index)
	if p != nil {
		xo.ParseVmessUri(p.GetUri())
		pStr = fmt.Sprintf("%s:%d", xo.Address, xo.Port)
		that.Client = NewXrayVmessClient(&XrayInbound{Port: that.Conf.Proxy.InboundPort})
		that.Client.StartVmessClient(that.Verifier.VmessResult.ChooseRandom())
		fmt.Printf("Xray started @socks5://127.0.0.1:%v", that.Conf.Proxy.InboundPort)
	}
	return
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

func (that *XrayRunner) ShowVmessVerifiedList() {
	vl := that.Verifier.GetVmessVerifiedList()
	color.Red.Println(vl.Date)
	color.Yellow.Println(fmt.Sprintf("Total: %d", vl.Total))
	for idx, p := range vl.Proxies {
		rawUrl := p.GetUri()
		xo := &XrayVmessOutbound{}
		xo.ParseVmessUri(rawUrl)
		if xo.Address != "" {
			color.Cyan.Println(fmt.Sprintf("%d. %s:%d (rtt: %dms)", idx, xo.Address, xo.Port, p.RTT))
		}
	}
}
