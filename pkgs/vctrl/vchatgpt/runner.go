package vchatgpt

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/views"
	"github.com/moqsien/gvc/pkgs/vctrl/vchatgpt/vtui"
)

type Runner struct {
	M *vtui.Model
}

func NewRunner() (r *Runner) {
	r = &Runner{
		M: vtui.NewModel(),
	}
	return
}

func (that *Runner) Run() {
	that.M.RegisterView(views.NewDefaultView())
	p := tea.NewProgram(that.M)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (that *Runner) RegisterView(v vtui.IView) {
	that.M.RegisterView(v)
}
