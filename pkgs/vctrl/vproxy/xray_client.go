package vproxy

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
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

type XrayVmessOutbound struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	UserId   string `json:"id"`
	Network  string `json:"network"`
	Security string `json:"security"`
	Path     string `json:"path"`
	Raw      string `json:"raw"`
}

type VO struct {
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
}

func (that *XrayVmessOutbound) ParseVmessUri(rawUri string) {
	if strings.HasPrefix(rawUri, "vmess://") {
		jsonStr := utils.DecodeBase64(strings.ReplaceAll(rawUri, "vmess://", ""))
		j := gjson.New(jsonStr)
		that.Security = j.GetString("tls")
		// if that.Security != "" {
		// 	that.Security = "tls"
		// }
		that.Address = j.GetString("add")
		that.Port = j.GetInt("port")
		that.UserId = j.GetString("id")
		that.Network = j.GetString("net")
		that.Path = j.GetString("path")
	}
}

type XrayClient struct {
	*core.Instance
	Config            *xconf.Config
	ConfigVmessStr    string
	VerifierIsRunning bool
	Verifier          Verifier
	Inbound           *XrayInbound
}

func NewXrayVmessClient(xi *XrayInbound, verifier ...Verifier) (r *XrayClient) {
	r = &XrayClient{
		ConfigVmessStr: XrayVmessConfStr,
		Inbound:        xi,
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
	if xo != nil {
		j := gjson.New(that.ConfigVmessStr)
		j.Set("outbounds.0.settings.vnext.0.address", xo.Address)
		j.Set("outbounds.0.settings.vnext.0.port", xo.Port)
		j.Set("outbounds.0.settings.vnext.0.users.0.id", xo.UserId)
		j.Set("outbounds.0.streamSettings.network", xo.Network)
		j.Set("outbounds.0.streamSettings.security", xo.Security)
		j.Set("outbounds.0.streamSettings.wsSettings.path", xo.Path)
		that.ConfigVmessStr = j.MustToJsonString()
	}
}

func (that *XrayClient) StartVmessClient(p RawProxy) error {
	xo := &XrayVmessOutbound{}
	xo.ParseVmessUri(p.GetUri())
	that.FormatVmessOutbound(xo)
	if config, err := serial.DecodeJSONConfig(utils.ConvertStrToReader(that.ConfigVmessStr)); err == nil {
		var f *core.Config
		f, err = config.Build()
		if err != nil {
			fmt.Println("[Build config for Xray failed] ", err)
			return err
		}
		that.Instance, err = core.New(f)
		if err != nil {
			fmt.Println("[Init Xray Instance Failed] ", err)
			return err
		}
		that.Start()
	} else {
		fmt.Println("[Start Client Failed] ", err)
		return err
	}
	return nil
}

func (that *XrayClient) CloseClient() {
	if that.Instance != nil {
		that.Close()
		that.Instance = nil
		runtime.GC()
	}
}

func (that *XrayClient) VerifyProxy() (ok bool, timeLag int64) {
	node := &point.Point{
		Protocols: []*protocol.Protocol{
			{
				Protocol: &protocol.Protocol_Simple{
					Simple: &protocol.Simple{
						Host:             "127.0.0.1",
						Port:             int32(that.Inbound.Port),
						PacketConnDirect: true,
					},
				},
			},
			{
				Protocol: &protocol.Protocol_Socks5{
					Socks5: &protocol.Socks5{},
				},
			},
		},
	}

	pro, err := register.Dialer(node)
	if err != nil {
		fmt.Println("[Dialer error] ", err)
		return
	}
	t := that.Verifier.GetConf().Proxy.MaxRTT
	if t == 0 {
		t = 5
	}
	c := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				add, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
				if err != nil {
					return nil, fmt.Errorf("parse address failed: %w", err)
				}
				add.WithContext(ctx)
				return pro.Conn(add)
			}}, Timeout: time.Duration(t) * time.Second,
	}
	vUrl := that.Verifier.GetConf().Proxy.VerifyUrl
	if vUrl == "" {
		vUrl = "https://www.google.com"
	}
	startTime := time.Now()
	resp, err := c.Get(vUrl)
	timeLag = time.Since(startTime).Milliseconds()
	if err != nil {
		fmt.Println("[Verify url failed] ", err)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	ok = strings.Contains(buf.String(), "</html>")
	return
}

func (that *XrayClient) RunVerifier(typ ...ProxyType) {
	ProxyChan := that.Verifier.GetProxyChan()
	if len(typ) == 0 || typ[0] == Vmess {
		that.VerifierIsRunning = true
	Outter:
		for {
			select {
			case p, ok := <-ProxyChan:
				if p != nil {
					if err := that.StartVmessClient(p); err == nil {
						if ok, timeLag := that.VerifyProxy(); ok {
							p.SetRTT(timeLag)
							that.Verifier.GetVmessCollector() <- p
						}
						that.CloseClient()
					} else {
						fmt.Println("[Start client failed] ", err)
					}
				}
				if !ok {
					that.VerifierIsRunning = false
					if !that.Verifier.IsAllClientsRunning() {
						close(that.Verifier.GetVmessCollector())
					}
					break Outter
				}
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}
