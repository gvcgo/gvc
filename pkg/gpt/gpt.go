package gpt

import (
	"os"
	"path/filepath"

	"github.com/gvcgo/gogpt/pkgs/config"
	"github.com/gvcgo/gogpt/pkgs/gpt"
	gptui "github.com/gvcgo/gogpt/pkgs/tui"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gvc/conf"
	"github.com/postfinance/single"
)

/*
Chatgpt & flyTek xinghuo
*/
const (
	GPTWorkDirName string = "gpt"
)

func GetGPTWorkDir() string {
	return filepath.Join(conf.GetGVCWorkDir(), GPTWorkDirName)
}

type GPT struct {
	gptConf *config.Config
}

func NewGPT() *GPT {
	g := &GPT{
		gptConf: config.NewConf(GetGPTWorkDir()),
	}
	g.gptConf.Reload()
	return g
}

func (g *GPT) Run() {
	promptFilePath := filepath.Join(g.gptConf.GetWorkDir(), gpt.PromptFileName)
	if ok, _ := gutils.PathIsExist(promptFilePath); !ok {
		prompt := gpt.NewGPTPrompt(g.gptConf)
		prompt.DownloadPrompt()
	}

	lockFile, _ := single.New("chatgpt")
	if err := lockFile.Lock(); err != nil {
		gprint.PrintError("Another gogpt program is running: %s", lockFile.Lockfile())
		os.Exit(1)
	}
	defer func() {
		lockFile.Unlock()
	}()

	ui := gptui.NewGPTUI(g.gptConf)
	ui.Run()
}
