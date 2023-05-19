package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/ctrl"
)

var cmdName string = func() string {
	epath, err := os.Executable()
	if err != nil {
		tui.PrintError(fmt.Sprintf("cannot find executable path: %+v", err))
		os.Exit(1)
	}
	return epath
}()

const (
	XtrayStarterCmd = "xtray-runner"
	XtrayKeeperCmd  = "xtray-keeper"
)

var (
	starterBatPath = filepath.Join(config.ProxyFilesDir, "starter.bat")
	keeperBatPath  = filepath.Join(config.ProxyFilesDir, "keeper.bat")
)

// Start-Process -WindowStyle hidden -FilePath "executable path"
func GenStarter() (starter *exec.Cmd) {
	if runtime.GOOS == utils.Windows {
		if ok, _ := utils.PathIsExist(starterBatPath); !ok {
			content := fmt.Sprintf("%s %s", cmdName, XtrayStarterCmd)
			os.WriteFile(starterBatPath, []byte(content), 0777)
		}
		starter = exec.Command("powershell", "Start-Process", "-WindowStyle", "hidden", "-FilePath", starterBatPath)
	} else {
		starter = exec.Command(cmdName, XtrayStarterCmd)
	}
	return
}

func GenKeeper() (keeper *exec.Cmd) {
	if runtime.GOOS == utils.Windows {
		if ok, _ := utils.PathIsExist(keeperBatPath); !ok {
			content := fmt.Sprintf("%s %s", cmdName, XtrayKeeperCmd)
			os.WriteFile(keeperBatPath, []byte(content), 0777)
		}
		keeper = exec.Command("powershell", "Start-Process", "-WindowStyle", "hidden", "-FilePath", keeperBatPath)
	} else {
		keeper = exec.Command(cmdName, XtrayKeeperCmd)
	}
	return
}

type XtrayExa struct {
	GVConf *config.GVConfig // gvc configuration
	Conf   *conf.Conf       // xtray configuration
	Runner *ctrl.XRunner
	Keeper *ctrl.XKeeper
}

func NewXtrayExa() *XtrayExa {
	xe := &XtrayExa{
		GVConf: config.New(),
	}

	xe.Conf = conf.NewConf()
	xe.Conf.FetcherUrl = xe.GVConf.Xtray.FetcherUrl
	xe.Conf.WorkDir = xe.GVConf.Xtray.WorkDir
	xe.Conf.RawProxyFile = xe.GVConf.Xtray.RawProxyFile
	xe.Conf.ProxyFile = xe.GVConf.Xtray.PorxyFile
	xe.Conf.PortRange.Start = xe.GVConf.Xtray.PortRange.Start
	xe.Conf.PortRange.End = xe.GVConf.Xtray.PortRange.End
	xe.Conf.Port = xe.GVConf.Xtray.Port
	xe.Conf.TestUrl = xe.GVConf.Xtray.TestUrl
	xe.Conf.SwitchyOmegaUrl = xe.GVConf.Xtray.SwitchyOmegaUrl
	xe.Conf.GeoInfoUrl = xe.GVConf.Xtray.GeoInfoUrl
	// "@every 1h30m10s" https://pkg.go.dev/github.com/robfig/cron
	xe.Conf.VerifierCron = xe.GVConf.Xtray.VerifierCron
	xe.Conf.KeeperCron = xe.GVConf.Xtray.KeeperCron
	xe.Conf.StorageSqlitePath = xe.GVConf.Xtray.StorageSqlitePath
	xe.Conf.StorageExportPath = xe.GVConf.Xtray.StorageExportPath
	xe.Conf.SockFileDir = xe.GVConf.Xtray.SockFileDir

	xe.Runner = ctrl.NewXRunner(xe.Conf)
	xe.Runner.RegisterStarter(GenStarter())
	xe.Runner.RegisterKeeper(GenKeeper())
	xe.Keeper = ctrl.NewXKeeper(xe.Conf, xe.Runner)
	return xe
}
