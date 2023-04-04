package utils

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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

func GetPathForWindows() string {
	content := os.Getenv("PATH")
	return strings.ReplaceAll(content, ";;", ";")
}

func FormatPathForWindows(newPath string) string {
	l := []string{}
	old := GetPathForWindows()
	for _, v := range strings.Split(newPath, ";") {
		if !strings.Contains(old, v) {
			l = append(l, v)
		}
	}
	return fmt.Sprintf("%s;%s", old, strings.Join(l, ";"))
}

func SetWinEnv(key, value string, isSys ...bool) (err error) {
	if runtime.GOOS == Windows {
		// handle path for windows.
		k := strings.ToLower(key)
		if k == "path" {
			value = FormatPathForWindows(value)
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
	FlushPathEnvForUnix()
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

func ReplaceFileContent(filePath, old, new string, perm fs.FileMode) {
	if ok, _ := PathIsExist(filePath); ok {
		if oldContent, err := os.ReadFile(filePath); err == nil {
			if !strings.Contains(string(oldContent), new) {
				newContent := strings.Replace(string(oldContent), old, new, -1)
				if len(newContent) > 0 {
					os.WriteFile(filePath, []byte(newContent), perm)
				}
			}
		}
	}
}

func ConvertStrToReader(str string) io.Reader {
	return bytes.NewReader([]byte(str))
}

func DecodeBase64(str string) (res string) {
	s, _ := base64.StdEncoding.DecodeString(str)
	res = string(s)
	return
}

func ParseArch(name string) string {
	name = strings.ToLower(name)
	for k, v := range ArchMap {
		if strings.Contains(name, k) {
			return v
		}
	}
	return ""
}

func ParsePlatform(name string) string {
	name = strings.ToLower(name)
	for k, v := range PlatformMap {
		if strings.Contains(name, k) {
			return v
		}
	}
	return ""
}
