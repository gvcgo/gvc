package vproxy

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
	"github.com/moqsien/goktrl"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type XrayCtrl struct {
	Ktrl           *goktrl.Ktrl
	Runner         *XrayRunner
	sockName       string
	d              *downloader.Downloader
	scriptPath     string
	keepScriptPath string
	batPath        string
	keepBatPath    string
}

func NewXrayCtrl() (xc *XrayCtrl) {
	xc = &XrayCtrl{
		Ktrl:           goktrl.NewKtrl(),
		Runner:         NewXrayRunner(),
		sockName:       "gvc_xray",
		d:              &downloader.Downloader{},
		scriptPath:     filepath.Join(config.ProxyFilesDir, "run_xray.sh"),
		keepScriptPath: filepath.Join(config.ProxyFilesDir, "keep_run.sh"),
		batPath:        filepath.Join(config.ProxyFilesDir, "run_xray.bat"),
		keepBatPath:    filepath.Join(config.ProxyFilesDir, "keep_run.bat"),
	}
	xc.initXrayCtrl()
	return
}

func (that *XrayCtrl) getConf() *config.GVConfig {
	return that.Runner.Conf
}

func (that *XrayCtrl) runPingServer() {
	cnf := that.getConf()
	if cnf.Proxy.XrayPingSock != "" {
		s := utils.NewUServer(cnf.Proxy.XrayPingSock)
		s.AddHandler("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		s.Start()
	}
}

func (that *XrayCtrl) IsXrayRunning() bool {
	cnf := that.getConf()
	c := utils.NewUClient(cnf.Proxy.XrayPingSock)
	if result, err := c.GetResp("/ping", map[string]string{}); err == nil {
		return strings.Contains(result, "pong")
	} else {
		fmt.Println("[Error detail] ", err)
	}
	return false
}

func (that *XrayCtrl) runByScript(batPath, scriptPath string) {
	binPath := filepath.Join(config.GVCWorkDir, "gvc")
	if ok, _ := utils.PathIsExist(binPath); !ok {
		fmt.Println("[gvc] is not found.")
		return
	}
	fmt.Println("Starting Xray Client...")
	that.writeScript()
	var cmd *exec.Cmd
	if runtime.GOOS == utils.Windows {
		cmd = exec.Command("powershell", batPath)
	} else {
		cmd = exec.Command("sh", scriptPath)
	}
	if cmd != nil {
		if err := cmd.Run(); err != nil {
			fmt.Println("[Start Xray Client Errored] ", err)
		}
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Xray Client Started.")
}

func (that *XrayCtrl) initXrayCtrl() {
	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "start",
		Help: "Start an Xray Client.",
		Func: func(c *goktrl.Context) {
			that.runByScript(that.keepBatPath, that.keepScriptPath)
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
				that.sendQuitSignal()
			}
		},
		KtrlHandler: func(c *goktrl.Context) {
			that.Runner.Stop()
			c.Send("Xray client stopped.", 200)
		},
		SocketName: that.sockName,
	})
	type restartOpts struct {
		Enable bool `alias:"e" descr:"Enable fixed vmess list or not."`
	}
	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "restart",
		Help: "restart xray client.",
		Opts: &restartOpts{},
		Func: func(c *goktrl.Context) {
			result, err := c.GetResult()
			that.hints(err)
			if len(result) > 0 {
				fmt.Println(string(result))
			}
		},
		ArgsDescription: "choose a specified proxy by index.",
		KtrlHandler: func(c *goktrl.Context) {
			opts := c.Options.(*restartOpts)
			pStr := that.Runner.RestartClient(opts.Enable, c.Args...)
			c.Send(fmt.Sprintf("Xray client restarted @ proxy: %s | %s", pStr, strings.Join(c.Args, ",")), 200)
		},
		SocketName: that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name:            "add",
		Help:            "Add vmesses to fixed list. ",
		ArgsRequired:    true,
		ArgsDescription: "legal vmess uris.",
		Func: func(c *goktrl.Context) {
			for _, vm := range c.Args {
				that.Runner.Verifier.SetVmessForFixed(vm)
			}
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.sockName,
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

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "show",
		Help: "Show available proxy list. ",
		Func: func(c *goktrl.Context) {
			fmt.Println("Free VPNs: ")
			that.Runner.ShowVmessVerifiedList()
			fmt.Println("Fixed VPNs: ")
			that.Runner.ShowVmessFixedList()
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "status",
		Help: "Show running status. ",
		Func: func(c *goktrl.Context) {
			if that.IsXrayRunning() {
				fmt.Println("[gvc] xray client is running.")
				return
			}
			fmt.Println("[gvc] xray server is stopped.")
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.sockName,
	})

	that.Ktrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "omega",
		Help: "Download Switchy-Omega for GoogleChrome. ",
		Func: func(c *goktrl.Context) {
			that.DownloadSwithOmega()
		},
		KtrlHandler: func(c *goktrl.Context) {},
		SocketName:  that.sockName,
	})
}

