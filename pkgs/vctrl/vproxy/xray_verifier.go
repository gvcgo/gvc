package vproxy

import (
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
)

type ChanRawProxy chan RawProxy

type XrayVerifier struct {
	VmessFetcher *ProxyFetcher
	VmessResult  []*XrayVmessOutbound
	ClientList   []*XrayClient
	Ports        []int
	Conf         *config.GVConfig
	ProxyChan    ChanRawProxy
}

func NewVerifier() (xv *XrayVerifier) {
	xv = &XrayVerifier{
		VmessFetcher: NewProxyFetcher(Vmess),
		VmessResult:  make([]*XrayVmessOutbound, 0),
		ClientList:   make([]*XrayClient, 0),
		Conf:         config.New(),
	}
	xv.Ports = xv.Conf.Proxy.GetVerifyPorts()
	for _, p := range xv.Ports {
		xv.ClientList = append(xv.ClientList, NewXrayVmessClient(&XrayInbound{Port: p}, xv))
	}
	return
}

func (that *XrayVerifier) GetProxyChan() ChanRawProxy {
	return that.ProxyChan
}

func (that *XrayVerifier) GetConf() *config.GVConfig {
	return that.Conf
}

func (that *XrayVerifier) sendProxy(force bool) {
	if that.VmessFetcher == nil {
		that.VmessFetcher = NewProxyFetcher(Vmess)
	}
	that.VmessFetcher.GetProxyList(force)
	p := that.VmessFetcher.ProxyList.GetProxyList()
	if len(p) > 0 {
		i := 0
		for i < len(p) {
			select {
			case that.ProxyChan <- p[i]:
				i++
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
		close(that.ProxyChan)
	}
}

func (that *XrayVerifier) RunVmess(force ...bool) {
	f := false
	if len(force) > 0 && force[0] {
		f = true
	}
	that.ProxyChan = make(ChanRawProxy, len(that.Ports))
	go that.sendProxy(f)
	for _, client := range that.ClientList {
		go client.RunVerifier(Vmess)
	}
	c := make(chan struct{})
	<-c
}
