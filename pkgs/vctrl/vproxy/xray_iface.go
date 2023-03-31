package vproxy

type Proxies interface {
	Today() string
	GetDate() string
	Reload()
	Update(any)
	GetProxyList() []*Proxy
}

type RawProxy interface {
	GetUri() string
}

type Verifier interface {
	GetProxyChan() ChanRawProxy
}
