package vproxy

import (
	config "github.com/moqsien/gvc/pkgs/confs"
)

type XrayVerifier struct {
	VmessFetcher *ProxyFetcher
	VmessResult  []*XrayVmessOutbound
	ClientList   []*XrayClient
	Ports        []int
	Conf         *config.GVConfig
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
		xv.ClientList = append(xv.ClientList, NewXrayVmessClient(&XrayInbound{Port: p}))
	}
	return
}
