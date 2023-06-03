package vctrl

import (
	"os/exec"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/neobox/pkgs/run"
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
	if that.conf.NeoBox.NeoConf != nil {
		that.runner = run.NewRunner(that.conf.NeoBox.NeoConf)
		run.SetNeoBoxEnvs(that.conf.NeoBox.NeoConf)
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
