package fetcher

import (
	"context"
	"fmt"
	"io"
	"io/fs"
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
	"github.com/gocolly/colly/v2"
	"github.com/k0kubun/go-ansi"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	Url              string
	Timeout          time.Duration
	PostBody         []byte
	Headers          map[string]string
	ManuallyRedirect bool
	Proxy            string
	SetUA            bool
}

func (that *Downloader) setHeaders(req *http.Request) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	if (that.Headers == nil || len(that.Headers) == 0) && that.SetUA {
		that.Headers = map[string]string{
			"User-Agent": ua,
		}
	}
	if req != nil {
		for key, value := range that.Headers {
			req.Header.Add(key, value)
		}
	}
}

func (that *Downloader) parseProxy() (scheme, host string, port int) {
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

func (that *Downloader) getClient() *http.Client {
	httpTransport, _ := http.DefaultTransport.(*http.Transport)
	scheme, host, port := that.parseProxy()
	switch scheme {
	case "http", "https":
		if pUri, err := url.Parse(that.Proxy); err == nil {
			httpTransport.Proxy = http.ProxyURL(pUri)
		}
	case "socks5":
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
			httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				address, err := proxy.ParseAddress(proxy.PaseNetwork(network), addr)
				if err != nil {
					return nil, fmt.Errorf("parse address failed: %w", err)
				}
				return proxy_result.Conn(ctx, address)
			}
		}
	default:
	}
	client := &http.Client{
		Timeout: that.Timeout,
	}
	if that.Timeout != 0 {
		if that.ManuallyRedirect {
			httpTransport.TLSHandshakeTimeout = that.Timeout
		}
		client = &http.Client{
			Transport: httpTransport,
			Timeout:   that.Timeout,
		}
	} else {
		client = &http.Client{Transport: httpTransport}
	}
	return client
}

func (that *Downloader) GetUrl() *http.Response {
	if that.Url != "" && utils.VerifyUrls(that.Url) {
		var (
			r   *http.Response
			err error
		)
		httpReq, err := http.NewRequest(http.MethodGet, that.Url, nil)
		if err != nil {
			fmt.Println("[New Request Failed] ", err)
			return nil
		}

		that.setHeaders(httpReq)
		client := that.getClient()
		if !that.ManuallyRedirect {
			r, err = client.Do(httpReq)
		} else {
			r, err = client.Transport.RoundTrip(httpReq)
		}
		if err != nil {
			fmt.Println("[Request Errored] URL: ", that.Url, ", err: ", err)
			return nil
		}
		return r
	} else {
		fmt.Println("[Illegal URL] URL: ", that.Url)
		return nil
	}
}

func (that *Downloader) GetWithColly() (resp []byte) {
	c := colly.NewCollector()
	if that.Proxy == "" {
		c.SetProxy(that.Proxy)
	}
	if that.Timeout != 0 {
		c.SetRequestTimeout(that.Timeout)
	}
	c.OnResponse(func(r *colly.Response) {
		resp = r.Body
	})
	c.Visit(that.Url)
	return
}

func (that *Downloader) parseFilename(fPath string) (fName string) {
	dirPath := filepath.Dir(fPath)
	fName = strings.ReplaceAll(fPath, dirPath, "")
	fName = strings.TrimPrefix(fName, "/")
	fName = strings.TrimPrefix(fName, `\`)
	return fmt.Sprintf("<%s>", fName)
}

func (that *Downloader) GetFile(fPath string, flag int, perm fs.FileMode, force ...bool) (size int64) {
	forceToDownload := false
	if len(force) > 0 && force[0] {
		forceToDownload = true
	}
	if ok, _ := utils.PathIsExist(fPath); ok && !forceToDownload {
		fmt.Println("[Downloader] File already exists.")
		return 100
	}

	resp := that.GetUrl()
	if resp == nil {
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	f, err := os.OpenFile(fPath, flag, perm)
	if err != nil {
		fmt.Println("[open file failed] ", err)
		return
	}
	defer f.Close()
	var dst io.Writer
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        color.InGreen("="),
			SaucerHead:    color.InGreen(">"),
			SaucerPadding: " ",
			BarStart:      color.InGreen("["),
			BarEnd:        color.InGreen("]"),
		}),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s", color.InYellow(that.parseFilename(fPath)))),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(ansi.NewAnsiStdout(), "\n")
		}),
	)
	_ = bar.RenderBlank()
	dst = io.MultiWriter(f, bar)
	size, _ = io.Copy(dst, resp.Body)
	return
}

func (that *Downloader) PostUrl() (httpRsp *http.Response) {
	if that.Url != "" && len(that.PostBody) > 0 {

		reqBody := strings.NewReader(string(that.PostBody))
		httpReq, err := http.NewRequest(http.MethodPost, that.Url, reqBody)
		if err != nil {
			fmt.Println("[New Request Failed] ", err)
			return nil
		}
		that.setHeaders(httpReq)
		client := that.getClient()
		if that.ManuallyRedirect {
			httpRsp, err = client.Transport.RoundTrip(httpReq)
		} else {
			httpRsp, err = client.Do(httpReq)
		}
		if err != nil {
			fmt.Println("[Request Failed] ", err)
			return nil
		}
		return httpRsp
	} else {
		fmt.Println("[Illegal URL] URL: ", that.Url)
		return nil
	}
}

func (that *Downloader) PostWithColly() (resp []byte) {
	c := colly.NewCollector()
	if that.Proxy == "" {
		c.SetProxy(that.Proxy)
	}
	if that.Timeout != 0 {
		c.SetRequestTimeout(that.Timeout)
	}
	c.OnResponse(func(r *colly.Response) {
		resp = r.Body
	})
	c.PostRaw(that.Url, that.PostBody)
	return
}
