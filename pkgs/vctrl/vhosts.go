package vctrl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/panjf2000/ants/v2"
	ping "github.com/prometheus-community/pro-bing"
	"github.com/schollz/progressbar/v3"
)

const (
	HEAD        = "# FromGhosts Start"
	TAIL        = "# FromGhosts End"
	TIME        = "# UpdateTime: %s"
	LinePattern = "%s\t\t\t%s # %s"
)

type taskArgs struct {
	IP  string
	URL string
}

type host struct {
	IP     string        // ip address
	AvgRTT time.Duration // average RTT
}

type hostList map[string]host // key: host name, value: host

type Hosts struct {
	Conf     *config.GVConfig
	filePath string
	rawList  map[string]string
	hList    hostList
	lineReg  *regexp.Regexp
	hostReg  *regexp.Regexp
	lock     *sync.Mutex
	wg       sync.WaitGroup
	pool     *ants.PoolWithFunc
	bar      *progressbar.ProgressBar
	*downloader.Downloader
}

func NewHosts() *Hosts {
	conf := config.New()
	lineReg := `((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`
	hostsReg := fmt.Sprintf(`%s[\s\S]*%s`, HEAD, TAIL)
	return &Hosts{
		Conf:       conf,
		filePath:   config.GetHostsFilePath(),
		rawList:    make(map[string]string, 200),
		hList:      make(hostList, 200),
		lineReg:    regexp.MustCompile(lineReg),
		hostReg:    regexp.MustCompile(hostsReg),
		lock:       &sync.Mutex{},
		wg:         sync.WaitGroup{},
		Downloader: &downloader.Downloader{},
	}
}

func (that *Hosts) extractHostUrl(text, ip string) string {
	raw := strings.Replace(text, ip, "", -1)
	return strings.TrimSpace(raw)
}

func (that *Hosts) ParseHosts(content []byte) {
	sc := bufio.NewScanner(strings.NewReader(string(content)))
	for sc.Scan() {
		text := sc.Text()
		ipList := that.lineReg.FindAllString(text, -1)
		if len(ipList) == 1 {
			ip_ := ipList[0]
			url := that.extractHostUrl(text, ip_)
			if url == "" {
				continue
			}
			if _, ok := that.rawList[ip_]; !ok {
				that.rawList[ip_] = url
			}
		}
	}
}

func (that *Hosts) GetHosts() {
	resps := make(chan *http.Response, 10)
	that.Timeout = time.Duration(that.Conf.Hosts.ReqTimeout) * time.Second
	for _, url := range that.Conf.Hosts.SourceUrls {
		that.wg.Add(1)
		var url_ string = url
		go func() {
			defer that.wg.Done()
			that.Url = url_
			resp := that.GetUrl()
			resps <- resp
		}()
	}
	that.wg.Wait()
	close(resps)
	for r := range resps {
		if r != nil {
			content, err := io.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				fmt.Println("[Read Body Errored] ", err)
				return
			}
			that.ParseHosts(content)
		}
	}
}

func (that *Hosts) toSave(url string) bool {
	filters := that.Conf.Hosts.HostFilters
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if strings.Contains(url, filter) {
			return true
		}
	}
	return false
}

