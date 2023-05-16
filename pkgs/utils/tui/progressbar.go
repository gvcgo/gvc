package tui

import (
	"fmt"

	"github.com/pterm/pterm"
)

const (
	BarLength int = 192
)

type ProgressBar struct {
	Bar             *pterm.ProgressbarPrinter
	Title           string
	ContentLength   int
	CurrentReceived int
	divideBy        int
}

func NewProgressBar(title string, length int) (p *ProgressBar) {
	if length <= 0 {
		return
	}
	p = &ProgressBar{
		ContentLength: length,
	}
	p.Bar = pterm.DefaultProgressbar.WithTitle(fmt.Sprintf("%s|<%sMB>", title, p.calcSize())).WithTotal(BarLength).WithShowCount(false).WithShowPercentage(true)
	p.divideBy = length / BarLength
	return
}

func (that *ProgressBar) calcSize() string {
	mb := float64(that.ContentLength) / float64(1024*1024)
	return fmt.Sprintf("%.2f", mb)
}

func (that *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	that.CurrentReceived += n
	var increasement int
	if that.CurrentReceived == that.ContentLength {
		increasement = that.Bar.Total - that.Bar.Current
	} else {
		increasement = (that.CurrentReceived / that.divideBy) - that.Bar.Current
	}
	if increasement > 0 {
		that.Bar.Add(increasement)
	}
	return
}

func (that *ProgressBar) Start() (err error) {
	bar := *that.Bar
	that.Bar, err = bar.Start()
	return
}
