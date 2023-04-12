package confs

type ChatgptConf struct {
	AppKey         string `koanf:"app_key"`
	LocalProxyPort int    `koanf:"local_proxy_port"`
	ProxyTimeout   int    `koanf:"proxy_timeout"`
}

func NewGptConf() (r *ChatgptConf) {
	r = &ChatgptConf{}
	return
}

func (that *ChatgptConf) Reset() {
	that.LocalProxyPort = 2019
	that.ProxyTimeout = 120
}
