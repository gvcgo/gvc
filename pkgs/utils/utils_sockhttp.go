package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gfile"
	"github.com/moqsien/processes/logger"
)

/*
For communication between different processes.
*/

type USocket struct {
	Name string
	Path string
}

func (that *USocket) SetSock(name string) {
	if name == "" {
		name = "default.sock"
	}
	if !strings.HasSuffix(name, ".sock") {
		name = fmt.Sprintf("%s.sock", name)
	}
	that.Name = name
	that.Path = gfile.TempDir(name)
}

func (that *USocket) CheckSock() {
	_, err := os.Stat(that.Path)
	if !os.IsNotExist(err) {
		_ = gfile.Remove(that.Path)
	}
}

type UServer struct {
	*USocket
	Engine *gin.Engine
}

func NewUServer(name string) (us *UServer) {
	us = &UServer{
		USocket: &USocket{},
		Engine:  gin.New(),
	}
	us.SetSock(name)
	return
}

func (that *UServer) Start() (err error) {
	if that.Path == "" {
		err = errors.New("no unix socket path specified")
		return
	}
	that.CheckSock()
	unixAddr, err := net.ResolveUnixAddr("unix", that.Path)
	if err != nil {
		return err
	}
	listener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		logger.Error("listening error:", err)
		return err
	}
	return http.Serve(listener, that.Engine)
}

func (that *UServer) AddHandler(route string, handler func(c *gin.Context)) {
	that.Engine.GET(route, handler)
}

type UClient struct {
	*USocket
	Client *http.Client
	params string
}

func NewUClient(name string) (uc *UClient) {
	uc = &UClient{
		USocket: &USocket{},
		Client:  &http.Client{},
	}
	uc.SetSock(name)
	uc.initClient()
	return
}

func (that *UClient) initClient() {
	that.Client.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", that.Path)
		},
	}
}

func (that *UClient) parseParms(params map[string]string) {
	that.params = ""
	for k, v := range params {
		if len(that.params) == 0 {
			that.params += fmt.Sprintf("?%s=%s", k, v)
		} else {
			that.params += fmt.Sprintf("&%s=%s", k, v)
		}
	}
}

func (that *UClient) GetResp(route string, params map[string]string) (string, error) {
	that.parseParms(params)
	url := fmt.Sprintf("http://%s/%s/%s",
		that.Name,
		strings.Trim(route, "/"),
		that.params)

	resp, err := that.Client.Get(url)
	if err != nil {
		return "", err
	}
	if result, err := io.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}
