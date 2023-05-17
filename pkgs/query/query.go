package query

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/node/register"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
	"github.com/k0kubun/go-ansi"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/schollz/progressbar/v3"
)

type Fetcher struct {
	Url          string
	PostBody     map[string]interface{}
	Timeout      time.Duration
	RetryTimes   int
	Headers      map[string]string
	Proxy        string
	NoRedirect   bool
	client       *resty.Client
	proxyEnvName string
}

func NewFetcher() *Fetcher {
	return &Fetcher{client: resty.New(), proxyEnvName: "GVC_DEFAULT_PROXY"}
}

func (that *Fetcher) setHeaders() {
	if that.client != nil || len(that.Headers) > 0 {
		for k, v := range that.Headers {
			that.client = that.client.SetHeader(k, v)
		}
	}
}

func (that *Fetcher) parseProxy() (scheme, host string, port int) {
	if that.Proxy == "" {
		that.Proxy = os.Getenv(that.proxyEnvName)
	}
	if that.Proxy == "" {
		return
	}
	if u, err := url.Parse(that.Proxy); err == nil {
		scheme = u.Scheme
		host = u.Hostname()
		port, _ = strconv.Atoi(u.Port())
		if port == 0 {
			port = 80
		}
	}
	return
}

func (that *Fetcher) SetProxyEnvName(name string) {
	if name != "" {
		that.proxyEnvName = name
	}
}

func (that *Fetcher) setProxy() {
	if that.client != nil && that.Proxy != "" {
		scheme, host, port := that.parseProxy()
		switch scheme {
		case "http", "https":
			that.client = that.client.SetProxy(that.Proxy)
		case "socks5":
			httpClient := that.client.GetClient()
			if transport, ok := httpClient.Transport.(*http.Transport); ok {
				node := &point.Point{
					Protocols: []*protocol.Protocol{
						{
							Protocol: &protocol.Protocol_Simple{
								Simple: &protocol.Simple{
									Host:             host,
									Port:             int32(port),
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
				if proxy_result, err := register.Dialer(node); err == nil {
					transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
						address, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
						if err != nil {
							return nil, fmt.Errorf("parse address failed: %w", err)
						}
						return proxy_result.Conn(ctx, address)
					}
				}
			}
		default:
			fmt.Println(color.InRed(fmt.Sprintf("Unsupported proxy: %s", that.Proxy)))
		}
	}
}

func (that *Fetcher) setMisc() {
	that.setHeaders()
	that.setProxy()
	if that.Timeout > 0 {
		that.client = that.client.SetTimeout(that.Timeout)
	}
	if that.RetryTimes > 0 {
		that.client = that.client.SetRetryCount(that.RetryTimes)
	}
	if that.NoRedirect {
		that.client = that.client.SetRedirectPolicy(resty.NoRedirectPolicy())
	}
}

func (that *Fetcher) RemoveProxy() {
	if that.client != nil {
		that.client.RemoveProxy()
	}
}

func (that *Fetcher) Get() (r *resty.Response) {
	if that.client == nil {
		fmt.Println(color.InRed("client is nil."))
		return
	} else {
		that.setMisc()
	}
	if resp, err := that.client.R().SetDoNotParseResponse(true).Get(that.Url); err != nil {
		fmt.Println(err)
	} else {
		r = resp
	}
	return
}

func (that *Fetcher) parseFilename(fPath string) (fName string) {
	dirPath := filepath.Dir(fPath)
	fName = strings.ReplaceAll(fPath, dirPath, "")
	fName = strings.TrimPrefix(fName, "/")
	fName = strings.TrimPrefix(fName, `\`)
	return
}

func (that *Fetcher) DownloadFile(fPath string, force ...bool) (size int64) {
	if that.client == nil {
		fmt.Println(color.InRed("client is nil."))
		return
	} else {
		that.setMisc()
	}
	forceToDownload := false
	if len(force) > 0 && force[0] {
		forceToDownload = true
	}
	if ok, _ := utils.PathIsExist(fPath); ok && !forceToDownload {
		fmt.Println(color.InYellow("[Downloader] File already exists."))
		return 100
	}
	if forceToDownload {
		os.RemoveAll(fPath)
	}
	if res, err := that.client.R().SetDoNotParseResponse(true).Get(that.Url); err == nil {
		outFile, err := os.Create(fPath)
		if err != nil {
			fmt.Println(color.InRed("Cannot open file"), err)
			return
		}
		defer utils.Closeq(outFile)

		defer utils.Closeq(res.RawResponse.Body)
		var dst io.Writer
		bar := progressbar.NewOptions64(
			res.RawResponse.ContentLength,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        color.InCyan("="),
				SaucerHead:    color.InCyan(">"),
				SaucerPadding: " ",
				BarStart:      color.InYellow("["),
				BarEnd:        color.InYellow("]"),
			}),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s", color.InYellow(that.parseFilename(fPath)))),
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionShowBytes(true),
			progressbar.OptionThrottle(20*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(ansi.NewAnsiStdout(), "\n")
			}),
		)
		_ = bar.RenderBlank()
		dst = io.MultiWriter(outFile, bar)

		// io.Copy reads maximum 32kb size, it is perfect for large file download too
		written, err := io.Copy(dst, res.RawResponse.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		size = written
	} else {
		fmt.Println(err)
	}
	return
}

func (that *Fetcher) GetAndSaveFile(localPath string, force ...bool) (size int64) {
	if that.client == nil {
		fmt.Println(color.InRed("client is nil."))
		return
	} else {
		that.setMisc()
	}
	forceToDownload := false
	if len(force) > 0 && force[0] {
		forceToDownload = true
	}
	if ok, _ := utils.PathIsExist(localPath); ok && !forceToDownload {
		fmt.Println(color.InYellow("[Downloader] File already exists."))
		return 100
	}
	if forceToDownload {
		os.RemoveAll(localPath)
	}
	if res, err := that.client.R().SetDoNotParseResponse(true).Get(that.Url); err == nil {
		outFile, err := os.Create(localPath)
		if err != nil {
			fmt.Println(color.InRed("Cannot open file"), err)
			return
		}
		defer utils.Closeq(outFile)

		defer utils.Closeq(res.RawResponse.Body)
		var dst io.Writer
		bar := tui.NewProgressBar(that.parseFilename(localPath), int(res.RawResponse.ContentLength))
		bar.Start()
		dst = io.MultiWriter(outFile, bar)
		// io.Copy reads maximum 32kb size, it is perfect for large file download too
		written, err := io.Copy(dst, res.RawResponse.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		size = written
	} else {
		fmt.Println(err)
	}
	return
}

func (that *Fetcher) Post() (r *resty.Response) {
	if that.client == nil {
		fmt.Println(color.InRed("client is nil."))
		return
	} else {
		that.setMisc()
	}
	if resp, err := that.client.SetDoNotParseResponse(true).R().SetBody(that.PostBody).Post(that.Url); err != nil {
		fmt.Println(err)
		return
	} else {
		r = resp
	}
	return
}
