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

var ArchOSs map[string]string = map[string]string{
	"x86-64":     "amd64",
	"x86":        "386",
	"arm64":      "arm64",
	"armv6":      "arm",
	"ppc64le":    "ppc64le",
	"macos":      "darwin",
	"os x 10.8+": "darwin",
	"os x 10.6+": "darwin",
	"linux":      "linux",
	"windows":    "windows",
	"freebsd":    "freebsd",
}

func MapArchAndOS(ArchOrOS string) (result string) {
	result, ok := ArchOSs[strings.ToLower(ArchOrOS)]
	if !ok {
		result = ArchOrOS
	}
	return
}

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

func PathIsExist(path string) (bool, error) {
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
