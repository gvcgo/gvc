package confs

type ChatgptConf struct {
	AppKey       string `koanf:"app_key"`
	HywwwLoveUrl string `koanf:"hywww_love_url"`
}

func NewGptConf() (r *ChatgptConf) {
	r = &ChatgptConf{
		AppKey: "",
	}
	return
}

func (that *ChatgptConf) Reset() {
	that.AppKey = ""
	that.HywwwLoveUrl = "https://chat.hywwwlove.top/v1.0/chat/?code=%s"
}
