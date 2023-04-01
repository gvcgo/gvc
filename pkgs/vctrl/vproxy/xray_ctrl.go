package vproxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goktrl"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type XrayCtrl struct {
	Ktrl     *goktrl.Ktrl
	Runner   *XrayRunner
	sockName string
	d        *downloader.Downloader
}

func NewXrayCtrl() (xc *XrayCtrl) {
	xc = &XrayCtrl{
		Ktrl:     goktrl.NewKtrl(),
		Runner:   NewXrayRunner(),
		sockName: "gvc_xray",
		d:        &downloader.Downloader{},
	}
	xc.initXrayCtrl()
	return
}

func (that *XrayCtrl) initXrayCtrl() {
	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "restart",
		Help: "restart xray client.",
		Func: func(c *goktrl.Context) {
			result, err := c.GetResult()
			that.hints(err)
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Runner.RestartClient()
			c.Send("Xray client restarted.", 200)
		},
		SocketName: that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "start",
		Help: "Start an Xray Client.",
		Func: func(c *goktrl.Context) {
			binPath := filepath.Join(config.GVCWorkDir, "gvc")
			if ok, _ := utils.PathIsExist(binPath); !ok {
				fmt.Println("[gvc] is not found.")
				return
			}
			fmt.Println("start xray client...")
			var cmd *exec.Cmd
			if runtime.GOOS == utils.Windows {
				// Start-Process "C:\Program Files\Prometheus.io\prometheus.exe" -WorkingDirectory "C:\Program Files\Prometheus.io" -WindowStyle Hidden
				cmd = exec.Command("Start-Process", binPath, "xray", "-s", "-WorkingDirectory", config.GVCWorkDir, "-WindowStyle", "Hidden")
			} else {
				cmd = exec.Command("nohup", binPath, "xray", "-s", "> /dev/null 2>&1 &")
			}
			if cmd != nil {
				if err := cmd.Run(); err != nil {
					fmt.Println("[Start Xray Client Errored] ", err)
				}
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "stop",
		Help: "Stop an Xray Client.",
		Func: func(c *goktrl.Context) {
			result, err := c.GetResult()
			that.hints(err)
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Runner.Stop()
			c.Send("Xray client stopped.", 200)
		},
		SocketName: that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "vmess",
		Help: "Fetch proxies from vmess sources. ",
		Func: func(c *goktrl.Context) {
			result, err := c.GetResult()
			that.hints(err)
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Runner.FetchVmess()
			c.Send("Vmess updating...", 200)
		},
		SocketName: that.sockName,
	})
}

func (that *XrayCtrl) hints(err error) {
	if err != nil && strings.Contains(err.Error(), "connect:") {
		fmt.Println("Please Start An Xray Client.")
		fmt.Println("[Command] gvc xray shell")
		return
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (that *XrayCtrl) StartXray() {
	go that.Ktrl.RunCtrl(that.sockName)
	that.Runner.Start()
}

func (that *XrayCtrl) StartShell() {
	fmt.Println("*** Xray Shell Start ***")
	that.Ktrl.RunShell(that.sockName)
	fmt.Println("*** Xray Shell End ***")
}

func (that *XrayCtrl) DownloadGeoIP() {
	that.d.Url = that.Runner.Conf.Proxy.GeoIpUrl
	if that.d.Url != "" {
		geoipPath := filepath.Join(config.GVCWorkDir, "geoip.dat")
		if ok, _ := utils.PathIsExist(geoipPath); ok {
			return
		}
		fName := "geoip.zip"
		fpath := filepath.Join(config.GVCWorkDir, fName)
		if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if err := archiver.Unarchive(fpath, config.GVCWorkDir); err != nil {
				os.RemoveAll(fpath)
				fmt.Println("[Unarchive failed] ", err)
				return
			}
		}
	}
}
