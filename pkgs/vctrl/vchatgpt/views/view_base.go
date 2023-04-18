package views

import "github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"

type ViewBase struct {
	Model    vtui.IModel
	ViewName string
	Enabled  bool
}

func (that *ViewBase) SetModel(m vtui.IModel) {
	that.Model = m
}

func (that *ViewBase) Name() string {
	return that.ViewName
}

func (that *ViewBase) IsEnabled() bool {
	return that.Enabled
}

func (that *ViewBase) Disable() {
	that.Enabled = false
}

func (that *ViewBase) Enable() {
	that.Enabled = true
}

func (that *ViewBase) ExtraCmdHandlers() []vtui.CmdHandler {
	return []vtui.CmdHandler{}
}

func NewBase(name string) *ViewBase {
	return &ViewBase{
		ViewName: name,
	}
}
