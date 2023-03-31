package vproxy

import (
	"fmt"
	"time"

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
	Verifier       Verifier
}

func NewXrayVmessClient(xi *XrayInbound, verifier ...Verifier) (r *XrayClient) {
	r = &XrayClient{
		ConfigVmessStr: XrayVmessConfStr,
	}
	if len(verifier) > 0 {
		r.Verifier = verifier[0]
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

func (that *XrayClient) RunVerifier(typ ...ProxyType) {
	ProxyChan := that.Verifier.GetProxyChan()
	if len(typ) == 0 || typ[0] == Vmess {
		for {
			select {
			case p, ok := <-ProxyChan:
				if !ok {
					break
				}
				if p != nil {
					fmt.Println(p.GetUri())
				}
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
