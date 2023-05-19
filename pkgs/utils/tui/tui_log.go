package tui

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pterm/pterm"
)

type TuiLog struct {
	content any
}

func NewLog(content any) *TuiLog {
	return &TuiLog{
		content: content,
	}
}

func (that *TuiLog) String() string {
	return gconv.String(that.content)
}

func (that *TuiLog) Debug() {
	pterm.EnableDebugMessages()
	pterm.Debug.Println(that.String())
}

func (that *TuiLog) Info() {
	pterm.Info.Println(that.String())
}

func (that *TuiLog) Success() {
	pterm.Success.Println(that.String())
}

func (that *TuiLog) Warning() {
	pterm.Warning.Println(that.String())
}

func (that *TuiLog) Error() {
	pterm.Error.Println(that.String())
}

func (that *TuiLog) Fatal() {
	pterm.Fatal.WithFatal(false).Println(that.String())
}

func PrintSuccess(info any) {
	l := NewLog(info)
	l.Success()
}

func PrintError(err any) {
	l := NewLog(err)
	l.Error()
}

func PrintInfo(info any) {
	l := NewLog(info)
	l.Info()
}

func PrintWarning(warn any) {
	l := NewLog(warn)
	l.Warning()
}
