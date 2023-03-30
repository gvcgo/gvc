package vproxy

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/xtls/xray-core/core"
	xconf "github.com/xtls/xray-core/infra/conf"
	_ "github.com/xtls/xray-core/main/confloader/external"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type XrayInbound struct {
	Port   int    `json:"port"`
	Listen string `json:"listen"`
}

type XrayVmessOutbound struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	UserId   string `json:"id"`
	Network  string `json:"network"`
	Security string `json:"security"`
	Path     string `json:"path"`
	Raw      string `json:"raw"`
}

type XrayClient struct {
	*core.Instance
	Config         *xconf.Config
	ConfigVmessStr string
	IsRunning      bool
}

func NewXrayVmessClient(xi *XrayInbound) (r *XrayClient) {
	r = &XrayClient{
		ConfigVmessStr: XrayVmessConfStr,
	}
	r.formatInbound(xi)
	return
}

func (that *XrayClient) formatInbound(xi *XrayInbound) {
	j := gjson.New(that.ConfigVmessStr)
	if xi.Port != 0 {
		j.Set("inbounds.0.port", xi.Port)
	}
	if xi.Listen != "" {
		j.Set("inbounds.0.listen", xi.Listen)
	}
	that.ConfigVmessStr = j.MustToJsonString()
}

func (that *XrayClient) FormatVmessOutbound(xo *XrayVmessOutbound) {
	j := gjson.New(that.ConfigVmessStr)
	j.Set("outbounds.0.settings.vnext.0.address", xo.Address)
	j.Set("outbounds.0.settings.vnext.0.port", xo.Port)
	j.Set("outbounds.0.settings.vnext.0.users.0.id", xo.UserId)
	j.Set("outbounds.0.streamSettings.network", xo.Network)
	j.Set("outbounds.0.streamSettings.security", xo.Security)
	j.Set("outbounds.0.streamSettings.wsSettings.path", xo.Path)
	that.ConfigVmessStr = j.MustToJsonString()
}

func (that *XrayClient) ParseVmessUri(rawUri string) {

}
