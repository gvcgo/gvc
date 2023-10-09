package utils

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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
	"unicode"

	"github.com/gogf/gf/os/genv"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
)

func WinIsAdmin() bool {
	if _, err := os.Open("C:\\Program Files\\WindowsApps"); err != nil {
		return false
	}
	return true
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
		gprint.PrintError("Cannot find current user.")
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
		gprint.PrintError(fmt.Sprintf("Cannot open file: %+v", err))
		return
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		gprint.PrintError(fmt.Sprintf("Cannot open file: %+v", err))
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func CopyDir(srcPath, desPath string) error {
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", srcPath)
		}
	}

	os.MkdirAll(desPath, os.ModePerm)
	if desInfo, err := os.Stat(desPath); err != nil {
		return err
	} else {
		if !desInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", desPath)
		}
	}

	if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
		return fmt.Errorf("same path for src and des")
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == srcPath {
			return nil
		}
		destNewPath := strings.Replace(path, srcPath, desPath, -1)

		if !f.IsDir() {
			if _, err := CopyFile(path, destNewPath); err != nil {
				return err
			}
		} else {
			if ok, _ := PathIsExist(destNewPath); !ok {
				return os.MkdirAll(destNewPath, os.ModePerm)
			}
		}
		return nil
	})
	return err
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
		gprint.PrintWarning("Some disk formats do not support symbol link: %s", "extFAT, FAT32")
		return exec.Command("cmd", "/c", "mklink", "/j", newfile, target).Run()
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

func GetExt(filename string) (ext string) {
	if strings.Contains(filename, ".tar.gz") {
		return ".tar.gz"
	}
	if strings.Contains(filename, ".zip") {
		return ".zip"
	}
	if strings.Contains(filename, "tar.xz") {
		return ".tar.xz"
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
		gprint.PrintError(fmt.Sprintf("Open file failed: %+v", err))
		return false
	}
	defer f.Close()

	var h hash.Hash
	switch strings.ToLower(cType) {
	case "sha256":
		h = sha256.New()
	case "sha1":
		h = sha1.New()
	case "sha512":
		h = sha512.New()
	default:
		gprint.PrintError(fmt.Sprintf("[Crypto] %s is not supported.", cType))
		return
	}

	if _, err = io.Copy(h, f); err != nil {
		gprint.PrintError(fmt.Sprintf("Copy file failed: %+v", err))
		return
	}

	if cSum != hex.EncodeToString(h.Sum(nil)) {
		gprint.PrintError("Checksum failed.")
		return
	}
	gprint.PrintSuccess("Checksum succeeded.")
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

func ExecuteSysCommand(collectOutput bool, args ...string) (*bytes.Buffer, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == Windows {
		args = append([]string{"/c"}, args...)
		cmd = exec.Command("cmd", args...)
	} else {
		FlushPathEnvForUnix()
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Env = os.Environ()
	var output bytes.Buffer
	if collectOutput {
		cmd.Stdout = &output
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return &output, cmd.Run()
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

func RecordVersion(version, dir string) {
	vf := filepath.Join(dir, "version")
	if ok, _ := PathIsExist(dir); !ok {
		return
	}
	if ok, _ := PathIsExist(vf); !ok {
		os.WriteFile(vf, []byte(version), 0644)
	}
}

func ReadVersion(dir string) (v string) {
	vf := filepath.Join(dir, "version")
	if content, err := os.ReadFile(vf); err == nil {
		v = string(content)
	}
	return
}

func BatchReplaceAll(str string, oldNew map[string]string) (r string) {
	r = str
	for old, new := range oldNew {
		r = strings.ReplaceAll(r, old, new)
	}
	return
}

func ContainsCJK(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			return true
		}
	}
	return false
}

func EnsureTrailingNewline(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

func FindMaxLengthOfStringList(sl []string) (max int) {
	for _, s := range sl {
		if len(s) > max {
			max = len(s)
		}
	}
	return
}

func Closeq(v interface{}) {
	if c, ok := v.(io.Closer); ok {
		silently(c.Close())
	}
}

func silently(_ ...interface{}) {}

func MakeDirs(dirs ...string) {
	for _, d := range dirs {
		if ok, _ := PathIsExist(d); !ok {
			if err := os.MkdirAll(d, os.ModePerm); err != nil {
				gprint.PrintError(fmt.Sprintf("Make dir [%s] failed: %+v", d, err))
			}
		}
	}
}
