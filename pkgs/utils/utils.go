package utils

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/gogf/gf/os/genv"
)

func VerifyUrls(rawUrl string) (r bool) {
	r = true
	_, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		r = false
		return
	}
	url, err := url.Parse(rawUrl)
	if err != nil || url.Scheme == "" || url.Host == "" {
		r = false
		return
	}
	return
}

func GetShell() (shell string) {
	if strings.Contains(runtime.GOOS, "window") {
		return "win"
	}
	s := os.Getenv("SHELL")
	if strings.Contains(s, "zsh") {
		return "zsh"
	}
	return "bash"
}

func GetHomeDir() (homeDir string) {
	u, err := user.Current()
	if err != nil {
		fmt.Println("[CurrentUser]", err)
		return
	}
	return u.HomeDir
}

func PahtIsExist(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

func CopyFileOnUnixSudo(from, to string) error {
	cmd := exec.Command("sudo", "cp", from, to)
	cmd.Env = genv.All()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
