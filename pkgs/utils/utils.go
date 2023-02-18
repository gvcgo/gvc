package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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

const (
	Win  string = "win"
	Zsh  string = "zsh"
	Bash string = "bash"
)

func GetShell() (shell string) {
	if strings.Contains(runtime.GOOS, "window") {
		return Win
	}
	s := os.Getenv("SHELL")
	if strings.Contains(s, "zsh") {
		return Zsh
	}
	return Bash
}

func GetShellRcFile() (rc string) {
	shell := GetShell()
	switch shell {
	case Zsh:
		rc = filepath.Join(GetHomeDir(), ".zshrc")
	case Bash:
		rc = filepath.Join(GetHomeDir(), ".bashrc")
	default:
		rc = Win
	}
	return
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
	cmd := exec.Command("sudo", "cp", "-R", from, to)
	cmd.Env = genv.All()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func WindowsSetEnv(key, value string) error {
	cmd := exec.Command("setx", key, value)
	cmd.Env = genv.All()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func CopyFile(src, dst string) (written int64, err error) {
	srcFile, err := os.Open(src)

	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func MkSymLink(target, newfile string) (err error) {
	if runtime.GOOS == "windows" {
		if err = exec.Command("cmd", "/c", "mklink", "/j", newfile, target).Run(); err == nil {
			return nil
		}
	}
	return os.Symlink(target, newfile)
}

func SetWinEnv(key, value string, isSys ...bool) (err error) {
	if runtime.GOOS == "windows" {
		var args []string
		if len(isSys) > 0 && isSys[0] {
			args = []string{"/c", "setx", key, value, "/m"}
		} else {
			args = []string{"/c", "setx", key, value}
		}
		err = exec.Command("cmd", args...).Run()
	}
	return
}

func GetExt(filename string) (ext string) {
	if strings.Contains(filename, ".tar.gz") {
		return ".tar.gz"
	}
	if strings.Contains(filename, ".zip") {
		return ".zip"
	}
	if strings.Contains(filename, ".") {
		l := strings.Split(filename, ".")
		return l[len(l)-1]
	}
	return
}

func CheckFile(fpath, cType, cSum string) (r bool) {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println("[Open file failed] ", err)
		return false
	}
	defer f.Close()

	var h hash.Hash
	switch strings.ToLower(cType) {
	case "sha256":
		h = sha256.New()
	case "sha1":
		h = sha1.New()
	default:
		fmt.Println("[Crypto] ", cType, " not supported.")
		return
	}

	if _, err = io.Copy(h, f); err != nil {
		fmt.Println("[Copy file failed] ", err)
		return
	}

	if cSum != hex.EncodeToString(h.Sum(nil)) {
		fmt.Println("Checksum failed.")
		return
	}
	fmt.Println("Checksum successed.")
	return true
}
