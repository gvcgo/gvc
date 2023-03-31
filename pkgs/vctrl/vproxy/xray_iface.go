package vproxy

import config "github.com/moqsien/gvc/pkgs/confs"

type Proxies interface {
	Today() string
	GetDate() string
	Reload()
	Update(any)
	GetProxyList() []*Proxy
}

type RawProxy interface {
	GetUri() string
	SetRTT(rtt int64)
}

type Verifier interface {
	GetProxyChan() ChanRawProxy
	GetVmessCollector() ChanRawProxy
	GetConf() *config.GVConfig
	IsAllClientsRunning() bool
}