func (that *XrayCtrl) writeScript() {
	if runtime.GOOS != utils.Windows && that.scriptPath != "" {
		if ok, _ := utils.PathIsExist(that.scriptPath); !ok {
			os.WriteFile(that.scriptPath, []byte(config.ProxyXrayShellScript), 0777)
		}
	}
	if runtime.GOOS != utils.Windows && that.keepScriptPath != "" {
		if ok, _ := utils.PathIsExist(that.keepScriptPath); !ok {
			os.WriteFile(that.keepScriptPath, []byte(config.ProxyXrayKeepRunningScript), 0777)
		}
	}

	if runtime.GOOS == utils.Windows && that.batPath != "" {
		if ok, _ := utils.PathIsExist(that.batPath); !ok {
			binPath := filepath.Join(config.GVCWorkDir, "gvc.exe")
			batContent := fmt.Sprintf(config.ProxyXrayBatScript, binPath, config.GVCWorkDir)
			os.WriteFile(that.batPath, []byte(batContent), 0777)
		}
	}
	if runtime.GOOS == utils.Windows && that.keepBatPath != "" {
		if ok, _ := utils.PathIsExist(that.keepBatPath); !ok {
			binPath := filepath.Join(config.GVCWorkDir, "gvc.exe")
			batContent := fmt.Sprintf(config.ProxyXrayKeepRunningBat, binPath, config.GVCWorkDir)
			os.WriteFile(that.keepBatPath, []byte(batContent), 0777)
		}
	}
}

func (that *XrayCtrl) hints(err error) {
	if err != nil && strings.Contains(err.Error(), "connect:") {
		fmt.Println("Please Start An Xray Client.")
		fmt.Println("[Use Command like] gvc xray")
		return
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (that *XrayCtrl) StartXray() {
	go that.Ktrl.RunCtrl(that.sockName)
	go that.runPingServer()
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

func (that *XrayCtrl) DownloadSwithOmega() {
	that.d.Url = that.Runner.Conf.Proxy.SwitchOmegaUrl
	if that.d.Url != "" {
		omegaPath := filepath.Join(config.ProxyFilesDir, "switchy_omega")
		if ok, _ := utils.PathIsExist(omegaPath); ok {
			fmt.Println("[Archive Path] ", omegaPath)
			return
		}
		fName := "switchy-omega.zip"
		fpath := filepath.Join(config.ProxyFilesDir, fName)
		if size := that.d.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
			if err := archiver.Unarchive(fpath, omegaPath); err != nil {
				os.RemoveAll(fpath)
				os.RemoveAll(omegaPath)
				fmt.Println("[Unarchive failed] ", err)
				return
			} else {
				fmt.Println("Swithy-Omega Download Succeeded.")
				fmt.Println("[Archive Path] ", omegaPath)
			}
		}
	}
}

func (that *XrayCtrl) keepRunningQuitSignal() {
	cnf := that.getConf()
	if cnf.Proxy.XrayKeeperSock != "" {
		s := utils.NewUServer(cnf.Proxy.XrayKeeperSock)
		s.AddHandler("/quit", func(c *gin.Context) {
			go func() {
				StopSignal <- struct{}{}
			}()
			c.String(http.StatusOK, "OK")
		})
		s.Start()
	}
}

func (that *XrayCtrl) sendQuitSignal() {
	cnf := that.getConf()
	c := utils.NewUClient(cnf.Proxy.XrayKeeperSock)
	if result, err := c.GetResp("/quit", map[string]string{}); err == nil {
		if strings.Contains(result, "OK") {
			fmt.Println("gvc [xray keeper] is stopped.")
			return
		}
	}
	fmt.Println("Nothing happened, please try again.")
}

func (that *XrayCtrl) KeepRunning() {
	if that.IsXrayRunning() {
		fmt.Println("[gvc xray client] is already running")
		return
	}
	go that.keepRunningQuitSignal()
	go func() {
		for {
			if !that.IsXrayRunning() {
				that.runByScript(that.batPath, that.scriptPath)
			}
			time.Sleep(time.Second * 60)
		}
	}()
	<-StopSignal
}
