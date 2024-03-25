package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/goutils/pkgs/request"
	"github.com/gvcgo/gvc/conf"
)

/*
Modifies hosts file.
*/
var remoteList = []string{
	"https://raw.githubusercontent.com/JohyC/Hosts/main/MicrosoftHosts.txt",
	"https://raw.githubusercontent.com/JohyC/Hosts/main/EpicHosts.txt",
	"https://raw.githubusercontent.com/JohyC/Hosts/main/SteamDomains.txt",
	"https://raw.githubusercontent.com/JohyC/Hosts/main/hosts.txt",
	"https://raw.githubusercontent.com/ineo6/hosts/master/next-hosts",
	"https://raw.githubusercontent.com/sengshinlee/hosts/main/hosts",
}

const (
	HostFilePathForNix = "/etc/hosts"
	HostFilePathForWin = `C:\Windows\System32\drivers\etc\hosts`
)

var StrWrapper string = `# FromGhosts Start
# UpdateTime: %s
%s
# FromGhosts End`

var StrRegExp = regexp.MustCompile(`# FromGhosts Start[\w\W]*# FromGhosts End`)

var WinScript string = `COPY %s %s`

func getHostsFilePath() string {
	switch runtime.GOOS {
	case gutils.Windows:
		return HostFilePathForWin
	default:
		return HostFilePathForNix
	}
}

func getTempFilePath() string {
	return filepath.Join(conf.GetGVCWorkDir(), "hosts.temp.txt")
}

type HostsModifier struct {
	items   map[string]string
	fetcher *request.Fetcher
	cfg     *conf.GVConfig
}

func NewModifier() *HostsModifier {
	return &HostsModifier{
		items:   make(map[string]string, 30),
		fetcher: request.NewFetcher(),
		cfg:     conf.NewGVConfig(),
	}
}

func (h *HostsModifier) GetHostsFiles() {
	for _, u := range remoteList {
		rp := h.cfg.GetReverseProxy()
		u = strings.TrimRight(rp, "/") + "/" + u
		h.fetcher.SetUrl(u)
		h.fetcher.Timeout = time.Second * 20
		r, _ := h.fetcher.GetString()
		h.Parse(r)
	}
}

func (h *HostsModifier) Parse(resp string) {
	lines := strings.Split(resp, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		h.add(line)
	}
}

func (h *HostsModifier) add(line string) {
	sList := strings.Split(line, " ")
	if len(sList) < 2 {
		return
	}
	key := sList[0]
	for _, s := range sList[1:] {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		h.items[key] = s
		break
	}
}

func (h *HostsModifier) PrepareTempFile() (ok bool) {
	if len(h.items) == 0 {
		return
	}
	lines := []string{}
	for k, v := range h.items {
		lines = append(lines, fmt.Sprintf("%s\t\t\t\t\t\t%s", k, v))
	}
	content, _ := os.ReadFile(getHostsFilePath())
	newStr := StrRegExp.ReplaceAllString(
		string(content),
		fmt.Sprintf(StrWrapper, time.Now().Format("2006-01-02 15:04:05"), strings.Join(lines, "\n")),
	)
	err := os.WriteFile(getTempFilePath(), []byte(newStr), os.ModePerm)
	if err != nil {
		return
	}
	return true
}

func (h *HostsModifier) copyAsSudo(src, dst string) {
	if runtime.GOOS != gutils.Windows {
		gutils.ExecuteSysCommand(false, "", "sudo", "cp", "-rf", src, dst)
	} else {
		script := fmt.Sprintf(WinScript, src, dst)
		scriptPath := filepath.Join(conf.GetGVCWorkDir(), "win_script_temp.ps1")
		if err := os.WriteFile(scriptPath, []byte(script), os.ModePerm); err == nil {
			gutils.ExecuteSysCommand(false, "",
				"powershell", "Start-Process", "-verb", "runas", scriptPath)
		}
		os.RemoveAll(scriptPath)
	}
}

func (h *HostsModifier) BackupOldFile() {
	hostFile := getHostsFilePath()
	backupFile := filepath.Join(filepath.Dir(hostFile), "hosts.backup")
	h.copyAsSudo(hostFile, backupFile)
}

func (h *HostsModifier) CopyTempFile() {
	h.BackupOldFile()
	h.copyAsSudo(getTempFilePath(), getHostsFilePath())
}

func (h *HostsModifier) Run() {
	h.GetHostsFiles()
	if ok := h.PrepareTempFile(); ok {
		h.CopyTempFile()
	}
}
