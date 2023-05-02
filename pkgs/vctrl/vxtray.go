package vctrl

import (
	"os"
	"os/exec"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/xtray/pkgs/conf"
	"github.com/moqsien/xtray/pkgs/ctrl"
)

var cmdName string = func() string {
	epath, err := os.Executable()
	if err != nil {
		panic("cannot find executable path")
	}
	return epath
}()

const (
	XtrayStarterCmd = "xtray-runner"
	XtrayKeeperCmd  = "xtray-keeper"
)

var (
	Starter *exec.Cmd = exec.Command(cmdName, XtrayStarterCmd)
	Keeper  *exec.Cmd = exec.Command(cmdName, XtrayKeeperCmd)
)

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
	xe.Conf.PorxyFile = xe.GVConf.Xtray.PorxyFile
	xe.Conf.PortRange.Start = xe.GVConf.Xtray.PortRange.Start
	xe.Conf.PortRange.End = xe.GVConf.Xtray.PortRange.End
	xe.Conf.Port = xe.GVConf.Xtray.Port
	xe.Conf.TestUrl = xe.GVConf.Xtray.TestUrl
	xe.Conf.SwitchyOmegaUrl = xe.GVConf.Xtray.SwitchyOmegaUrl
	xe.Conf.GeoInfoUrl = xe.GVConf.Xtray.GeoInfoUrl
	// "@every 1h30m10s" https://pkg.go.dev/github.com/robfig/cron
	xe.Conf.VerifierCron = xe.GVConf.Xtray.VerifierCron
	xe.Conf.KeeperCron = xe.GVConf.Xtray.KeeperCron

	xe.Runner = ctrl.NewXRunner(xe.Conf)
	xe.Runner.RegisterStarter(Starter)
	xe.Runner.RegisterKeeper(Keeper)
	xe.Keeper = ctrl.NewXKeeper(xe.Conf, xe.Runner)
	return xe
}