func (that *Hosts) pingHosts(args interface{}) {
	var url, ip string
	defer that.wg.Done()
	defer that.bar.Add(1)
	if tArgs, ok := args.(*taskArgs); !ok {
		return
	} else {
		url = tArgs.URL
		ip = tArgs.IP
	}
	if !that.toSave(url) {
		return
	}
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Println("Ping hosts errored: ", err)
		return
	}
	pinger.Count = that.Conf.Hosts.PingCount
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}
	pinger.Timeout = time.Duration(that.Conf.Hosts.ReqTimeout) * time.Millisecond
	pinger.OnFinish = func(statics *ping.Statistics) {
		if len(statics.Rtts) > 0 {
			that.lock.Lock()
			if old, ok := that.hList[url]; !ok {
				that.hList[url] = host{IP: ip, AvgRTT: statics.AvgRtt}
			} else {
				if old.AvgRTT > statics.AvgRtt {
					that.hList[url] = host{IP: ip, AvgRTT: statics.AvgRtt}
				}
			}
			that.lock.Unlock()
		}
	}
	err = pinger.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func (that *Hosts) ReadAndBackupHosts(hPath, hBackupPath string) (content []byte) {
	var (
		err  error
		file *os.File
	)
	file, err = os.Open(hPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err = io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	if utils.GetShell() == "win" {
		err = os.WriteFile(hBackupPath, content, 0644)
	} else {
		err = utils.CopyFileOnUnixSudo(hPath, hBackupPath)
	}
	if err != nil {
		fmt.Println("Hosts file backup failed: ", err)
		return
	}
	return
}

func (that *Hosts) replace(oldContent []byte, newHostStr string) string {
	old := string(oldContent)
	if strings.Contains(old, HEAD) {
		return that.hostReg.ReplaceAllString(old, newHostStr)
	} else {
		if newHostStr != "" {
			return fmt.Sprintf("%s\n%s", oldContent, newHostStr)
		}
		return old
	}
}

func (that *Hosts) FormatAndSaveHosts(oldContent []byte) {
	if len(that.hList) > 0 {
		lineList := []string{}
		for url, h := range that.hList {
			lineList = append(lineList, fmt.Sprintf(LinePattern, h.IP, url, h.AvgRTT))
		}
		loc, _ := time.LoadLocation("Asia/Shanghai")
		if len(oldContent) < 1 {
			return
		}
		newHostStr := fmt.Sprintf("%s\n%s\n%s\n%s",
			HEAD,
			fmt.Sprintf(TIME, time.Now().In(loc).Format("2006-01-02 15:04:05")),
			strings.Join(lineList, "\n"),
			TAIL)
		newStr := that.replace(oldContent, newHostStr)
		if newStr == "" {
			return
		}
		var err error
		if utils.GetShell() == "win" {
			err = os.WriteFile(config.GetHostsFilePath(), []byte(newStr), 0666)
		} else {
			err = os.WriteFile(config.TempHostsFilePath, []byte(newStr), 0666)
			if err == nil {
				err = utils.CopyFileOnUnixSudo(config.TempHostsFilePath, config.GetHostsFilePath())
			}
		}
		if err != nil {
			fmt.Println("\nWrite file errored: ", err)
			return
		}
		fmt.Println("\nSuccessed!")
	}
}

func (that *Hosts) Run() {
	that.GetHosts()
	hostFilePath := config.GetHostsFilePath()
	hostBackupFilePath := filepath.Join(filepath.Dir(hostFilePath), "hosts.backup")
	oldContent := that.ReadAndBackupHosts(hostFilePath, hostBackupFilePath)

	length := len(that.rawList)
	bar := progressbar.NewOptions(length,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionSetDescription("[gvc] ping hosts..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	that.bar = bar

	if pool, err := ants.NewPoolWithFunc(that.Conf.Hosts.WorkerNum, func(args interface{}) {
		that.pingHosts(args)
	}); err != nil {
		return
	} else {
		that.pool = pool
	}

	defer that.pool.Release()
	for k, v := range that.rawList {
		that.wg.Add(1)
		err := that.pool.Invoke(&taskArgs{
			IP:  k,
			URL: v,
		})
		if err != nil {
			fmt.Println("[Invoke task failed] ", err)
		}
	}
	that.wg.Wait()
	time.Sleep(1 * time.Second)
	fmt.Printf("Find available hosts: <%v/%v(raw)>", len(that.hList), len(that.rawList))
	that.FormatAndSaveHosts(oldContent)
}

func (that *Hosts) ShowFilePath() {
	fmt.Println("HostsFile: ", config.GetHostsFilePath())
}
