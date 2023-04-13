package vproxy

import (
	"fmt"
	"os"
	"strconv"
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

func (that *XrayRunner) Start(idx ...string) {
	that.Verifier.RunVmess(true)
	that.Cron.AddFunc(that.Conf.Proxy.GetCrontabStr(), func() {
		that.Verifier.RunVmess(false)
	})
	that.Cron.Start()

	that.RestartClient(false, idx...)
	<-StopSignal
	fmt.Println("Exiting xray...")
	os.Exit(0)
}

func (that *XrayRunner) Stop() {
	StopSignal <- struct{}{}
}

func (that *XrayRunner) RestartClient(enableFixed bool, idx ...string) (pStr string) {
	if that.Client != nil {
		that.Client.CloseClient()
		time.Sleep(500 * time.Millisecond)
	}
	index := -1
	if len(idx) > 0 {
		index, _ = strconv.Atoi(idx[0])
	}
	vmessList := that.Verifier.VmessResult
	if len(vmessList.GetProxyList()) == 0 || (enableFixed && len(that.Verifier.VmessFixed.GetProxyList()) > 0) {
		vmessList = that.Verifier.VmessFixed
	}
	vmessList.CheckFilePath()

	xo := &XrayVmessOutbound{}
	p := vmessList.ChooseByIndex(index)
	if p != nil {
		xo.ParseVmessUri(p.GetUri())
		pStr = fmt.Sprintf("%s:%d", xo.Address, xo.Port)
		that.Client = NewXrayVmessClient(&XrayInbound{Port: that.Conf.Proxy.InboundPort})
		that.Client.StartVmessClient(p)
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
	vl.ShowProxyList()
}

func (that *XrayRunner) ShowVmessFixedList() {
	vl := that.Verifier.GetVmessFixedList()
	vl.ShowProxyList()
}
