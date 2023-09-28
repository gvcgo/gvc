package vctrl

import (
	"os"
	"os/exec"

	"github.com/moqsien/goutils/pkgs/logs"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/neobox/pkgs/run"
	"github.com/moqsien/neobox/pkgs/storage/model"
	nutils "github.com/moqsien/neobox/pkgs/utils"
)

type NeoBox struct {
	conf   *config.GVConfig
	runner *run.Runner
}

func NewBox(starter, keeperStarter *exec.Cmd) (n *NeoBox) {
	n = &NeoBox{
		conf: config.New(),
	}
	n.Initiate()
	n.registerStarter(starter)
	n.registerKeeperStarter(keeperStarter)
	return
}

func (that *NeoBox) Initiate() {
	if that.conf.NeoBox.NeoConf.LogDir != "" {
		os.MkdirAll(that.conf.NeoBox.NeoConf.LogDir, 0666)
	}
	if that.conf.NeoBox.NeoConf.GeoInfoDir != "" {
		os.MkdirAll(that.conf.NeoBox.NeoConf.GeoInfoDir, 0666)
	}
	if that.conf.NeoBox.NeoConf.SocketDir != "" {
		os.MkdirAll(that.conf.NeoBox.NeoConf.SocketDir, 0666)
	}
	if that.conf.NeoBox.NeoConf != nil {
		that.runner = run.NewRunner(that.conf.NeoBox.NeoConf)
		// set envs for neobox
		nutils.SetNeoboxEnvs(that.conf.NeoBox.NeoConf.GeoInfoDir, that.conf.NeoBox.NeoConf.SocketDir)
		// set logs
		logs.SetLogger(that.conf.NeoBox.NeoConf.LogDir)
		// init sqlitedb for neobox
		model.NewDBEngine(that.conf.NeoBox.NeoConf)
	}
}

func (that *NeoBox) registerStarter(cmd *exec.Cmd) {
	if that.runner != nil {
		that.runner.SetStarter(cmd)
	}
}

func (that *NeoBox) registerKeeperStarter(cmd *exec.Cmd) {
	if that.runner != nil {
		that.runner.SetKeeperStarter(cmd)
	}
}

func (that *NeoBox) StartShell() {
	if that.runner != nil {
		that.runner.OpenShell()
	}
}

func (that *NeoBox) StartClient() {
	if that.runner != nil {
		that.runner.Start()
	}
}

func (that *NeoBox) StartKeeper() {
	if that.runner != nil {
		that.runner.StartKeeper()
	}
}
