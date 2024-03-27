package git

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
)

var (
	GitSSHProxyCommandWin  string = `ProxyCommand connect -S %s %s %s`
	GitSSHProxyCommandNix  string = `ProxyCommand nc -v -x %s %s %s`
	GitSSHProxyCommandHttp string = `ProxyCommand g g cs -a %s -p %s`
)

var GitSSHConfigStr string = `Host github.com
  User git
  Port 22
  Hostname github.com
  IdentityFile "%s"
  TCPKeepAlive yes
  %s

Host ssh.github.com
  User git
  Port 443
  Hostname ssh.github.com
  IdentityFile "%s"
  TCPKeepAlive yes
  %s
`

/*
Adds proxy to hosts file.
*/
func SetProxyForSSH() {
	cfg := conf.NewGVConfig()
	pURI := cfg.GetLocalProxy()
	if pURI == "" {
		gprint.PrintError("No legal proxy is specified.")
		return
	}
	u, err := url.Parse(pURI)
	if err != nil {
		return
	}
	homeDir, _ := os.UserHomeDir()
	dotSSHPath := filepath.Join(homeDir, ".ssh")
	idRSAPath := filepath.Join(dotSSHPath, "id_rsa")
	if ok, _ := gutils.PathIsExist(idRSAPath); !ok {
		gprint.PrintError("Cannot find ~/.ssh/id_rsa.")
		return
	}
	uStr := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	pxyCmd := ""

	switch runtime.GOOS {
	case gutils.Windows:
		if strings.Contains(u.Scheme, "sock") {
			pxyCmd = fmt.Sprintf(
				GitSSHProxyCommandWin,
				uStr,
				`%h`,
				`%p`,
			)
		} else {
			pxyCmd = fmt.Sprintf(
				GitSSHProxyCommandHttp,
				`%h`,
				`%p`,
			)
		}
	case gutils.Linux, gutils.Darwin:
		if strings.Contains(u.Scheme, "sock") {
			pxyCmd = fmt.Sprintf(
				GitSSHProxyCommandNix,
				uStr,
				`%h`,
				`%p`,
			)
		} else {
			pxyCmd = fmt.Sprintf(
				GitSSHProxyCommandHttp,
				`%h`,
				`%p`,
			)
		}
	default:
		gprint.PrintError("Unsupported OS.")
		return
	}

	content := fmt.Sprintf(
		GitSSHConfigStr,
		idRSAPath,
		pxyCmd,
		idRSAPath,
		pxyCmd,
	)
	setProxyForSSH(dotSSHPath, content)
}

func setProxyForSSH(dotSSHPath, content string) {
	confPath := filepath.Join(dotSSHPath, "config")
	if ok, _ := gutils.PathIsExist(confPath); !ok {
		os.WriteFile(confPath, []byte(content), 0o666)
	} else {
		oldContentByte, _ := os.ReadFile(confPath)
		oldContent := string(oldContentByte)
		if !strings.Contains(oldContent, "ProxyCommand") && len(oldContent) > 0 {
			os.WriteFile(confPath, []byte(oldContent+"\n"+content), os.ModePerm)
		} else {
			os.WriteFile(confPath, []byte(content), 0o666)
		}
	}
}

func ToggleProxyForSSH() {
	homeDir, _ := os.UserHomeDir()
	confPath := filepath.Join(homeDir, ".ssh", "config")
	backupConfPath := filepath.Join(homeDir, ".ssh", "config.bak")

	ok1, _ := gutils.PathIsExist(confPath)
	ok2, _ := gutils.PathIsExist(backupConfPath)

	if !ok1 && !ok2 {
		gprint.PrintWarning("Set a proxy for ssh...")
		SetProxyForSSH()
	} else if ok1 && !ok2 {
		if err := os.Rename(confPath, backupConfPath); err == nil {
			gprint.PrintInfo("Proxy disabled.")
		}
	} else if !ok1 && ok2 {
		if err := os.Rename(backupConfPath, confPath); err == nil {
			gprint.PrintSuccess("Proxy enabled.")
		}
	} else {
		os.RemoveAll(backupConfPath)
		if err := os.Rename(confPath, backupConfPath); err == nil {
			gprint.PrintInfo("Proxy disabled.")
		}
	}
}
