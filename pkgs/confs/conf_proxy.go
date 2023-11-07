package confs

/*
GVC Reverse Proxy
*/
type GVCReverseProxyConf struct {
	ReverseProxyUrl string `koanf:"reverse_proxy_url"`
}

func NewReverseProxyConf() (r *GVCReverseProxyConf) {
	r = &GVCReverseProxyConf{}
	return
}

func (that *GVCReverseProxyConf) Reset() {
	that.ReverseProxyUrl = "https://worker-github.moqsien2022.workers.dev/proxy/"
}
