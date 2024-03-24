package asciinema

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cirocosta/asciinema-edit/cast"
	"github.com/cirocosta/asciinema-edit/commands/transformer"
	"github.com/cirocosta/asciinema-edit/editor"
	acmd "github.com/gvcgo/asciinema/cmd"
	autil "github.com/gvcgo/asciinema/util"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
)

func getName(base string) string {
	if base == "" {
		return base
	}
	return strings.Split(base, ".")[0]
}

func handleFilePath(fpath string) (title, result string) {
	cwd, _ := os.Getwd()
	if fpath == "" {
		return "default_cast", filepath.Join(cwd, "default.cast")
	}
	base := filepath.Base(fpath)
	if base == fpath {
		return getName(base), filepath.Join(cwd, base)
	}
	return getName(base), fpath
}

/*
Github: https://github.com/asciinema
Homepage: https://asciinema.org/
*/
type Asciinema struct {
	cmd *acmd.Runner
}

func NewAsciinema() *Asciinema {
	os.Setenv(autil.DefaultHomeEnv, conf.GetGVCWorkDir())
	a := &Asciinema{
		cmd: acmd.New(),
	}
	return a
}

// Authorization to asciinema.org
func (a *Asciinema) Auth() error {
	authUrl, info := a.cmd.Auth()
	gprint.PrintInfo(info)
	var cmd *exec.Cmd
	if runtime.GOOS == gutils.Darwin {
		cmd = exec.Command("open", authUrl)
	} else if runtime.GOOS == gutils.Linux {
		cmd = exec.Command("x-www-browser", authUrl)
	} else if runtime.GOOS == gutils.Windows {
		cmd = exec.Command("cmd", "/c", "start", authUrl)
	} else {
		return fmt.Errorf("unsupported os")
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Records an asciinema cast.
func (a *Asciinema) Record(fPath string) error {
	a.cmd.Title, a.cmd.FilePath = handleFilePath(fPath)
	return a.cmd.Rec()
}

// Plays an asciinema cast.
func (a *Asciinema) Play(fPath string) error {
	a.cmd.Title, a.cmd.FilePath = handleFilePath(fPath)
	return a.cmd.Play()
}

// Uploads an asciinema cast to asciinema.org.
func (a *Asciinema) Upload(fPath string) error {
	a.cmd.Title, a.cmd.FilePath = handleFilePath(fPath)
	if respStr, err := a.cmd.Upload(); err == nil {
		gprint.PrintInfo(respStr)
		return err
	}
	return nil
}

/*
Cut an asciinema cast.
*/

type cutTransformation struct {
	from float64
	to   float64
}

func (t *cutTransformation) Transform(c *cast.Cast) (err error) {
	err = editor.Cut(c, t.from, t.to)
	return
}

func (a *Asciinema) Cut(fPath, outFilePath string, start, end float64) (err error) {
	transformation := &cutTransformation{
		from: start,
		to:   end,
	}
	t, err := transformer.New(transformation, fPath, outFilePath)
	if err != nil {
		return err
	}
	defer t.Close()
	return t.Transform()
}

/*
Modifies the playing speed of an asciinema cast.
*/
type speedTransformation struct {
	from   float64
	to     float64
	factor float64
}

func (t *speedTransformation) Transform(c *cast.Cast) (err error) {
	if t.from == 0 && t.to == 0 {
		t.from = c.EventStream[0].Time
		t.to = c.EventStream[len(c.EventStream)-1].Time
	}

	err = editor.Speed(c, t.factor, t.from, t.to)
	return
}

func (a *Asciinema) Speed(fPath, outFilePath string, factor, start, end float64) (err error) {
	transformation := &speedTransformation{
		factor: factor,
		from:   start,
		to:     end,
	}
	t, err := transformer.New(transformation, fPath, outFilePath)
	if err != nil {
		return err
	}
	defer t.Close()
	return t.Transform()
}
