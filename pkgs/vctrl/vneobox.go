package vctrl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/logs"
	"github.com/moqsien/gshell/pkgs/ktrl"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
	neoconf "github.com/moqsien/neobox/pkgs/conf"
	"github.com/moqsien/neobox/pkgs/run"
	"github.com/moqsien/neobox/pkgs/storage/model"
	"github.com/spf13/cobra"
)

type NeoBox struct {
	conf   *config.GVConfig
	nconf  *neoconf.NeoConf
	runner *run.Runner
}

func NewBox(starter, keeperStarter *exec.Cmd) (n *NeoBox) {
	n = &NeoBox{
		conf:  config.New(),
		nconf: &neoconf.NeoConf{},
	}
	n.Initiate()
	n.registerStarter(starter)
	n.registerKeeperStarter(keeperStarter)
	return
}

func (that *NeoBox) Initiate() {
	// fixbugs： cannot use backup files for different platforms
	if ok, _ := utils.PathIsExist(that.nconf.WorkDir); !ok {
		that.conf.Reset()
	}
	that.nconf = that.conf.NeoBox.GetNeoConf()

	utils.MakeDirs(
		that.nconf.LogDir,
		that.nconf.GeoInfoDir,
		that.nconf.SocketDir,
	)

	that.runner = run.NewRunner(that.nconf)
	// set envs for neobox
	// nutils.SetNeoboxEnvs(that.nconf.GeoInfoDir, that.nconf.SocketDir)
	// set logs
	logs.SetLogger(that.nconf.LogDir)
	// init sqlitedb for neobox
	model.NewDBEngine(that.nconf)
}

func (that *NeoBox) registerStarter(cmd *exec.Cmd) {
	if that.runner != nil {
		that.runner.SetStarter(cmd)
	} else {
		gprint.PrintError("No runner found.")
	}
}

func (that *NeoBox) registerKeeperStarter(cmd *exec.Cmd) {
	if that.runner != nil {
		that.runner.SetKeeperStarter(cmd)
	} else {
		gprint.PrintError("No runner found.")
	}
}

func (that *NeoBox) StartShell() {
	if that.runner != nil {
		that.runner.OpenShell()
	} else {
		gprint.PrintError("No runner found.")
	}
}

func (that *NeoBox) StartClient() {
	if that.runner != nil {
		that.runner.Start()
	} else {
		gprint.PrintError("No runner found.")
	}
}

func (that *NeoBox) StartKeeper() {
	if that.runner != nil {
		that.runner.StartKeeper()
	} else {
		gprint.PrintError("No runner found.")
	}
}

const (
	NeoBoxCmdName          string = "neobox"
	NeoBoxAutoStartCmdName string = "autostart"
)

func (that *NeoBox) AutoStart(cmd *cobra.Command, args ...string) {
	if that.runner == nil {
		gprint.PrintError("No runner found.")
		return
	}
	sh := that.runner.GetShell()
	opts := []string{
		run.RestartUseDomain,
		run.RestartForceSingbox,
		run.RestartShowProxy,
		run.RestartShowConfig,
	}
	optStr := ""
	for _, o := range opts {
		if ok, _ := cmd.Flags().GetBool(o); ok {
			optStr += o
		}
	}
	ctx := &ktrl.KtrlContext{}
	ctx.SetArgs(args...)
	sh.Restart(ctx, optStr)
}

func (that *NeoBox) GenAutoStartScript() {
	scriptName := "neobox-autostart.sh"
	if runtime.GOOS == utils.Windows {
		scriptName = "neobox-autostart.bat"
	}
	scriptPath := filepath.Join(config.GVCDir, scriptName)
	binPath, _ := os.Executable()
	content := fmt.Sprintf("%s %s %s", binPath, NeoBoxCmdName, NeoBoxAutoStartCmdName)
	os.WriteFile(scriptPath, []byte(content), 0777)
	gprint.PrintInfo("Autostart script path: %s", scriptPath)
}

/*
synchronize local files to remote repo.
*/
func (that *NeoBox) HandleNeoboxConf(toDownload bool) {
	fPath := config.ProxyFilesDir
	remoteFileName := "neobox_confs.zip"
	repoSyncer := NewSynchronizer()
	if toDownload {
		// download and deploy.
		repoSyncer.DownloadFile(
			fPath,
			remoteFileName,
			EncryptByZip,
		)
	} else {
		repoSyncer.UploadFile(
			fPath,
			remoteFileName,
			EncryptByZip,
		)
	}
}
