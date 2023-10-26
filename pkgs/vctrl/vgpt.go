package vctrl

import (
	"os"
	"path/filepath"

	gconf "github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	gptui "github.com/moqsien/gogpt/pkgs/tui"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gutils"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/postfinance/single"
)

/*
ChatGPT

IFlytek Spark
*/
type VGpt struct {
	Conf    *config.GVConfig
	GPTConf *gconf.Config
}

func NewVGPT() (vg *VGpt) {
	vg = &VGpt{
		Conf: config.New(),
	}
	vg.GPTConf = gconf.NewConf(vg.Conf.GPT.WorkDir)
	vg.GPTConf.Reload()
	return
}

func (that *VGpt) Run() {
	promptFilePath := filepath.Join(that.GPTConf.GetWorkDir(), gpt.PromptFileName)
	if ok, _ := gutils.PathIsExist(promptFilePath); !ok {
		prompt := gpt.NewGPTPrompt(that.GPTConf)
		prompt.DownloadPrompt()
	}

	// fmt.Printf("%+v", that.GPTConf.OpenAI)
	lockFile, _ := single.New("chatgpt")
	if err := lockFile.Lock(); err != nil {
		gprint.PrintError("Another gogpt program is running: %s", lockFile.Lockfile())
		os.Exit(1)
	}
	defer func() {
		lockFile.Unlock()
	}()

	ui := gptui.NewGPTUI(that.GPTConf)
	ui.Run()
}
