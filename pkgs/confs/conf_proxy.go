package confs

import "strings"

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
	/*
		https://gh.flyinbug.top/gh/
		https://gvc.1710717.xyz/proxy/
		https://gh.chapro.xyz/
	*/
	that.ReverseProxyUrl = "https://gvc.1710717.xyz/proxy/" // only for gvc related
}

func (that *GVCReverseProxyConf) WrapUrl(origUrl string) (finUrl string) {
	if !strings.HasPrefix(origUrl, "http") || that.ReverseProxyUrl == "" {
		return origUrl
	}
	finUrl = that.ReverseProxyUrl + origUrl
	return
}
