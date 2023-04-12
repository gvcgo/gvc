package vproxy

import (
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
)

type ChanRawProxy chan RawProxy

type XrayVerifier struct {
	VmessFetcher   *ProxyFetcher
	VmessResult    *VmessList
	ClientList     []*XrayClient
	Ports          []int
	Conf           *config.GVConfig
	ProxyChan      ChanRawProxy
	VmessCollector ChanRawProxy
}

func NewVerifier() (xv *XrayVerifier) {
	xv = &XrayVerifier{
		VmessFetcher:   NewProxyFetcher(Vmess),
		VmessResult:    NewVmessList("proxies-verified-vmess.yml"),
		ClientList:     make([]*XrayClient, 0),
		Conf:           config.New(),
		VmessCollector: make(ChanRawProxy, 100),
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

func (that *XrayVerifier) GetVmessCollector() ChanRawProxy {
	return that.VmessCollector
}

func (that *XrayVerifier) GetConf() *config.GVConfig {
	return that.Conf
}

func (that *XrayVerifier) GetVmessVerifiedList() *VmessList {
	that.VmessResult.Reload()
	return that.VmessResult
}

func (that *XrayVerifier) IsAllClientsRunning() bool {
	for _, client := range that.ClientList {
		if client.VerifierIsRunning {
			return true
		}
	}
	return false
}

func (that *XrayVerifier) receiveResult() {
	res := []RawProxy{}
OUTTER:
	for {
		select {
		case p, ok := <-that.VmessCollector:
			if p != nil {
				res = append(res, p)
			}
			if !ok {
				break OUTTER
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
	that.VmessResult.Update(res)
}

func (that *XrayVerifier) sendProxy(force bool) {
	if that.VmessFetcher == nil {
		that.VmessFetcher = NewProxyFetcher(Vmess)
	}
	// wether force to fetch a new proxy list
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
	go that.receiveResult()
	// c := make(chan struct{})
	// <-c
}
