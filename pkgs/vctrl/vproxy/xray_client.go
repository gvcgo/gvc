package vproxy

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/xtls/xray-core/core"
	xconf "github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	_ "github.com/xtls/xray-core/main/confloader/external"
	_ "github.com/xtls/xray-core/main/distro/all"
)

type XrayInbound struct {
	Port   int    `json:"port"`
	Listen string `json:"listen"`
}

type XrayVmessOutbound struct{}

type XrayClient struct {
	*core.Instance
	Config    *xconf.Config
	IsRunning bool
}

func NewXrayClient(xi *XrayInbound) (r *XrayClient) {
	r = &XrayClient{}
	r.initConf(r.formatVmessConf(xi))
	return
}

func (that *XrayClient) formatVmessConf(xi *XrayInbound) (r string) {
	j := gjson.New(XrayVmessConfStr)
	if xi.Port != 0 {
		j.Set("inbounds.0.port", xi.Port)
	}
	if xi.Listen != "" {
		j.Set("inbounds.0.listen", xi.Listen)
	}
	return j.MustToJsonString()
}

func (that *XrayClient) initConf(confStr string) (err error) {
	if confStr != "" {
		r := utils.ConvertStrToReader(confStr)
		that.Config, err = serial.DecodeJSONConfig(r)
	}
	return
}
