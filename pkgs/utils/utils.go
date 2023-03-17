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

const (
	Windows string = "windows"
	MacOS   string = "darwin"
	Linux   string = "linux"
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
	if runtime.GOOS == Windows {
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

func RunCommand(args ...string) {
	var cmd *exec.Cmd
	if runtime.GOOS == Windows {
		args = append([]string{"/c"}, args...)
		cmd = exec.Command("cmd", args...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func MkSymLink(target, newfile string) (err error) {
	if runtime.GOOS == Windows {
		if err = exec.Command("cmd", "/c", "mklink", "/j", newfile, target).Run(); err == nil {
			return nil
		}
	}
	return os.Symlink(target, newfile)
}

func GetWinAppdataEnv() string {
	return os.Getenv("APPDATA")
}

func SetWinEnv(key, value string, isSys ...bool) (err error) {
	if runtime.GOOS == Windows {
		// handle path for windows.
		k := strings.ToLower(key)
		if k == "path" {
			oldPath := strings.Trim(os.Getenv("Path"), ";")
			if strings.Contains(oldPath, value) {
				return
			}
			value = fmt.Sprintf("%s;%s", oldPath, value)
		}
		var c *exec.Cmd
		if len(isSys) > 0 && isSys[0] {
			c = exec.Command("setx", key, value, "/m")
		} else {
			c = exec.Command("setx", key, value)
		}
		c.Env = os.Environ()
		err = c.Run()
		if k == "path" {
			fmt.Println("!!!Close current window to make envs valid!!!")
		}
	}
	return
}

func WinCmdExit() {
	if runtime.GOOS == Windows {
		exec.Command("exit()").Run()
	}
}

func SetUnixEnv(envars string) {
	shellrc := GetShellRcFile()
	if file, err := os.Open(shellrc); err == nil {
		defer file.Close()
		content, err := io.ReadAll(file)
		if err == nil {
			c := string(content)
			os.WriteFile(fmt.Sprintf("%s.backup", shellrc), content, 0644)
			flag := strings.Split(envars, "\n")[0]
			if strings.Contains(c, flag) {
				return
			}
			s := fmt.Sprintf("%v\n%s", c, envars)
			os.WriteFile(shellrc, []byte(strings.ReplaceAll(s, GetHomeDir(), "$HOME")), 0644)
		}
	}
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

func JoinUnixFilePath(pathList ...string) (r string) {
	newList := []string{}
	for _, p := range pathList {
		newList = append(newList, strings.Trim(p, "/"))
	}
	r = strings.Join(newList, "/")
	if !strings.HasPrefix(r, "/") {
		r = fmt.Sprintf("%s%s", "/", r)
	}
	return
}

func FlushPathEnvForUnix() (err error) {
	if runtime.GOOS != Windows {
		return exec.Command("source", GetShellRcFile()).Run()
	}
	return
}

func ExecuteCommand(args ...string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == Windows {
		args = append([]string{"/c"}, args...)
		cmd = exec.Command("cmd", args...)
	} else {
		FlushPathEnvForUnix()
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func ClearDir(dirPath string) {
	fList, _ := os.ReadDir(dirPath)
	for _, f := range fList {
		os.RemoveAll(filepath.Join(dirPath, f.Name()))
	}
}
